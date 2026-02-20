# Testing Guide

Panduan testing untuk aplikasi Ecommerce BE.

## Prerequisites

1. **Dependencies** (sudah termasuk di go.mod)
   - `github.com/stretchr/testify` - assertions & mocks
   - `github.com/DATA-DOG/go-sqlmock` - database mocking

2. **Database Test**
   - PostgreSQL server harus running
   - Database dibuat **otomatis** saat integration test
   - Tidak perlu PostgreSQL client tools (psql, createdb, dll)

## Tipe Testing

### Unit Tests
- ✅ Tidak perlu database (menggunakan mock)
- ✅ Test: Repository, Service, Controller
- ✅ Bisa dijalankan di komputer manapun

### Integration Tests
- ✅ Perlu PostgreSQL server running
- ✅ Database dibuat otomatis
- ✅ Test end-to-end API endpoints

## Menjalankan Tests

```bash
# Semua tests
make test

# Unit tests only
make test-unit

# Integration tests only
make test-integration

# Dengan coverage
make test-coverage

# Test spesifik package
go test -v ./services/user/...
```

## Struktur Testing

```
be_eco/
├── repositories/*/..._test.go    # Unit test repository
├── services/*/..._test.go        # Unit test service
├── controllers/*/..._test.go     # Unit test controller
└── tests/
    ├── helpers/                   # Test utilities
    └── integration/               # Integration tests
```

## Best Practices

1. **Naming**: `TestFunctionName_Scenario`
2. **Mock Expectations**: Selalu assert `mockRepo.AssertExpectations(t)`
3. **Cleanup**: Gunakan `helpers.CleanupTestDB()` setelah integration test
4. **Isolation**: Setiap test harus independent

## Coverage Target

- Minimum: 70%
- Ideal: 80%+

```bash
make test-coverage
# Buka coverage.html di browser
```

## Testing dengan Docker (Optional)

Jika tidak ada PostgreSQL di host:

```bash
# Start PostgreSQL di container
make test-db-up

# Set environment
export TEST_DB_HOST=localhost
export TEST_DB_PORT=5433

# Run tests
make test-integration

# Stop container
make test-db-down
```

## Troubleshooting

**Error: "Failed to connect to test database"**
- Pastikan PostgreSQL server running
- Check permission CREATE DATABASE untuk user postgres

**Error: "database does not exist"**
- Database seharusnya dibuat otomatis
- Pastikan PostgreSQL server accessible
- Check environment variables

**Unit test error tentang database?**
- Unit test tidak perlu database
- Periksa mock setup

## Contoh Test

### Unit Test
```go
func TestUserService_Login(t *testing.T) {
    mockRepo := new(MockUserRepository)
    service := NewUserService(mockRepo)
    
    mockRepo.On("FindByEmail", "test@example.com").Return(user, nil)
    
    result, err := service.Login(req)
    assert.NoError(t, err)
    mockRepo.AssertExpectations(t)
}
```

### Integration Test
```go
func TestLoginUser(t *testing.T) {
    helpers.CleanupTestDB(t, testDB)
    
    payload := map[string]interface{}{
        "email": "test@example.com",
        "password": "password123",
    }
    
    req := createRequest("POST", "/api/v1/auth/login", payload)
    w := httptest.NewRecorder()
    testRouter.ServeHTTP(w, req)
    
    assert.Equal(t, http.StatusOK, w.Code)
}
```

## Resources

- [Go Testing](https://pkg.go.dev/testing)
- [Testify](https://github.com/stretchr/testify)
- [SQLMock](https://github.com/DATA-DOG/go-sqlmock)
