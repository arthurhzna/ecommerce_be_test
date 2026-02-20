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

func TestOrderController_CreateOrder(t *testing.T) {
	productUUID := uuid.New()

	tests := []struct {
		name           string
		payload        map[string]interface{}
		setup          func(*MockOrderService, *MockUserService)
		expectedStatus int
	}{
		{
			name: "success create order",
			payload: map[string]interface{}{
				"items": []map[string]interface{}{
					{
						"productUUID": productUUID.String(),
						"quantity":    2,
					},
				},
			},
			setup: func(orderSvc *MockOrderService, userSvc *MockUserService) {
				user := &models.User{
					ID:    1,
					Email: "test@example.com",
				}
				userSvc.On("GetUserByEmail", mock.Anything, "test@example.com").Return(user, nil)

				orderResponse := &dto.OrderResponse{
					UUID:   uuid.New(),
					UserID: 1,
					Amount: 200.00,
					Status: "pending",
					Items: []dto.OrderItemResponse{
						{
							UUID:        uuid.New(),
							ProductID:   1,
							ProductName: "Product 1",
							Quantity:    2,
							Price:       100.00,
							Subtotal:    200.00,
						},
					},
				}
				orderSvc.On("CreateOrder", mock.Anything, uint(1), mock.AnythingOfType("*dto.CreateOrderRequest")).Return(orderResponse, nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "invalid payload - missing items",
			payload: map[string]interface{}{
				"items": []map[string]interface{}{},
			},
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name: "invalid payload - missing quantity",
			payload: map[string]interface{}{
				"items": []map[string]interface{}{
					{
						"productUUID": productUUID.String(),
					},
				},
			},
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name: "service error",
			payload: map[string]interface{}{
				"items": []map[string]interface{}{
					{
						"productUUID": productUUID.String(),
						"quantity":    2,
					},
				},
			},
			setup: func(orderSvc *MockOrderService, userSvc *MockUserService) {
				user := &models.User{
					ID:    1,
					Email: "test@example.com",
				}
				userSvc.On("GetUserByEmail", mock.Anything, "test@example.com").Return(user, nil)
				orderSvc.On("CreateOrder", mock.Anything, uint(1), mock.AnythingOfType("*dto.CreateOrderRequest")).Return(nil, errors.New("product not found"))
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			router := gin.New()

			mockOrderService := new(MockOrderService)
			mockUserService := new(MockUserService)

			if tt.setup != nil {
				tt.setup(mockOrderService, mockUserService)
			}

			mockServiceRegistry := &MockServiceRegistryForOrder{
				userService:  mockUserService,
				orderService: mockOrderService,
			}

			controller := NewOrderController(mockServiceRegistry)
			router.POST("/orders", func(c *gin.Context) {
				// Set user context
				userResponse := &dto.UserResponse{
					UUID:  uuid.New(),
					Email: "test@example.com",
					Role:  "customer",
				}
				c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), constants.UserLogin, userResponse))
				controller.CreateOrder(c)
			})

			jsonValue, _ := json.Marshal(tt.payload)
			req, _ := http.NewRequest("POST", "/orders", bytes.NewBuffer(jsonValue))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if mockOrderService != nil {
				mockOrderService.AssertExpectations(t)
			}
			if mockUserService != nil {
				mockUserService.AssertExpectations(t)
			}
		})
	}
}

func TestOrderController_GetMyOrders(t *testing.T) {
	tests := []struct {
		name           string
		setup          func(*MockOrderService, *MockUserService)
		expectedStatus int
	}{
		{
			name: "success get my orders",
			setup: func(orderSvc *MockOrderService, userSvc *MockUserService) {
				user := &models.User{
					ID:    1,
					Email: "test@example.com",
				}
				userSvc.On("GetUserByEmail", mock.Anything, "test@example.com").Return(user, nil)

				orderList := &dto.OrderListResponse{
					Orders: []dto.OrderResponse{
						{
							UUID:   uuid.New(),
							UserID: 1,
							Amount: 100.00,
							Status: "pending",
						},
					},
				}
				orderSvc.On("GetMyOrders", mock.Anything, uint(1)).Return(orderList, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "service error",
			setup: func(orderSvc *MockOrderService, userSvc *MockUserService) {
				user := &models.User{
					ID:    1,
					Email: "test@example.com",
				}
				userSvc.On("GetUserByEmail", mock.Anything, "test@example.com").Return(user, nil)
				orderSvc.On("GetMyOrders", mock.Anything, uint(1)).Return(nil, errors.New("database error"))
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			router := gin.New()

			mockOrderService := new(MockOrderService)
			mockUserService := new(MockUserService)

			if tt.setup != nil {
				tt.setup(mockOrderService, mockUserService)
			}

			mockServiceRegistry := &MockServiceRegistryForOrder{
				userService:  mockUserService,
				orderService: mockOrderService,
			}

			controller := NewOrderController(mockServiceRegistry)
			router.GET("/orders/my", func(c *gin.Context) {
				// Set user context
				userResponse := &dto.UserResponse{
					UUID:  uuid.New(),
					Email: "test@example.com",
					Role:  "customer",
				}
				c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), constants.UserLogin, userResponse))
				controller.GetMyOrders(c)
			})

			req, _ := http.NewRequest("GET", "/orders/my", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if mockOrderService != nil {
				mockOrderService.AssertExpectations(t)
			}
			if mockUserService != nil {
				mockUserService.AssertExpectations(t)
			}
		})
	}
}

// Mock implementations
type MockOrderService struct {
	mock.Mock
}

func (m *MockOrderService) CreateOrder(ctx context.Context, userID uint, req *dto.CreateOrderRequest) (*dto.OrderResponse, error) {
	args := m.Called(ctx, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.OrderResponse), args.Error(1)
}

func (m *MockOrderService) GetMyOrders(ctx context.Context, userID uint) (*dto.OrderListResponse, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.OrderListResponse), args.Error(1)
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

type MockServiceRegistryForOrder struct {
	userService    *MockUserService
	productService productService.IProductService
	orderService   *MockOrderService
	paymentService paymentService.IPaymentService
}

func (m *MockServiceRegistryForOrder) GetUser() userService.IUserService {
	return m.userService
}

func (m *MockServiceRegistryForOrder) GetProduct() productService.IProductService {
	return m.productService
}

func (m *MockServiceRegistryForOrder) GetOrder() orderService.IOrderService {
	return m.orderService
}

func (m *MockServiceRegistryForOrder) GetPayment() paymentService.IPaymentService {
	return m.paymentService
}
