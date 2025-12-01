package models

type Filter struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	FilterName string `gorm:"size:50;not null;unique" json:"filter_name"` //names like gender usetype brand
}

type FilterOption struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	FilterID uint   `gorm:"not null" json:"filter_id"`
	Label    string `gorm:"size:50;not null" json:"label"` //what user see like male,daily,meta
}

type ProductFilterOption struct { //links products to filter options
	ProductID      uint `gorm:"primaryKey" json:"product_id"`
	FilterOptionID uint `gorm:"primaryKey" json:"filter_option_id"`
}
