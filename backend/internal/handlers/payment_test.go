package handlers

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/uit/vinhomes-management/internal/models"
)

func TestCallbackQueryRejectsDuplicateParameters(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest("GET", "/return?vnp_TxnRef=one&vnp_TxnRef=two", nil)
	if params := queryParamsToMap(c); params != nil {
		t.Fatal("ambiguous callback was accepted")
	}
}
func TestPublicPaymentSummaryOmitsPrivateFields(t *testing.T) {
	x := &models.PaymentTransaction{ID: "transaction", UserID: "private-user", PaymentURL: "private-payment-url", Amount: 100000}
	summary := paymentSummary(x)
	for _, key := range []string{"user_id", "payment_url", "raw_callback"} {
		if _, ok := summary[key]; ok {
			t.Fatalf("private field %s exposed", key)
		}
	}
	if summary["id"] != "transaction" || summary["amount"] != float64(100000) {
		t.Fatal("invalid summary")
	}
}
