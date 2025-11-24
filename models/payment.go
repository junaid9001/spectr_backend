package models

import (
	"time"

	"gorm.io/gorm"
)

type Payment struct {
	ID            uint    `gorm:"primaryKey" json:"id"`
	OrderID       uint    `gorm:"index;not null" json:"order_id"`
	Amount        float64 `gorm:"not null" json:"amount"`
	PaymentStatus string  `gorm:"size:32;not null" json:"payment_status"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}
