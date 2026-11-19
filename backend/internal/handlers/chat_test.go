package handlers

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/uit/vinhomes-management/internal/database"
	"github.com/uit/vinhomes-management/internal/models"
	"github.com/uit/vinhomes-management/internal/testutil"
)

func TestChatPrivacyReceiptsAndRetries(t *testing.T) {
	database.DB = testutil.OpenDB(t)
	gin.SetMode(gin.TestMode)
	users := []models.User{}
	for _, role := range []string{"admin", "resident", "resident"} {
		id := uuid.NewString()
		u := models.User{ID: id, Username: id, Email: id + "@example.com", Role: role, Status: "active"}
		if err := database.DB.Create(&u).Error; err != nil {
			t.Fatal(err)
		}
		users = append(users, u)
	}
	r := gin.New()
	r.Use(func(c *gin.Context) {
		for _, u := range users {
			if u.ID == c.GetHeader("X-Test-User") {
				c.Set("user_id", u.ID)
				c.Set("role", u.Role)
			}
		}
	})
	r.GET("/people", ChatPeople)
	r.GET("/unread", ChatUnread)
	r.POST("/start", StartChat)
	r.GET("/thread/:id/messages", ChatMessages)
	r.POST("/thread/:id/messages", SendChatMessage)
	r.PUT("/thread/:id/read", ReadChat)
	r.PUT("/thread/:id/status", ChatStatus)
	call := func(user int, method, path string, payload any, code int) map[string]json.RawMessage {
		t.Helper()
		raw, _ := json.Marshal(payload)
		req := httptest.NewRequest(method, path, bytes.NewReader(raw))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Test-User", users[user].ID)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != code {
			t.Fatalf("%s %s: want %d got %d %s", method, path, code, w.Code, w.Body.String())
		}
		var out map[string]json.RawMessage
		if err := json.Unmarshal(w.Body.Bytes(), &out); err != nil {
			t.Fatal(err)
		}
		return out
	}
	start := func(user int) models.ChatConversation {
		t.Helper()
		out := call(user, "POST", "/start", map[string]string{"user_id": users[2].ID}, 200)
		var c models.ChatConversation
		json.Unmarshal(out["data"], &c)
		return c
	}
	thread := start(1)
	if thread.UserID != users[1].ID || start(1).ID != thread.ID {
		t.Fatal("owner spoof or duplicate conversation")
	}
	other := start(2)
	path := "/thread/" + thread.ID
	call(1, "GET", "/people", nil, 403)
	call(2, "GET", path+"/messages", nil, 403)
	call(2, "POST", path+"/messages", nil, 403)
	call(2, "PUT", path+"/read", nil, 403)
	call(1, "PUT", path+"/status", map[string]string{"status": "resolved"}, 403)
	nonce := uuid.NewString()
	payload := map[string]any{"body": "Hello <script>plain text</script>", "client_id": nonce, "from_admin": true, "sender_id": users[0].ID}
	out := call(1, "POST", path+"/messages", payload, 200)
	var message models.ChatMessage
	json.Unmarshal(out["data"], &message)
	if message.FromAdmin || message.SenderID != users[1].ID {
		t.Fatal("sender spoofed")
	}
	out = call(1, "POST", path+"/messages", payload, 200)
	var duplicate models.ChatMessage
	json.Unmarshal(out["data"], &duplicate)
	if duplicate.ID != message.ID {
		t.Fatal("retry duplicated message")
	}
	call(1, "POST", path+"/messages", map[string]string{"body": "different", "client_id": nonce}, 409)
	for _, body := range []string{" ", strings.Repeat("ắ", 2001)} {
		call(1, "POST", path+"/messages", map[string]string{"body": body, "client_id": uuid.NewString()}, 400)
	}
	unread := func(user int, want int) {
		t.Helper()
		res := call(user, "GET", "/unread", nil, 200)
		var n int
		json.Unmarshal(res["unread_count"], &n)
		if n != want {
			t.Fatalf("unread=%d want %d", n, want)
		}
	}
	unread(0, 1)
	unread(1, 0)
	unread(2, 0)
	board := call(0, "GET", "/people", nil, 200)
	if strings.Contains(string(board["data"]), "password") || !strings.Contains(string(board["data"]), "plain text") {
		t.Fatal("incorrect inbox data")
	}
	call(0, "PUT", path+"/read", map[string]string{"message_id": message.ID}, 200)
	unread(0, 0)
	out = call(0, "POST", path+"/messages", map[string]string{"body": "We can help", "client_id": uuid.NewString()}, 200)
	var reply models.ChatMessage
	json.Unmarshal(out["data"], &reply)
	if !reply.FromAdmin {
		t.Fatal("admin reply missing role")
	}
	unread(1, 1)
	unread(2, 0)
	call(2, "PUT", "/thread/"+other.ID+"/read", map[string]string{"message_id": reply.ID}, 400)
	call(1, "PUT", path+"/read", map[string]string{"message_id": reply.ID}, 200)
	unread(1, 0)
	call(1, "PUT", path+"/read", map[string]string{"message_id": message.ID}, 200)
	unread(1, 0)
	call(0, "PUT", path+"/status", map[string]string{"status": "resolved"}, 200)
	call(1, "POST", path+"/messages", map[string]string{"body": "Another question", "client_id": uuid.NewString()}, 200)
	var stored models.ChatConversation
	database.DB.First(&stored, "id = ?", thread.ID)
	if stored.Status != "open" {
		t.Fatal("new message should reopen conversation")
	}
	// History is bounded and an unrelated cursor cannot expose another resident's messages.
	now := time.Now()
	batch := []models.ChatMessage{}
	for i := 0; i < 105; i++ {
		batch = append(batch, models.ChatMessage{ID: uuid.NewString(), ClientID: uuid.NewString(), ConversationID: thread.ID, SenderID: users[1].ID, Body: "History", CreatedAt: now.Add(time.Duration(i) * time.Millisecond)})
	}
	if err := database.DB.Create(&batch).Error; err != nil {
		t.Fatal(err)
	}
	out = call(1, "GET", path+"/messages", nil, 200)
	var page []models.ChatMessage
	json.Unmarshal(out["data"], &page)
	if len(page) != 100 || string(out["has_more"]) != "true" {
		t.Fatal("history not bounded")
	}
	out = call(1, "GET", path+"/messages?before="+page[0].ID, nil, 200)
	json.Unmarshal(out["data"], &page)
	if len(page) != 8 || string(out["has_more"]) != "false" {
		t.Fatal("incorrect older history page", len(page))
	}
	call(2, "GET", "/thread/"+other.ID+"/messages?before="+message.ID, nil, 400)
	database.DB.Model(&users[1]).Update("status", "locked")
	call(0, "POST", path+"/messages", map[string]string{"body": "blocked", "client_id": uuid.NewString()}, 409)
	call(0, "GET", path+"/messages", nil, 200)
}
