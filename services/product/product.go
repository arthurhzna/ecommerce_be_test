package product

import (
	"context"

	"github.com/arthurhzna/ecommerce_be_test/domain/dto"
	"github.com/arthurhzna/ecommerce_be_test/domain/models"
	"github.com/arthurhzna/ecommerce_be_test/repositories"
)

type ProductService struct {
	repository repositories.IRepositoryRegistry
}

type IProductService interface {
	GetProductsWithoutPagination(context.Context) (*dto.ProductListResponse, error)
	CreateProduct(context.Context, *dto.CreateProductRequest) (*dto.ProductResponse, error)
}

func NewProductService(repository repositories.IRepositoryRegistry) IProductService {
	return &ProductService{repository: repository}
}

func (p *ProductService) GetProductsWithoutPagination(ctx context.Context) (*dto.ProductListResponse, error) {
	products, err := p.repository.GetProduct().FindAllWithoutPagination(ctx)
	if err != nil {
		return nil, err
	}

	productsResponse := make([]dto.ProductResponse, len(products))
	for i, product := range products {
		productsResponse[i] = dto.ProductResponse{
			UUID:        product.UUID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			Stock:       product.Stock,
		}
	}
	response := dto.ProductListResponse{
		Products: productsResponse,
	}

	return &response, nil
}

func (p *ProductService) CreateProduct(ctx context.Context, req *dto.CreateProductRequest) (*dto.ProductResponse, error) {

	product, err := p.repository.GetProduct().Create(ctx, &models.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
	})

	if err != nil {
		return nil, err
	}

	response := dto.ProductResponse{
		UUID:        product.UUID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
	}

	return &response, nil
}
