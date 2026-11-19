package config

import (
	"os"
	"testing"
)

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

func TestGetVNPayEnvTrimsWhitespaceAndMatchingQuotes(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  string
	}{
		{name: "plain", value: "KD4N9ECY", want: "KD4N9ECY"},
		{name: "double quoted", value: `"KD4N9ECY"`, want: "KD4N9ECY"},
		{name: "single quoted", value: `'KD4N9ECY'`, want: "KD4N9ECY"},
		{name: "surrounding spaces", value: "  KD4N9ECY  ", want: "KD4N9ECY"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("TEST_VNPAY_VALUE", tt.value)
			if got := getVNPayEnv("TEST_VNPAY_VALUE"); got != tt.want {
				t.Fatalf("getVNPayEnv() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestVNPayCheckoutDefaultsToQRAndAllowsOverride(t *testing.T) {
	t.Chdir(t.TempDir())
	t.Setenv("VNPAY_BANK_CODE", "")
	if err := os.Unsetenv("VNPAY_BANK_CODE"); err != nil {
		t.Fatal(err)
	}
	if got := LoadConfig().VNPayBankCode; got != "VNPAYQR" {
		t.Fatalf("default checkout channel: %q", got)
	}
	t.Setenv("VNPAY_BANK_CODE", "VNBANK")
	if got := LoadConfig().VNPayBankCode; got != "VNBANK" {
		t.Fatalf("override checkout channel: %q", got)
	}
	t.Setenv("VNPAY_BANK_CODE", "")
	if got := LoadConfig().VNPayBankCode; got != "" {
		t.Fatalf("empty setting must show all methods: %q", got)
	}
}
