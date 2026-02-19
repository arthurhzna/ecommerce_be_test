package usergo

import (
	"context"
	"strings"
	"time"

	"github.com/arthurhzna/ecommerce_be_test/config"
	"github.com/arthurhzna/ecommerce_be_test/constants"
	errUser "github.com/arthurhzna/ecommerce_be_test/constants/error/user"
	"github.com/arthurhzna/ecommerce_be_test/domain/dto"
	"github.com/arthurhzna/ecommerce_be_test/repositories"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repository repositories.IRepositoryRegistry
}

type IUserService interface {
	Login(context.Context, *dto.LoginRequest) (*dto.LoginResponse, error)
	Register(context.Context, *dto.RegisterRequest) (*dto.RegisterResponse, error)
}

func NewUserService(repository repositories.IRepositoryRegistry) IUserService {
	return &UserService{repository: repository}
}

func (u *UserService) Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := u.repository.GetUser().FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, errUser.ErrPasswordInCorrect
	}

	expirationTime := time.Now().Add(time.Duration(config.Config.JwtExpirationTime) * time.Minute).Unix()
	data := &dto.UserResponse{
		UUID:  user.UUID,
		Name:  user.Name,
		Email: user.Email,
		Role:  strings.ToLower(user.Role.Code),
	}

	claims := &dto.Claims{
		User: data,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Unix(expirationTime, 0)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.Config.JwtSecretKey))
	if err != nil {
		return nil, err
	}

	response := &dto.LoginResponse{
		User:  *data,
		Token: tokenString,
	}

	return response, nil
}

func (u *UserService) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.RegisterResponse, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	if u.isEmailExist(ctx, req.Email) {
		return nil, errUser.ErrEmailExist
	}

	if req.Password != req.ConfirmPassword {
		return nil, errUser.ErrPasswordDoesNotMatch
	}

	user, err := u.repository.GetUser().Register(ctx, &dto.RegisterRequest{
		Name:     req.Name,
		Password: string(hashedPassword),
		Email:    req.Email,
		RoleID:   constants.Customer,
	})
	if err != nil {
		return nil, err
	}

	response := &dto.RegisterResponse{
		User: dto.UserResponse{
			UUID:  user.UUID,
			Name:  user.Name,
			Email: user.Email,
		},
	}

	return response, nil
}

func (u *UserService) isEmailExist(ctx context.Context, email string) bool {
	user, err := u.repository.GetUser().FindByEmail(ctx, email)
	if err != nil {
		return false
	}

	if user != nil {
		return true
	}

	return false
}
