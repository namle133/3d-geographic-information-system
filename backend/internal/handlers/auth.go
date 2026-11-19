package handlers

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/uit/vinhomes-management/config"
	"github.com/uit/vinhomes-management/internal/database"
	"github.com/uit/vinhomes-management/internal/middleware"
	"github.com/uit/vinhomes-management/internal/models"
	"github.com/uit/vinhomes-management/internal/services"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token,omitempty"`
	UserID       string `json:"user_id"`
	Role         string `json:"role"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	Name         string `json:"name"`
}

type RegisterRequest struct {
	Username         string `json:"username" binding:"required"`
	Email            string `json:"email" binding:"required,email"`
	Password         string `json:"password" binding:"required,min=8"`
	Role             string `json:"role" binding:"omitempty,oneof=resident"`
	FullName         string `json:"full_name"`
	IDCard           string `json:"id_card"`
	DateOfBirth      string `json:"date_of_birth"`
	Gender           string `json:"gender"`
	PhoneNumber      string `json:"phone_number"`
	PermanentAddress string `json:"permanent_address"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPasswordRequest struct {
	Email       string `json:"email" binding:"required,email"`
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

func Login(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		loginID := req.Username
		if loginID == "" {
			loginID = req.Email
		}
		if loginID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email hoặc username là bắt buộc"})
			return
		}

		var user models.User
		if err := database.DB.Where("username = ? OR email = ?", loginID, loginID).First(&user).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Email hoặc mật khẩu không đúng"})
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Email hoặc mật khẩu không đúng"})
			return
		}

		if user.Status != "active" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Tài khoản đã bị khóa hoặc vô hiệu hóa"})
			return
		}

		token, err := middleware.GenerateToken(&user, cfg)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": LoginResponse{
				Token:    token,
				UserID:   user.ID,
				Role:     user.Role,
				Username: user.Username,
				Email:    user.Email,
				Name:     user.Username,
			},
		})
	}
}

func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Role == "" {
		req.Role = "resident"
	}

	// For resident role, validate required fields
	if req.Role == "resident" {
		if req.FullName == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Họ tên là bắt buộc"})
			return
		}
		if len(req.FullName) < 2 || len(req.FullName) > 100 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Họ tên phải từ 2-100 ký tự"})
			return
		}
		if req.IDCard == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "CCCD là bắt buộc"})
			return
		}
		if req.DateOfBirth == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Ngày sinh là bắt buộc"})
			return
		}
		if req.Gender == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Giới tính là bắt buộc"})
			return
		}
		if req.Gender != "M" && req.Gender != "F" && req.Gender != "O" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Giới tính không hợp lệ"})
			return
		}
		if req.PhoneNumber == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Số điện thoại là bắt buộc"})
			return
		}
		if req.PermanentAddress == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Địa chỉ thường trú là bắt buộc"})
			return
		}
		if len(req.PermanentAddress) < 5 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Địa chỉ phải ít nhất 5 ký tự"})
			return
		}
	}

	// Check if user already exists
	var existingUser models.User
	if err := database.DB.Where("username = ? OR email = ?", req.Username, req.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Người dùng đã tồn tại"})
		return
	}

	// Check if ID card already exists (for resident role)
	if req.Role == "resident" {
		var existingResident models.Resident
		if err := database.DB.Where("id_card = ?", req.IDCard).First(&existingResident).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "CCCD đã được đăng ký"})
			return
		}
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể mã hóa mật khẩu"})
		return
	}

	// Create user
	user := models.User{
		ID:       uuid.New().String(),
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     req.Role,
		Status:   "active",
	}

	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể tạo người dùng"})
		return
	}

	// Cư dân: tạo 1 bản ghi Resident đầy đủ (không gọi createRoleRecord — tránh bản ghi rỗng trùng user_id)
	if req.Role == "resident" {
		dateOfBirth, err := time.Parse("2006-01-02", req.DateOfBirth)
		if err != nil {
			database.DB.Delete(&user)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Định dạng ngày sinh không hợp lệ (YYYY-MM-DD)"})
			return
		}

		gender := req.Gender
		if gender == "O" {
			gender = "M" // DB chỉ lưu M/F; form có thêm "Khác"
		}

		resident := models.Resident{
			ID:               uuid.New().String(),
			UserID:           user.ID,
			FullName:         req.FullName,
			IDCard:           req.IDCard,
			DateOfBirth:      dateOfBirth,
			Gender:           gender,
			PhoneNumber:      req.PhoneNumber,
			Email:            req.Email,
			PermanentAddress: req.PermanentAddress,
		}

		if err := database.DB.Create(&resident).Error; err != nil {
			database.DB.Delete(&user)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể tạo hồ sơ cư dân: " + err.Error()})
			return
		}
	} else {
		createRoleRecord(user.ID, req.Role)
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Đăng ký thành công",
		"user_id": user.ID,
	})
}

// ForgotPassword sends reset tokens only to the account email, never to the caller.
func ForgotPassword(cfg *config.Config) gin.HandlerFunc {
	return forgotPasswordWithSender(services.NewSMTPPasswordResetSender(cfg))
}

func resetTokenHash(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func forgotPasswordWithSender(sender services.PasswordResetSender) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req ForgotPasswordRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email không hợp lệ"})
			return
		}
		if !sender.Configured() {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Khôi phục qua email chưa được cấu hình. Vui lòng liên hệ ban quản lý."})
			return
		}
		const message = "Nếu email tồn tại, mã đặt lại mật khẩu sẽ được gửi tới email đã đăng ký."
		var user models.User
		if err := database.DB.Where("email = ? AND status = ?", req.Email, "active").First(&user).Error; err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Chưa thể xử lý yêu cầu. Vui lòng thử lại sau."})
				return
			}
			c.JSON(http.StatusOK, gin.H{"message": message})
			return
		}
		tokenBytes := make([]byte, 32)
		if _, err := rand.Read(tokenBytes); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Chưa thể tạo mã khôi phục"})
			return
		}
		token := hex.EncodeToString(tokenBytes)
		tokenHash := resetTokenHash(token)
		now := time.Now()
		// Do not repeatedly invalidate a code or send email more than once per minute.
		update := database.DB.Model(&models.User{}).Where("id = ? AND (reset_token_expiry IS NULL OR reset_token_expiry < ?)", user.ID, now.Add(29*time.Minute)).Updates(map[string]interface{}{"reset_token": tokenHash, "reset_token_expiry": now.Add(30 * time.Minute)})
		if update.Error != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Chưa thể xử lý yêu cầu. Vui lòng thử lại sau."})
			return
		}
		if update.RowsAffected > 0 {
			if err := sender.Send(user.Email, token); err != nil {
				database.DB.Model(&models.User{}).Where("id = ? AND reset_token = ?", user.ID, tokenHash).Updates(map[string]interface{}{"reset_token": "", "reset_token_expiry": nil})
				log.Print("Password reset email delivery failed")
			}
		}
		c.JSON(http.StatusOK, gin.H{"message": message})
	}
}

func ResetPassword(c *gin.Context) {
	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Thông tin đặt lại mật khẩu không hợp lệ"})
		return
	}
	tokenBytes, err := hex.DecodeString(req.Token)
	if err != nil || len(tokenBytes) != 32 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Mã reset không hợp lệ hoặc đã hết hạn"})
		return
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Chưa thể cập nhật mật khẩu"})
		return
	}
	// Atomically consume a valid token, including the expiry check, exactly once.
	result := database.DB.Model(&models.User{}).Where("email = ? AND reset_token = ? AND reset_token_expiry > ? AND status = ?", req.Email, resetTokenHash(req.Token), time.Now(), "active").Updates(map[string]interface{}{"password": string(hashed), "reset_token": "", "reset_token_expiry": nil})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Chưa thể cập nhật mật khẩu"})
		return
	}
	if result.RowsAffected != 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Mã reset không hợp lệ hoặc đã hết hạn"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Đặt lại mật khẩu thành công"})
}

func GetProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var user models.User

	if err := database.DB.Preload("Resident").Preload("Manager").Preload("Staff").Preload("Admin").
		First(&user, "id = ?", userID.(string)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

// UpdateProfile — UC 2.7.5/2.7.6
func UpdateProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var req struct {
		Email       string `json:"email"`
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := database.DB.First(&user, "id = ?", userID.(string)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if req.Email != "" {
		database.DB.Model(&user).Update("email", req.Email)
	}

	if req.NewPassword != "" {
		if req.OldPassword == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Cần nhập mật khẩu cũ"})
			return
		}
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Mật khẩu cũ không đúng"})
			return
		}
		hashed, _ := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
		database.DB.Model(&user).Update("password", string(hashed))
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cập nhật hồ sơ thành công"})
}
