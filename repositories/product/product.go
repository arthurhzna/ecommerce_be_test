package product

import (
	"context"

	"github.com/arthurhzna/ecommerce_be_test/domain/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductRepository struct {
	db *gorm.DB
}

type IProductRepository interface {
	FindAllWithoutPagination(ctx context.Context) ([]models.Product, error)
	Create(ctx context.Context, product *models.Product) (*models.Product, error)
}

func NewProductRepository(db *gorm.DB) IProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) FindAllWithoutPagination(ctx context.Context) ([]models.Product, error) {
	var products []models.Product
	err := r.db.WithContext(ctx).Find(&products).Error
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (r *ProductRepository) Create(ctx context.Context, product *models.Product) (*models.Product, error) {

	product.UUID = uuid.New()

	err := r.db.WithContext(ctx).Create(product).Error
	if err != nil {
		return nil, err
	}
	return product, nil
}
