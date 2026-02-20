package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/arthurhzna/ecommerce_be_test/domain/dto"
	"github.com/arthurhzna/ecommerce_be_test/domain/models"
	orderService "github.com/arthurhzna/ecommerce_be_test/services/order"
	paymentService "github.com/arthurhzna/ecommerce_be_test/services/payment"
	productService "github.com/arthurhzna/ecommerce_be_test/services/product"
	userService "github.com/arthurhzna/ecommerce_be_test/services/user"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestUserController_Login(t *testing.T) {
	tests := []struct {
		name           string
		payload        map[string]interface{}
		setup          func(*MockUserService)
		expectedStatus int
		expectToken    bool
	}{
		{
			name: "success login",
			payload: map[string]interface{}{
				"email":    "test@example.com",
				"password": "password123",
			},
			setup: func(m *MockUserService) {
				response := &dto.LoginResponse{
					User: dto.UserResponse{
						UUID:  uuid.New(),
						Name:  "Test User",
						Email: "test@example.com",
						Role:  "customer",
					},
					Token: "test-token",
				}
				m.On("Login", mock.Anything, mock.AnythingOfType("*dto.LoginRequest")).Return(response, nil)
			},
			expectedStatus: http.StatusOK,
			expectToken:    true,
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
		{
			name: "missing required fields",
			payload: map[string]interface{}{
				"email": "test@example.com",
			},
			expectedStatus: http.StatusUnprocessableEntity,
			expectToken:    false,
		},
		{
			name: "login service error",
			payload: map[string]interface{}{
				"email":    "test@example.com",
				"password": "password123",
			},
			setup: func(m *MockUserService) {
				m.On("Login", mock.Anything, mock.AnythingOfType("*dto.LoginRequest")).Return(nil, errors.New("user not found"))
			},
			expectedStatus: http.StatusBadRequest,
			expectToken:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := setupRouter()

			mockUserService := new(MockUserService)
			if tt.setup != nil {
				tt.setup(mockUserService)
			}

			mockServiceRegistry := &MockServiceRegistry{
				userService: mockUserService,
			}

			controller := NewUserController(mockServiceRegistry)

			router.POST("/login", controller.Login)

			jsonValue, _ := json.Marshal(tt.payload)
			req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonValue))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectToken {
				var response map[string]interface{}
				json.Unmarshal(w.Body.Bytes(), &response)
				assert.NotNil(t, response["token"])
			}

			mockUserService.AssertExpectations(t)
		})
	}
}

func TestUserController_Register(t *testing.T) {
	tests := []struct {
		name           string
		payload        map[string]interface{}
		setup          func(*MockUserService)
		expectedStatus int
	}{
		{
			name: "success register",
			payload: map[string]interface{}{
				"name":            "New User",
				"email":           "newuser@example.com",
				"password":        "password123",
				"confirmPassword": "password123",
			},
			setup: func(m *MockUserService) {
				response := &dto.RegisterResponse{
					User: dto.UserResponse{
						UUID:  uuid.New(),
						Name:  "New User",
						Email: "newuser@example.com",
					},
				}
				m.On("Register", mock.Anything, mock.AnythingOfType("*dto.RegisterRequest")).Return(response, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "invalid email format",
			payload: map[string]interface{}{
				"name":            "New User",
				"email":           "invalid-email",
				"password":        "password123",
				"confirmPassword": "password123",
			},
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name: "missing required fields",
			payload: map[string]interface{}{
				"email": "test@example.com",
			},
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name: "register service error",
			payload: map[string]interface{}{
				"name":            "New User",
				"email":           "newuser@example.com",
				"password":        "password123",
				"confirmPassword": "password123",
			},
			setup: func(m *MockUserService) {
				m.On("Register", mock.Anything, mock.AnythingOfType("*dto.RegisterRequest")).Return(nil, errors.New("email already exists"))
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := setupRouter()

			mockUserService := new(MockUserService)
			if tt.setup != nil {
				tt.setup(mockUserService)
			}

			mockServiceRegistry := &MockServiceRegistry{
				userService: mockUserService,
			}

			controller := NewUserController(mockServiceRegistry)

			router.POST("/register", controller.Register)

			jsonValue, _ := json.Marshal(tt.payload)
			req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonValue))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if mockUserService != nil {
				mockUserService.AssertExpectations(t)
			}
		})
	}
}

// Mock implementations for testing
type MockServiceRegistry struct {
	userService    *MockUserService
	productService productService.IProductService
	orderService   orderService.IOrderService
	paymentService paymentService.IPaymentService
}

func (m *MockServiceRegistry) GetUser() userService.IUserService {
	return m.userService
}

func (m *MockServiceRegistry) GetProduct() productService.IProductService {
	return m.productService
}

func (m *MockServiceRegistry) GetOrder() orderService.IOrderService {
	return m.orderService
}

func (m *MockServiceRegistry) GetPayment() paymentService.IPaymentService {
	return m.paymentService
}

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.LoginResponse), args.Error(1)
}

func (m *MockUserService) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.RegisterResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.RegisterResponse), args.Error(1)
}

func (m *MockUserService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}
