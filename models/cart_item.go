package models

import "gorm.io/gorm"

type CartItem struct {
	gorm.Model
	UserId     uint    `gorm:"not null;index" json:"user_id"`
	ProductId  uint    `gorm:"not null;index" json:"product_id"`
	Quantity   int     `gorm:"default:1" json:"quantity"`
	UnitPrice  float64 `gorm:"type:decimal(10,2);not null" json:"unit_price"`
	TotalPrice float64 `gorm:"type:decimal(10,2);not null" json:"total_price"`
	Product    Product `gorm:"foreignKey:ProductId;constraint:onUpdate:CASCADE,onDelete:CASCADE" json:"product"`
	//clean cart if Pro or USe deleted
}
