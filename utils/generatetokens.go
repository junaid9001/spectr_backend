package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/junaid9001/spectr_backend/models"
	"gorm.io/gorm"
)

// access token
func GenerateAccessToken(userId uint, role string) (string, error) {
	secretKey := os.Getenv("JWT_SECRETKEY")

	claims := jwt.MapClaims{
		"userId": userId,
		"role":   role,
		"exp":    time.Now().Add(200 * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secretKey))
}

//random refresh token plain and  hashed

func GenerateRefreshToken() (string, string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", "", err
	}
	token := hex.EncodeToString(b)
	hash := sha256.Sum256([]byte(token))
	return token, hex.EncodeToString(hash[:]), nil

}

//save refresh token in db

func SaveRefreshToken(db *gorm.DB, userId uint, hashedToken string, expirestAt time.Time) error {
	RT := models.RefreshToken{
		UserId:    userId,
		Token:     hashedToken,
		ExpiresAt: expirestAt,
	}
	return db.Create(&RT).Error
}

//validate refresh token

func ValidateRefreshToken(db *gorm.DB, token string) (*models.RefreshToken, error) {
	hash := sha256.Sum256([]byte(token))
	hashedToken := hex.EncodeToString(hash[:])

	var rt models.RefreshToken

	err := db.Where("token=? AND expires_at > ?", hashedToken, time.Now()).First(&rt).Error
	if err != nil {
		return nil, errors.New("invalid or expired refresh token")
	}
	return &rt, nil

}

// delete rt from db
func DeleteRefreshToken(db *gorm.DB, token string) error {

	hash := sha256.Sum256([]byte(token))
	hashedToken := hex.EncodeToString(hash[:])
	return db.Where("token=?", hashedToken).Delete(&models.RefreshToken{}).Error
}
