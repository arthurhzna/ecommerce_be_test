package models

import (
	"time"

	"gorm.io/gorm"

	"github.com/google/uuid"
)

type Product struct {
	ID          uint      `gorm:"primaryKey;autoIncrement"`
	UUID        uuid.UUID `gorm:"type:uuid;not null"`
	Name        string    `gorm:"varchar(100);not null"`
	Description string    `gorm:"varchar(255);not null"`
	Price       float64   `gorm:"type:numeric(10,2);not null"`
	Stock       int       `gorm:"not null"`

	CreatedAt *time.Time
	UpdatedAt *time.Time

	DeletedAt *gorm.DeletedAt
}
