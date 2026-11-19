package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/uit/vinhomes-management/config"
	"github.com/uit/vinhomes-management/internal/database"
	"github.com/uit/vinhomes-management/internal/models"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	ErrInvoiceNotFound     = errors.New("invoice not found")
	ErrInvoiceNotPayable   = errors.New("invoice is not payable")
	ErrInvalidAmount       = errors.New("invalid payment amount")
	ErrPaymentNotFound     = errors.New("payment transaction not found")
	ErrPaymentAlreadyDone  = errors.New("payment already completed")
	ErrVNPayNotConfigured  = errors.New("vnpay not configured")
	ErrInvoiceAccessDenied = errors.New("access denied to invoice")
	ErrPaymentInProgress   = errors.New("Hóa đơn đang có giao dịch VNPay chờ xử lý. Vui lòng đợi giao dịch hoàn tất hoặc hết hạn trước khi thanh toán lại")
	ErrInvalidCallback     = errors.New("invalid VNPay callback")
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
	if !validPaymentAmount(amount) {
		return nil, ErrInvalidAmount
	}
	var invoice models.Invoice
	payment := models.InvoicePayment{ID: uuid.NewString(), InvoiceID: invoiceID, Amount: amount, Method: method, Reference: reference, PaidAt: time.Now(), RecordedBy: recordedBy}
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&invoice, "id = ?", invoiceID).Error; err != nil {
			return err
		}
		if invoice.Status == "paid" || invoice.Status == "cancelled" {
			return ErrInvoiceNotPayable
		}
		if amount > outstandingAmount(&invoice) {
			return ErrInvalidAmount
		}
		var pending int64
		if err := activeVNPayTransactions(tx, invoiceID, time.Now()).Count(&pending).Error; err != nil {
			return err
		}
		if pending > 0 {
			return ErrPaymentInProgress
		}
		if err := tx.Create(&payment).Error; err != nil {
			return err
		}
		invoice.PaidAmount += amount
		if invoice.PaidAmount >= invoice.TotalAmount {
			invoice.Status = "paid"
		}
		return tx.Omit(clause.Associations).Save(&invoice).Error
	})
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrInvoiceNotFound
	}
	if err != nil {
		return nil, err
	}
	NewNotificationService().NotifyPaymentReceived(&invoice, amount)
	return &payment, nil
}

func activeVNPayTransactions(tx *gorm.DB, invoiceID string, now time.Time) *gorm.DB {
	return tx.Model(&models.PaymentTransaction{}).Where("invoice_id = ? AND provider = ? AND status = ?", invoiceID, "vnpay", "pending").
		Where("(expires_at > ?) OR (expires_at IS NULL AND created_at > ?)", now, now.Add(-vnPayPaymentTTL))
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
	if _, err := s.CanAccessInvoice(invoiceID, userID, role); err != nil {
		return nil, err
	}
	var txn models.PaymentTransaction
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		var invoice models.Invoice
		// All payment writers lock the invoice before the transaction to avoid deadlocks.
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&invoice, "id = ?", invoiceID).Error; err != nil {
			return err
		}
		if invoice.Status == "paid" || invoice.Status == "cancelled" || outstandingAmount(&invoice) <= 0 {
			return ErrInvoiceNotPayable
		}
		now := time.Now()
		err := activeVNPayTransactions(tx, invoiceID, now).Order("created_at DESC").First(&txn).Error
		if err == nil {
			if txn.UserID != userID || vnpayAmountInMinorUnits(txn.Amount) != vnpayAmountInMinorUnits(outstandingAmount(&invoice)) {
				return ErrPaymentInProgress
			}
			return nil
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		txn = models.PaymentTransaction{}
		amount := float64(vnpayAmountInMinorUnits(outstandingAmount(&invoice))) / 100
		txnID := strings.ReplaceAll(uuid.NewString(), "-", "")
		paymentURL, err := s.vnpay.BuildPaymentURL(txnID, fmt.Sprintf("Thanh toan hoa don %s", invoice.InvoiceCode), clientIP, amount)
		if err != nil {
			return err
		}
		expires := now.Add(vnPayPaymentTTL).Truncate(time.Second)
		txn = models.PaymentTransaction{ID: txnID, InvoiceID: invoiceID, UserID: userID, Amount: amount, Provider: "vnpay", Status: "pending", PaymentURL: paymentURL, CreatedAt: now, ExpiresAt: &expires}
		return tx.Create(&txn).Error
	})
	if err != nil {
		return nil, err
	}
	return &InitiateVNPayResult{TransactionID: txn.ID, PaymentURL: txn.PaymentURL, Amount: txn.Amount, InvoiceID: txn.InvoiceID}, nil
}

type VNPayIPNResult struct {
	RspCode string `json:"RspCode"`
	Message string `json:"Message"`
}

func (s *PaymentService) ProcessVNPayIPN(params map[string]string) VNPayIPNResult {
	if !s.vnpay.VerifyCallback(params) {
		return VNPayIPNResult{"97", "Invalid signature"}
	}
	var initial models.PaymentTransaction
	if err := database.DB.First(&initial, "id = ? AND provider = ?", params["vnp_TxnRef"], "vnpay").Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return VNPayIPNResult{"01", "Order not found"}
		}
		return VNPayIPNResult{"99", "Database error"}
	}
	result := VNPayIPNResult{"00", "Confirm Success"}
	var credited bool
	var invoice models.Invoice
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&invoice, "id = ?", initial.InvoiceID).Error; err != nil {
			return err
		}
		var txn models.PaymentTransaction
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&txn, "id = ?", initial.ID).Error; err != nil {
			return err
		}
		if err := validateVNPayAmount(params["vnp_Amount"], txn.Amount); err != nil {
			result = VNPayIPNResult{"04", "Invalid amount"}
			return nil
		}
		if txn.Status == "success" || txn.Status == "failed" {
			result = VNPayIPNResult{"02", "Order already confirmed"}
			return nil
		}
		// Missing or malformed status must never finalize the order.
		if len(params["vnp_ResponseCode"]) != 2 || len(params["vnp_TransactionStatus"]) != 2 {
			return ErrInvalidCallback
		}
		raw, err := json.Marshal(params)
		if err != nil {
			return err
		}
		txn.RawCallback = datatypes.JSON(raw)
		now := time.Now()
		txn.CompletedAt = &now
		if params["vnp_ResponseCode"] != "00" || params["vnp_TransactionStatus"] != "00" {
			txn.Status = "failed"
			return tx.Omit(clause.Associations).Save(&txn).Error
		}
		if params["vnp_TransactionNo"] == "" || params["vnp_TransactionNo"] == "0" {
			return ErrInvalidCallback
		}
		txn.ProviderTxnID = params["vnp_TransactionNo"]
		txn.Status = "success"
		txn.RequiresReview = invoice.Status == "cancelled" || txn.Amount > outstandingAmount(&invoice)
		payment := models.InvoicePayment{ID: uuid.NewString(), PaymentTransactionID: &txn.ID, InvoiceID: txn.InvoiceID, Amount: txn.Amount, Method: "vnpay", Reference: txn.ProviderTxnID, PaidAt: now, RecordedBy: "system"}
		if err := tx.Create(&payment).Error; err != nil {
			return err
		}
		invoice.PaidAmount += txn.Amount
		if invoice.Status != "cancelled" && invoice.PaidAmount >= invoice.TotalAmount {
			invoice.Status = "paid"
		}
		if err := tx.Omit(clause.Associations).Save(&invoice).Error; err != nil {
			return err
		}
		if err := tx.Omit(clause.Associations).Save(&txn).Error; err != nil {
			return err
		}
		credited = true
		return nil
	})
	if err != nil {
		log.Printf("VNPay IPN could not commit transaction %s", initial.ID)
		return VNPayIPNResult{"99", "Unable to confirm transaction"}
	}
	if credited {
		NewNotificationService().NotifyPaymentReceived(&invoice, initial.Amount)
	}
	return result
}

func validateVNPayAmount(raw string, expected float64) error {
	if raw == "" || len(raw) > 12 || !validPaymentAmount(expected) {
		return ErrInvalidAmount
	}
	for _, c := range raw {
		if c < '0' || c > '9' {
			return ErrInvalidAmount
		}
	}
	amount, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || amount != vnpayAmountInMinorUnits(expected) {
		return ErrInvalidAmount
	}
	return nil
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

// SyncReturnCallback verifies the return URL and reads the status committed by IPN.
// Visiting or reloading a browser URL never records a payment.
func (s *PaymentService) SyncReturnCallback(params map[string]string) (*models.PaymentTransaction, error) {
	if !s.vnpay.VerifyCallback(params) {
		return nil, ErrInvalidCallback
	}
	var txn models.PaymentTransaction
	if err := database.DB.First(&txn, "id = ? AND provider = ?", params["vnp_TxnRef"], "vnpay").Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPaymentNotFound
		}
		return nil, err
	}
	if err := validateVNPayAmount(params["vnp_Amount"], txn.Amount); err != nil {
		return nil, err
	}
	if len(params["vnp_ResponseCode"]) != 2 || len(params["vnp_TransactionStatus"]) != 2 {
		return nil, ErrInvalidCallback
	}
	return &txn, nil
}
