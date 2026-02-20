package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/arthurhzna/ecommerce_be_test/constants"
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

func TestPaymentController_ProcessPayment(t *testing.T) {
	orderUUID := uuid.New()

	tests := []struct {
		name           string
		orderUUID      string
		payload        map[string]interface{}
		setup          func(*MockPaymentService, *MockUserServiceForPayment)
		expectedStatus int
	}{
		{
			name:      "success process payment",
			orderUUID: orderUUID.String(),
			payload:   map[string]interface{}{},
			setup: func(paymentSvc *MockPaymentService, userSvc *MockUserServiceForPayment) {
				user := &models.User{
					ID:    1,
					Email: "test@example.com",
				}
				userSvc.On("GetUserByEmail", mock.Anything, "test@example.com").Return(user, nil)

				paymentResponse := &dto.PaymentResponse{
					UUID:    uuid.New(),
					OrderID: 1,
					Amount:  100.00,
					Status:  "paid",
				}
				paymentSvc.On("ProcessPayment", mock.Anything, uint(1), orderUUID, mock.AnythingOfType("*dto.PaymentRequest")).Return(paymentResponse, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid uuid format",
			orderUUID:      "invalid-uuid",
			payload:        map[string]interface{}{},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:      "service error",
			orderUUID: orderUUID.String(),
			payload:   map[string]interface{}{},
			setup: func(paymentSvc *MockPaymentService, userSvc *MockUserServiceForPayment) {
				user := &models.User{
					ID:    1,
					Email: "test@example.com",
				}
				userSvc.On("GetUserByEmail", mock.Anything, "test@example.com").Return(user, nil)
				paymentSvc.On("ProcessPayment", mock.Anything, uint(1), orderUUID, mock.AnythingOfType("*dto.PaymentRequest")).Return(nil, errors.New("order not found"))
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			router := gin.New()

			mockPaymentService := new(MockPaymentService)
			mockUserService := new(MockUserServiceForPayment)

			if tt.setup != nil {
				tt.setup(mockPaymentService, mockUserService)
			}

			mockServiceRegistry := &MockServiceRegistryForPayment{
				userService:    mockUserService,
				paymentService: mockPaymentService,
			}

			controller := NewPaymentController(mockServiceRegistry)
			router.POST("/payments/:orderUUID/pay", func(c *gin.Context) {
				// Set user context
				userResponse := &dto.UserResponse{
					UUID:  uuid.New(),
					Email: "test@example.com",
					Role:  "customer",
				}
				c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), constants.UserLogin, userResponse))
				controller.ProcessPayment(c)
			})

			jsonValue, _ := json.Marshal(tt.payload)
			req, _ := http.NewRequest("POST", "/payments/"+tt.orderUUID+"/pay", bytes.NewBuffer(jsonValue))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if mockPaymentService != nil {
				mockPaymentService.AssertExpectations(t)
			}
			if mockUserService != nil {
				mockUserService.AssertExpectations(t)
			}
		})
	}
}

// Mock implementations
type MockPaymentService struct {
	mock.Mock
}

func (m *MockPaymentService) ProcessPayment(ctx context.Context, userID uint, orderUUID uuid.UUID, req *dto.PaymentRequest) (*dto.PaymentResponse, error) {
	args := m.Called(ctx, userID, orderUUID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.PaymentResponse), args.Error(1)
}

type MockUserServiceForPayment struct {
	mock.Mock
}

func (m *MockUserServiceForPayment) Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.LoginResponse), args.Error(1)
}

func (m *MockUserServiceForPayment) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.RegisterResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.RegisterResponse), args.Error(1)
}

func (m *MockUserServiceForPayment) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

type MockServiceRegistryForPayment struct {
	userService    *MockUserServiceForPayment
	paymentService *MockPaymentService
}

func (m *MockServiceRegistryForPayment) GetUser() userService.IUserService {
	return m.userService
}

func (m *MockServiceRegistryForPayment) GetProduct() productService.IProductService {
	return nil
}

func (m *MockServiceRegistryForPayment) GetOrder() orderService.IOrderService {
	return nil
}

func (m *MockServiceRegistryForPayment) GetPayment() paymentService.IPaymentService {
	return m.paymentService
}
