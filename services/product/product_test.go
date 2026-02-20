package product

import (
	"context"
	"errors"
	"testing"

	errProduct "github.com/arthurhzna/ecommerce_be_test/constants/error/product"
	"github.com/arthurhzna/ecommerce_be_test/domain/dto"
	"github.com/arthurhzna/ecommerce_be_test/domain/models"
	"github.com/arthurhzna/ecommerce_be_test/tests/helpers"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestProductService_GetProductsWithoutPagination(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*helpers.MockProductRepository)
		wantErr bool
	}{
		{
			name: "success get products",
			setup: func(m *helpers.MockProductRepository) {
				products := []models.Product{
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
				}
				m.On("FindAllWithoutPagination", context.Background()).Return(products, nil)
			},
			wantErr: false,
		},
		{
			name: "empty products",
			setup: func(m *helpers.MockProductRepository) {
				products := []models.Product{}
				m.On("FindAllWithoutPagination", context.Background()).Return(products, nil)
			},
			wantErr: false,
		},
		{
			name: "database error",
			setup: func(m *helpers.MockProductRepository) {
				m.On("FindAllWithoutPagination", context.Background()).Return(nil, errors.New("database error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(helpers.MockProductRepository)
			mockRegistry := helpers.NewMockRepositoryRegistry()
			mockRegistry.ProductRepo = mockRepo

			if tt.setup != nil {
				tt.setup(mockRepo)
			}

			service := NewProductService(mockRegistry)
			result, err := service.GetProductsWithoutPagination(context.Background())

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				if result != nil {
					assert.NotNil(t, result.Products)
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestProductService_CreateProduct(t *testing.T) {
	tests := []struct {
		name    string
		req     *dto.CreateProductRequest
		setup   func(*helpers.MockProductRepository)
		wantErr bool
		errType error
	}{
		{
			name: "success create product",
			req: &dto.CreateProductRequest{
				Name:        "New Product",
				Description: "New Product Description",
				Price:       150.00,
				Stock:       15,
			},
			setup: func(m *helpers.MockProductRepository) {
				m.On("FindByName", context.Background(), "New Product").Return(nil, nil)
				product := &models.Product{
					UUID:        uuid.New(),
					Name:        "New Product",
					Description: "New Product Description",
					Price:       150.00,
					Stock:       15,
				}
				m.On("Create", context.Background(), mock.MatchedBy(func(p *models.Product) bool {
					return p.Name == "New Product" && p.Price == 150.00
				})).Return(product, nil)
			},
			wantErr: false,
		},
		{
			name: "product name already exists",
			req: &dto.CreateProductRequest{
				Name:        "Existing Product",
				Description: "Description",
				Price:       100.00,
				Stock:       10,
			},
			setup: func(m *helpers.MockProductRepository) {
				existingProduct := &models.Product{
					Name: "Existing Product",
				}
				m.On("FindByName", context.Background(), "Existing Product").Return(existingProduct, nil)
			},
			wantErr: true,
			errType: errProduct.ErrProductNameAlreadyExist,
		},
		{
			name: "database error on create",
			req: &dto.CreateProductRequest{
				Name:        "New Product",
				Description: "Description",
				Price:       100.00,
				Stock:       10,
			},
			setup: func(m *helpers.MockProductRepository) {
				m.On("FindByName", context.Background(), "New Product").Return(nil, nil)
				m.On("Create", context.Background(), mock.Anything).Return(nil, errors.New("database error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(helpers.MockProductRepository)
			mockRegistry := helpers.NewMockRepositoryRegistry()
			mockRegistry.ProductRepo = mockRepo

			if tt.setup != nil {
				tt.setup(mockRepo)
			}

			service := NewProductService(mockRegistry)
			result, err := service.CreateProduct(context.Background(), tt.req)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
				if tt.errType != nil {
					assert.Equal(t, tt.errType, err)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.req.Name, result.Name)
				assert.Equal(t, tt.req.Price, result.Price)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
