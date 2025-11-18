package models

import "time"

//users table
type User struct {
	ID             uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name           string    `gorm:"size:50;not null" json:"name"`
	Email          string    `gorm:"size:50;uniqueIndex;not null" json:"email"`
	HashedPassword string    `gorm:"size:255" json:"-"`
	Role           string    `gorm:"size:10;default:user" json:"role"`
	IsBlocked      bool      `gorm:"default:false" json:"is_blocked"`
	IsVerified     bool      `gorm:"default:false" json:"is_verified"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
