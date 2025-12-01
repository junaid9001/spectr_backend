package models

type Category struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	CategoryName string `gorm:"size:50;not null;unique" json:"category_name"`
	ParentID     *uint  `gorm:"constraint:OnDelete:CASCADE" json:"parent_id"` //nil = root
}
