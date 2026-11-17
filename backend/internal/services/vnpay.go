package services

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"math"
	"net"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/uit/vinhomes-management/config"
)

type VNPayClient struct {
	TmnCode    string
	HashSecret string
	PayURL     string
	ReturnURL  string
	IPNURL     string
}

func NewVNPayClient(cfg *config.Config) *VNPayClient {
	if cfg == nil {
		return &VNPayClient{
			PayURL: "https://sandbox.vnpayment.vn/paymentv2/vpcpay.html",
		}
	}
	return &VNPayClient{
		TmnCode:    cfg.VNPayTmnCode,
		HashSecret: cfg.VNPayHashSecret,
		PayURL:     cfg.VNPayPayURL,
		ReturnURL:  cfg.VNPayReturnURL,
		IPNURL:     cfg.VNPayIPNURL,
	}
}

func (v *VNPayClient) Configured() bool {
	return v.TmnCode != "" && v.HashSecret != ""
}

// phpURLEncode matches PHP urlencode() used by VNPay reference implementations.
func phpURLEncode(value string) string {
	encoded := url.QueryEscape(value)
	return strings.ReplaceAll(encoded, "%20", "+")
}

func (v *VNPayClient) buildHashData(params map[string]string) string {
	keys := make([]string, 0, len(params))
	for k, val := range params {
		if val != "" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		parts = append(parts, phpURLEncode(k)+"="+phpURLEncode(params[k]))
	}
	return strings.Join(parts, "&")
}

func (v *VNPayClient) sign(params map[string]string) string {
	hashData := v.buildHashData(params)
	mac := hmac.New(sha512.New, []byte(v.HashSecret))
	mac.Write([]byte(hashData))
	return hex.EncodeToString(mac.Sum(nil))
}

func copyVNPayParams(src map[string]string) map[string]string {
	dst := make(map[string]string, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

func (v *VNPayClient) VerifyCallback(params map[string]string) bool {
	receivedHash, ok := params["vnp_SecureHash"]
	if !ok {
		return false
	}
	verifyParams := copyVNPayParams(params)
	delete(verifyParams, "vnp_SecureHash")
	delete(verifyParams, "vnp_SecureHashType")
	expected := v.sign(verifyParams)
	return hmac.Equal([]byte(expected), []byte(receivedHash))
}

func (v *VNPayClient) BuildPaymentURL(txnRef, orderInfo, clientIP string, amountVND float64) (string, error) {
	if !v.Configured() {
		return "", fmt.Errorf("VNPay chưa được cấu hình — đăng ký sandbox tại https://sandbox.vnpayment.vn/devreg/ và cập nhật VNPAY_TMN_CODE, VNPAY_HASH_SECRET")
	}

	loc, err := time.LoadLocation("Asia/Ho_Chi_Minh")
	if err != nil {
		loc = time.FixedZone("ICT", 7*3600)
	}
	now := time.Now().In(loc)
	expire := now.Add(15 * time.Minute)

	amountInt := int64(math.Round(amountVND * 100))
	if amountInt <= 0 {
		return "", fmt.Errorf("số tiền thanh toán không hợp lệ")
	}

	params := map[string]string{
		"vnp_Version":    "2.1.0",
		"vnp_Command":    "pay",
		"vnp_TmnCode":    v.TmnCode,
		"vnp_Amount":     fmt.Sprintf("%d", amountInt),
		"vnp_CurrCode":   "VND",
		"vnp_TxnRef":     txnRef,
		"vnp_OrderInfo":  orderInfo,
		"vnp_OrderType":  "other",
		"vnp_Locale":     "vn",
		"vnp_ReturnUrl":  v.ReturnURL,
		"vnp_IpAddr":     clientIP,
		"vnp_CreateDate": now.Format("20060102150405"),
		"vnp_ExpireDate": expire.Format("20060102150405"),
	}

	hash := v.sign(params)
	query := v.buildHashData(params)
	return fmt.Sprintf("%s?%s&vnp_SecureHash=%s", v.PayURL, query, hash), nil
}

func ClientIPFromAddr(addr string) string {
	host := addr
	if parsedHost, _, err := net.SplitHostPort(addr); err == nil {
		host = parsedHost
	}

	host = strings.Trim(host, "[]")
	ip := net.ParseIP(host)
	if ip == nil || ip.IsLoopback() {
		return "127.0.0.1"
	}
	return ip.String()
}
