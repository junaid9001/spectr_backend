package models

import (
	"time"

	"gorm.io/gorm"
)

type OrderItem struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	OrderID    uint           `gorm:"index;not null" json:"order_id"`
	ProductID  uint           `gorm:"index;not null" json:"product_id"`
	UnitPrice  float64        `gorm:"not null" json:"unit_price"`
	Quantity   int            `gorm:"not null" json:"quantity"`
	TotalPrice float64        `gorm:"not null" json:"total_price"`
	CreatedAt  time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
