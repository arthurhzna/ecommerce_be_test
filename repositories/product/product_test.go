package product

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/arthurhzna/ecommerce_be_test/domain/models"
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

func TestProductRepository_FindAllWithoutPagination(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(sqlmock.Sqlmock)
		wantErr bool
	}{
		{
			name: "success find all products",
			setup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "uuid", "name", "description", "price", "stock", "created_at", "updated_at", "deleted_at"}).
					AddRow(1, uuid.New(), "Product 1", "Description 1", 100.00, 10, time.Now(), time.Now(), nil).
					AddRow(2, uuid.New(), "Product 2", "Description 2", 200.00, 20, time.Now(), time.Now(), nil)
				mock.ExpectQuery(`SELECT (.+) FROM "products"`).
					WillReturnRows(rows)
			},
			wantErr: false,
		},
		{
			name: "empty products",
			setup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "uuid", "name", "description", "price", "stock", "created_at", "updated_at"})
				mock.ExpectQuery(`SELECT (.+) FROM "products"`).
					WillReturnRows(rows)
			},
			wantErr: false,
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

			repo := NewProductRepository(db)
			result, err := repo.FindAllWithoutPagination(context.Background())

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("There were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestProductRepository_Create(t *testing.T) {
	tests := []struct {
		name    string
		product *models.Product
		setup   func(sqlmock.Sqlmock)
		wantErr bool
	}{
		{
			name: "success create product",
			product: &models.Product{
				Name:        "New Product",
				Description: "New Description",
				Price:       150.00,
				Stock:       15,
			},
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(`INSERT INTO "products"`).
					WithArgs(sqlmock.AnyArg(), "New Product", "New Description", 150.00, 15, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mock.ExpectCommit()
			},
			wantErr: false,
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

			repo := NewProductRepository(db)
			result, err := repo.Create(context.Background(), tt.product)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.product.Name, result.Name)
				assert.NotEqual(t, uuid.Nil, result.UUID)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("There were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestProductRepository_FindByName(t *testing.T) {
	tests := []struct {
		name     string
		prodName string
		setup    func(sqlmock.Sqlmock)
		wantErr  bool
		wantNil  bool
	}{
		{
			name:     "success find product by name",
			prodName: "Test Product",
			setup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "uuid", "name", "description", "price", "stock", "created_at", "updated_at", "deleted_at"}).
					AddRow(1, uuid.New(), "Test Product", "Description", 100.00, 10, time.Now(), time.Now(), nil)
				mock.ExpectQuery(`SELECT (.+) FROM "products"`).
					WithArgs("Test Product", 1).
					WillReturnRows(rows)
			},
			wantErr: false,
			wantNil: false,
		},
		{
			name:     "product not found",
			prodName: "Non Existent",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT (.+) FROM "products"`).
					WithArgs("Non Existent", 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			wantErr: false,
			wantNil: true,
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

			repo := NewProductRepository(db)
			result, err := repo.FindByName(context.Background(), tt.prodName)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if tt.wantNil {
				assert.Nil(t, result)
			} else {
				assert.NotNil(t, result)
				if result != nil {
					assert.Equal(t, tt.prodName, result.Name)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("There were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestProductRepository_UpdateStock(t *testing.T) {
	tests := []struct {
		name      string
		productID uint
		quantity  int
		setup     func(sqlmock.Sqlmock)
		wantErr   bool
	}{
		{
			name:      "success update stock",
			productID: 1,
			quantity:  5,
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`UPDATE "products"`).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
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

			repo := NewProductRepository(db)
			err := repo.UpdateStock(context.Background(), tt.productID, tt.quantity)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("There were unfulfilled expectations: %s", err)
			}
		})
	}
}
