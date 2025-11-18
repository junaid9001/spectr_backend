package models

import (
	"time"

	"gorm.io/gorm"
)

type RefreshToken struct {
	ID        uint      `gorm:"primaryKey"`
	UserId    uint      `gorm:"not null"`
	Token     string    `gorm:"not null;unique"` //hashed token
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
