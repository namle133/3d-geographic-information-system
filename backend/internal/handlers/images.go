package handlers

import (
	"bytes"
	"context"
	"errors"
	"image"
	"image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
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

func imageBackend() string {
	account, bucket := os.Getenv("R2_ACCOUNT_ID"), os.Getenv("R2_BUCKET")
	access, secret := os.Getenv("R2_ACCESS_KEY_ID"), os.Getenv("R2_SECRET_ACCESS_KEY")
	configured := account != "" || bucket != "" || access != "" || secret != ""
	if configured {
		if account == "" || bucket == "" || access == "" || secret == "" {
			return ""
		}
		return "r2"
	}
	if imageDirectory() != "" {
		return "filesystem"
	}
	return ""
}

func r2Client() (*s3.Client, string, error) {
	account, bucket := os.Getenv("R2_ACCOUNT_ID"), os.Getenv("R2_BUCKET")
	access, secret := os.Getenv("R2_ACCESS_KEY_ID"), os.Getenv("R2_SECRET_ACCESS_KEY")
	if account == "" || bucket == "" || access == "" || secret == "" || strings.ContainsAny(account+bucket, "/\\") {
		return nil, "", errors.New("R2 image storage is not configured")
	}
	cfg, err := awsconfig.LoadDefaultConfig(context.Background(),
		awsconfig.WithRegion("auto"),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(access, secret, "")),
	)
	if err != nil {
		return nil, "", err
	}
	endpoint := "https://" + account + ".r2.cloudflarestorage.com"
	return s3.NewFromConfig(cfg, func(options *s3.Options) {
		options.BaseEndpoint = aws.String(endpoint)
		options.UsePathStyle = true
	}), bucket, nil
}

func saveImage(backend, id string, data []byte) error {
	if backend == "r2" {
		client, bucket, err := r2Client()
		if err != nil {
			return err
		}
		_, err = client.PutObject(context.Background(), &s3.PutObjectInput{Bucket: aws.String(bucket), Key: aws.String("images/" + id + ".jpg"), Body: bytes.NewReader(data), ContentType: aws.String("image/jpeg")})
		return err
	}
	if backend != "filesystem" {
		return errors.New("image storage is unavailable")
	}
	if err := os.MkdirAll(imageDirectory(), 0700); err != nil {
		return err
	}
	return os.WriteFile(imagePath(id), data, 0600)
}

func readImage(backend, id string) ([]byte, error) {
	if backend == "r2" {
		client, bucket, err := r2Client()
		if err != nil {
			return nil, err
		}
		object, err := client.GetObject(context.Background(), &s3.GetObjectInput{Bucket: aws.String(bucket), Key: aws.String("images/" + id + ".jpg")})
		if err != nil {
			return nil, err
		}
		defer object.Body.Close()
		return io.ReadAll(io.LimitReader(object.Body, 3<<20))
	}
	if backend != "filesystem" {
		return nil, errors.New("image storage is unavailable")
	}
	return os.ReadFile(imagePath(id))
}

func removeImage(backend, id string) error {
	if backend == "r2" {
		client, bucket, err := r2Client()
		if err != nil {
			return err
		}
		_, err = client.DeleteObject(context.Background(), &s3.DeleteObjectInput{Bucket: aws.String(bucket), Key: aws.String("images/" + id + ".jpg")})
		return err
	}
	if backend != "filesystem" {
		return errors.New("image storage is unavailable")
	}
	err := os.Remove(imagePath(id))
	if os.IsNotExist(err) {
		return nil
	}
	return err
}

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
	backend := imageBackend()
	if backend == "" {
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
	asset := models.ImageAsset{ID: uuid.NewString(), OwnerID: c.GetString("user_id"), Purpose: purpose, ConversationID: conversation, StorageBackend: backend, Size: int64(output.Len()), CreatedAt: time.Now()}
	if saveImage(backend, asset.ID, output.Bytes()) != nil {
		c.JSON(503, gin.H{"error": "Image storage is unavailable"})
		return
	}
	if database.DB.Create(&asset).Error != nil {
		_ = removeImage(backend, asset.ID)
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
	c.Header("Cache-Control", "private, no-store")
	c.Header("X-Content-Type-Options", "nosniff")
	c.Header("Content-Type", "image/jpeg")
	backend := asset.StorageBackend
	if backend == "" {
		backend = "filesystem"
	} // Existing rows used the original volume backend.
	data, err := readImage(backend, id)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	c.Data(http.StatusOK, "image/jpeg", data)
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
	backend := asset.StorageBackend
	if backend == "" {
		backend = "filesystem"
	}
	if backend == "" {
		c.AbortWithStatus(503)
		return
	}
	if err := removeImage(backend, asset.ID); err != nil {
		c.AbortWithStatus(500)
		return
	}
	if database.DB.Delete(&asset).Error != nil {
		c.AbortWithStatus(500)
		return
	}
	c.Status(204)
}
