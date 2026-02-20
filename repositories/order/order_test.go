package order

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	statusConstants "github.com/arthurhzna/ecommerce_be_test/constants/status"
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

func TestOrderRepository_Create(t *testing.T) {
	tests := []struct {
		name    string
		order   *models.Order
		setup   func(sqlmock.Sqlmock)
		wantErr bool
	}{
		{
			name: "success create order",
			order: &models.Order{
				UserID: 1,
				Amount: 100.00,
				Status: statusConstants.OrderStatusPending,
			},
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(`INSERT INTO "orders"`).
					WithArgs(sqlmock.AnyArg(), uint(1), sqlmock.AnyArg(), 100.00, "pending", sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
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

			repo := NewOrderRepository(db)
			result, err := repo.Create(context.Background(), tt.order)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.NotEqual(t, uuid.Nil, result.UUID)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("There were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestOrderRepository_FindByUserID(t *testing.T) {
	tests := []struct {
		name    string
		userID  uint
		setup   func(sqlmock.Sqlmock)
		wantErr bool
	}{
		{
			name:   "success find orders by user id",
			userID: 1,
			setup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "uuid", "user_id", "payment_id", "amount", "status", "paid_at", "created_at", "updated_at"}).
					AddRow(1, uuid.New(), 1, nil, 100.00, "pending", nil, time.Now(), time.Now()).
					AddRow(2, uuid.New(), 1, nil, 200.00, "pending", nil, time.Now(), time.Now())
				mock.ExpectQuery(`SELECT (.+) FROM "orders"`).
					WithArgs(1).
					WillReturnRows(rows)
				orderItemsRows := sqlmock.NewRows([]string{"id", "uuid", "order_id", "product_id", "quantity", "price", "created_at", "updated_at"})
				mock.ExpectQuery(`SELECT (.+) FROM "order_items"`).
					WillReturnRows(orderItemsRows)
			},
			wantErr: false,
		},
		{
			name:   "empty orders",
			userID: 1,
			setup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "uuid", "user_id", "payment_id", "amount", "status", "paid_at", "created_at", "updated_at"})
				mock.ExpectQuery(`SELECT (.+) FROM "orders"`).
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

			repo := NewOrderRepository(db)
			result, err := repo.FindByUserID(context.Background(), tt.userID)

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

func TestOrderRepository_FindByID(t *testing.T) {
	tests := []struct {
		name    string
		orderID uint
		setup   func(sqlmock.Sqlmock)
		wantErr bool
		wantNil bool
	}{
		{
			name:    "success find order by id",
			orderID: 1,
			setup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "uuid", "user_id", "payment_id", "amount", "status", "paid_at", "created_at", "updated_at"}).
					AddRow(1, uuid.New(), 1, nil, 100.00, "pending", nil, time.Now(), time.Now())
				mock.ExpectQuery(`SELECT (.+) FROM "orders"`).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnRows(rows)
				orderItemsRows := sqlmock.NewRows([]string{"id", "uuid", "order_id", "product_id", "quantity", "price", "created_at", "updated_at"})
				mock.ExpectQuery(`SELECT (.+) FROM "order_items"`).
					WillReturnRows(orderItemsRows)
			},
			wantErr: false,
			wantNil: false,
		},
		{
			name:    "order not found",
			orderID: 999,
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT (.+) FROM "orders"`).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			wantErr: true,
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

			repo := NewOrderRepository(db)
			result, err := repo.FindByID(context.Background(), tt.orderID)

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
					assert.Equal(t, tt.orderID, result.ID)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("There were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestOrderRepository_FindByUUID(t *testing.T) {
	tests := []struct {
		name      string
		orderUUID uuid.UUID
		setup     func(sqlmock.Sqlmock)
		wantErr   bool
		wantNil   bool
	}{
		{
			name:      "success find order by uuid",
			orderUUID: uuid.New(),
			setup: func(mock sqlmock.Sqlmock) {
				orderUUID := uuid.New()
				rows := sqlmock.NewRows([]string{"id", "uuid", "user_id", "payment_id", "amount", "status", "paid_at", "created_at", "updated_at"}).
					AddRow(1, orderUUID, 1, nil, 100.00, "pending", nil, time.Now(), time.Now())
				mock.ExpectQuery(`SELECT (.+) FROM "orders"`).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnRows(rows)
				orderItemsRows := sqlmock.NewRows([]string{"id", "uuid", "order_id", "product_id", "quantity", "price", "created_at", "updated_at"})
				mock.ExpectQuery(`SELECT (.+) FROM "order_items"`).
					WillReturnRows(orderItemsRows)
			},
			wantErr: false,
			wantNil: false,
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

			repo := NewOrderRepository(db)
			result, err := repo.FindByUUID(context.Background(), tt.orderUUID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if tt.wantNil {
				assert.Nil(t, result)
			} else {
				assert.NotNil(t, result)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("There were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestOrderRepository_Update(t *testing.T) {
	tests := []struct {
		name    string
		order   *models.Order
		setup   func(sqlmock.Sqlmock)
		wantErr bool
	}{
		{
			name: "success update order",
			order: &models.Order{
				ID:     1,
				UUID:   uuid.New(),
				UserID: 1,
				Amount: 150.00,
				Status: statusConstants.OrderStatusPaid,
			},
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`UPDATE "orders"`).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
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

			repo := NewOrderRepository(db)
			err := repo.Update(context.Background(), tt.order)

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
