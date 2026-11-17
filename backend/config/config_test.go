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
