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
	BankCode   string
}

const vnPayPaymentTTL = 15 * time.Minute

func validPaymentAmount(amount float64) bool {
	return !math.IsNaN(amount) && !math.IsInf(amount, 0) && amount > 0 && amount*100 <= 999999999999
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
		BankCode:   cfg.VNPayBankCode,
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

func vnpayAmountInMinorUnits(amountVND float64) int64 {
	return int64(math.Round(amountVND * 100))
}

func copyVNPayParams(src map[string]string) map[string]string {
	dst := make(map[string]string, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

func (v *VNPayClient) VerifyCallback(params map[string]string) bool {
	if !v.Configured() || params["vnp_TmnCode"] != v.TmnCode {
		return false
	}
	receivedHash, ok := params["vnp_SecureHash"]
	if !ok || len(receivedHash) != sha512.Size*2 {
		return false
	}
	verifyParams := copyVNPayParams(params)
	delete(verifyParams, "vnp_SecureHash")
	delete(verifyParams, "vnp_SecureHashType")
	expected, _ := hex.DecodeString(v.sign(verifyParams))
	received, err := hex.DecodeString(receivedHash)
	return err == nil && hmac.Equal(expected, received)
}

func (v *VNPayClient) BuildPaymentURL(txnRef, orderInfo, clientIP string, amountVND float64) (string, error) {
	if !v.Configured() {
		return "", fmt.Errorf("VNPay chưa được cấu hình — đăng ký sandbox tại https://sandbox.vnpayment.vn/devreg/ và thiết lập VNPAY_TMN_CODE, VNPAY_HASH_SECRET trong biến môi trường")
	}

	loc, err := time.LoadLocation("Asia/Ho_Chi_Minh")
	if err != nil {
		loc = time.FixedZone("ICT", 7*3600)
	}
	now := time.Now().In(loc)
	expire := now.Add(vnPayPaymentTTL)
	payURL, err := url.Parse(v.PayURL)
	if err != nil || payURL.Host == "" || payURL.Scheme != "https" || payURL.RawQuery != "" {
		return "", fmt.Errorf("URL cổng thanh toán VNPay không hợp lệ")
	}
	returnURL, err := url.Parse(v.ReturnURL)
	if err != nil || returnURL.Host == "" || (returnURL.Scheme != "https" && returnURL.Scheme != "http") {
		return "", fmt.Errorf("URL nhận kết quả VNPay không hợp lệ")
	}

	amountInt := vnpayAmountInMinorUnits(amountVND)
	if !validPaymentAmount(amountVND) || amountInt < 500000 {
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
	if v.BankCode != "" {
		params["vnp_BankCode"] = v.BankCode
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
