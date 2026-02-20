package payment

import (
	"context"

	"github.com/arthurhzna/ecommerce_be_test/domain/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentRepository struct {
	db *gorm.DB
}

type IPaymentRepository interface {
	Create(ctx context.Context, payment *models.Payment) (*models.Payment, error)
	FindByOrderID(ctx context.Context, orderID uint) (*models.Payment, error)
	Update(ctx context.Context, payment *models.Payment) error
}

func NewPaymentRepository(db *gorm.DB) IPaymentRepository {
	return &PaymentRepository{db: db}
}

func (r *PaymentRepository) Create(ctx context.Context, payment *models.Payment) (*models.Payment, error) {
	payment.UUID = uuid.New()

	err := r.db.WithContext(ctx).Create(payment).Error
	if err != nil {
		return nil, err
	}
	return payment, nil
}

func (r *PaymentRepository) FindByOrderID(ctx context.Context, orderID uint) (*models.Payment, error) {
	var payment models.Payment
	err := r.db.WithContext(ctx).
		Where("order_id = ?", orderID).
		First(&payment).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &payment, nil
}

func (r *PaymentRepository) Update(ctx context.Context, payment *models.Payment) error {
	return r.db.WithContext(ctx).Save(payment).Error
}
