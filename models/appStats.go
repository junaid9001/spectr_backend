package models

type AppStats struct {
	ID                uint    `gorm:"primaryKey" json:"id"`
	TotalUsers        int     `gorm:"not null;default:0" json:"total_users"`
	TotalRevenue      float64 `gorm:"not null;default:0" json:"total_revenue"`
	TotalProductsSold int     `gorm:"not null;default:0" json:"total_products_sold"`
	TotalSales        int     `gorm:"not null;default:0" json:"total_sales"`
}
