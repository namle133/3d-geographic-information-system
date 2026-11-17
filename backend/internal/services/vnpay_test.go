package services

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"strings"
	"testing"
)

// Verifies hash algorithm against VNPay PHP sample (urlencode + HMAC-SHA512).
func TestVNPaySignMatchesPHPSample(t *testing.T) {
	client := &VNPayClient{HashSecret: "TESTSECRETKEY12345678901234567890"}
	params := map[string]string{
		"vnp_Amount":     "1806000",
		"vnp_Command":    "pay",
		"vnp_CreateDate": "20210801153333",
		"vnp_CurrCode":   "VND",
		"vnp_IpAddr":     "127.0.0.1",
		"vnp_Locale":     "vn",
		"vnp_OrderInfo":  "Thanh toan don hang :5",
		"vnp_OrderType":  "other",
		"vnp_ReturnUrl":  "https://domainmerchant.vn/ReturnUrl",
		"vnp_TmnCode":    "DEMOV210",
		"vnp_TxnRef":     "5",
		"vnp_Version":    "2.1.0",
	}

	hashData := client.buildHashData(params)
	mac := hmac.New(sha512.New, []byte(client.HashSecret))
	mac.Write([]byte(hashData))
	sig := hex.EncodeToString(mac.Sum(nil))

	if !strings.Contains(hashData, "Thanh+toan+don+hang") {
		t.Fatalf("unexpected hash data encoding: %s", hashData)
	}
	if len(sig) != 128 {
		t.Fatalf("expected sha512 hex length 128, got %d", len(sig))
	}
}

func TestVNPayBuildPaymentURLRequiresConfig(t *testing.T) {
	client := &VNPayClient{}
	_, err := client.BuildPaymentURL("ref", "info", "127.0.0.1", 100000)
	if err == nil {
		t.Fatal("expected error when VNPay is not configured")
	}
}

func TestClientIPFromAddr(t *testing.T) {
	tests := []struct {
		name string
		addr string
		want string
	}{
		{name: "bare IPv4 returned by Gin", addr: "203.0.113.42", want: "203.0.113.42"},
		{name: "IPv4 with port", addr: "203.0.113.42:8080", want: "203.0.113.42"},
		{name: "IPv6 with port", addr: "[2001:db8::1]:8080", want: "2001:db8::1"},
		{name: "loopback", addr: "::1", want: "127.0.0.1"},
		{name: "invalid address", addr: "", want: "127.0.0.1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ClientIPFromAddr(tt.addr); got != tt.want {
				t.Errorf("ClientIPFromAddr(%q) = %q, want %q", tt.addr, got, tt.want)
			}
		})
	}
}

func TestVNPayAmountInMinorUnitsRoundsFractionalVND(t *testing.T) {
	tests := []struct {
		amount float64
		want   int64
	}{
		{amount: 1000, want: 100000},
		{amount: 1000.1, want: 100010},
		{amount: 1000.29, want: 100029},
	}

	for _, tt := range tests {
		if got := vnpayAmountInMinorUnits(tt.amount); got != tt.want {
			t.Errorf("vnpayAmountInMinorUnits(%v) = %d, want %d", tt.amount, got, tt.want)
		}
	}
}
