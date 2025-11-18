package services

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/junaid9001/spectr_backend/config"
	"github.com/junaid9001/spectr_backend/models"
)

// creates, stores, and emails a 6-digit OTP
func GenerateOtp(userId uint, email, purpose string) (string, error) {
	otp, err := GenerateRandomOtp(6)
	if err != nil {
		return "", err
	}

	otpEnrty := models.Otp{
		UserId:    userId,
		OtpCode:   otp,
		ExpiresAt: time.Now().Add(10 * time.Minute),
		Purpose:   purpose,
		IsUsed:    false,
	}
	if err := config.DB.Create(&otpEnrty).Error; err != nil {
		return "", err
	}
	subject := "Your otp code "

	body := fmt.Sprintf("Your OTP for %s is: %s. It expires in 10 minutes.", purpose, otp)

	if err := SendEmail(email, subject, body); err != nil {
		return "", err
	}
	return otp, nil

}

// generate random otp
func GenerateRandomOtp(length int) (string, error) {
	const digits = "0123456789"
	otp := make([]byte, length)

	for i := range otp {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			return "", err
		}

		otp[i] = digits[num.Int64()]

	}
	return string(otp), nil

}

//validate otp

func ValidateOtp(userId uint, otp, purpose string) (bool, error) {

	var entry models.Otp

	err := config.DB.Where("user_id=? AND purpose=? AND is_used=?", userId, purpose, false).Order("created_at DESC").First(&entry).Error
	if err != nil {
		return false, fmt.Errorf("otp not found or already used")
	}

	if time.Now().After(entry.ExpiresAt) {
		return false, fmt.Errorf("otp expired")
	}

	if entry.OtpCode != otp {
		return false, fmt.Errorf("invalid otp")
	}

	entry.IsUsed = true
	if err := config.DB.Save(&entry).Error; err != nil {
		return false, fmt.Errorf("failed to update otp status: %w", err)
	}

	if purpose == "signup" {
		if err := config.DB.Model(&models.User{}).Where("id=?", userId).Update("is_verified", true).Error; err != nil {
			return false, fmt.Errorf("failed to verify user: %w", err)
		}
	}

	return true, nil

}
