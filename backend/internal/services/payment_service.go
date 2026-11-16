package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/uit/vinhomes-management/config"
	"github.com/uit/vinhomes-management/internal/database"
	"github.com/uit/vinhomes-management/internal/models"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

var (
	ErrInvoiceNotFound      = errors.New("invoice not found")
	ErrInvoiceNotPayable    = errors.New("invoice is not payable")
	ErrInvalidAmount        = errors.New("invalid payment amount")
	ErrPaymentNotFound      = errors.New("payment transaction not found")
	ErrPaymentAlreadyDone   = errors.New("payment already completed")
	ErrVNPayNotConfigured   = errors.New("vnpay not configured")
	ErrInvoiceAccessDenied  = errors.New("access denied to invoice")
)

type PaymentService struct {
	vnpay *VNPayClient
}

func NewPaymentService(cfg *config.Config) *PaymentService {
	return &PaymentService{vnpay: NewVNPayClient(cfg)}
}

func outstandingAmount(invoice *models.Invoice) float64 {
	outstanding := invoice.TotalAmount - invoice.PaidAmount
	if outstanding < 0 {
		return 0
	}
	return outstanding
}

func (s *PaymentService) CanAccessInvoice(invoiceID, userID, role string) (*models.Invoice, error) {
	var invoice models.Invoice
	query := database.DB.Preload("Contract.Resident")
	if role == "resident" {
		var resident models.Resident
		if err := database.DB.Where("user_id = ?", userID).First(&resident).Error; err != nil {
			return nil, ErrInvoiceAccessDenied
		}
		err := query.Joins("JOIN contracts ON contracts.id = invoices.contract_id").
			Where("invoices.id = ? AND contracts.resident_id = ?", invoiceID, resident.ID).
			First(&invoice).Error
		if err != nil {
			return nil, ErrInvoiceAccessDenied
		}
		return &invoice, nil
	}

	if err := query.First(&invoice, "id = ?", invoiceID).Error; err != nil {
		return nil, ErrInvoiceNotFound
	}
	return &invoice, nil
}

func (s *PaymentService) ApplyPaymentToInvoice(invoiceID string, amount float64, method, reference, recordedBy string) (*models.InvoicePayment, error) {
	if amount <= 0 {
		return nil, ErrInvalidAmount
	}

	var invoice models.Invoice
	if err := database.DB.First(&invoice, "id = ?", invoiceID).Error; err != nil {
		return nil, ErrInvoiceNotFound
	}

	if invoice.Status == "paid" || invoice.Status == "cancelled" {
		return nil, ErrInvoiceNotPayable
	}

	payment := models.InvoicePayment{
		ID:         uuid.New().String(),
		InvoiceID:  invoiceID,
		Amount:     amount,
		Method:     method,
		Reference:  reference,
		PaidAt:     time.Now(),
		RecordedBy: recordedBy,
	}

	if err := database.DB.Create(&payment).Error; err != nil {
		return nil, err
	}

	invoice.PaidAmount += amount
	if invoice.PaidAmount >= invoice.TotalAmount {
		invoice.Status = "paid"
	} else if invoice.Status == "overdue" || invoice.Status == "pending" {
		// keep current status until fully paid
	}

	if err := database.DB.Save(&invoice).Error; err != nil {
		return nil, err
	}

	NewNotificationService().NotifyPaymentReceived(&invoice, amount)
	return &payment, nil
}

type InitiateVNPayResult struct {
	TransactionID string  `json:"transaction_id"`
	PaymentURL    string  `json:"payment_url"`
	Amount        float64 `json:"amount"`
	InvoiceID     string  `json:"invoice_id"`
}

func (s *PaymentService) InitiateVNPay(invoiceID, userID, role, clientIP string) (*InitiateVNPayResult, error) {
	if !s.vnpay.Configured() {
		return nil, ErrVNPayNotConfigured
	}

	invoice, err := s.CanAccessInvoice(invoiceID, userID, role)
	if err != nil {
		return nil, err
	}

	if invoice.Status == "paid" || invoice.Status == "cancelled" {
		return nil, ErrInvoiceNotPayable
	}

	amount := outstandingAmount(invoice)
	if amount <= 0 {
		return nil, ErrInvoiceNotPayable
	}

	txnID := strings.ReplaceAll(uuid.New().String(), "-", "")[:16]
	orderInfo := fmt.Sprintf("Thanh toan hoa don %s", invoice.InvoiceCode)

	paymentURL, err := s.vnpay.BuildPaymentURL(txnID, orderInfo, clientIP, amount)
	if err != nil {
		return nil, err
	}

	txn := models.PaymentTransaction{
		ID:         txnID,
		InvoiceID:  invoiceID,
		UserID:     userID,
		Amount:     amount,
		Provider:   "vnpay",
		Status:     "pending",
		PaymentURL: paymentURL,
		CreatedAt:  time.Now(),
	}

	if err := database.DB.Create(&txn).Error; err != nil {
		return nil, err
	}

	return &InitiateVNPayResult{
		TransactionID: txn.ID,
		PaymentURL:    paymentURL,
		Amount:        amount,
		InvoiceID:     invoiceID,
	}, nil
}

type VNPayIPNResult struct {
	RspCode string `json:"RspCode"`
	Message string `json:"Message"`
}

func (s *PaymentService) ProcessVNPayIPN(params map[string]string) VNPayIPNResult {
	if !s.vnpay.VerifyCallback(copyParams(params)) {
		return VNPayIPNResult{RspCode: "97", Message: "Invalid signature"}
	}

	txnRef := params["vnp_TxnRef"]
	if txnRef == "" {
		return VNPayIPNResult{RspCode: "01", Message: "Order not found"}
	}

	var txn models.PaymentTransaction
	if err := database.DB.First(&txn, "id = ?", txnRef).Error; err != nil {
		return VNPayIPNResult{RspCode: "01", Message: "Order not found"}
	}

	if txn.Status == "success" {
		return VNPayIPNResult{RspCode: "02", Message: "Order already confirmed"}
	}

	raw, _ := json.Marshal(params)
	txn.RawCallback = datatypes.JSON(raw)

	responseCode := params["vnp_ResponseCode"]
	if responseCode != "00" {
		txn.Status = "failed"
		now := time.Now()
		txn.CompletedAt = &now
		database.DB.Save(&txn)
		return VNPayIPNResult{RspCode: "00", Message: "Confirm Success"}
	}

	// Validate amount (VNPay sends amount * 100)
	var vnpAmount int64
	fmt.Sscanf(params["vnp_Amount"], "%d", &vnpAmount)
	expectedAmount := int64(txn.Amount * 100)
	if vnpAmount != expectedAmount {
		log.Printf("VNPay amount mismatch: got %d want %d for txn %s", vnpAmount, expectedAmount, txnRef)
		return VNPayIPNResult{RspCode: "04", Message: "Invalid amount"}
	}

	txn.ProviderTxnID = params["vnp_TransactionNo"]
	txn.Status = "success"
	now := time.Now()
	txn.CompletedAt = &now

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&txn).Error; err != nil {
			return err
		}

		payment := models.InvoicePayment{
			ID:         uuid.New().String(),
			InvoiceID:  txn.InvoiceID,
			Amount:     txn.Amount,
			Method:     "vnpay",
			Reference:  txn.ProviderTxnID,
			PaidAt:     now,
			RecordedBy: "system",
		}
		if err := tx.Create(&payment).Error; err != nil {
			return err
		}

		var invoice models.Invoice
		if err := tx.First(&invoice, "id = ?", txn.InvoiceID).Error; err != nil {
			return err
		}
		invoice.PaidAmount += txn.Amount
		if invoice.PaidAmount >= invoice.TotalAmount {
			invoice.Status = "paid"
		}
		return tx.Save(&invoice).Error
	})

	if err != nil {
		log.Printf("VNPay IPN processing error: %v", err)
		return VNPayIPNResult{RspCode: "99", Message: "Unknown error"}
	}

	var invoice models.Invoice
	if err := database.DB.First(&invoice, "id = ?", txn.InvoiceID).Error; err == nil {
		NewNotificationService().NotifyPaymentReceived(&invoice, txn.Amount)
	}

	return VNPayIPNResult{RspCode: "00", Message: "Confirm Success"}
}

func copyParams(src map[string]string) map[string]string {
	dst := make(map[string]string, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

func (s *PaymentService) GetTransaction(id, userID, role string) (*models.PaymentTransaction, error) {
	var txn models.PaymentTransaction
	if err := database.DB.First(&txn, "id = ?", id).Error; err != nil {
		return nil, ErrPaymentNotFound
	}

	if role == "resident" && txn.UserID != userID {
		return nil, ErrInvoiceAccessDenied
	}

	return &txn, nil
}

func (s *PaymentService) SyncReturnCallback(params map[string]string) (*models.PaymentTransaction, error) {
	if !s.vnpay.VerifyCallback(copyParams(params)) {
		return nil, fmt.Errorf("invalid signature")
	}

	txnRef := params["vnp_TxnRef"]
	var txn models.PaymentTransaction
	if err := database.DB.First(&txn, "id = ?", txnRef).Error; err != nil {
		return nil, ErrPaymentNotFound
	}

	if txn.Status == "success" {
		return &txn, nil
	}

	responseCode := params["vnp_ResponseCode"]
	if responseCode == "00" {
		result := s.ProcessVNPayIPN(params)
		if result.RspCode != "00" && result.RspCode != "02" {
			return nil, fmt.Errorf("%s", result.Message)
		}
		_ = database.DB.First(&txn, "id = ?", txnRef)
		return &txn, nil
	}

	if txn.Status == "pending" {
		txn.Status = "failed"
		now := time.Now()
		txn.CompletedAt = &now
		raw, _ := json.Marshal(params)
		txn.RawCallback = datatypes.JSON(raw)
		database.DB.Save(&txn)
	}

	return &txn, nil
}
