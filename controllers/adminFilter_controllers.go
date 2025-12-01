package controllers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/junaid9001/spectr_backend/models"
	"github.com/junaid9001/spectr_backend/utils"
	"gorm.io/gorm"
)

//done postman

// add filter
func AddFilter(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			FilterName string `json:"filter_name" binding:"required"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": err.Error()})
			return
		}

		filter := models.Filter{
			FilterName: input.FilterName,
		}

		if err := db.Create(&filter).Error; err != nil {
			//error handling for duplicate name
			if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "UNIQUE") {
				c.JSON(http.StatusBadRequest, gin.H{
					"status": "failed",
					"error":  "filter name already exists",
				})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "error": "db error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "success", "message": filter.FilterName + "filter_added", "filter_id": filter.ID})

	}
}

//----------------*-----------------------

func AddFilterOption(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			FilterId uint   `json:"filter_id" binding:"required"`
			Label    string `json:"label" binding:"required"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": err.Error()})
			return
		}

		filter := models.Filter{}

		if err := db.First(&filter, input.FilterId).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"status": "failed", "error": err.Error()})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "error": "db error"})
			return
		}

		var filterOption = models.FilterOption{
			FilterID: input.FilterId,
			Label:    input.Label,
		}

		result := db.Create(&filterOption)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "error": result.Error.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "success", "filter_option_id": filterOption.ID})

	}
}

//----------------*-----------------------

// admin can link filter option to product
func AddFilterOptionToProduct(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		var input struct {
			FilterOptionID uint `json:"filter_option_id" binding:"required"`
			ProductID      uint `json:"product_id" binding:"required"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": err.Error()})
			return
		}

		var product models.Product

		if err := db.First(&product, input.ProductID).Error; err != nil {

			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"status": "failed", "error": err.Error()})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "error": "db error"})
			return
		}

		var filterOption models.FilterOption
		if err := db.First(&filterOption, input.FilterOptionID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"status": "failed", "error": "filter option not found"})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "error": "db error"})
			return
		}

		if err := db.Transaction(func(tx *gorm.DB) error {

			//link product to filter option

			linkProductToFilter := models.ProductFilterOption{
				FilterOptionID: input.FilterOptionID,
				ProductID:      product.ID,
			}

			if err := tx.Create(&linkProductToFilter).Error; err != nil {
				return err
			}

			return nil
		}); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "error": "db error"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"status": "success"})

	}
}

func ViewProductsByFilterOptionId(db *gorm.DB) gin.HandlerFunc {
	db.Debug()
	return func(c *gin.Context) {
		FilterOptionId, err := utils.StringToUint(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": "invalid id format"})
		}

		var products []models.Product

		if err := db.Joins("JOIN product_filter_options pfo ON pfo.product_id=products.id").
			Where("pfo.filter_option_id=?", FilterOptionId).Find(&products).
			Error; err != nil {

			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"status": "failed", "error": "filter option not found"})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "error": "db error"})
			return

		}

		if len(products) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"status": "failed", "error": "no products found for this filter option"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "success", "products": products})
	}
}
