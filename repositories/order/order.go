package order

import (
	"context"

	"github.com/arthurhzna/ecommerce_be_test/domain/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

type IOrderRepository interface {
	Create(ctx context.Context, order *models.Order) (*models.Order, error)
	FindByUserID(ctx context.Context, userID uint) ([]models.Order, error)
	FindByID(ctx context.Context, id uint) (*models.Order, error)
	FindByUUID(ctx context.Context, orderUUID uuid.UUID) (*models.Order, error)
	Update(ctx context.Context, order *models.Order) error
}

func NewOrderRepository(db *gorm.DB) IOrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Create(ctx context.Context, order *models.Order) (*models.Order, error) {
	order.UUID = uuid.New()

	err := r.db.WithContext(ctx).Create(order).Error
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (r *OrderRepository) FindByUserID(ctx context.Context, userID uint) ([]models.Order, error) {
	var orders []models.Order

	err := r.db.WithContext(ctx).
		Preload("OrderItems").
		Preload("OrderItems.Product").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&orders).Error
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *OrderRepository) FindByID(ctx context.Context, id uint) (*models.Order, error) {
	var order models.Order
	err := r.db.WithContext(ctx).
		Preload("OrderItems").
		Preload("OrderItems.Product").
		Where("id = ?", id).
		First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepository) FindByUUID(ctx context.Context, orderUUID uuid.UUID) (*models.Order, error) {
	var order models.Order
	err := r.db.WithContext(ctx).
		Preload("OrderItems").
		Preload("OrderItems.Product").
		Where("uuid = ?", orderUUID).
		First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepository) Update(ctx context.Context, order *models.Order) error {
	return r.db.WithContext(ctx).Save(order).Error
}
