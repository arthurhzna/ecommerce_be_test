package repositories

import (
	userRepo "github.com/arthurhzna/ecommerce_be_test/repositories/user"
	"gorm.io/gorm"
)

type Registry struct {
	db *gorm.DB
}

type IRepositoryRegistry interface {
	GetUser() userRepo.IUserRepository
}

func NewRepositoryRegistry(db *gorm.DB) IRepositoryRegistry {
	return &Registry{db: db}
}

func (r *Registry) GetUser() userRepo.IUserRepository {
	return userRepo.NewUserRepository(r.db)
}
