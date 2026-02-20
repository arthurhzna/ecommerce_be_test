package controllers

import (
	orderController "github.com/arthurhzna/ecommerce_be_test/controllers/order"
	paymentController "github.com/arthurhzna/ecommerce_be_test/controllers/payment"
	productController "github.com/arthurhzna/ecommerce_be_test/controllers/product"
	userController "github.com/arthurhzna/ecommerce_be_test/controllers/user"
	"github.com/arthurhzna/ecommerce_be_test/services"
)

type Registry struct {
	service services.IServiceRegistry
}

type IControllerRegistry interface {
	GetUser() userController.IUserController
	GetProduct() productController.IProductController
	GetOrder() orderController.IOrderController
	GetPayment() paymentController.IPaymentController
}

func NewControllerRegistry(service services.IServiceRegistry) IControllerRegistry {
	return &Registry{service: service}
}

func (u *Registry) GetUser() userController.IUserController {
	return userController.NewUserController(u.service)
}

func (u *Registry) GetProduct() productController.IProductController {
	return productController.NewProductController(u.service)
}

func (u *Registry) GetOrder() orderController.IOrderController {
	return orderController.NewOrderController(u.service)
}

func (u *Registry) GetPayment() paymentController.IPaymentController {
	return paymentController.NewPaymentController(u.service)
}
