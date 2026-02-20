package repositories

import (
	orderRepo "github.com/arthurhzna/ecommerce_be_test/repositories/order"
	orderItemRepo "github.com/arthurhzna/ecommerce_be_test/repositories/order_item"
	paymentRepo "github.com/arthurhzna/ecommerce_be_test/repositories/payment"
	productRepo "github.com/arthurhzna/ecommerce_be_test/repositories/product"
	userRepo "github.com/arthurhzna/ecommerce_be_test/repositories/user"
	"gorm.io/gorm"
)

type Registry struct {
	db *gorm.DB
}

type IRepositoryRegistry interface {
	GetUser() userRepo.IUserRepository
	GetProduct() productRepo.IProductRepository
	GetOrder() orderRepo.IOrderRepository
	GetOrderItem() orderItemRepo.IOrderItemRepository
	GetPayment() paymentRepo.IPaymentRepository
}

func NewRepositoryRegistry(db *gorm.DB) IRepositoryRegistry {
	return &Registry{db: db}
}

func (r *Registry) GetUser() userRepo.IUserRepository {
	return userRepo.NewUserRepository(r.db)
}

func (r *Registry) GetProduct() productRepo.IProductRepository {
	return productRepo.NewProductRepository(r.db)
}

func (r *Registry) GetOrder() orderRepo.IOrderRepository {
	return orderRepo.NewOrderRepository(r.db)
}

func (r *Registry) GetOrderItem() orderItemRepo.IOrderItemRepository {
	return orderItemRepo.NewOrderItemRepository(r.db)
}

func (r *Registry) GetPayment() paymentRepo.IPaymentRepository {
	return paymentRepo.NewPaymentRepository(r.db)
}
