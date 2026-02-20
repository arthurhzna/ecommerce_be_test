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
	orderService "github.com/arthurhzna/ecommerce_be_test/services/order"
	paymentService "github.com/arthurhzna/ecommerce_be_test/services/payment"
	productService "github.com/arthurhzna/ecommerce_be_test/services/product"
	userService "github.com/arthurhzna/ecommerce_be_test/services/user"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestProductController_GetProductsWithoutPagination(t *testing.T) {
	tests := []struct {
		name           string
		setup          func(*MockProductService)
		expectedStatus int
	}{
		{
			name: "success get products",
			setup: func(m *MockProductService) {
				response := &dto.ProductListResponse{
					Products: []dto.ProductResponse{
						{
							UUID:        uuid.New(),
							Name:        "Product 1",
							Description: "Description 1",
							Price:       100.00,
							Stock:       10,
						},
						{
							UUID:        uuid.New(),
							Name:        "Product 2",
							Description: "Description 2",
							Price:       200.00,
							Stock:       20,
						},
					},
				}
				m.On("GetProductsWithoutPagination", mock.Anything).Return(response, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "empty products",
			setup: func(m *MockProductService) {
				response := &dto.ProductListResponse{
					Products: []dto.ProductResponse{},
				}
				m.On("GetProductsWithoutPagination", mock.Anything).Return(response, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "service error",
			setup: func(m *MockProductService) {
				m.On("GetProductsWithoutPagination", mock.Anything).Return(nil, errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			router := gin.New()

			mockProductService := new(MockProductService)
			if tt.setup != nil {
				tt.setup(mockProductService)
			}

			mockServiceRegistry := &MockServiceRegistryForProduct{
				productService: mockProductService,
			}

			controller := NewProductController(mockServiceRegistry)
			router.GET("/products", controller.GetProductsWithoutPagination)

			req, _ := http.NewRequest("GET", "/products", nil)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusOK {
				var response map[string]interface{}
				json.Unmarshal(w.Body.Bytes(), &response)
				assert.Equal(t, "success", response["status"])
				assert.NotNil(t, response["data"])
			}

			mockProductService.AssertExpectations(t)
		})
	}
}

func TestProductController_CreateProduct(t *testing.T) {
	tests := []struct {
		name           string
		payload        map[string]interface{}
		setup          func(*MockProductService)
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
			setup: func(m *MockProductService) {
				response := &dto.ProductResponse{
					UUID:        uuid.New(),
					Name:        "New Product",
					Description: "New Product Description",
					Price:       150.00,
					Stock:       15,
				}
				m.On("CreateProduct", mock.Anything, mock.AnythingOfType("*dto.CreateProductRequest")).Return(response, nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "invalid payload - missing required fields",
			payload: map[string]interface{}{
				"name": "Product",
			},
			setup:          nil,
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name: "invalid payload - invalid name length",
			payload: map[string]interface{}{
				"name":        "AB",
				"description": "Valid description",
				"price":       100.00,
				"stock":       10,
			},
			setup:          nil,
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name: "invalid payload - invalid description length",
			payload: map[string]interface{}{
				"name":        "Valid Product Name",
				"description": "Short",
				"price":       100.00,
				"stock":       10,
			},
			setup:          nil,
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name: "invalid payload - negative price",
			payload: map[string]interface{}{
				"name":        "Valid Product Name",
				"description": "Valid product description",
				"price":       -100.00,
				"stock":       10,
			},
			setup:          nil,
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name: "invalid payload - negative stock",
			payload: map[string]interface{}{
				"name":        "Valid Product Name",
				"description": "Valid product description",
				"price":       100.00,
				"stock":       -10,
			},
			setup:          nil,
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name: "service error",
			payload: map[string]interface{}{
				"name":        "New Product",
				"description": "New Product Description",
				"price":       150.00,
				"stock":       15,
			},
			setup: func(m *MockProductService) {
				m.On("CreateProduct", mock.Anything, mock.AnythingOfType("*dto.CreateProductRequest")).Return(nil, errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			router := gin.New()

			mockProductService := new(MockProductService)
			if tt.setup != nil {
				tt.setup(mockProductService)
			}

			mockServiceRegistry := &MockServiceRegistryForProduct{
				productService: mockProductService,
			}

			controller := NewProductController(mockServiceRegistry)
			router.POST("/products/create", controller.CreateProduct)

			jsonValue, _ := json.Marshal(tt.payload)
			req, _ := http.NewRequest("POST", "/products/create", bytes.NewBuffer(jsonValue))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusCreated {
				var response map[string]interface{}
				json.Unmarshal(w.Body.Bytes(), &response)
				assert.Equal(t, "success", response["status"])
				assert.NotNil(t, response["data"])
			}

			if mockProductService != nil {
				mockProductService.AssertExpectations(t)
			}
		})
	}
}

// Mock implementations
type MockProductService struct {
	mock.Mock
}

func (m *MockProductService) GetProductsWithoutPagination(ctx context.Context) (*dto.ProductListResponse, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ProductListResponse), args.Error(1)
}

func (m *MockProductService) CreateProduct(ctx context.Context, req *dto.CreateProductRequest) (*dto.ProductResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ProductResponse), args.Error(1)
}

type MockServiceRegistryForProduct struct {
	userService    userService.IUserService
	productService *MockProductService
	orderService   orderService.IOrderService
	paymentService paymentService.IPaymentService
}

func (m *MockServiceRegistryForProduct) GetUser() userService.IUserService {
	return m.userService
}

func (m *MockServiceRegistryForProduct) GetProduct() productService.IProductService {
	return m.productService
}

func (m *MockServiceRegistryForProduct) GetOrder() orderService.IOrderService {
	return m.orderService
}

func (m *MockServiceRegistryForProduct) GetPayment() paymentService.IPaymentService {
	return m.paymentService
}
