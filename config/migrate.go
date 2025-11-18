package config

import (
	"fmt"
	"log"

	"github.com/junaid9001/spectr_backend/models"
)

func MigrateAll() {
	err := DB.AutoMigrate(&models.User{}, &models.Otp{}, &models.RefreshToken{})

	if err != nil {
		log.Fatal("Table migration failed", err.Error())
		return
	}

	fmt.Print("All models migrated")
}
