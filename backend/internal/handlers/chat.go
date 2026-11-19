package handlers

import (
	"errors"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/uit/vinhomes-management/internal/database"
	"github.com/uit/vinhomes-management/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func chatAdmin(c *gin.Context) bool { return c.GetString("role") == "admin" }
func chatThread(c *gin.Context) (*models.ChatConversation, bool) {
	var thread models.ChatConversation
	if err := database.DB.First(&thread, "id = ?", c.Param("id")).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(404, gin.H{"error": "Conversation not found"})
		} else {
			c.JSON(500, gin.H{"error": "Unable to load conversation"})
		}
		return nil, false
	}
	if !chatAdmin(c) && thread.UserID != c.GetString("user_id") {
		c.JSON(403, gin.H{"error": "Access denied"})
		return nil, false
	}
	return &thread, true
}

// ChatPeople is the administrator's board. Passwords and personal profile records are never returned.
func ChatPeople(c *gin.Context) {
	if !chatAdmin(c) {
		c.JSON(403, gin.H{"error": "Access denied"})
		return
	}
	type person struct {
		ID             string     `json:"id"`
		Username       string     `json:"username"`
		Email          string     `json:"email"`
		Role           string     `json:"role"`
		AccountStatus  string     `json:"account_status"`
		ConversationID *string    `json:"conversation_id"`
		Status         *string    `json:"status"`
		LastMessageAt  *time.Time `json:"last_message_at"`
		Preview        string     `json:"preview"`
		Unread         int64      `json:"unread"`
	}
	var people []person
	err := database.DB.Raw(`SELECT u.id,u.username,u.email,u.role,u.status AS account_status,c.id AS conversation_id,c.status,c.last_message_at,
 COALESCE((SELECT body FROM chat_messages m WHERE m.conversation_id=c.id ORDER BY m.created_at DESC,m.id DESC LIMIT 1),'') AS preview,
 (SELECT COUNT(*) FROM chat_messages m WHERE m.conversation_id=c.id AND m.from_admin=false AND (c.admin_read_at IS NULL OR m.created_at>c.admin_read_at)) AS unread
 FROM users u LEFT JOIN chat_conversations c ON c.user_id=u.id WHERE u.role <> 'admin'
 ORDER BY c.last_message_at DESC NULLS LAST,u.username,u.id`).Scan(&people).Error
	if err != nil {
		c.JSON(500, gin.H{"error": "Unable to load conversations"})
		return
	}
	if people == nil {
		people = []person{}
	}
	c.JSON(200, gin.H{"data": people})
}
func StartChat(c *gin.Context) {
	owner := c.GetString("user_id")
	if chatAdmin(c) {
		var req struct {
			UserID string `json:"user_id" binding:"required"`
		}
		if c.ShouldBindJSON(&req) != nil {
			c.JSON(400, gin.H{"error": "Select a user"})
			return
		}
		owner = req.UserID
	}
	var user models.User
	if err := database.DB.First(&user, "id = ? AND role <> ?", owner, "admin").Error; err != nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}
	if user.Status != "active" {
		c.JSON(409, gin.H{"error": "This account is inactive"})
		return
	}
	thread := models.ChatConversation{ID: uuid.NewString(), UserID: owner, Status: "open", CreatedAt: time.Now()}
	if err := database.DB.Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "user_id"}}, DoNothing: true}).Create(&thread).Error; err != nil {
		c.JSON(500, gin.H{"error": "Unable to open conversation"})
		return
	}
	thread = models.ChatConversation{}
	if err := database.DB.First(&thread, "user_id = ?", owner).Error; err != nil {
		c.JSON(500, gin.H{"error": "Unable to open conversation"})
		return
	}
	c.JSON(200, gin.H{"data": thread})
}
func ChatMessages(c *gin.Context) {
	thread, ok := chatThread(c)
	if !ok {
		return
	}
	query := database.DB.Where("conversation_id = ?", thread.ID)
	if before := c.Query("before"); before != "" {
		var cursor models.ChatMessage
		if database.DB.First(&cursor, "id = ? AND conversation_id = ?", before, thread.ID).Error != nil {
			c.JSON(400, gin.H{"error": "Invalid message cursor"})
			return
		}
		query = query.Where("(created_at < ?) OR (created_at = ? AND id < ?)", cursor.CreatedAt, cursor.CreatedAt, cursor.ID)
	}
	var messages []models.ChatMessage
	if err := query.Order("created_at DESC,id DESC").Limit(101).Find(&messages).Error; err != nil {
		c.JSON(500, gin.H{"error": "Unable to load messages"})
		return
	}
	more := len(messages) > 100
	if more {
		messages = messages[:100]
	}
	if messages == nil {
		messages = []models.ChatMessage{}
	}
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}
	c.JSON(200, gin.H{"data": messages, "conversation": thread, "has_more": more})
}
func SendChatMessage(c *gin.Context) {
	thread, ok := chatThread(c)
	if !ok {
		return
	}
	var req struct {
		Body     string `json:"body"`
		ClientID string `json:"client_id"`
	}
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 16384)
	if c.ShouldBindJSON(&req) != nil {
		c.JSON(400, gin.H{"error": "Invalid message"})
		return
	}
	req.Body = strings.TrimSpace(req.Body)
	if req.Body == "" || !utf8.ValidString(req.Body) || utf8.RuneCountInString(req.Body) > 2000 {
		c.JSON(400, gin.H{"error": "Messages must contain 1–2000 characters"})
		return
	}
	if _, err := uuid.Parse(req.ClientID); err != nil {
		c.JSON(400, gin.H{"error": "Invalid message reference"})
		return
	}
	var msg models.ChatMessage
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(thread, "id = ?", thread.ID).Error; err != nil {
			return err
		}
		var owner models.User
		if err := tx.First(&owner, "id = ?", thread.UserID).Error; err != nil {
			return err
		}
		if owner.Status != "active" {
			return errors.New("inactive")
		}
		err := tx.First(&msg, "client_id = ?", req.ClientID).Error
		if err == nil {
			if msg.ConversationID != thread.ID || msg.SenderID != c.GetString("user_id") || msg.Body != req.Body {
				return errors.New("conflict")
			}
			return nil
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		msg = models.ChatMessage{ID: uuid.NewString(), ConversationID: thread.ID, SenderID: c.GetString("user_id"), ClientID: req.ClientID, Body: req.Body, FromAdmin: chatAdmin(c), CreatedAt: time.Now()}
		if err := tx.Create(&msg).Error; err != nil {
			return err
		}
		return tx.Model(thread).Updates(map[string]interface{}{"last_message_at": msg.CreatedAt, "status": "open"}).Error
	})
	if err != nil {
		c.JSON(409, gin.H{"error": "Unable to send message. Refresh and try again."})
		return
	}
	publishChat(thread)
	c.JSON(200, gin.H{"data": msg})
}
func ReadChat(c *gin.Context) {
	thread, ok := chatThread(c)
	if !ok {
		return
	}
	var req struct {
		MessageID string `json:"message_id"`
	}
	if c.ShouldBindJSON(&req) != nil {
		c.JSON(400, gin.H{"error": "Invalid read receipt"})
		return
	}
	var msg models.ChatMessage
	if database.DB.First(&msg, "id = ? AND conversation_id = ?", req.MessageID, thread.ID).Error != nil {
		c.JSON(400, gin.H{"error": "Invalid read receipt"})
		return
	}
	field := "user_read_at"
	if chatAdmin(c) {
		field = "admin_read_at"
	}
	result := database.DB.Model(thread).Where(field+" IS NULL OR "+field+" < ?", msg.CreatedAt).Update(field, msg.CreatedAt)
	if result.Error != nil {
		c.JSON(500, gin.H{"error": "Unable to mark messages as read"})
		return
	}
	if result.RowsAffected > 0 {
		publishChat(thread)
	}
	c.JSON(200, gin.H{"success": true})
}
func ChatStatus(c *gin.Context) {
	if !chatAdmin(c) {
		c.JSON(403, gin.H{"error": "Access denied"})
		return
	}
	thread, ok := chatThread(c)
	if !ok {
		return
	}
	var req struct {
		Status string `json:"status" binding:"required,oneof=open resolved"`
	}
	if c.ShouldBindJSON(&req) != nil {
		c.JSON(400, gin.H{"error": "Invalid conversation status"})
		return
	}
	if err := database.DB.Model(thread).Update("status", req.Status).Error; err != nil {
		c.JSON(500, gin.H{"error": "Unable to update conversation"})
		return
	}
	publishChat(thread)
	c.JSON(200, gin.H{"success": true})
}

func ChatUnread(c *gin.Context) {
	var count int64
	query := database.DB.Table("chat_messages m").Joins("JOIN chat_conversations c ON c.id=m.conversation_id")
	if chatAdmin(c) {
		query = query.Where("m.from_admin=false AND (c.admin_read_at IS NULL OR m.created_at>c.admin_read_at)")
	} else {
		query = query.Where("c.user_id=? AND m.from_admin=true AND (c.user_read_at IS NULL OR m.created_at>c.user_read_at)", c.GetString("user_id"))
	}
	if err := query.Count(&count).Error; err != nil {
		c.JSON(500, gin.H{"error": "Unable to load unread messages"})
		return
	}
	c.JSON(200, gin.H{"unread_count": count})
}
