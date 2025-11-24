package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/junaid9001/spectr_backend/models"
	"github.com/junaid9001/spectr_backend/utils"
	"gorm.io/gorm"
)

//done

// list all users
func AllUsers(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		var users []models.User

		if err := db.Find(&users).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "failed",
				"error":  err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data":   users,
			"count":  len(users),
		})

	}
}

//update user role

func UpdateUserRole(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idstr := c.Param("id")

		var input struct {
			Role string `json:"role" binding:"required"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}

		if input.Role != "user" && input.Role != "admin" {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "failed",
				"error":  "Invalid role provided. Allowed: 'user', 'admin'",
			})

			return
		}

		userId, err := utils.StringToUint(idstr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "failed",
				"error":  "invalid id format",
			})
			return
		}

		result := db.Model(&models.User{}).Where("id=?", userId).Update("role", input.Role)

		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "error": result.Error.Error()})
			return
		}

		if result.RowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"status": "failed", "error": "user not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "success", "message": "user role is successfully changed"})

	}
}

//update user status (is_blocked)

func UpdateUserStatus(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		idstr := c.Param("id")
		userId, err := utils.StringToUint(idstr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "failed",
				"error":  "invalid id format",
			})
			return
		}

		var input struct {
			IsBlocked *bool `json:"is_blocked" binding:"required"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}

		result := db.Model(&models.User{}).Where("id=?", userId).Update("is_blocked", input.IsBlocked)

		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "error": result.Error.Error()})
			return
		}

		if result.RowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"status": "failed", "error": "user not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "success", "message": "user's status is changed succesfully"})

	}
}
