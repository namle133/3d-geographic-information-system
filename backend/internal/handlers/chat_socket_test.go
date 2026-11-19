package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/uit/vinhomes-management/config"
	"github.com/uit/vinhomes-management/internal/database"
	"github.com/uit/vinhomes-management/internal/middleware"
	"github.com/uit/vinhomes-management/internal/models"
	"github.com/uit/vinhomes-management/internal/testutil"
)

func TestChatSocketAuthenticationAndIsolation(t *testing.T) {
	database.DB = testutil.OpenDB(t)
	cfg := &config.Config{JWTSecret: "test-only-socket-key", JWTExpiryHours: 1}
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/ws", ChatSocket(cfg))
	r.POST("/thread/:id/messages", middleware.AuthMiddleware(cfg), SendChatMessage)
	server := httptest.NewServer(r)
	defer server.Close()
	url := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"
	headers := http.Header{"Origin": []string{"https://3d-geographic-information-system.vercel.app"}}
	if conn, resp, err := websocket.DefaultDialer.Dial(url, http.Header{"Origin": []string{"https://untrusted.example"}}); err == nil {
		conn.Close()
		t.Fatal("untrusted origin accepted")
	} else if resp.StatusCode != 403 {
		t.Fatal(resp.StatusCode)
	}
	bad, _, err := websocket.DefaultDialer.Dial(url, headers)
	if err != nil {
		t.Fatal(err)
	}
	bad.WriteJSON(gin.H{"token": "invalid"})
	bad.SetReadDeadline(time.Now().Add(time.Second))
	if _, _, err := bad.ReadMessage(); err == nil {
		t.Fatal("invalid session accepted")
	}
	bad.Close()
	users := []models.User{}
	for _, role := range []string{"resident", "admin", "resident"} {
		id := uuid.NewString()
		u := models.User{ID: id, Username: id, Email: id + "@example.com", Role: role, Status: "active"}
		if err := database.DB.Create(&u).Error; err != nil {
			t.Fatal(err)
		}
		users = append(users, u)
	}
	sockets := []*websocket.Conn{}
	for _, user := range users {
		token, err := middleware.GenerateToken(&user, cfg)
		if err != nil {
			t.Fatal(err)
		}
		conn, _, err := websocket.DefaultDialer.Dial(url, headers)
		if err != nil {
			t.Fatal(err)
		}
		defer conn.Close()
		conn.WriteJSON(gin.H{"token": token})
		conn.SetReadDeadline(time.Now().Add(2 * time.Second))
		var event map[string]string
		if err := conn.ReadJSON(&event); err != nil || event["type"] != "ready" {
			t.Fatal("missing authenticated ready", err, event)
		}
		sockets = append(sockets, conn)
	}
	thread := models.ChatConversation{ID: uuid.NewString(), UserID: users[0].ID}
	thread.Status = "open"
	if err := database.DB.Create(&thread).Error; err != nil {
		t.Fatal(err)
	}
	token, _ := middleware.GenerateToken(&users[0], cfg)
	payload, _ := json.Marshal(gin.H{"body": "WebSocket integration test", "client_id": uuid.NewString()})
	request := httptest.NewRequest("POST", "/thread/"+thread.ID+"/messages", bytes.NewReader(payload))
	request.Header.Set("Authorization", "Bearer "+token)
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	r.ServeHTTP(response, request)
	if response.Code != 200 {
		t.Fatal(response.Code, response.Body.String())
	}
	for _, socket := range sockets[:2] {
		var event map[string]string
		if err := socket.ReadJSON(&event); err != nil || event["conversation_id"] != thread.ID {
			t.Fatal("owner/admin did not get update", err, event)
		}
	}
	sockets[2].SetReadDeadline(time.Now().Add(100 * time.Millisecond))
	if _, _, err := sockets[2].ReadMessage(); err == nil {
		t.Fatal("other resident received private event")
	}
	// Revoked accounts cannot authenticate a new connection.
	database.DB.Model(&users[0]).Update("status", "locked")
	token, _ = middleware.GenerateToken(&users[0], cfg)
	if _, err := socketUser(token, cfg); err == nil {
		t.Fatal("locked account authorized")
	}
	// Wait for every read pump to exit before this test's database is removed.
	for _, socket := range sockets {
		socket.Close()
	}
	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		chatHub.Lock()
		n := len(chatHub.clients)
		chatHub.Unlock()
		if n == 0 {
			return
		}
		time.Sleep(time.Millisecond)
	}
	t.Fatal("socket subscribers leaked")
}
