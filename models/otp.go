package models

import "time"

type Otp struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserId    uint      `gorm:"not null" json:"user_id"`
	OtpCode   string    `gorm:"size:10;not null" json:"otp_code"`
	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`
	Purpose   string    `gorm:"size:20;not null" json:"purpose"`
	IsUsed    bool      `gorm:"default:false" json:"is_used"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}
