package payment

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

func TestPaymentRepository_Create(t *testing.T) {
	tests := []struct {
		name    string
		payment *models.Payment
		setup   func(sqlmock.Sqlmock)
		wantErr bool
	}{
		{
			name: "success create payment",
			payment: &models.Payment{
				OrderID: 1,
				Amount:  100.00,
				Status:  statusConstants.PaymentStatusPending,
			},
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(`INSERT INTO "payments"`).
					WithArgs(sqlmock.AnyArg(), uint(1), 100.00, "pending", sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"id", "transaction_id", "description"}).AddRow(1, nil, nil))
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

			repo := NewPaymentRepository(db)
			result, err := repo.Create(context.Background(), tt.payment)

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

func TestPaymentRepository_FindByOrderID(t *testing.T) {
	tests := []struct {
		name    string
		orderID uint
		setup   func(sqlmock.Sqlmock)
		wantErr bool
		wantNil bool
	}{
		{
			name:    "success find payment by order id",
			orderID: 1,
			setup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "uuid", "order_id", "amount", "status", "transaction_id", "description", "paid_at", "expired_at", "created_at", "updated_at"}).
					AddRow(1, uuid.New(), 1, 100.00, "pending", nil, nil, nil, nil, time.Now(), time.Now())
				mock.ExpectQuery(`SELECT (.+) FROM "payments"`).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnRows(rows)
			},
			wantErr: false,
			wantNil: false,
		},
		{
			name:    "payment not found",
			orderID: 999,
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT (.+) FROM "payments"`).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
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

			repo := NewPaymentRepository(db)
			result, err := repo.FindByOrderID(context.Background(), tt.orderID)

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
					assert.Equal(t, tt.orderID, result.OrderID)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("There were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestPaymentRepository_Update(t *testing.T) {
	tests := []struct {
		name    string
		payment *models.Payment
		setup   func(sqlmock.Sqlmock)
		wantErr bool
	}{
		{
			name: "success update payment",
			payment: &models.Payment{
				ID:      1,
				UUID:    uuid.New(),
				OrderID: 1,
				Amount:  100.00,
				Status:  statusConstants.PaymentStatusPaid,
			},
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`UPDATE "payments"`).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
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

			repo := NewPaymentRepository(db)
			err := repo.Update(context.Background(), tt.payment)

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
