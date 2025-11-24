package models

import "gorm.io/gorm"

type Wishlist struct {
	gorm.Model
	UserId    uint    `gorm:"not null;index:idx_user_product,unique" json:"user_id"`
	ProductId uint    `gorm:"not null;index:idx_user_product,unique" json:"product_id"`
	Product   Product `gorm:"foreignKey:ProductId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"product"`
}
