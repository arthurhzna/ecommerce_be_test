package user

import (
	"context"
	"testing"
	"time"

	"errors"

	"github.com/DATA-DOG/go-sqlmock"
	errUser "github.com/arthurhzna/ecommerce_be_test/constants/error/user"
	"github.com/arthurhzna/ecommerce_be_test/domain/dto"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to open mock sql db: %v", err)
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open gorm db: %v", err)
	}

	return gormDB, mock
}

func TestUserRepository_Register(t *testing.T) {
	tests := []struct {
		name    string
		req     *dto.RegisterRequest
		setup   func(sqlmock.Sqlmock)
		wantErr bool
	}{
		{
			name: "success register",
			req: &dto.RegisterRequest{
				Name:     "Test User",
				Email:    "test@example.com",
				Password: "hashedpassword",
				RoleID:   2,
			},
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(`INSERT INTO "users"`).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "duplicate email",
			req: &dto.RegisterRequest{
				Name:     "Test User",
				Email:    "existing@example.com",
				Password: "hashedpassword",
				RoleID:   2,
			},
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(`INSERT INTO "users"`).
					WillReturnError(errors.New("duplicate key"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := setupMockDB(t)
			sqlDB, _ := db.DB()
			defer sqlDB.Close()

			if tt.setup != nil {
				tt.setup(mock)
			}

			repo := NewUserRepository(db)
			result, err := repo.Register(context.Background(), tt.req)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.req.Name, result.Name)
				assert.Equal(t, tt.req.Email, result.Email)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("There were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestUserRepository_FindByEmail(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		setup   func(sqlmock.Sqlmock)
		wantErr bool
		errType error
	}{
		{
			name:  "success find user",
			email: "test@example.com",
			setup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "uuid", "name", "email", "password", "role_id", "created_at", "updated_at"}).
					AddRow(1, uuid.New(), "Test User", "test@example.com", "hashedpassword", 2, time.Now(), time.Now())
				mock.ExpectQuery(`SELECT (.+) FROM "users"`).
					WithArgs("test@example.com", 1).
					WillReturnRows(rows)

				roleRows := sqlmock.NewRows([]string{"id", "code", "name", "created_at", "updated_at"}).
					AddRow(2, "customer", "Customer", time.Now(), time.Now())
				mock.ExpectQuery(`SELECT (.+) FROM "roles"`).
					WillReturnRows(roleRows)
			},
			wantErr: false,
		},
		{
			name:  "user not found",
			email: "notfound@example.com",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT (.+) FROM "users"`).
					WithArgs("notfound@example.com", 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			wantErr: true,
			errType: errUser.ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := setupMockDB(t)
			sqlDB, _ := db.DB()
			defer sqlDB.Close()

			if tt.setup != nil {
				tt.setup(mock)
			}

			repo := NewUserRepository(db)
			result, err := repo.FindByEmail(context.Background(), tt.email)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
				if tt.errType != nil {
					assert.Equal(t, tt.errType, err)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.email, result.Email)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("There were unfulfilled expectations: %s", err)
			}
		})
	}
}
