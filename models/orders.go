package models

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	UserID        uint           `gorm:"not null" json:"user_id"`
	TotalAmount   float64        `gorm:"not null" json:"total_amount"`
	Address       string         `gorm:"type:text;not null" json:"address"`
	Status        string         `gorm:"default:'pending';not null" json:"status"`
	PaymentStatus string         `gorm:"size:30;default:'pending;not null" json:"payment_status"`
	CreatedAt     time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	//relation
	OrderItems []OrderItem `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE" json:"order_items"`
}
