package controllers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/junaid9001/spectr_backend/models"
	"github.com/junaid9001/spectr_backend/utils"
	"gorm.io/gorm"
)

//done

// add product to wishlist
func AddToWishlist(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, ok := utils.GetUserId(c)
		if !ok {
			return //response are send in func
		}

		var input struct {
			ProductId uint `json:"product_id" binding:"required"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": err.Error()})
			return
		}

		var product models.Product

		if err := db.First(&product, input.ProductId).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"status": "failed", "error": err.Error()})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "error": err.Error()})
			return
		}

		var checkAlreadyExists models.Wishlist
		//if already in wish list
		if err := db.Where("user_id=? AND product_id=?", userId, input.ProductId).
			First(&checkAlreadyExists).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{"status": "failed", "error": "product already in wishlist"})
			return
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {

			c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "error": "db error"})
			return
		}

		usersWishlist := models.Wishlist{
			UserId:    userId,
			ProductId: product.ID,
		}

		if err := db.Create(&usersWishlist).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"status": "success", "data": usersWishlist})

	}
}

//view wish list(get)

func GetWishList(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, ok := utils.GetUserId(c)
		if !ok {
			return
		}

		var usersWishList []models.Wishlist

		if err := db.Where("user_id=?", userId).Find(&usersWishList).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "success", "data": usersWishList})
	}
}

//delete prod from wishlist by id

func DeleteFromWishList(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idstr := c.Param("product_id")
		ProdId, err := utils.StringToUint(idstr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": "invalid id format"})
			return
		}

		userId, ok := utils.GetUserId(c)
		if !ok {
			return
		}

		result := db.Unscoped().Where("user_id=? AND product_id=?", userId, ProdId).Delete(&models.Wishlist{})
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "error": result.Error.Error()})
			return
		}

		if result.RowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"status": "failed", "error": "wishlist item not found"})
			return
		}

		c.Status(http.StatusNoContent)
	}
}
