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
