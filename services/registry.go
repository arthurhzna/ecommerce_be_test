package services

import (
	"github.com/arthurhzna/ecommerce_be_test/repositories"
	orderService "github.com/arthurhzna/ecommerce_be_test/services/order"
	paymentService "github.com/arthurhzna/ecommerce_be_test/services/payment"
	productService "github.com/arthurhzna/ecommerce_be_test/services/product"
	userService "github.com/arthurhzna/ecommerce_be_test/services/user"
)

type Registry struct {
	repository repositories.IRepositoryRegistry
}

type IServiceRegistry interface {
	GetUser() userService.IUserService
	GetProduct() productService.IProductService
	GetOrder() orderService.IOrderService
	GetPayment() paymentService.IPaymentService
}

func NewServiceRegistry(repository repositories.IRepositoryRegistry) IServiceRegistry {
	return &Registry{repository: repository}
}

func (r *Registry) GetUser() userService.IUserService {
	return userService.NewUserService(r.repository)
}

func (r *Registry) GetProduct() productService.IProductService {
	return productService.NewProductService(r.repository)
}

func (r *Registry) GetOrder() orderService.IOrderService {
	return orderService.NewOrderService(r.repository)
}

func (r *Registry) GetPayment() paymentService.IPaymentService {
	return paymentService.NewPaymentService(r.repository)
}
