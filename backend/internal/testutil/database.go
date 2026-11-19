package testutil

import (
	"net/url"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/uit/vinhomes-management/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func OpenDB(t *testing.T) *gorm.DB {
	t.Helper()
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("set TEST_DATABASE_URL or run scripts/test-integration.sh for PostgreSQL tests")
	}
	base, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatal(err)
	}
	schema := "test_" + uuid.New().String()[:8]
	if err := base.Exec("CREATE SCHEMA " + schema).Error; err != nil {
		t.Fatal(err)
	}
	u, err := url.Parse(dsn)
	if err != nil {
		t.Fatal(err)
	}
	q := u.Query()
	q.Set("search_path", schema)
	u.RawQuery = q.Encode()
	db, err := gorm.Open(postgres.Open(u.String()), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
		base.Exec("DROP SCHEMA " + schema + " CASCADE")
		sqlBase, _ := base.DB()
		sqlBase.Close()
	})
	if err := db.AutoMigrate(&models.User{}, &models.Building{}, &models.Floor{}, &models.Apartment{}, &models.Resident{}, &models.Contract{}, &models.Invoice{}, &models.InvoicePayment{}, &models.PaymentTransaction{}, &models.Notification{}, &models.ChatConversation{}, &models.ChatMessage{}); err != nil {
		t.Fatal(err)
	}
	return db
}

func Invoice(t *testing.T, db *gorm.DB, amount float64) models.Invoice {
	t.Helper()
	id := uuid.NewString()
	rows := []interface{}{
		&models.Building{ID: id},
		&models.Floor{ID: id, BuildingID: id},
		&models.Apartment{ID: id, FloorID: id, ApartmentCode: id},
		&models.User{ID: id, Username: id, Email: id + "@example.com", Role: "resident", Status: "active"},
		&models.Resident{ID: id, UserID: id, IDCard: id},
		&models.Contract{ID: id, ApartmentID: id, ResidentID: id, ContractCode: id},
	}
	for _, row := range rows {
		if err := db.Create(row).Error; err != nil {
			t.Fatal(err)
		}
	}
	inv := models.Invoice{ID: id, ContractID: id, ApartmentID: id, InvoiceCode: id, TotalAmount: amount, Status: "pending"}
	if err := db.Create(&inv).Error; err != nil {
		t.Fatal(err)
	}
	return inv
}
