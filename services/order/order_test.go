package order

import (
	"context"
	"errors"
	"testing"

	errProduct "github.com/arthurhzna/ecommerce_be_test/constants/error/product"
	statusConstants "github.com/arthurhzna/ecommerce_be_test/constants/status"
	"github.com/arthurhzna/ecommerce_be_test/domain/dto"
	"github.com/arthurhzna/ecommerce_be_test/domain/models"
	"github.com/arthurhzna/ecommerce_be_test/tests/helpers"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestOrderService_CreateOrder(t *testing.T) {
	productUUID1 := uuid.New()
	productUUID2 := uuid.New()

	tests := []struct {
		name    string
		userID  uint
		req     *dto.CreateOrderRequest
		setup   func(*helpers.MockRepositoryRegistry)
		wantErr bool
		errType error
	}{
		{
			name:   "success create order",
			userID: 1,
			req: &dto.CreateOrderRequest{
				Items: []dto.OrderItemRequest{
					{
						ProductUUID: productUUID1,
						Quantity:    2,
					},
					{
						ProductUUID: productUUID2,
						Quantity:    1,
					},
				},
			},
			setup: func(m *helpers.MockRepositoryRegistry) {
				product1 := &models.Product{
					ID:    1,
					UUID:  productUUID1,
					Name:  "Product 1",
					Price: 100.00,
					Stock: 10,
				}
				product2 := &models.Product{
					ID:    2,
					UUID:  productUUID2,
					Name:  "Product 2",
					Price: 200.00,
					Stock: 5,
				}

				m.ProductRepo.On("FindByUUID", context.Background(), productUUID1).Return(product1, nil)
				m.ProductRepo.On("FindByUUID", context.Background(), productUUID2).Return(product2, nil)

				order := &models.Order{
					ID:     1,
					UUID:   uuid.New(),
					UserID: 1,
					Amount: 400.00,
					Status: statusConstants.OrderStatusPending,
				}
				m.OrderRepo.On("Create", context.Background(), mock.AnythingOfType("*models.Order")).Return(order, nil)
				m.OrderItemRepo.On("CreateBulk", context.Background(), mock.Anything).Return(nil)

				m.ProductRepo.On("UpdateStock", context.Background(), uint(1), 2).Return(nil)
				m.ProductRepo.On("UpdateStock", context.Background(), uint(2), 1).Return(nil)

				orderWithItems := &models.Order{
					ID:     1,
					UUID:   order.UUID,
					UserID: 1,
					Amount: 400.00,
					Status: statusConstants.OrderStatusPending,
					OrderItems: []models.OrderItem{
						{
							OrderID:   1,
							ProductID: 1,
							Quantity:  2,
							Price:     100.00,
							Product:   *product1,
						},
						{
							OrderID:   1,
							ProductID: 2,
							Quantity:  1,
							Price:     200.00,
							Product:   *product2,
						},
					},
				}
				m.OrderRepo.On("FindByID", context.Background(), uint(1)).Return(orderWithItems, nil)
			},
			wantErr: false,
		},
		{
			name:   "product not found",
			userID: 1,
			req: &dto.CreateOrderRequest{
				Items: []dto.OrderItemRequest{
					{
						ProductUUID: productUUID1,
						Quantity:    2,
					},
				},
			},
			setup: func(m *helpers.MockRepositoryRegistry) {
				m.ProductRepo.On("FindByUUID", context.Background(), productUUID1).Return(nil, errors.New("product not found"))
			},
			wantErr: true,
			errType: errProduct.ErrProductNotFound,
		},
		{
			name:   "insufficient stock",
			userID: 1,
			req: &dto.CreateOrderRequest{
				Items: []dto.OrderItemRequest{
					{
						ProductUUID: productUUID1,
						Quantity:    20,
					},
				},
			},
			setup: func(m *helpers.MockRepositoryRegistry) {
				product1 := &models.Product{
					ID:    1,
					UUID:  productUUID1,
					Name:  "Product 1",
					Price: 100.00,
					Stock: 10,
				}
				m.ProductRepo.On("FindByUUID", context.Background(), productUUID1).Return(product1, nil)
			},
			wantErr: true,
			errType: errProduct.ErrInsufficientStock,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRegistry := helpers.NewMockRepositoryRegistry()

			if tt.setup != nil {
				tt.setup(mockRegistry)
			}

			service := NewOrderService(mockRegistry)
			result, err := service.CreateOrder(context.Background(), tt.userID, tt.req)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
				if tt.errType != nil {
					assert.Equal(t, tt.errType, err)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.NotEqual(t, uuid.Nil, result.UUID)
			}

			mockRegistry.ProductRepo.AssertExpectations(t)
			mockRegistry.OrderRepo.AssertExpectations(t)
			mockRegistry.OrderItemRepo.AssertExpectations(t)
		})
	}
}

func TestOrderService_GetMyOrders(t *testing.T) {
	tests := []struct {
		name    string
		userID  uint
		setup   func(*helpers.MockRepositoryRegistry)
		wantErr bool
	}{
		{
			name:   "success get my orders",
			userID: 1,
			setup: func(m *helpers.MockRepositoryRegistry) {
				orders := []models.Order{
					{
						ID:     1,
						UUID:   uuid.New(),
						UserID: 1,
						Amount: 100.00,
						Status: statusConstants.OrderStatusPending,
						OrderItems: []models.OrderItem{
							{
								OrderID:   1,
								ProductID: 1,
								Quantity:  1,
								Price:     100.00,
								Product: models.Product{
									ID:   1,
									Name: "Product 1",
								},
							},
						},
					},
				}
				m.OrderRepo.On("FindByUserID", context.Background(), uint(1)).Return(orders, nil)
			},
			wantErr: false,
		},
		{
			name:   "empty orders",
			userID: 1,
			setup: func(m *helpers.MockRepositoryRegistry) {
				orders := []models.Order{}
				m.OrderRepo.On("FindByUserID", context.Background(), uint(1)).Return(orders, nil)
			},
			wantErr: false,
		},
		{
			name:   "database error",
			userID: 1,
			setup: func(m *helpers.MockRepositoryRegistry) {
				m.OrderRepo.On("FindByUserID", context.Background(), uint(1)).Return(nil, errors.New("database error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRegistry := helpers.NewMockRepositoryRegistry()

			if tt.setup != nil {
				tt.setup(mockRegistry)
			}

			service := NewOrderService(mockRegistry)
			result, err := service.GetMyOrders(context.Background(), tt.userID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.NotNil(t, result.Orders)
			}

			mockRegistry.OrderRepo.AssertExpectations(t)
		})
	}
}
