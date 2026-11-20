package handlers

import (
	"bytes"
	"errors"
	"image"
	"image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/uit/vinhomes-management/internal/database"
	"github.com/uit/vinhomes-management/internal/models"
	"gorm.io/gorm"
)

var imageWriteLock sync.Mutex
var imageDecodeSlots = make(chan struct{}, 2)

func imageDirectory() string {
	if dir := os.Getenv("UPLOAD_DIR"); dir != "" {
		return dir
	}
	if dir := os.Getenv("RAILWAY_VOLUME_MOUNT_PATH"); dir != "" {
		return filepath.Join(dir, "images")
	}
	return "" // Never silently store production files on an ephemeral filesystem.
}
func imagePath(id string) string { return filepath.Join(imageDirectory(), id+".jpg") }
func validImageIDs(tx *gorm.DB, ids []string, owner, purpose, conversation string) error {
	if len(ids) > 5 {
		return errors.New("Choose up to 5 images")
	}
	seen := map[string]bool{}
	for _, id := range ids {
		if _, err := uuid.Parse(id); err != nil || seen[id] {
			return errors.New("Invalid image reference")
		}
		seen[id] = true
		var asset models.ImageAsset
		if tx.First(&asset, "id = ? AND owner_id = ? AND purpose = ? AND conversation_id = ?", id, owner, purpose, conversation).Error != nil {
			return errors.New("Invalid image reference")
		}
		if purpose == "maintenance" && asset.MaintenanceID != "" {
			return errors.New("Image is already attached")
		}
	}
	return nil
}
func UploadImage(c *gin.Context) {
	if imageDirectory() == "" {
		c.JSON(503, gin.H{"error": "Image storage is unavailable"})
		return
	}
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 3<<20)
	if c.Request.ParseMultipartForm(3<<20) != nil {
		c.JSON(400, gin.H{"error": "Upload a JPEG or PNG image up to 2 MB"})
		return
	}
	defer c.Request.MultipartForm.RemoveAll()
	purpose, conversation := c.PostForm("purpose"), c.PostForm("conversation_id")
	if purpose != "maintenance" && purpose != "chat" {
		c.JSON(400, gin.H{"error": "Invalid upload purpose"})
		return
	}
	if purpose == "chat" {
		var thread models.ChatConversation
		if database.DB.First(&thread, "id = ?", conversation).Error != nil || (!chatAdmin(c) && thread.UserID != c.GetString("user_id")) {
			c.JSON(403, gin.H{"error": "Access denied"})
			return
		}
	} else {
		conversation = ""
	}
	select {
	case imageDecodeSlots <- struct{}{}:
		defer func() { <-imageDecodeSlots }()
	default:
		c.JSON(429, gin.H{"error": "Image uploads are busy. Try again shortly."})
		return
	}
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{"error": "Choose an image"})
		return
	}
	defer file.Close()
	raw, err := io.ReadAll(io.LimitReader(file, (2<<20)+1))
	if err != nil || len(raw) > 2<<20 {
		c.JSON(400, gin.H{"error": "Upload a JPEG or PNG image up to 2 MB"})
		return
	}
	info, kind, err := image.DecodeConfig(bytes.NewReader(raw))
	if err != nil || (kind != "jpeg" && kind != "png") || info.Width < 1 || info.Height < 1 || int64(info.Width)*int64(info.Height) > 16_000_000 {
		c.JSON(400, gin.H{"error": "Use a valid JPEG or PNG image below 16 megapixels"})
		return
	}
	decoded, _, err := image.Decode(bytes.NewReader(raw))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid image"})
		return
	}
	var output bytes.Buffer
	// Re-encoding removes EXIF/location metadata and any trailing non-image payload.
	if jpeg.Encode(&output, decoded, &jpeg.Options{Quality: 85}) != nil || output.Len() > 3<<20 {
		c.JSON(400, gin.H{"error": "Image is too large"})
		return
	}
	imageWriteLock.Lock()
	defer imageWriteLock.Unlock()
	var total, owned int64
	if database.DB.Model(&models.ImageAsset{}).Select("COALESCE(SUM(size),0)").Scan(&total).Error != nil {
		c.AbortWithStatus(503)
		return
	}
	if database.DB.Model(&models.ImageAsset{}).Where("owner_id = ?", c.GetString("user_id")).Select("COALESCE(SUM(size),0)").Scan(&owned).Error != nil {
		c.AbortWithStatus(503)
		return
	}
	if total+int64(output.Len()) > 400<<20 || owned+int64(output.Len()) > 50<<20 {
		c.JSON(409, gin.H{"error": "Image storage limit reached"})
		return
	}
	asset := models.ImageAsset{ID: uuid.NewString(), OwnerID: c.GetString("user_id"), Purpose: purpose, ConversationID: conversation, Size: int64(output.Len()), CreatedAt: time.Now()}
	if os.MkdirAll(imageDirectory(), 0700) != nil || os.WriteFile(imagePath(asset.ID), output.Bytes(), 0600) != nil {
		c.JSON(503, gin.H{"error": "Image storage is unavailable"})
		return
	}
	if database.DB.Create(&asset).Error != nil {
		os.Remove(imagePath(asset.ID))
		c.JSON(500, gin.H{"error": "Unable to save image"})
		return
	}
	c.JSON(201, gin.H{"data": asset})
}
func imageUsed(asset models.ImageAsset) bool {
	if asset.MaintenanceID != "" {
		return true
	}
	var count int64
	database.DB.Model(&models.ChatMessage{}).Where("? = ANY(image_ids)", asset.ID).Count(&count)
	return count > 0
}
func GetImage(c *gin.Context) {
	id := c.Param("imageId")
	if _, err := uuid.Parse(id); err != nil {
		c.AbortWithStatus(404)
		return
	}
	var asset models.ImageAsset
	if database.DB.First(&asset, "id = ?", id).Error != nil {
		c.AbortWithStatus(404)
		return
	}
	allowed := asset.OwnerID == c.GetString("user_id") || chatAdmin(c)
	if !allowed && asset.Purpose == "chat" && imageUsed(asset) {
		var thread models.ChatConversation
		allowed = database.DB.First(&thread, "id = ? AND user_id = ?", asset.ConversationID, c.GetString("user_id")).Error == nil
	}
	if !allowed && asset.Purpose == "maintenance" && asset.MaintenanceID != "" {
		if c.GetString("role") == "manager" || c.GetString("role") == "staff" {
			allowed = true
		} else {
			var count int64
			database.DB.Table("maintenance_requests m").Joins("JOIN residents r ON r.id=m.resident_id").Where("m.id = ? AND r.user_id = ?", asset.MaintenanceID, c.GetString("user_id")).Count(&count)
			allowed = count > 0
		}
	}
	if !allowed {
		c.AbortWithStatus(403)
		return
	}
	if imageDirectory() == "" {
		c.AbortWithStatus(503)
		return
	}
	c.Header("Cache-Control", "private, no-store")
	c.Header("X-Content-Type-Options", "nosniff")
	c.Header("Content-Type", "image/jpeg")
	c.File(imagePath(id))
}
func DeleteUnusedImage(c *gin.Context) {
	imageWriteLock.Lock()
	defer imageWriteLock.Unlock()
	var asset models.ImageAsset
	if database.DB.First(&asset, "id = ? AND owner_id = ?", c.Param("imageId"), c.GetString("user_id")).Error != nil {
		c.AbortWithStatus(404)
		return
	}
	if imageUsed(asset) {
		c.JSON(409, gin.H{"error": "Attached images cannot be removed"})
		return
	}
	if imageDirectory() == "" {
		c.AbortWithStatus(503)
		return
	}
	if err := os.Remove(imagePath(asset.ID)); err != nil && !os.IsNotExist(err) {
		c.AbortWithStatus(500)
		return
	}
	if database.DB.Delete(&asset).Error != nil {
		c.AbortWithStatus(500)
		return
	}
	c.Status(204)
}
