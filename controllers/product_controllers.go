package controllers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/junaid9001/spectr_backend/models"
	"github.com/junaid9001/spectr_backend/utils"
	"gorm.io/gorm"
)

// create new product (admin)
func CreateProduct(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		var product models.Product

		//bind text fields 1st then files using c.PostForm(key name of the file)

		if err := c.ShouldBind(&product); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "failed",
				"error":  err.Error(),
			})
			return
		}

		file, err := c.FormFile("image")
		if err == nil {
			//create a unique file name

			fileName := fmt.Sprintf("%d_%s", time.Now().Unix(), filepath.Base(file.Filename))

			//define where to save

			uploadPath := "uploads/" + fileName

			//save
			if err := c.SaveUploadedFile(file, uploadPath); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
				return
			}

			//Store ONLY the path string in the database
			product.ImageUrl = "/" + uploadPath

		}

		result := db.Create(&product)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "error": result.Error.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"status":  "success",
			"product": product,
		})

	}
}

//all products (public endpoint)

func GetAllProducts(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		products := []models.Product{}

		if err := db.Find(&products).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "failed",
				"error":  err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data":   products,
		})
	}
}

// get product by id (public endpoint)
func GetProductByID(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idstr := c.Param("id")

		id, err := utils.StringToUint(idstr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "failed",
				"error":  "invalid id format",
			})
			return
		}

		var product models.Product

		if err := db.First(&product, id).Error; err != nil {

			if err == gorm.ErrRecordNotFound {
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

		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data":   product,
		})
	}
}

// struct for update product
type ProductUpdateInput struct {
	Name          *string  `json:"name"`
	Description   *string  `json:"description"`
	Price         *float64 `json:"price"`
	StockQuantity *int     `json:"stock_quantity"`
	Category      *string  `json:"category"`
}

//update product info by id

func UpdateProductByID(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		idstr := c.Param("id")

		id, err := utils.StringToUint(idstr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "failed",
				"error":  "invalid id format",
			})
			return
		}

		var input ProductUpdateInput

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "failed",
				"error":  err.Error(),
			})
			return
		}

		var product models.Product

		if err := db.First(&product, id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
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

		//transaction (returning nil =commit/ err=rollback)
		if err := db.Transaction(func(tx *gorm.DB) error {
			if input.Name != nil {
				product.Name = strings.TrimSpace(*input.Name)
			}

			if input.Description != nil {
				product.Description = strings.TrimSpace(*input.Description)
			}

			if input.Price != nil {
				product.Price = *input.Price
			}

			if input.StockQuantity != nil {
				product.StockQuantity = *input.StockQuantity
			}

			if input.Category != nil {
				product.Category = *input.Category
			}

			if product.Price < 0 {
				return fmt.Errorf("price cant be less that zero")
			}

			if err := tx.Save(&product).Error; err != nil {
				return err
			}

			return nil

		}); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "success", "data": product})
	}
}

//soft delete product by id

func DeleteProductByID(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		idstr := c.Param("id")

		id, err := utils.StringToUint(idstr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "failed",
				"error":  "invalid id format",
			})
			return
		}

		result := db.Delete(&models.Product{}, id)

		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "error": result.Error.Error()})
			return
		}

		if result.RowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"status": "failed", "error": "product not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "success"})
		//restore by db.unscoped update deleted_at to nil

	}
}
