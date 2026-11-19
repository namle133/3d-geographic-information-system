package handlers

import (
	"errors"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
	"github.com/uit/vinhomes-management/config"
	"github.com/uit/vinhomes-management/internal/database"
	"github.com/uit/vinhomes-management/internal/middleware"
	"github.com/uit/vinhomes-management/internal/models"
)

type chatSubscriber struct {
	userID, role string
	events       chan string
}
type chatEventHub struct {
	sync.Mutex
	clients map[*chatSubscriber]bool
}

var chatHub = chatEventHub{clients: make(map[*chatSubscriber]bool)}

// Events contain only a conversation ID. History is always fetched through the
// authenticated HTTP endpoints, which recheck ownership and account status.
func publishChat(thread *models.ChatConversation) {
	chatHub.Lock()
	defer chatHub.Unlock()
	for client := range chatHub.clients {
		if client.role == "admin" || client.userID == thread.UserID {
			select {
			case client.events <- thread.ID:
			default: /* Events are invalidations; the next refresh also catches up. */
			}
		}
	}
}
func chatSocketOrigin(r *http.Request) bool {
	u, err := url.Parse(r.Header.Get("Origin"))
	if err != nil {
		return false
	}
	if u.Scheme == "https" && u.Host == "3d-geographic-information-system.vercel.app" {
		return true
	}
	return u.Scheme == "http" && (u.Hostname() == "localhost" || u.Hostname() == "127.0.0.1")
}
func socketUser(token string, cfg *config.Config) (*models.User, error) {
	claims := &middleware.Claims{}
	_, err := jwt.ParseWithClaims(token, claims, func(_ *jwt.Token) (interface{}, error) { return []byte(cfg.JWTSecret), nil }, jwt.WithValidMethods([]string{"HS256"}))
	if err != nil {
		return nil, err
	}
	if claims.ExpiresAt == nil {
		return nil, errors.New("expiration required")
	}
	var user models.User
	err = database.DB.First(&user, "id = ? AND status = ?", claims.UserID, "active").Error
	return &user, err
}
func ChatSocket(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		upgrader := websocket.Upgrader{CheckOrigin: chatSocketOrigin, HandshakeTimeout: 5 * time.Second}
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			return
		}
		defer conn.Close()
		// Browser WebSockets cannot set Authorization headers. Authenticate in the
		// first encrypted frame, never in a URL or access log. Nothing is subscribed yet.
		conn.SetReadLimit(8192)
		conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		var auth struct {
			Token string `json:"token"`
		}
		if conn.ReadJSON(&auth) != nil {
			return
		}
		user, err := socketUser(auth.Token, cfg)
		if err != nil {
			conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(1008, "Authentication required"), time.Now().Add(time.Second))
			return
		}
		client := &chatSubscriber{userID: user.ID, role: user.Role, events: make(chan string, 32)}
		chatHub.Lock()
		chatHub.clients[client] = true
		chatHub.Unlock()
		defer func() { chatHub.Lock(); delete(chatHub.clients, client); chatHub.Unlock() }()
		conn.SetReadDeadline(time.Now().Add(45 * time.Second))
		conn.SetPongHandler(func(string) error { return conn.SetReadDeadline(time.Now().Add(45 * time.Second)) })
		done := make(chan struct{})
		go func() {
			defer close(done)
			for {
				if _, _, err := conn.ReadMessage(); err != nil {
					return
				}
			}
		}()
		write := func(value any) error {
			conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
			return conn.WriteJSON(value)
		}
		if write(gin.H{"type": "ready"}) != nil {
			return
		}
		ticker := time.NewTicker(15 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-done:
				return
			case id := <-client.events:
				if write(gin.H{"type": "chat.updated", "conversation_id": id}) != nil {
					return
				}
			case <-ticker.C:
				current, err := socketUser(auth.Token, cfg)
				if err != nil || current.Role != client.role {
					conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(1008, "Session expired"), time.Now().Add(time.Second))
					return
				}
				if conn.WriteControl(websocket.PingMessage, nil, time.Now().Add(5*time.Second)) != nil {
					return
				}
			}
		}
	}
}
