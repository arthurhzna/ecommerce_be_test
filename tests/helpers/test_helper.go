package helpers

import (
	"os"
	"testing"

	"github.com/arthurhzna/ecommerce_be_test/config"
	"github.com/arthurhzna/ecommerce_be_test/domain/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var TestDB *gorm.DB

func SetupTestDB(t *testing.T) *gorm.DB {
	config.Init()

	// Use test database
	testDBName := os.Getenv("TEST_DB_NAME")
	if testDBName == "" {
		testDBName = "ecommerce_test"
	}

	testDBHost := os.Getenv("TEST_DB_HOST")
	if testDBHost == "" {
		testDBHost = "localhost"
	}

	testDBPort := os.Getenv("TEST_DB_PORT")
	if testDBPort == "" {
		testDBPort = "5432"
	}

	testDBUser := os.Getenv("TEST_DB_USER")
	if testDBUser == "" {
		testDBUser = "postgres"
	}

	testDBPassword := os.Getenv("TEST_DB_PASSWORD")
	if testDBPassword == "" {
		testDBPassword = "postgres"
	}

	adminDSN := "host=" + testDBHost + " user=" + testDBUser + " password=" + testDBPassword + " dbname=postgres port=" + testDBPort + " sslmode=disable"
	adminDB, err := gorm.Open(postgres.Open(adminDSN), &gorm.Config{})
	if err == nil {
		var exists bool
		adminDB.Raw("SELECT 1 FROM pg_database WHERE datname = ?", testDBName).Scan(&exists)

		if !exists {
			sqlDB, _ := adminDB.DB()
			if sqlDB != nil {
				_, err = sqlDB.Exec("CREATE DATABASE " + testDBName)
				if err != nil {
					t.Logf("Warning: Failed to create test database (might already exist): %v", err)
				}
				sqlDB.Close()
			}
		} else {
			sqlDB, _ := adminDB.DB()
			if sqlDB != nil {
				sqlDB.Close()
			}
		}
	}

	dsn := "host=" + testDBHost + " user=" + testDBUser + " password=" + testDBPassword + " dbname=" + testDBName + " port=" + testDBPort + " sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v. Make sure PostgreSQL is running and credentials are correct.", err)
	}

	err = db.AutoMigrate(
		&models.Role{},
		&models.User{},
		&models.Product{},
		&models.Order{},
		&models.OrderItem{},
		&models.Payment{},
	)
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	TestDB = db
	return db
}

func CleanupTestDB(t *testing.T, db *gorm.DB) {
	if db == nil {
		return
	}

	db.Exec("SET session_replication_role = 'replica';")

	tables := []string{"payments", "order_items", "orders", "products", "users", "roles"}
	for _, table := range tables {
		if err := db.Exec("TRUNCATE TABLE " + table + " CASCADE").Error; err != nil {
			t.Logf("Warning: Failed to truncate table %s: %v", table, err)
		}
	}

	db.Exec("SET session_replication_role = 'origin';")
}

func SeedTestData(t *testing.T, db *gorm.DB) {
	roles := []models.Role{
		{ID: 1, Code: "admin", Name: "Admin"},
		{ID: 2, Code: "customer", Name: "Customer"},
	}
	for _, role := range roles {
		if err := db.FirstOrCreate(&role, models.Role{Code: role.Code}).Error; err != nil {
			t.Logf("Warning: Failed to seed role %s: %v", role.Code, err)
		}
	}
}
