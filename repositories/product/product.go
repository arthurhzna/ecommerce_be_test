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
	FindByID(ctx context.Context, id uint) (*models.Product, error)
	FindByUUID(ctx context.Context, productUUID uuid.UUID) (*models.Product, error)
	FindByName(ctx context.Context, name string) (*models.Product, error)
	UpdateStock(ctx context.Context, productID uint, quantity int) error
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

func (r *ProductRepository) FindByID(ctx context.Context, id uint) (*models.Product, error) {
	var product models.Product
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepository) FindByUUID(ctx context.Context, productUUID uuid.UUID) (*models.Product, error) {
	var product models.Product
	err := r.db.WithContext(ctx).Where("uuid = ?", productUUID).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepository) FindByName(ctx context.Context, name string) (*models.Product, error) {
	var product models.Product
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&product).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepository) UpdateStock(ctx context.Context, productID uint, quantity int) error {
	return r.db.WithContext(ctx).
		Model(&models.Product{}).
		Where("id = ?", productID).
		Update("stock", gorm.Expr("stock - ?", quantity)).Error
}
