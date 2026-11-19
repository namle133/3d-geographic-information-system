package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	// Database
	DatabaseURL string
	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string

	// Server
	ServerPort string
	GinMode    string

	// JWT
	JWTSecret      string
	JWTExpiryHours int

	// Email
	SMTPHost     string
	SMTPPort     string
	SMTPUser     string
	SMTPPassword string
	SMTPFrom     string

	// App
	AppEnv string

	// VNPay
	VNPayTmnCode    string
	VNPayHashSecret string
	VNPayPayURL     string
	VNPayReturnURL  string
	VNPayIPNURL     string
}

func LoadConfig() *Config {
	loadEnvFiles()

	jwtExpiry, _ := strconv.Atoi(getEnv("JWT_EXPIRY_HOURS", "24"))

	return &Config{
		DatabaseURL: getEnv("DATABASE_URL", ""),
		// Mặc định: PostgreSQL trong Docker (docker-compose port 5433)
		DBHost:          getEnv("DB_HOST", "127.0.0.1"),
		DBPort:          getEnv("DB_PORT", "5433"),
		DBUser:          getEnv("DB_USER", "vinhomes_user"),
		DBPassword:      getEnv("DB_PASSWORD", "vinhomes_pass"),
		DBName:          getEnv("DB_NAME", "vinhomes_db"),
		ServerPort:      getEnv("PORT", getEnv("SERVER_PORT", "8080")),
		GinMode:         getEnv("GIN_MODE", "debug"),
		JWTSecret:       getEnv("JWT_SECRET", "your-secret-key"),
		JWTExpiryHours:  jwtExpiry,
		SMTPHost:        getEnv("SMTP_HOST", ""),
		SMTPPort:        getEnv("SMTP_PORT", "587"),
		SMTPUser:        getEnv("SMTP_USER", ""),
		SMTPPassword:    getEnv("SMTP_PASSWORD", ""),
		SMTPFrom:        getEnv("SMTP_FROM", getEnv("SMTP_USER", "")),
		AppEnv:          getEnv("APP_ENV", "development"),
		VNPayTmnCode:    getVNPayEnv("VNPAY_TMN_CODE"),
		VNPayHashSecret: getVNPayEnv("VNPAY_HASH_SECRET"),
		VNPayPayURL:     getVNPayEnv("VNPAY_PAY_URL", "https://sandbox.vnpayment.vn/paymentv2/vpcpay.html"),
		VNPayReturnURL:  getVNPayEnv("VNPAY_RETURN_URL", "https://3d-geographic-information-system.vercel.app/pages/payment-result.html"),
		VNPayIPNURL:     getVNPayEnv("VNPAY_IPN_URL", "https://3d-geographic-information-system-production.up.railway.app/api/v1/payments/vnpay/ipn"),
	}
}

func getVNPayEnv(key string, defaults ...string) string {
	fallback := ""
	if len(defaults) > 0 {
		fallback = defaults[0]
	}
	value := strings.TrimSpace(getEnv(key, fallback))
	if len(value) >= 2 && ((value[0] == '"' && value[len(value)-1] == '"') || (value[0] == '\'' && value[len(value)-1] == '\'')) {
		value = strings.TrimSpace(value[1 : len(value)-1])
	}
	return value
}

// loadEnvFiles ưu tiên config/.env rồi .env ở thư mục backend.
func loadEnvFiles() {
	for _, path := range []string{"config/.env", ".env"} {
		// Process environment > config/.env > .env; load only missing values.
		_ = godotenv.Load(path)
	}
}

func (c *Config) GetDSN() string {
	if c.DatabaseURL != "" {
		return c.DatabaseURL
	}
	host := c.DBHost
	if host == "localhost" {
		// Tránh kết nối IPv6 [::1] trùng instance Postgres khác trên máy
		host = "127.0.0.1"
	}
	// Railway PostgreSQL requires SSL; local Docker does not
	sslMode := "disable"
	if c.AppEnv == "production" {
		sslMode = "require"
	}
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, c.DBPort, c.DBUser, c.DBPassword, c.DBName, sslMode)
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
