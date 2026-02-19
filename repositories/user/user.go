package user

import (
	"context"
	"errors"

	errWrap "github.com/arthurhzna/ecommerce_be_test/common/error"
	errConstant "github.com/arthurhzna/ecommerce_be_test/constants/error"
	errUser "github.com/arthurhzna/ecommerce_be_test/constants/error/user"

	"github.com/arthurhzna/ecommerce_be_test/domain/dto"
	"github.com/arthurhzna/ecommerce_be_test/domain/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

type IUserRepository interface {
	Register(context.Context, *dto.RegisterRequest) (*models.User, error)
	FindByEmail(context.Context, string) (*models.User, error)
}

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Register(ctx context.Context, req *dto.RegisterRequest) (*models.User, error) {
	user := models.User{
		UUID:     uuid.New(),
		Name:     req.Name,
		Password: req.Password,
		Email:    req.Email,
		RoleID:   req.RoleID,
	}

	err := r.db.WithContext(ctx).Create(&user).Error
	if err != nil {
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}

	return &user, nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).
		Preload("Role").
		Where("email = ?", email).
		First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errUser.ErrUserNotFound
		}
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}
	return &user, nil
}
