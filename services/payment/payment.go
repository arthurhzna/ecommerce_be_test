package payment

import (
	"context"
	"fmt"
	"time"

	errOrder "github.com/arthurhzna/ecommerce_be_test/constants/error/order"
	errPayment "github.com/arthurhzna/ecommerce_be_test/constants/error/payment"
	statusConstants "github.com/arthurhzna/ecommerce_be_test/constants/status"
	"github.com/arthurhzna/ecommerce_be_test/domain/dto"
	"github.com/arthurhzna/ecommerce_be_test/domain/models"
	"github.com/arthurhzna/ecommerce_be_test/repositories"
	"github.com/google/uuid"
)

type PaymentService struct {
	repository repositories.IRepositoryRegistry
}

type IPaymentService interface {
	ProcessPayment(ctx context.Context, userID uint, orderUUID uuid.UUID, req *dto.PaymentRequest) (*dto.PaymentResponse, error)
}

func NewPaymentService(repository repositories.IRepositoryRegistry) IPaymentService {
	return &PaymentService{repository: repository}
}

func (s *PaymentService) ProcessPayment(ctx context.Context, userID uint, orderUUID uuid.UUID, req *dto.PaymentRequest) (*dto.PaymentResponse, error) {

	order, err := s.repository.GetOrder().FindByUUID(ctx, orderUUID)
	if err != nil {
		return nil, errOrder.ErrOrderNotFound
	}

	if order.UserID != userID {
		return nil, errOrder.ErrOrderUnauthorized
	}

	if order.Status != statusConstants.OrderStatusPending {
		return nil, errOrder.ErrOrderNotPending
	}

	existingPayment, err := s.repository.GetPayment().FindByOrderID(ctx, order.ID)
	if err != nil {
		return nil, err
	}

	if existingPayment != nil && existingPayment.Status == statusConstants.PaymentStatusPaid {
		return nil, errPayment.ErrPaymentAlreadyPaid
	}

	var payment *models.Payment
	if existingPayment != nil {
		existingPayment.Status = statusConstants.PaymentStatusPending
		payment = existingPayment
	} else {
		payment = &models.Payment{
			OrderID: order.ID,
			Amount:  order.Amount,
			Status:  statusConstants.PaymentStatusPending,
		}
		payment, err = s.repository.GetPayment().Create(ctx, payment)
		if err != nil {
			return nil, err
		}
	}

	transactionID := s.processPaymentGateway(payment, req)

	payment.Status = statusConstants.PaymentStatusPaid
	payment.TransactionID = &transactionID
	now := time.Now()
	payment.PaidAt = &now

	err = s.repository.GetPayment().Update(ctx, payment)
	if err != nil {
		return nil, err
	}

	order.Status = statusConstants.OrderStatusPaid
	order.PaymentID = &payment.ID
	order.PaidAt = &now

	err = s.repository.GetOrder().Update(ctx, order)
	if err != nil {
		return nil, err
	}

	return &dto.PaymentResponse{
		UUID:          payment.UUID,
		OrderID:       payment.OrderID,
		Amount:        payment.Amount,
		Status:        payment.Status.GetStatusString(),
		TransactionID: payment.TransactionID,
		PaidAt:        payment.PaidAt,
		CreatedAt:     payment.CreatedAt,
	}, nil
}

func (s *PaymentService) processPaymentGateway(payment *models.Payment, req *dto.PaymentRequest) string {
	//simulasi payment gateway

	transactionID := fmt.Sprintf("TXN-%s-%d", time.Now().Format("20060102"), payment.ID)

	return transactionID
}
