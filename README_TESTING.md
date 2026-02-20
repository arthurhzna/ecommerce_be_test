# Testing Guide

Complete guide for testing Ecommerce BE application.

## Test Types

### Unit Tests
- ✅ No database required (uses mocks)
- ✅ Test: Repository, Service, Controller
- ✅ Can be run on any computer
- ✅ Uses `sqlmock` to mock database queries
- ✅ Uses `testify/mock` to mock services

### Integration Tests
- ✅ Requires PostgreSQL server running
- ✅ Database is created automatically by Go code
- ✅ Test end-to-end API endpoints
- ✅ Uses a separate test database (`ecommerce_test`)

## Running Tests

### All Tests
```bash
# From project root
cd be_eco
go test -v ./...
```

### Unit Tests Only
```bash
# Test all unit tests (repositories, services, controllers)
go test -v ./repositories/... ./services/... ./controllers/...

# Or test per layer
go test -v ./repositories/...
go test -v ./services/...
go test -v ./controllers/...
```

### Integration Tests Only
```bash
# Make sure PostgreSQL is running
go test -v ./tests/integration/...
```

### Test Specific Package
```bash
# Test specific controller
go test -v ./controllers/product
go test -v ./controllers/order
go test -v ./controllers/payment
go test -v ./controllers/user

# Test specific service
go test -v ./services/product
go test -v ./services/order

# Test specific repository
go test -v ./repositories/product
go test -v ./repositories/order
```

### Test with Coverage
```bash
# Generate coverage report
go test -v -coverprofile=coverage.out ./...

# Generate HTML coverage report
go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

# Open coverage.html in browser
```

### Test with Specific Pattern
```bash
# Test whose name contains "Order"
go test -v -run TestOrder ./...

# Test whose name contains "Product"
go test -v -run TestProduct ./...
```

## Test Structure

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
   - Always assert `mockRepo.AssertExpectations(t)` for unit tests
   - Setup mock before calling the function being tested

3. **Test Isolation**: 
   - Each test must be independent
   - Don't depend on state from other tests
   - Use `cleanupTestDatabase()` after integration test

4. **Test Data**:
   - Use clear and consistent data
   - Avoid unclear hardcoded values
   - Use factory functions to create test data

5. **Error Testing**:
   - Test both success and error cases
   - Test edge cases (empty data, invalid input, etc.)

## Test Examples

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

## Testing with Docker (Optional)

If PostgreSQL is not available on the host, use Docker:

```bash
# Start PostgreSQL in container (if docker-compose.test.yml exists)
docker compose -f docker-compose.test.yml up -d postgres-test

# Or use existing docker-compose.yml
docker compose up -d postgres

# Set environment variables
export TEST_DB_HOST=localhost
export TEST_DB_PORT=5433  # or appropriate port
export TEST_DB_USER=postgres
export TEST_DB_PASSWORD=postgres
export TEST_DB_NAME=ecommerce_test

# Run integration tests
go test -v ./tests/integration/...

# Stop container
docker compose down
```

## Environment Variables for Testing

Integration tests use the following environment variables (with default values):

```bash
TEST_DB_HOST=localhost
TEST_DB_PORT=5432
TEST_DB_USER=postgres
TEST_DB_PASSWORD=postgres
TEST_DB_NAME=ecommerce_test
```

If not set, will use default values from `config/config.go`.

## Test Statistics

Run tests to see statistics:

```bash
go test -v ./... 2>&1 | grep -E "(PASS|FAIL|ok)"
```

## Resources

- [Go Testing Package](https://pkg.go.dev/testing)
- [Testify Documentation](https://github.com/stretchr/testify)
- [SQLMock Documentation](https://github.com/DATA-DOG/go-sqlmock)
- [GORM Testing](https://gorm.io/docs/testing.html)
