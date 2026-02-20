package usergo

import (
	"context"
	"testing"

	"github.com/arthurhzna/ecommerce_be_test/constants"
	errUser "github.com/arthurhzna/ecommerce_be_test/constants/error/user"
	"github.com/arthurhzna/ecommerce_be_test/domain/dto"
	"github.com/arthurhzna/ecommerce_be_test/domain/models"
	"github.com/arthurhzna/ecommerce_be_test/tests/helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestUserService_Login(t *testing.T) {
	tests := []struct {
		name    string
		req     *dto.LoginRequest
		setup   func(*helpers.MockUserRepository)
		wantErr bool
		errType error
	}{
		{
			name: "success login",
			req: &dto.LoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			setup: func(m *helpers.MockUserRepository) {
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
				user := &models.User{
					Email:    "test@example.com",
					Password: string(hashedPassword),
					Name:     "Test User",
					Role: models.Role{
						Code: "customer",
					},
				}
				m.On("FindByEmail", context.Background(), "test@example.com").Return(user, nil)
			},
			wantErr: false,
		},
		{
			name: "user not found",
			req: &dto.LoginRequest{
				Email:    "notfound@example.com",
				Password: "password123",
			},
			setup: func(m *helpers.MockUserRepository) {
				m.On("FindByEmail", context.Background(), "notfound@example.com").Return(nil, errUser.ErrUserNotFound)
			},
			wantErr: true,
			errType: errUser.ErrUserNotFound,
		},
		{
			name: "wrong password",
			req: &dto.LoginRequest{
				Email:    "test@example.com",
				Password: "wrongpassword",
			},
			setup: func(m *helpers.MockUserRepository) {
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
				user := &models.User{
					Email:    "test@example.com",
					Password: string(hashedPassword),
					Name:     "Test User",
					Role: models.Role{
						Code: "customer",
					},
				}
				m.On("FindByEmail", context.Background(), "test@example.com").Return(user, nil)
			},
			wantErr: true,
			errType: errUser.ErrPasswordInCorrect,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(helpers.MockUserRepository)
			mockRegistry := helpers.NewMockRepositoryRegistry()
			mockRegistry.UserRepo = mockRepo

			if tt.setup != nil {
				tt.setup(mockRepo)
			}

			service := NewUserService(mockRegistry)
			result, err := service.Login(context.Background(), tt.req)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
				if tt.errType != nil {
					assert.Equal(t, tt.errType, err)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.NotEmpty(t, result.Token)
				assert.Equal(t, tt.req.Email, result.User.Email)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserService_Register(t *testing.T) {
	tests := []struct {
		name    string
		req     *dto.RegisterRequest
		setup   func(*helpers.MockUserRepository)
		wantErr bool
		errType error
	}{
		{
			name: "success register",
			req: &dto.RegisterRequest{
				Name:            "New User",
				Email:           "newuser@example.com",
				Password:        "password123",
				ConfirmPassword: "password123",
			},
			setup: func(m *helpers.MockUserRepository) {
				m.On("FindByEmail", context.Background(), "newuser@example.com").Return(nil, errUser.ErrUserNotFound)
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
				user := &models.User{
					Name:     "New User",
					Email:    "newuser@example.com",
					Password: string(hashedPassword),
					RoleID:   constants.Customer,
				}
				m.On("Register", context.Background(), mock.MatchedBy(func(req *dto.RegisterRequest) bool {
					return req.Email == "newuser@example.com" && req.Name == "New User"
				})).Return(user, nil)
			},
			wantErr: false,
		},
		{
			name: "email already exists",
			req: &dto.RegisterRequest{
				Name:            "Existing User",
				Email:           "existing@example.com",
				Password:        "password123",
				ConfirmPassword: "password123",
			},
			setup: func(m *helpers.MockUserRepository) {
				user := &models.User{
					Email: "existing@example.com",
				}
				m.On("FindByEmail", context.Background(), "existing@example.com").Return(user, nil)
			},
			wantErr: true,
			errType: errUser.ErrEmailExist,
		},
		{
			name: "password does not match",
			req: &dto.RegisterRequest{
				Name:            "New User",
				Email:           "newuser@example.com",
				Password:        "password123",
				ConfirmPassword: "differentpassword",
			},
			setup: func(m *helpers.MockUserRepository) {
				m.On("FindByEmail", context.Background(), "newuser@example.com").Return(nil, errUser.ErrUserNotFound)
			},
			wantErr: true,
			errType: errUser.ErrPasswordDoesNotMatch,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(helpers.MockUserRepository)
			mockRegistry := helpers.NewMockRepositoryRegistry()
			mockRegistry.UserRepo = mockRepo

			if tt.setup != nil {
				tt.setup(mockRepo)
			}

			service := NewUserService(mockRegistry)
			result, err := service.Register(context.Background(), tt.req)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
				if tt.errType != nil {
					assert.Equal(t, tt.errType, err)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.req.Email, result.User.Email)
				assert.Equal(t, tt.req.Name, result.User.Name)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserService_GetUserByEmail(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		setup   func(*helpers.MockUserRepository)
		wantErr bool
	}{
		{
			name:  "success get user",
			email: "test@example.com",
			setup: func(m *helpers.MockUserRepository) {
				user := &models.User{
					Email: "test@example.com",
					Name:  "Test User",
				}
				m.On("FindByEmail", context.Background(), "test@example.com").Return(user, nil)
			},
			wantErr: false,
		},
		{
			name:  "user not found",
			email: "notfound@example.com",
			setup: func(m *helpers.MockUserRepository) {
				m.On("FindByEmail", context.Background(), "notfound@example.com").Return(nil, errUser.ErrUserNotFound)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(helpers.MockUserRepository)
			mockRegistry := helpers.NewMockRepositoryRegistry()
			mockRegistry.UserRepo = mockRepo

			if tt.setup != nil {
				tt.setup(mockRepo)
			}

			service := NewUserService(mockRegistry)
			result, err := service.GetUserByEmail(context.Background(), tt.email)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.email, result.Email)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
