package controllers

import (
	userController "github.com/arthurhzna/ecommerce_be_test/controllers/user"
	"github.com/arthurhzna/ecommerce_be_test/services"
)

type Registry struct {
	service services.IServiceRegistry
}

type IControllerRegistry interface {
	GetUser() userController.IUserController
}

func NewControllerRegistry(service services.IServiceRegistry) IControllerRegistry {
	return &Registry{service: service}
}

func (u *Registry) GetUser() userController.IUserController {
	return userController.NewUserController(u.service)
}
