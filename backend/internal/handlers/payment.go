package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/uit/vinhomes-management/config"
	"github.com/uit/vinhomes-management/internal/models"
	"github.com/uit/vinhomes-management/internal/services"
)

func queryParamsToMap(c *gin.Context) map[string]string {
	m := make(map[string]string)
	for k, values := range c.Request.URL.Query() {
		if len(values) != 1 {
			return nil
		}
		if len(values) > 0 {
			m[k] = values[0]
		}
	}
	return m
}

func InitiateVNPayPayment(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		invoiceID := c.Param("id")
		userID, _ := c.Get("user_id")
		role, _ := c.Get("role")

		svc := services.NewPaymentService(cfg)
		result, err := svc.InitiateVNPay(
			invoiceID,
			userID.(string),
			role.(string),
			services.ClientIPFromAddr(c.ClientIP()),
		)
		if err != nil {
			status := http.StatusBadRequest
			switch err {
			case services.ErrInvoiceNotFound, services.ErrPaymentNotFound:
				status = http.StatusNotFound
			case services.ErrInvoiceAccessDenied:
				status = http.StatusForbidden
			case services.ErrVNPayNotConfigured:
				status = http.StatusServiceUnavailable
			case services.ErrPaymentInProgress:
				status = http.StatusConflict
			}
			c.JSON(status, gin.H{"error": err.Error(), "success": false})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": result, "success": true})
	}
}

func GetPaymentTransaction(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		role, _ := c.Get("role")

		svc := services.NewPaymentService(cfg)
		txn, err := svc.GetTransaction(c.Param("id"), userID.(string), role.(string))
		if err != nil {
			status := http.StatusNotFound
			if err == services.ErrInvoiceAccessDenied {
				status = http.StatusForbidden
			}
			c.JSON(status, gin.H{"error": err.Error(), "success": false})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": txn, "success": true})
	}
}

func VNPayReturnCallback(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		params := queryParamsToMap(c)
		svc := services.NewPaymentService(cfg)

		txn, err := svc.SyncReturnCallback(params)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":           err.Error(),
				"success":         false,
				"response_code":   params["vnp_ResponseCode"],
				"transaction_ref": params["vnp_TxnRef"],
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data":               paymentSummary(txn),
			"success":            txn.Status == "success",
			"verified":           true,
			"response_code":      params["vnp_ResponseCode"],
			"transaction_status": params["vnp_TransactionStatus"],
			"transaction_ref":    params["vnp_TxnRef"],
		})
	}
}

func VNPayIPNCallback(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		params := queryParamsToMap(c)
		svc := services.NewPaymentService(cfg)
		result := svc.ProcessVNPayIPN(params)
		c.JSON(http.StatusOK, result)
	}
}

// Return callbacks are public; expose only the verified transaction summary.
func paymentSummary(txn *models.PaymentTransaction) gin.H {
	return gin.H{"id": txn.ID, "invoice_id": txn.InvoiceID, "amount": txn.Amount, "status": txn.Status, "provider_txn_id": txn.ProviderTxnID, "requires_review": txn.RequiresReview}
}
