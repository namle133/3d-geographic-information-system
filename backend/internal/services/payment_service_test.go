package services

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/uit/vinhomes-management/config"
	"github.com/uit/vinhomes-management/internal/database"
	"github.com/uit/vinhomes-management/internal/models"
	"github.com/uit/vinhomes-management/internal/testutil"
	"gorm.io/gorm"
)

func paymentFixture(t *testing.T) (*PaymentService, models.PaymentTransaction) {
	t.Helper()
	database.DB = testutil.OpenDB(t)
	inv := testutil.Invoice(t, database.DB, 100000)
	expiry := time.Now().Add(vnPayPaymentTTL)
	txn := models.PaymentTransaction{ID: uuid.NewString(), InvoiceID: inv.ID, UserID: inv.ID, Amount: inv.TotalAmount, Provider: "vnpay", Status: "pending", CreatedAt: time.Now(), ExpiresAt: &expiry}
	if err := database.DB.Create(&txn).Error; err != nil {
		t.Fatal(err)
	}
	svc := NewPaymentService(&config.Config{VNPayTmnCode: "AUDIT123", VNPayHashSecret: "audit-secret", VNPayPayURL: "https://sandbox.vnpayment.vn/paymentv2/vpcpay.html", VNPayReturnURL: "https://example.com/result"})
	return svc, txn
}
func callback(s *PaymentService, txn models.PaymentTransaction, response, status string) map[string]string {
	p := map[string]string{"vnp_TmnCode": "AUDIT123", "vnp_TxnRef": txn.ID, "vnp_Amount": fmt.Sprint(vnpayAmountInMinorUnits(txn.Amount)), "vnp_ResponseCode": response, "vnp_TransactionStatus": status, "vnp_TransactionNo": "12345678"}
	p["vnp_SecureHash"] = s.vnpay.sign(p)
	return p
}
func paymentState(t *testing.T, txn models.PaymentTransaction) (models.Invoice, models.PaymentTransaction, int64) {
	t.Helper()
	var i models.Invoice
	var x models.PaymentTransaction
	var n int64
	if err := database.DB.First(&i, "id = ?", txn.InvoiceID).Error; err != nil {
		t.Fatal(err)
	}
	if err := database.DB.First(&x, "id = ?", txn.ID).Error; err != nil {
		t.Fatal(err)
	}
	if err := database.DB.Model(&models.InvoicePayment{}).Where("invoice_id = ?", txn.InvoiceID).Count(&n).Error; err != nil {
		t.Fatal(err)
	}
	return i, x, n
}
func TestIPNConcurrentCallbacksCreditExactlyOnce(t *testing.T) {
	s, txn := paymentFixture(t)
	var arrived atomic.Int32
	var ready sync.WaitGroup
	ready.Add(2)
	// Force both requests to finish the initial lookup before either can commit.
	database.DB.Callback().Query().After("gorm:query").Register("test:barrier", func(tx *gorm.DB) {
		if x, ok := tx.Statement.Dest.(*models.PaymentTransaction); ok && x.ID == txn.ID && arrived.Add(1) <= 2 {
			ready.Done()
			ready.Wait()
		}
	})
	var wg sync.WaitGroup
	wg.Add(2)
	results := make(chan string, 2)
	for i := 0; i < 2; i++ {
		go func() { defer wg.Done(); results <- s.ProcessVNPayIPN(callback(s, txn, "00", "00")).RspCode }()
	}
	wg.Wait()
	database.DB.Callback().Query().Remove("test:barrier")
	close(results)
	got := map[string]int{}
	for r := range results {
		got[r]++
	}
	if got["00"] != 1 || got["02"] != 1 {
		t.Fatalf("callbacks: %v", got)
	}
	invoice, x, n := paymentState(t, txn)
	if invoice.PaidAmount != 100000 || n != 1 || x.Status != "success" {
		t.Fatalf("invoice=%+v txn=%s rows=%d", invoice, x.Status, n)
	}
}
func TestIPNRejectsInvalidCallbacksWithoutCrediting(t *testing.T) {
	for _, tc := range []struct{ name, response, status, amount, want string }{
		{"failed transaction", "00", "02", "10000000", "00"},
		{"cancelled", "24", "02", "10000000", "00"},
		{"wrong amount", "00", "00", "10000001", "04"},
		{"trailing amount data", "00", "00", "10000000x", "04"},
		{"missing status", "00", "", "10000000", "99"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			s, txn := paymentFixture(t)
			p := callback(s, txn, tc.response, tc.status)
			p["vnp_Amount"] = tc.amount
			delete(p, "vnp_SecureHash")
			p["vnp_SecureHash"] = s.vnpay.sign(p)
			if r := s.ProcessVNPayIPN(p); r.RspCode != tc.want {
				t.Fatalf("%+v", r)
			}
			i, _, n := paymentState(t, txn)
			if i.PaidAmount != 0 || n != 0 {
				t.Fatalf("credited rejected callback: %v %d", i.PaidAmount, n)
			}
		})
	}
}
func TestReturnDoesNotRecordPaymentAndChecksAmount(t *testing.T) {
	s, txn := paymentFixture(t)
	p := callback(s, txn, "00", "00")
	x, err := s.SyncReturnCallback(p)
	if err != nil || x.Status != "pending" {
		t.Fatalf("return: %v %v", x, err)
	}
	i, _, n := paymentState(t, txn)
	if i.PaidAmount != 0 || n != 0 {
		t.Fatal("return recorded payment")
	}
	p["vnp_Amount"] = "1"
	delete(p, "vnp_SecureHash")
	p["vnp_SecureHash"] = s.vnpay.sign(p)
	if _, err := s.SyncReturnCallback(p); !errors.Is(err, ErrInvalidAmount) {
		t.Fatalf("amount: %v", err)
	}
	if r := s.ProcessVNPayIPN(callback(s, txn, "00", "00")); r.RspCode != "00" {
		t.Fatal(r)
	}
	x, err = s.SyncReturnCallback(callback(s, txn, "00", "00"))
	if err != nil || x.Status != "success" {
		t.Fatalf("confirmed: %v %v", x, err)
	}
}
func TestIPNDatabaseFailuresRollBackAndRequestRetry(t *testing.T) {
	for _, failed := range []bool{false, true} {
		t.Run(fmt.Sprint(failed), func(t *testing.T) {
			s, txn := paymentFixture(t)
			response, status := "00", "00"
			query := "ALTER TABLE invoice_payments ADD CHECK (amount <> 100000)"
			if failed {
				query = "ALTER TABLE payment_transactions ADD CHECK (status <> 'failed')"
				response, status = "24", "02"
			}
			if err := database.DB.Exec(query).Error; err != nil {
				t.Fatal(err)
			}
			if r := s.ProcessVNPayIPN(callback(s, txn, response, status)); r.RspCode != "99" {
				t.Fatalf("%+v", r)
			}
			i, x, n := paymentState(t, txn)
			if i.PaidAmount != 0 || n != 0 || x.Status != "pending" {
				t.Fatal("partial commit")
			}
		})
	}
}
func TestInitiateReusesPendingAndPreventsConcurrentOrders(t *testing.T) {
	s, txn := paymentFixture(t)
	var wg sync.WaitGroup
	wg.Add(2)
	results := make(chan *InitiateVNPayResult, 2)
	errs := make(chan error, 2)
	for i := 0; i < 2; i++ {
		go func() {
			defer wg.Done()
			r, e := s.InitiateVNPay(txn.InvoiceID, txn.UserID, "resident", "203.0.113.10")
			results <- r
			errs <- e
		}()
	}
	wg.Wait()
	close(results)
	close(errs)
	for err := range errs {
		if err != nil {
			t.Fatal(err)
		}
	}
	for r := range results {
		if r.TransactionID != txn.ID {
			t.Fatal("did not reuse pending")
		}
	}
	if _, err := s.InitiateVNPay(txn.InvoiceID, "another-user", "admin", "203.0.113.11"); !errors.Is(err, ErrPaymentInProgress) {
		t.Fatal(err)
	}
	if _, err := s.ApplyPaymentToInvoice(txn.InvoiceID, 100000, "cash", "", "staff"); !errors.Is(err, ErrPaymentInProgress) {
		t.Fatal(err)
	}
	// Expired links no longer prevent a new order; a paid invoice always does.
	database.DB.Model(&txn).Update("expires_at", time.Now().Add(-time.Minute))
	r, err := s.InitiateVNPay(txn.InvoiceID, txn.UserID, "resident", "203.0.113.10")
	if err != nil || r.TransactionID == txn.ID {
		t.Fatalf("expired: %v %v", r, err)
	}
}
func TestIPNOverpaymentIsRecordedForReview(t *testing.T) {
	s, txn := paymentFixture(t)
	other := txn
	other.ID = uuid.NewString()
	if err := database.DB.Create(&other).Error; err != nil {
		t.Fatal(err)
	}
	for _, x := range []models.PaymentTransaction{txn, other} {
		if r := s.ProcessVNPayIPN(callback(s, x, "00", "00")); r.RspCode != "00" {
			t.Fatal(r)
		}
	}
	i, x, n := paymentState(t, other)
	if i.PaidAmount != 200000 || n != 2 || !x.RequiresReview {
		t.Fatal("captured overpayment must be kept and flagged")
	}
}
func TestManualPaymentRollsBackOnInvoiceFailure(t *testing.T) {
	_, txn := paymentFixture(t)
	database.DB.Delete(&txn)
	if err := database.DB.Exec("ALTER TABLE invoices ADD CHECK (paid_amount = 0)").Error; err != nil {
		t.Fatal(err)
	}
	s := NewPaymentService(nil)
	if _, err := s.ApplyPaymentToInvoice(txn.InvoiceID, 100000, "cash", "", "staff"); err == nil {
		t.Fatal("expected failure")
	}
	var n int64
	database.DB.Model(&models.InvoicePayment{}).Count(&n)
	if n != 0 {
		t.Fatal("orphan payment")
	}
}

func TestConcurrentInitiationCreatesOnlyOneNewTransaction(t *testing.T) {
	s, txn := paymentFixture(t)
	if err := database.DB.Delete(&txn).Error; err != nil {
		t.Fatal(err)
	}
	var wg sync.WaitGroup
	wg.Add(2)
	ids := make(chan string, 2)
	errs := make(chan error, 2)
	for i := 0; i < 2; i++ {
		go func() {
			defer wg.Done()
			r, e := s.InitiateVNPay(txn.InvoiceID, txn.UserID, "resident", "203.0.113.10")
			errs <- e
			if r != nil {
				ids <- r.TransactionID
			}
		}()
	}
	wg.Wait()
	close(errs)
	close(ids)
	for e := range errs {
		if e != nil {
			t.Fatal(e)
		}
	}
	unique := map[string]bool{}
	for id := range ids {
		unique[id] = true
	}
	if len(unique) != 1 {
		t.Fatalf("created different orders: %v", unique)
	}
	var count int64
	database.DB.Model(&models.PaymentTransaction{}).Count(&count)
	if count != 1 {
		t.Fatal(count)
	}
}
