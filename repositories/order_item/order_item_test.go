package order_item

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

func TestOrderItemRepository_CreateBulk(t *testing.T) {
	tests := []struct {
		name    string
		items   []*models.OrderItem
		setup   func(sqlmock.Sqlmock)
		wantErr bool
	}{
		{
			name: "success create bulk order items",
			items: []*models.OrderItem{
				{
					OrderID:   1,
					ProductID: 1,
					Quantity:  2,
					Price:     100.00,
				},
				{
					OrderID:   1,
					ProductID: 2,
					Quantity:  3,
					Price:     200.00,
				},
			},
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(`INSERT INTO "order_items"`).
					WithArgs(
						sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
						sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
						sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
						sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
					).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1).AddRow(2))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name:  "empty items",
			items: []*models.OrderItem{},
			setup: func(mock sqlmock.Sqlmock) {
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

			repo := NewOrderItemRepository(db)
			err := repo.CreateBulk(context.Background(), tt.items)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				for _, item := range tt.items {
					assert.NotEqual(t, uuid.Nil, item.UUID)
				}
			}
			_ = mock
		})
	}
}

func TestOrderItemRepository_FindByOrderID(t *testing.T) {
	tests := []struct {
		name    string
		orderID uint
		setup   func(sqlmock.Sqlmock)
		wantErr bool
	}{
		{
			name:    "success find order items by order id",
			orderID: 1,
			setup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "uuid", "order_id", "product_id", "quantity", "price", "created_at", "updated_at"}).
					AddRow(1, uuid.New(), 1, 1, 2, 100.00, time.Now(), time.Now()).
					AddRow(2, uuid.New(), 1, 2, 3, 200.00, time.Now(), time.Now())
				mock.ExpectQuery(`SELECT (.+) FROM "order_items"`).
					WithArgs(1).
					WillReturnRows(rows)
				mock.ExpectQuery(`SELECT (.+) FROM "products"`).
					WillReturnRows(sqlmock.NewRows([]string{"id", "uuid", "name", "description", "price", "stock", "created_at", "updated_at", "deleted_at"}))
			},
			wantErr: false,
		},
		{
			name:    "empty order items",
			orderID: 1,
			setup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "uuid", "order_id", "product_id", "quantity", "price", "created_at", "updated_at"})
				mock.ExpectQuery(`SELECT (.+) FROM "order_items"`).
					WithArgs(1).
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

			repo := NewOrderItemRepository(db)
			result, err := repo.FindByOrderID(context.Background(), tt.orderID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}
			if tt.setup != nil && !tt.wantErr {
			}
		})
	}
}
