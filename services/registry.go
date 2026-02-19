package services

import (
	"github.com/arthurhzna/ecommerce_be_test/repositories"
	userService "github.com/arthurhzna/ecommerce_be_test/services/user.go"
)

type Registry struct {
	repository repositories.IRepositoryRegistry
}

type IServiceRegistry interface {
	GetUser() userService.IUserService
}

func NewServiceRegistry(repository repositories.IRepositoryRegistry) IServiceRegistry {
	return &Registry{repository: repository}
}

func (r *Registry) GetUser() userService.IUserService {
	return userService.NewUserService(r.repository)
}
