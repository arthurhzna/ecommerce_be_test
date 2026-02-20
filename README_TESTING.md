# Testing Guide

Panduan lengkap untuk testing aplikasi Ecommerce BE.

## Prerequisites

1. **Dependencies** (sudah termasuk di go.mod)
   - `github.com/stretchr/testify` - assertions & mocks
   - `github.com/DATA-DOG/go-sqlmock` - database mocking

2. **Database Test**
   - PostgreSQL server harus running (untuk integration tests)
   - Database dibuat **otomatis** saat integration test
   - Tidak perlu PostgreSQL client tools (psql, createdb, dll)

## Tipe Testing

### Unit Tests
- ✅ Tidak perlu database (menggunakan mock)
- ✅ Test: Repository, Service, Controller
- ✅ Bisa dijalankan di komputer manapun
- ✅ Menggunakan `sqlmock` untuk mock database queries
- ✅ Menggunakan `testify/mock` untuk mock services

### Integration Tests
- ✅ Perlu PostgreSQL server running
- ✅ Database dibuat otomatis oleh Go code
- ✅ Test end-to-end API endpoints
- ✅ Menggunakan test database yang terpisah (`ecommerce_test`)

## Menjalankan Tests

### Semua Tests
```bash
# Dari root project
cd be_eco
go test -v ./...
```

### Unit Tests Only
```bash
# Test semua unit tests (repositories, services, controllers)
go test -v ./repositories/... ./services/... ./controllers/...

# Atau test per layer
go test -v ./repositories/...
go test -v ./services/...
go test -v ./controllers/...
```

### Integration Tests Only
```bash
# Pastikan PostgreSQL running
go test -v ./tests/integration/...
```

### Test Package Tertentu
```bash
# Test controller tertentu
go test -v ./controllers/product
go test -v ./controllers/order
go test -v ./controllers/payment
go test -v ./controllers/user

# Test service tertentu
go test -v ./services/product
go test -v ./services/order

# Test repository tertentu
go test -v ./repositories/product
go test -v ./repositories/order
```

### Test dengan Coverage
```bash
# Generate coverage report
go test -v -coverprofile=coverage.out ./...

# Generate HTML coverage report
go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

# Buka coverage.html di browser
```

### Test dengan Pattern Tertentu
```bash
# Test yang namanya mengandung "Order"
go test -v -run TestOrder ./...

# Test yang namanya mengandung "Product"
go test -v -run TestProduct ./...
```

## Struktur Testing

```
be_eco/
├── controllers/
│   ├── order/order_test.go          ✅ Unit test OrderController
│   ├── payment/payment_test.go      ✅ Unit test PaymentController
│   ├── product/product_test.go      ✅ Unit test ProductController
│   └── user/user_test.go            ✅ Unit test UserController
├── services/
│   ├── order/order_test.go          ✅ Unit test OrderService
│   ├── payment/payment_test.go      ✅ Unit test PaymentService
│   ├── product/product_test.go      ✅ Unit test ProductService
│   └── user/user_test.go            ✅ Unit test UserService
├── repositories/
│   ├── order/order_test.go           ✅ Unit test OrderRepository
│   ├── order_item/order_item_test.go ✅ Unit test OrderItemRepository
│   ├── payment/payment_test.go      ✅ Unit test PaymentRepository
│   ├── product/product_test.go      ✅ Unit test ProductRepository
│   └── user/user_test.go            ✅ Unit test UserRepository
└── tests/
    ├── helpers/                      # Test utilities & mocks
    └── integration/
        └── api_test.go               ✅ Integration tests (all endpoints)
```

## Test Coverage

### Unit Tests Coverage

**Controllers:**
- ✅ `OrderController` - CreateOrder, GetMyOrders
- ✅ `PaymentController` - ProcessPayment
- ✅ `ProductController` - GetProductsWithoutPagination, CreateProduct
- ✅ `UserController` - Login, Register

**Services:**
- ✅ `OrderService` - CreateOrder, GetMyOrders
- ✅ `PaymentService` - ProcessPayment
- ✅ `ProductService` - GetProductsWithoutPagination, CreateProduct
- ✅ `UserService` - Login, Register, GetUserByEmail

**Repositories:**
- ✅ `OrderRepository` - Create, FindByUserID, FindByID, FindByUUID, Update
- ✅ `OrderItemRepository` - CreateBulk, FindByOrderID
- ✅ `PaymentRepository` - Create, FindByOrderID, Update
- ✅ `ProductRepository` - FindAllWithoutPagination, Create, FindByName, UpdateStock
- ✅ `UserRepository` - Register, FindByEmail

### Integration Tests Coverage

- ✅ `TestHealthCheck` - Health check endpoint
- ✅ `TestRegisterUser` - User registration
- ✅ `TestLoginUser` - User login
- ✅ `TestGetProducts` - Get all products
- ✅ `TestCreateProduct` - Create product (Admin only)
- ✅ `TestCreateOrder` - Create order
- ✅ `TestGetMyOrders` - Get user orders
- ✅ `TestProcessPayment` - Process payment

## Best Practices

1. **Naming Convention**: `TestFunctionName_Scenario`
   ```go
   TestUserService_Login
   TestUserService_Login/success_login
   TestUserService_Login/user_not_found
   ```

2. **Mock Expectations**: 
   - Selalu assert `mockRepo.AssertExpectations(t)` untuk unit tests
   - Setup mock sebelum memanggil function yang di-test

3. **Test Isolation**: 
   - Setiap test harus independent
   - Jangan bergantung pada state dari test lain
   - Gunakan `cleanupTestDatabase()` setelah integration test

4. **Test Data**:
   - Gunakan data yang jelas dan konsisten
   - Hindari hardcoded values yang tidak jelas
   - Gunakan factory functions untuk create test data

5. **Error Testing**:
   - Test both success dan error cases
   - Test edge cases (empty data, invalid input, dll)

## Contoh Test

### Unit Test - Repository
```go
func TestOrderRepository_Create(t *testing.T) {
    db, mock := setupMockDB(t)
    defer db.Close()
    
    mock.ExpectBegin()
    mock.ExpectQuery(`INSERT INTO "orders"`).
        WithArgs(sqlmock.AnyArg(), ...).
        WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
    mock.ExpectCommit()
    
    repo := NewOrderRepository(db)
    order, err := repo.Create(context.Background(), &models.Order{...})
    
    assert.NoError(t, err)
    assert.NotNil(t, order)
    assert.NoError(t, mock.ExpectationsWereMet())
}
```

### Unit Test - Service
```go
func TestUserService_Login(t *testing.T) {
    mockRepo := new(helpers.MockUserRepository)
    service := NewUserService(mockRepo)
    
    mockRepo.On("FindByEmail", "test@example.com").Return(user, nil)
    
    result, err := service.Login(context.Background(), req)
    assert.NoError(t, err)
    assert.NotNil(t, result)
    mockRepo.AssertExpectations(t)
}
```

### Unit Test - Controller
```go
func TestProductController_CreateProduct(t *testing.T) {
    gin.SetMode(gin.TestMode)
    router := gin.New()
    
    mockService := new(MockProductService)
    mockService.On("CreateProduct", mock.Anything, mock.Anything).
        Return(productResponse, nil)
    
    controller := NewProductController(mockServiceRegistry)
    router.POST("/products/create", controller.CreateProduct)
    
    jsonValue, _ := json.Marshal(payload)
    req, _ := http.NewRequest("POST", "/products/create", bytes.NewBuffer(jsonValue))
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)
    
    assert.Equal(t, http.StatusCreated, w.Code)
    mockService.AssertExpectations(t)
}
```

### Integration Test
```go
func TestLoginUser(t *testing.T) {
    cleanupTestDatabase(testDB)
    helpers.SeedTestData(t, testDB)
    
    payload := map[string]interface{}{
        "email":    "test@example.com",
        "password": "password123",
    }
    
    jsonValue, _ := json.Marshal(payload)
    req, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonValue))
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Api-Key", config.Config.ApiKey)
    
    w := httptest.NewRecorder()
    testRouter.ServeHTTP(w, req)
    
    assert.Equal(t, http.StatusOK, w.Code)
    
    var response map[string]interface{}
    json.Unmarshal(w.Body.Bytes(), &response)
    assert.Equal(t, "success", response["status"])
    assert.NotNil(t, response["data"])
}
```

## Testing dengan Docker (Optional)

Jika tidak ada PostgreSQL di host, gunakan Docker:

```bash
# Start PostgreSQL di container (jika ada docker-compose.test.yml)
docker compose -f docker-compose.test.yml up -d postgres-test

# Atau gunakan docker-compose.yml yang ada
docker compose up -d postgres

# Set environment variables
export TEST_DB_HOST=localhost
export TEST_DB_PORT=5433  # atau port yang sesuai
export TEST_DB_USER=postgres
export TEST_DB_PASSWORD=postgres
export TEST_DB_NAME=ecommerce_test

# Run integration tests
go test -v ./tests/integration/...

# Stop container
docker compose down
```

## Environment Variables untuk Testing

Integration tests menggunakan environment variables berikut (dengan default values):

```bash
TEST_DB_HOST=localhost
TEST_DB_PORT=5432
TEST_DB_USER=postgres
TEST_DB_PASSWORD=postgres
TEST_DB_NAME=ecommerce_test
```

Jika tidak di-set, akan menggunakan default values dari `config/config.go`.

## Troubleshooting

### Error: "Failed to connect to test database"
- ✅ Pastikan PostgreSQL server running
- ✅ Check permission CREATE DATABASE untuk user postgres
- ✅ Verify connection dengan: `psql -h localhost -U postgres -c "SELECT 1"`

### Error: "database does not exist"
- ✅ Database seharusnya dibuat otomatis oleh `setupTestDatabase()`
- ✅ Pastikan PostgreSQL server accessible
- ✅ Check environment variables (TEST_DB_*)
- ✅ Pastikan user postgres punya permission CREATE DATABASE

### Unit test error tentang database?
- ✅ Unit test tidak perlu database (menggunakan mock)
- ✅ Periksa mock setup di test
- ✅ Pastikan menggunakan `sqlmock` dengan benar
- ✅ Check `mock.ExpectationsWereMet()` untuk melihat unfulfilled expectations

### Error: "all expectations were already fulfilled"
- ✅ GORM melakukan Preload queries yang perlu di-mock
- ✅ Tambahkan mock expectations untuk Preload queries
- ✅ Atau gunakan `sqlmock.AnyArg()` untuk arguments yang fleksibel

### Error: "arguments do not match"
- ✅ GORM First() menggunakan LIMIT yang menghasilkan 2 arguments
- ✅ Gunakan `sqlmock.AnyArg()` untuk arguments yang tidak pasti
- ✅ Atau match exact arguments sesuai dengan query GORM

## Test Statistics

Jalankan test untuk melihat statistik:

```bash
go test -v ./... 2>&1 | grep -E "(PASS|FAIL|ok)"
```

## Resources

- [Go Testing Package](https://pkg.go.dev/testing)
- [Testify Documentation](https://github.com/stretchr/testify)
- [SQLMock Documentation](https://github.com/DATA-DOG/go-sqlmock)
- [GORM Testing](https://gorm.io/docs/testing.html)

## Notes

- Unit tests menggunakan mock, jadi tidak perlu database
- Integration tests memerlukan PostgreSQL server running
- Database test dibuat otomatis, tidak perlu setup manual
- Semua test harus independent dan bisa dijalankan secara parallel
- Gunakan `-short` flag untuk skip integration tests: `go test -short ./...`
