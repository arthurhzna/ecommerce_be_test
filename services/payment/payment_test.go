package payment

import (
	"context"
	"errors"
	"testing"
	"time"

	errOrder "github.com/arthurhzna/ecommerce_be_test/constants/error/order"
	errPayment "github.com/arthurhzna/ecommerce_be_test/constants/error/payment"
	statusConstants "github.com/arthurhzna/ecommerce_be_test/constants/status"
	"github.com/arthurhzna/ecommerce_be_test/domain/dto"
	"github.com/arthurhzna/ecommerce_be_test/domain/models"
	"github.com/arthurhzna/ecommerce_be_test/tests/helpers"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPaymentService_ProcessPayment(t *testing.T) {
	orderUUID := uuid.New()
	userID := uint(1)
	otherUserID := uint(2)

	tests := []struct {
		name      string
		userID    uint
		orderUUID uuid.UUID
		req       *dto.PaymentRequest
		setup     func(*helpers.MockRepositoryRegistry)
		wantErr   bool
		errType   error
	}{
		{
			name:      "success process payment - new payment",
			userID:    userID,
			orderUUID: orderUUID,
			req:       &dto.PaymentRequest{},
			setup: func(m *helpers.MockRepositoryRegistry) {
				order := &models.Order{
					ID:     1,
					UUID:   orderUUID,
					UserID: userID,
					Amount: 100.00,
					Status: statusConstants.OrderStatusPending,
				}
				m.OrderRepo.On("FindByUUID", context.Background(), orderUUID).Return(order, nil)

				m.PaymentRepo.On("FindByOrderID", context.Background(), uint(1)).Return(nil, nil)

				payment := &models.Payment{
					ID:      1,
					UUID:    uuid.New(),
					OrderID: 1,
					Amount:  100.00,
					Status:  statusConstants.PaymentStatusPending,
				}
				m.PaymentRepo.On("Create", context.Background(), mock.AnythingOfType("*models.Payment")).Return(payment, nil)
				m.PaymentRepo.On("Update", context.Background(), mock.AnythingOfType("*models.Payment")).Return(nil)
				m.OrderRepo.On("Update", context.Background(), mock.AnythingOfType("*models.Order")).Return(nil)
			},
			wantErr: false,
		},
		{
			name:      "order not found",
			userID:    userID,
			orderUUID: orderUUID,
			req:       &dto.PaymentRequest{},
			setup: func(m *helpers.MockRepositoryRegistry) {
				m.OrderRepo.On("FindByUUID", context.Background(), orderUUID).Return(nil, errors.New("order not found"))
			},
			wantErr: true,
			errType: errOrder.ErrOrderNotFound,
		},
		{
			name:      "unauthorized - order belongs to different user",
			userID:    otherUserID,
			orderUUID: orderUUID,
			req:       &dto.PaymentRequest{},
			setup: func(m *helpers.MockRepositoryRegistry) {
				order := &models.Order{
					ID:     1,
					UUID:   orderUUID,
					UserID: userID,
					Amount: 100.00,
					Status: statusConstants.OrderStatusPending,
				}
				m.OrderRepo.On("FindByUUID", context.Background(), orderUUID).Return(order, nil)
			},
			wantErr: true,
			errType: errOrder.ErrOrderUnauthorized,
		},
		{
			name:      "order not pending",
			userID:    userID,
			orderUUID: orderUUID,
			req:       &dto.PaymentRequest{},
			setup: func(m *helpers.MockRepositoryRegistry) {
				order := &models.Order{
					ID:     1,
					UUID:   orderUUID,
					UserID: userID,
					Amount: 100.00,
					Status: statusConstants.OrderStatusPaid,
				}
				m.OrderRepo.On("FindByUUID", context.Background(), orderUUID).Return(order, nil)
			},
			wantErr: true,
			errType: errOrder.ErrOrderNotPending,
		},
		{
			name:      "payment already paid",
			userID:    userID,
			orderUUID: orderUUID,
			req:       &dto.PaymentRequest{},
			setup: func(m *helpers.MockRepositoryRegistry) {
				order := &models.Order{
					ID:     1,
					UUID:   orderUUID,
					UserID: userID,
					Amount: 100.00,
					Status: statusConstants.OrderStatusPending,
				}
				m.OrderRepo.On("FindByUUID", context.Background(), orderUUID).Return(order, nil)

				existingPayment := &models.Payment{
					ID:      1,
					OrderID: 1,
					Amount:  100.00,
					Status:  statusConstants.PaymentStatusPaid,
				}
				m.PaymentRepo.On("FindByOrderID", context.Background(), uint(1)).Return(existingPayment, nil)
			},
			wantErr: true,
			errType: errPayment.ErrPaymentAlreadyPaid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRegistry := helpers.NewMockRepositoryRegistry()

			if tt.setup != nil {
				tt.setup(mockRegistry)
			}

			service := NewPaymentService(mockRegistry)
			result, err := service.ProcessPayment(context.Background(), tt.userID, tt.orderUUID, tt.req)

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
				assert.Equal(t, "paid", result.Status)
			}

			mockRegistry.OrderRepo.AssertExpectations(t)
			mockRegistry.PaymentRepo.AssertExpectations(t)
		})
	}
}

func stringPtr(s string) *string {
	return &s
}

func timePtr(t time.Time) *time.Time {
	return &t
}

func uintPtr(u uint) *uint {
	return &u
}
