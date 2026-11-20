package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/uit/vinhomes-management/internal/database"
	"github.com/uit/vinhomes-management/internal/middleware"
	"github.com/uit/vinhomes-management/internal/models"
	"github.com/uit/vinhomes-management/internal/testutil"
	"image"
	"image/png"
	"mime/multipart"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestPrivateImageStorageAndChatAttachments(t *testing.T) {
	database.DB = testutil.OpenDB(t)
	t.Setenv("UPLOAD_DIR", t.TempDir())
	gin.SetMode(gin.TestMode)
	users := []models.User{}
	for _, role := range []string{"resident", "admin", "resident", "staff"} {
		id := uuid.NewString()
		u := models.User{ID: id, Username: id, Email: id + "@example.com", Role: role, Status: "active"}
		if err := database.DB.Create(&u).Error; err != nil {
			t.Fatal(err)
		}
		users = append(users, u)
	}
	thread := models.ChatConversation{ID: uuid.NewString(), UserID: users[0].ID, Status: "open"}
	database.DB.Create(&thread)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		for _, u := range users {
			if c.GetHeader("X-Test-User") == u.ID {
				c.Set("user_id", u.ID)
				c.Set("role", u.Role)
				return
			}
		}
		c.AbortWithStatus(401)
	})
	r.POST("/images", UploadImage)
	r.POST("/maintenance", CreateMaintenanceRequest)
	r.POST("/maintenance/:id/complete", middleware.RequireRole("admin", "manager", "staff"), CompleteMaintenanceRequest)
	r.GET("/images/:imageId", GetImage)
	r.DELETE("/images/:imageId", DeleteUnusedImage)
	r.POST("/thread/:id/messages", SendChatMessage)
	call := func(user int, method, path string, raw []byte, contentType string, code int) *httptest.ResponseRecorder {
		t.Helper()
		req := httptest.NewRequest(method, path, bytes.NewReader(raw))
		req.Header.Set("Content-Type", contentType)
		if user >= 0 {
			req.Header.Set("X-Test-User", users[user].ID)
		}
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != code {
			t.Fatalf("%s %s want %d got %d %s", method, path, code, w.Code, w.Body.String())
		}
		return w
	}
	upload := func(user int, purpose, conversation string, raw []byte, code int) string {
		t.Helper()
		var body bytes.Buffer
		writer := multipart.NewWriter(&body)
		writer.WriteField("purpose", purpose)
		writer.WriteField("conversation_id", conversation)
		f, _ := writer.CreateFormFile("file", "untrusted.svg")
		f.Write(raw)
		writer.Close()
		w := call(user, "POST", "/images", body.Bytes(), writer.FormDataContentType(), code)
		if code != 201 {
			return ""
		}
		var result struct {
			Data models.ImageAsset `json:"data"`
		}
		json.Unmarshal(w.Body.Bytes(), &result)
		return result.Data.ID
	}
	var valid bytes.Buffer
	png.Encode(&valid, image.NewRGBA(image.Rect(0, 0, 16, 16)))
	valid.WriteString("PRIVATE-TRAILING-METADATA")
	upload(0, "chat", "not-a-conversation", valid.Bytes(), 403)
	upload(2, "chat", thread.ID, valid.Bytes(), 403)
	upload(0, "maintenance", "", []byte("<svg onload=alert(1)>"), 400)
	upload(0, "maintenance", "", make([]byte, (2<<20)+1), 400)
	id := upload(1, "chat", thread.ID, valid.Bytes(), 201)
	call(-1, "GET", "/images/"+id, nil, "", 401)
	call(0, "GET", "/images/"+id, nil, "", 403)
	call(2, "GET", "/images/"+id, nil, "", 403)
	call(3, "GET", "/images/"+id, nil, "", 403)
	imageResponse := call(1, "GET", "/images/"+id, nil, "", 200)
	if imageResponse.Header().Get("Cache-Control") != "private, no-store" || strings.Contains(imageResponse.Body.String(), "PRIVATE-TRAILING-METADATA") {
		t.Fatal("image metadata or cache privacy failure")
	}
	if _, _, err := image.Decode(bytes.NewReader(imageResponse.Body.Bytes())); err != nil {
		t.Fatal(err)
	}
	send := func(user int, ids []string, nonce string, code int) []byte {
		raw, _ := json.Marshal(gin.H{"body": "", "client_id": nonce, "image_ids": ids})
		return call(user, "POST", "/thread/"+thread.ID+"/messages", raw, "application/json", code).Body.Bytes()
	}
	nonce := uuid.NewString()
	first := send(1, []string{id}, nonce, 200)
	retry := send(1, []string{id}, nonce, 200)
	if !bytes.Equal(first, retry) {
		t.Fatal("attachment retry changed message")
	}
	call(0, "GET", "/images/"+id, nil, "", 200)
	call(2, "GET", "/images/"+id, nil, "", 403)
	send(0, []string{id}, uuid.NewString(), 409) // Cannot impersonate the image uploader.
	call(1, "DELETE", "/images/"+id, nil, "", 409)
	pending := upload(0, "maintenance", "", valid.Bytes(), 201)
	if validImageIDs(database.DB, []string{pending}, users[0].ID, "maintenance", "") != nil {
		t.Fatal("valid pending maintenance image rejected")
	}
	if validImageIDs(database.DB, []string{pending}, users[2].ID, "maintenance", "") == nil {
		t.Fatal("other user's image accepted")
	}
	call(0, "DELETE", "/images/"+pending, nil, "", 204)
	if _, err := os.Stat(imagePath(pending)); !os.IsNotExist(err) {
		t.Fatal("unused file not removed")
	}
	if validImageIDs(database.DB, []string{id}, users[1].ID, "maintenance", "") == nil {
		t.Fatal("chat image accepted as maintenance image")
	}
	building := models.Building{ID: uuid.NewString(), Name: "Test"}
	database.DB.Create(&building)
	floor := models.Floor{ID: uuid.NewString(), BuildingID: building.ID}
	database.DB.Create(&floor)
	apartment := models.Apartment{ID: uuid.NewString(), FloorID: floor.ID, ApartmentCode: uuid.NewString()}
	database.DB.Create(&apartment)
	resident := models.Resident{ID: uuid.NewString(), UserID: users[0].ID, IDCard: uuid.NewString()}
	database.DB.Create(&resident)
	maintenanceImage := upload(0, "maintenance", "", valid.Bytes(), 201)
	payload, _ := json.Marshal(gin.H{"apartment_id": apartment.ID, "resident_id": "spoofed", "issue_type": "water", "description": "Test image", "priority": "low", "photos": []string{maintenanceImage}})
	created := call(0, "POST", "/maintenance", payload, "application/json", 201)
	var maintenance struct {
		Data models.MaintenanceRequest `json:"data"`
	}
	json.Unmarshal(created.Body.Bytes(), &maintenance)
	if maintenance.Data.ResidentID != resident.ID {
		t.Fatal("resident ownership was spoofed")
	}
	call(0, "GET", "/images/"+maintenanceImage, nil, "", 200)
	call(3, "GET", "/images/"+maintenanceImage, nil, "", 200)
	call(2, "GET", "/images/"+maintenanceImage, nil, "", 403)
	completedImage := upload(3, "maintenance", "", valid.Bytes(), 201)
	completion, _ := json.Marshal(gin.H{"completion_notes": "Completed", "photos": []string{completedImage}})
	call(0, "POST", "/maintenance/"+maintenance.Data.ID+"/complete", completion, "application/json", 403)
	call(3, "POST", "/maintenance/"+maintenance.Data.ID+"/complete", completion, "application/json", 200)
	call(0, "GET", "/images/"+completedImage, nil, "", 200)
	call(2, "GET", "/images/"+completedImage, nil, "", 403)

	// The same directory survives handler/database reconnections; bytes are never
	// stored in a process-local map or the container's default filesystem.
	if bytes, err := os.ReadFile(imagePath(id)); err != nil || len(bytes) == 0 {
		t.Fatal("image was not persisted", err)
	}
	t.Setenv("UPLOAD_DIR", "")
	t.Setenv("RAILWAY_VOLUME_MOUNT_PATH", "")
	upload(0, "maintenance", "", valid.Bytes(), 503)
}
