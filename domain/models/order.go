package models

import (
	"time"

	constants "github.com/arthurhzna/ecommerce_be_test/constants/status"

	"github.com/google/uuid"
)

type Order struct {
	ID        uint                  `gorm:"primaryKey;autoIncrement"`
	UUID      uuid.UUID             `gorm:"type:uuid;not null"`
	UserID    uint                  `gorm:"not null"`
	PaymentID uint                  `gorm:"not null"`
	Amount    float64               `gorm:"type:decimal(10,2);not null"`
	Status    constants.OrderStatus `gorm:"type:varchar(20);not null"`
	PaidAt    *time.Time            `gorm:"type:timestamp"`

	CreatedAt *time.Time
	UpdatedAt *time.Time

	OrderItems []OrderItem `gorm:"foreignKey:OrderID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	Payment    Payment     `gorm:"foreignKey:PaymentID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
}
