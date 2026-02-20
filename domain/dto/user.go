package dto

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	User *UserResponse
	jwt.RegisteredClaims
}

type UserResponse struct {
	UUID  uuid.UUID `json:"uuid"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
	Role  string    `json:"role,omitempty"`
}

type RegisterRequest struct {
	Name            string `json:"name" validate:"required"`
	Password        string `json:"password" validate:"required"`
	ConfirmPassword string `json:"confirmPassword" validate:"required"`
	Email           string `json:"email" validate:"required,email"`
	RoleID          uint   `json:"-"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RegisterResponse struct {
	User UserResponse `json:"user"`
}

type LoginResponse struct {
	User  UserResponse `json:"user"`
	Token string       `json:"token"`
}
