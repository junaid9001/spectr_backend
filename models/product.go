package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model //adds id,Create,updated,deletedAt automatically

	Name          string  `gorm:"size:255;not null" json:"name" binding:"required" form:"name"`
	Description   string  `gorm:"type:text" json:"description" form:"description"`
	Price         float64 `gorm:"type:decimal(10,2);not null" json:"price" binding:"required" form:"price" `
	StockQuantity int     `gorm:"not null;default:0" json:"stock_quantity" binding:"required,gte=0" form:"stock_quantity"`
	ImageUrl      string  `gorm:"type:text" json:"image"`
	CategoryID    *uint   `gorm:"constraint:OnDelete:SET NULL;" json:"category_id" form:"category_id"`
	Brand         string  `gorm:"size:30;default:'spectr';index" json:"brand" form:"brand"`
}
