package models

import (
	"time"

	constants "github.com/arthurhzna/ecommerce_be_test/constants/status"

	"github.com/google/uuid"
)

type Payment struct {
	ID            uint                    `gorm:"primaryKey;autoIncrement"`
	UUID          uuid.UUID               `gorm:"type:uuid;not null"`
	OrderID       uint                    `gorm:"not null"`
	Amount        float64                 `gorm:"not null"`
	Status        constants.PaymentStatus `gorm:"type:varchar(20);not null"`
	TransactionID *string                 `gorm:"type:varchar(100);default:null"`
	Description   *string                 `gorm:"type:text;default:null"`
	PaidAt        *time.Time
	ExpiredAt     *time.Time
	CreatedAt     *time.Time
	UpdatedAt     *time.Time
}
