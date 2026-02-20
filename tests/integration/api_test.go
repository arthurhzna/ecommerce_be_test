package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/arthurhzna/ecommerce_be_test/common/response"
	"github.com/arthurhzna/ecommerce_be_test/config"
	"github.com/arthurhzna/ecommerce_be_test/constants"
	"github.com/arthurhzna/ecommerce_be_test/controllers"
	"github.com/arthurhzna/ecommerce_be_test/database/seeders"
	"github.com/arthurhzna/ecommerce_be_test/domain/models"
	"github.com/arthurhzna/ecommerce_be_test/middlewares"
	"github.com/arthurhzna/ecommerce_be_test/repositories"
	"github.com/arthurhzna/ecommerce_be_test/routes"
	"github.com/arthurhzna/ecommerce_be_test/services"
	"github.com/arthurhzna/ecommerce_be_test/tests/helpers"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var testDB *gorm.DB
var testRouter *gin.Engine

func setupTestRouter(db *gorm.DB) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middlewares.HandlePanic())

	seeders.NewSeederRegistry(db).Run()

	repository := repositories.NewRepositoryRegistry(db)
	service := services.NewServiceRegistry(repository)
	controller := controllers.NewControllerRegistry(service)

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, response.Response{
			Status:  constants.Success,
			Message: "Welcome to Ecommerce BE Test",
		})
	})

	group := router.Group("/api/v1")
	route := routes.NewRouteRegistry(controller, group)
	route.Serve()

	return router
}

func TestMain(m *testing.M) {
	config.Init()

	var err error
	testDB, err = setupTestDatabase()
	if err != nil {
		panic("Failed to setup test database: " + err.Error())
	}

	testRouter = setupTestRouter(testDB)

	code := m.Run()

	if testDB != nil {
		cleanupTestDatabase(testDB)
	}

	os.Exit(code)
}

func setupTestDatabase() (*gorm.DB, error) {
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
			_, err = sqlDB.Exec("CREATE DATABASE " + testDBName)
			if err != nil {
			}
		}

		sqlDB, _ := adminDB.DB()
		if sqlDB != nil {
			sqlDB.Close()
		}
	}

	dsn := "host=" + testDBHost + " user=" + testDBUser + " password=" + testDBPassword + " dbname=" + testDBName + " port=" + testDBPort + " sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
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
		return nil, err
	}

	return db, nil
}

func cleanupTestDatabase(db *gorm.DB) {
	if db == nil {
		return
	}

	db.Exec("SET session_replication_role = 'replica';")

	tables := []string{"payments", "order_items", "orders", "products", "users", "roles"}
	for _, table := range tables {
		db.Exec("TRUNCATE TABLE " + table + " CASCADE")
	}

	db.Exec("SET session_replication_role = 'origin';")
}

func TestHealthCheck(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "success", response["status"])
}

func TestRegisterUser(t *testing.T) {
	cleanupTestDatabase(testDB)
	helpers.SeedTestData(t, testDB)

	tests := []struct {
		name           string
		payload        map[string]interface{}
		expectedStatus int
	}{
		{
			name: "success register",
			payload: map[string]interface{}{
				"name":            "Test User",
				"email":           "testuser@example.com",
				"password":        "password123",
				"confirmPassword": "password123",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "duplicate email",
			payload: map[string]interface{}{
				"name":            "Another User",
				"email":           "testuser@example.com",
				"password":        "password123",
				"confirmPassword": "password123",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "invalid email format",
			payload: map[string]interface{}{
				"name":            "Test User",
				"email":           "invalid-email",
				"password":        "password123",
				"confirmPassword": "password123",
			},
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name: "password mismatch",
			payload: map[string]interface{}{
				"name":            "Test User",
				"email":           "newuser@example.com",
				"password":        "password123",
				"confirmPassword": "differentpassword",
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonValue, _ := json.Marshal(tt.payload)
			req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(jsonValue))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Api-Key", config.Config.ApiKey)

			w := httptest.NewRecorder()
			testRouter.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestLoginUser(t *testing.T) {
	cleanupTestDatabase(testDB)
	helpers.SeedTestData(t, testDB)

	registerPayload := map[string]interface{}{
		"name":            "Login Test User",
		"email":           "logintest@example.com",
		"password":        "password123",
		"confirmPassword": "password123",
	}
	jsonValue, _ := json.Marshal(registerPayload)
	req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Api-Key", config.Config.ApiKey)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	tests := []struct {
		name           string
		payload        map[string]interface{}
		expectedStatus int
		expectToken    bool
	}{
		{
			name: "success login",
			payload: map[string]interface{}{
				"email":    "logintest@example.com",
				"password": "password123",
			},
			expectedStatus: http.StatusOK,
			expectToken:    true,
		},
		{
			name: "wrong password",
			payload: map[string]interface{}{
				"email":    "logintest@example.com",
				"password": "wrongpassword",
			},
			expectedStatus: http.StatusBadRequest,
			expectToken:    false,
		},
		{
			name: "user not found",
			payload: map[string]interface{}{
				"email":    "notfound@example.com",
				"password": "password123",
			},
			expectedStatus: http.StatusBadRequest,
			expectToken:    false,
		},
		{
			name: "invalid email format",
			payload: map[string]interface{}{
				"email":    "invalid-email",
				"password": "password123",
			},
			expectedStatus: http.StatusUnprocessableEntity,
			expectToken:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonValue, _ := json.Marshal(tt.payload)
			req, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonValue))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Api-Key", config.Config.ApiKey)

			w := httptest.NewRecorder()
			testRouter.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectToken {
				var response map[string]interface{}
				json.Unmarshal(w.Body.Bytes(), &response)
				assert.NotNil(t, response["token"])
			}
		})
	}
}

func TestGetProducts(t *testing.T) {
	cleanupTestDatabase(testDB)
	helpers.SeedTestData(t, testDB)

	seeders.RunUserSeeder(testDB)

	product := models.Product{
		Name:        "Test Product",
		Description: "Test Description",
		Price:       100.00,
		Stock:       10,
	}
	testDB.Create(&product)

	loginPayload := map[string]interface{}{
		"email":    "admin@gmail.com",
		"password": "admin123",
	}
	jsonValue, _ := json.Marshal(loginPayload)
	req, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Api-Key", config.Config.ApiKey)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	var loginResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &loginResponse)
	token, ok := loginResponse["token"].(string)
	if !ok {
		t.Fatal("Failed to get token for product test")
	}

	req, _ = http.NewRequest("GET", "/api/v1/products", nil)
	req.Header.Set("Api-Key", config.Config.ApiKey)
	req.Header.Set("Authorization", "Bearer "+token)

	w = httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "success", response["status"])
	assert.NotNil(t, response["data"])
}

func TestCreateProduct(t *testing.T) {
	cleanupTestDatabase(testDB)
	helpers.SeedTestData(t, testDB)

	seeders.RunUserSeeder(testDB)

	loginPayload := map[string]interface{}{
		"email":    "admin@gmail.com",
		"password": "admin123",
	}
	jsonValue, _ := json.Marshal(loginPayload)
	req, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Api-Key", config.Config.ApiKey)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	var loginResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &loginResponse)
	token, ok := loginResponse["token"].(string)
	if !ok {
		t.Fatal("Admin login failed, cannot proceed with product creation test")
	}

	tests := []struct {
		name           string
		payload        map[string]interface{}
		expectedStatus int
	}{
		{
			name: "success create product",
			payload: map[string]interface{}{
				"name":        "New Product",
				"description": "New Product Description",
				"price":       150.00,
				"stock":       15,
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "invalid payload",
			payload: map[string]interface{}{
				"name": "AB",
			},
			expectedStatus: http.StatusUnprocessableEntity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonValue, _ := json.Marshal(tt.payload)
			req, _ := http.NewRequest("POST", "/api/v1/products/create", bytes.NewBuffer(jsonValue))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Api-Key", config.Config.ApiKey)
			req.Header.Set("Authorization", "Bearer "+token)

			w := httptest.NewRecorder()
			testRouter.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestCreateOrder(t *testing.T) {
	cleanupTestDatabase(testDB)
	helpers.SeedTestData(t, testDB)
	seeders.RunUserSeeder(testDB)

	// Create product
	product := models.Product{
		UUID:        uuid.New(),
		Name:        "Test Product",
		Description: "Test Description",
		Price:       100.00,
		Stock:       10,
	}
	testDB.Create(&product)

	// Register and login as customer
	registerPayload := map[string]interface{}{
		"name":            "Order Test User",
		"email":           "ordertest@example.com",
		"password":        "password123",
		"confirmPassword": "password123",
	}
	jsonValue, _ := json.Marshal(registerPayload)
	req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Api-Key", config.Config.ApiKey)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	loginPayload := map[string]interface{}{
		"email":    "ordertest@example.com",
		"password": "password123",
	}
	jsonValue, _ = json.Marshal(loginPayload)
	req, _ = http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Api-Key", config.Config.ApiKey)
	w = httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	var loginResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &loginResponse)
	token, ok := loginResponse["token"].(string)
	if !ok {
		t.Fatal("Failed to get token for order test")
	}

	tests := []struct {
		name           string
		payload        map[string]interface{}
		expectedStatus int
	}{
		{
			name: "success create order",
			payload: map[string]interface{}{
				"items": []map[string]interface{}{
					{
						"productUUID": product.UUID.String(),
						"quantity":    2,
					},
				},
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "invalid payload - empty items",
			payload: map[string]interface{}{
				"items": []map[string]interface{}{},
			},
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name: "invalid payload - invalid quantity",
			payload: map[string]interface{}{
				"items": []map[string]interface{}{
					{
						"productUUID": product.UUID.String(),
						"quantity":    0,
					},
				},
			},
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name: "insufficient stock",
			payload: map[string]interface{}{
				"items": []map[string]interface{}{
					{
						"productUUID": product.UUID.String(),
						"quantity":    100,
					},
				},
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonValue, _ := json.Marshal(tt.payload)
			req, _ := http.NewRequest("POST", "/api/v1/orders", bytes.NewBuffer(jsonValue))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Api-Key", config.Config.ApiKey)
			req.Header.Set("Authorization", "Bearer "+token)

			w := httptest.NewRecorder()
			testRouter.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusCreated {
				var response map[string]interface{}
				json.Unmarshal(w.Body.Bytes(), &response)
				assert.Equal(t, "success", response["status"])
				assert.NotNil(t, response["data"])
			}
		})
	}
}

func TestGetMyOrders(t *testing.T) {
	cleanupTestDatabase(testDB)
	helpers.SeedTestData(t, testDB)
	seeders.RunUserSeeder(testDB)

	product := models.Product{
		UUID:        uuid.New(),
		Name:        "Test Product",
		Description: "Test Description",
		Price:       100.00,
		Stock:       10,
	}
	testDB.Create(&product)

	registerPayload := map[string]interface{}{
		"name":            "Order List User",
		"email":           "orderlist@example.com",
		"password":        "password123",
		"confirmPassword": "password123",
	}
	jsonValue, _ := json.Marshal(registerPayload)
	req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Api-Key", config.Config.ApiKey)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	loginPayload := map[string]interface{}{
		"email":    "orderlist@example.com",
		"password": "password123",
	}
	jsonValue, _ = json.Marshal(loginPayload)
	req, _ = http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Api-Key", config.Config.ApiKey)
	w = httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	var loginResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &loginResponse)
	token, ok := loginResponse["token"].(string)
	if !ok {
		t.Fatal("Failed to get token for order list test")
	}

	orderPayload := map[string]interface{}{
		"items": []map[string]interface{}{
			{
				"productUUID": product.UUID.String(),
				"quantity":    1,
			},
		},
	}
	jsonValue, _ = json.Marshal(orderPayload)
	req, _ = http.NewRequest("POST", "/api/v1/orders", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Api-Key", config.Config.ApiKey)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	req, _ = http.NewRequest("GET", "/api/v1/orders/my", nil)
	req.Header.Set("Api-Key", config.Config.ApiKey)
	req.Header.Set("Authorization", "Bearer "+token)

	w = httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "success", response["status"])
	assert.NotNil(t, response["data"])
}

func TestProcessPayment(t *testing.T) {
	cleanupTestDatabase(testDB)
	helpers.SeedTestData(t, testDB)
	seeders.RunUserSeeder(testDB)

	product := models.Product{
		UUID:        uuid.New(),
		Name:        "Test Product",
		Description: "Test Description",
		Price:       100.00,
		Stock:       10,
	}
	testDB.Create(&product)

	registerPayload := map[string]interface{}{
		"name":            "Payment Test User",
		"email":           "paymenttest@example.com",
		"password":        "password123",
		"confirmPassword": "password123",
	}
	jsonValue, _ := json.Marshal(registerPayload)
	req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Api-Key", config.Config.ApiKey)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	loginPayload := map[string]interface{}{
		"email":    "paymenttest@example.com",
		"password": "password123",
	}
	jsonValue, _ = json.Marshal(loginPayload)
	req, _ = http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Api-Key", config.Config.ApiKey)
	w = httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	var loginResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &loginResponse)
	token, ok := loginResponse["token"].(string)
	if !ok {
		t.Fatal("Failed to get token for payment test")
	}

	orderPayload := map[string]interface{}{
		"items": []map[string]interface{}{
			{
				"productUUID": product.UUID.String(),
				"quantity":    1,
			},
		},
	}
	jsonValue, _ = json.Marshal(orderPayload)
	req, _ = http.NewRequest("POST", "/api/v1/orders", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Api-Key", config.Config.ApiKey)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	var orderResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &orderResponse)

	if orderResponse["status"] != "success" {
		t.Logf("Order creation failed: %v", orderResponse)
		t.Fatal("Failed to create order for payment test")
	}

	orderData, ok := orderResponse["data"].(map[string]interface{})
	if !ok {
		t.Logf("Order response: %v", orderResponse)
		t.Fatal("Failed to get order data")
	}
	orderUUID, ok := orderData["uuid"].(string)
	if !ok {
		t.Logf("Order data: %v", orderData)
		t.Fatal("Failed to get order UUID")
	}

	tests := []struct {
		name           string
		orderUUID      string
		payload        map[string]interface{}
		expectedStatus int
	}{
		{
			name:           "success process payment",
			orderUUID:      orderUUID,
			payload:        map[string]interface{}{},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid uuid format",
			orderUUID:      "invalid-uuid",
			payload:        map[string]interface{}{},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonValue, _ := json.Marshal(tt.payload)
			req, _ := http.NewRequest("POST", "/api/v1/payments/"+tt.orderUUID+"/pay", bytes.NewBuffer(jsonValue))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Api-Key", config.Config.ApiKey)
			req.Header.Set("Authorization", "Bearer "+token)

			w := httptest.NewRecorder()
			testRouter.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusOK {
				var response map[string]interface{}
				json.Unmarshal(w.Body.Bytes(), &response)
				assert.Equal(t, "success", response["status"])
				assert.NotNil(t, response["data"])
			}
		})
	}
}
