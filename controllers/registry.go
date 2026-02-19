package controllers

import (
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
