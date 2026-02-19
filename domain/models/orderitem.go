package models

import (
	"time"

	"github.com/google/uuid"
)

type OrderItem struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	UUID      uuid.UUID `gorm:"type:uuid;not null"`
	OrderID   uint      `gorm:"not null"`
	ProductID uint      `gorm:"not null"`
	Quantity  int       `gorm:"not null"`
	Price     float64   `gorm:"not null"`

	CreatedAt time.Time
	UpdatedAt time.Time

	Order   Order   `gorm:"foreignKey:OrderID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	Product Product `gorm:"foreignKey:ProductID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
}
