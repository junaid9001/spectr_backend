package controllers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/junaid9001/spectr_backend/models"
	"github.com/junaid9001/spectr_backend/utils"
	"gorm.io/gorm"
)

//done

func GetUserProfile(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		//get user if also safe accertion
		userId, ok := utils.GetUserId(c)
		if !ok {
			return
		}

		var user models.User

		if err := db.Where("id=?", userId).First(&user).Error; err != nil {

			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{
					"status": "failed",
					"error":  err.Error(),
				})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "failed",
				"error":  err.Error(),
			})
			return
		}

		user.HashedPassword = ""

		c.JSON(http.StatusOK, gin.H{
			"status":     "success",
			"id":         user.ID,
			"name":       user.Name,
			"email":      user.Email,
			"role":       user.Role,
			"created_at": user.CreatedAt,
		})
	}
}

func UpdateUserProfile(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			Name string `json:"name" binding:"required,min=1"`
		}

		//get user if also safe accertion
		userId, ok := utils.GetUserId(c)
		if !ok {
			return
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "failed",
				"error":  err.Error(),
			})
			return
		}

		name := strings.TrimSpace(input.Name)

		if len(name) < 1 {
			c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": "name cannot be empty"})
			return
		}

		if err := db.Model(&models.User{}).
			Where("id=?", userId).Updates(map[string]any{"name": name}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "failed",
				"error":  err.Error(),
			})
			return
		}

		user := models.User{}

		if err := db.Where("id=?", userId).First(&user).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"status": "failed", "error": "user not found"})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "failed",
				"error":  err.Error(),
			})
			return

		}

		user.HashedPassword = ""

		c.JSON(http.StatusOK, gin.H{
			"status":     "success",
			"id":         user.ID,
			"name":       user.Name,
			"email":      user.Email,
			"role":       user.Role,
			"created_at": user.CreatedAt,
		})

	}
	//---
}
