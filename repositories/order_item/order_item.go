package order_item

import (
	"context"

	"github.com/arthurhzna/ecommerce_be_test/domain/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderItemRepository struct {
	db *gorm.DB
}

type IOrderItemRepository interface {
	CreateBulk(ctx context.Context, items []*models.OrderItem) error
	FindByOrderID(ctx context.Context, orderID uint) ([]models.OrderItem, error)
}

func NewOrderItemRepository(db *gorm.DB) IOrderItemRepository {
	return &OrderItemRepository{db: db}
}

func (r *OrderItemRepository) CreateBulk(ctx context.Context, items []*models.OrderItem) error {
	for _, item := range items {
		item.UUID = uuid.New()
	}

	err := r.db.WithContext(ctx).CreateInBatches(items, 100).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *OrderItemRepository) FindByOrderID(ctx context.Context, orderID uint) ([]models.OrderItem, error) {
	var items []models.OrderItem
	err := r.db.WithContext(ctx).
		Preload("Product").
		Where("order_id = ?", orderID).
		Find(&items).Error
	if err != nil {
		return nil, err
	}
	return items, nil
}
