package utils

import (
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

// returs userId and Role
func ValidateJwt(tokenStr string) (int, string, error) {

	secretKey := os.Getenv("JWT_SECRETKEY")

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return 0, "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Extract userId
		userIDFloat, ok := claims["userId"].(float64)
		if !ok {
			return 0, "", fmt.Errorf("invalid userId in token")
		}

		role, ok := claims["role"].(string)
		if !ok {
			return 0, "", fmt.Errorf("invalid role in token")
		}

		return int(userIDFloat), role, nil

	}

	return 0, "", fmt.Errorf("invalid token")

}
