package models

import "time"

type Role struct {
	ID   uint   `gorm:"primaryKey;autoIncrement"`
	Name string `gorm:"type:varchar(20);not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
