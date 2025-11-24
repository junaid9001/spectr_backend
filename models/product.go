package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model //adds id,Create,updated,deletedAt automatically

	Name          string  `gorm:"size:255;not null" json:"name" binding:"required" form:"name"`
	Description   string  `gorm:"type:text" json:"description" form:"description"`
	Price         float64 `gorm:"type:decimal(10,2);not null" json:"price" binding:"required" form:"price" `
	StockQuantity int     `gorm:"not null;default:0" json:"stock_quantity" binding:"required,gte=0" form:"stock_quantity"`
	ImageUrl      string  `gorm:"type:text" json:"image"`
	Category      string  `gorm:"type:varchar(100);default:'smart';index" json:"category" form:"category"`
}
