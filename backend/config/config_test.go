package config

import "testing"

func TestLoadConfigDoesNotDefaultVNPayCredentials(t *testing.T) {
	t.Setenv("VNPAY_TMN_CODE", "")
	t.Setenv("VNPAY_HASH_SECRET", "")

	cfg := LoadConfig()
	if cfg.VNPayTmnCode != "" {
		t.Fatal("VNPay terminal code should require explicit configuration")
	}
	if cfg.VNPayHashSecret != "" {
		t.Fatal("VNPay hash secret should require explicit configuration")
	}
}
