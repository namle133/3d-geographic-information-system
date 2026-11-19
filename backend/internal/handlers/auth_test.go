package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/uit/vinhomes-management/internal/database"
	"github.com/uit/vinhomes-management/internal/models"
	"github.com/uit/vinhomes-management/internal/testutil"
	"golang.org/x/crypto/bcrypt"
)

type fakeResetSender struct {
	configured   bool
	email, token string
	err          error
	calls        int
}

func (f *fakeResetSender) Configured() bool { return f.configured }
func (f *fakeResetSender) Send(email, token string) error {
	f.email = email
	f.token = token
	f.calls++
	return f.err
}
func requestJSON(t *testing.T, handler gin.HandlerFunc, body map[string]string) *httptest.ResponseRecorder {
	t.Helper()
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/test", handler)
	raw, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader(raw))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
func TestPublicRegistrationRejectsElevatedRoles(t *testing.T) {
	for _, role := range []string{"admin", "manager", "staff"} {
		w := requestJSON(t, Register, map[string]string{"username": "test", "email": "test@example.com", "password": "test-password", "role": role})
		if w.Code != 400 {
			t.Fatalf("role %s: %d", role, w.Code)
		}
	}
}
func TestForgotPasswordUnavailableDoesNotReturnToken(t *testing.T) {
	w := requestJSON(t, forgotPasswordWithSender(&fakeResetSender{}), map[string]string{"email": "test@example.com"})
	if w.Code != 503 || strings.Contains(w.Body.String(), "reset_token") {
		t.Fatal(w.Code, w.Body.String())
	}
}
func TestPasswordResetTokenIsPrivateAndSingleUse(t *testing.T) {
	database.DB = testutil.OpenDB(t)
	id := uuid.NewString()
	user := models.User{ID: id, Username: id, Email: id + "@example.com", Role: "resident", Status: "active"}
	if err := database.DB.Create(&user).Error; err != nil {
		t.Fatal(err)
	}
	sender := &fakeResetSender{configured: true}
	handler := forgotPasswordWithSender(sender)
	w := requestJSON(t, handler, map[string]string{"email": user.Email})
	if w.Code != 200 || sender.email != user.Email || len(sender.token) != 64 || strings.Contains(w.Body.String(), sender.token) || strings.Contains(w.Body.String(), "reset_token") {
		t.Fatalf("reset response: %d %s", w.Code, w.Body.String())
	}
	var stored models.User
	database.DB.First(&stored, "id = ?", id)
	if stored.ResetToken != resetTokenHash(sender.token) || stored.ResetToken == sender.token {
		t.Fatal("token must be hashed")
	}
	unknown := requestJSON(t, handler, map[string]string{"email": "unknown@example.com"})
	if unknown.Body.String() != w.Body.String() {
		t.Fatal("response exposes account existence")
	}
	requestJSON(t, handler, map[string]string{"email": user.Email})
	if sender.calls != 1 {
		t.Fatal("reset emails were not throttled")
	}
	token := sender.token
	payload := map[string]string{"email": user.Email, "token": token, "new_password": "new-test-password"}
	if w := requestJSON(t, ResetPassword, payload); w.Code != 200 {
		t.Fatal(w.Code, w.Body.String())
	}
	if w := requestJSON(t, ResetPassword, payload); w.Code != 400 {
		t.Fatal("token reused", w.Code)
	}
	database.DB.First(&stored, "id = ?", id)
	if err := bcrypt.CompareHashAndPassword([]byte(stored.Password), []byte(payload["new_password"])); err != nil {
		t.Fatal("password was not saved")
	}
	// Neither old plaintext demo tokens nor expired tokens may be used.
	database.DB.Model(&stored).Updates(map[string]interface{}{"reset_token": token, "reset_token_expiry": time.Now().Add(time.Minute)})
	if w := requestJSON(t, ResetPassword, payload); w.Code != 400 {
		t.Fatal("accepted plaintext token")
	}
	database.DB.Model(&stored).Updates(map[string]interface{}{"reset_token": resetTokenHash(token), "reset_token_expiry": time.Now().Add(-time.Second)})
	if w := requestJSON(t, ResetPassword, payload); w.Code != 400 {
		t.Fatal("accepted expired token")
	}
}
func TestPasswordResetFailedDeliveryInvalidatesToken(t *testing.T) {
	database.DB = testutil.OpenDB(t)
	id := uuid.NewString()
	u := models.User{ID: id, Username: id, Email: id + "@example.com", Role: "resident", Status: "active"}
	database.DB.Create(&u)
	sender := &fakeResetSender{configured: true, err: errors.New("delivery failed")}
	w := requestJSON(t, forgotPasswordWithSender(sender), map[string]string{"email": u.Email})
	if w.Code != 200 {
		t.Fatal(w.Code)
	}
	var stored models.User
	database.DB.First(&stored, "id = ?", id)
	if stored.ResetToken != "" || stored.ResetTokenExpiry != nil {
		t.Fatal("undelivered token remains active")
	}
}
