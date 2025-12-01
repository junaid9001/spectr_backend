package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/junaid9001/spectr_backend/models"
	"github.com/junaid9001/spectr_backend/utils"
	"gorm.io/gorm"
)

//done postman

// search product by name (public) (done by query  parameter) tested
func SearchProduct(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		name := c.Query("name")

		if name == "" {
			c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": "search something"})
			return
		}

		//partial match 2*% == 1*%

		search := fmt.Sprintf("%%%s%%", name)

		var products []models.Product

		// ILIKE is caseintensice

		if err := db.Where("name ILIKE ?", search).Find(&products).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "error": "db error"})
			return
		}

		if len(products) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"status": "failed", "error": "no products found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "success", "data": products})

	}
}

//filter product by category(smart/luxury) query based ?category=

func FilterProductByCategoryID(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		categoryIdStr := c.Query("category_id")

		if categoryIdStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": "choose a category"})
			return
		}

		categoryId, err := utils.StringToUint(categoryIdStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id format"})
			return
		}

		var selected models.Category
		if err := db.First(&selected, categoryId).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{
					"status": "failed",
					"error":  "category not found",
				})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "error": "db error"})
			return
		}

		var Categories []models.Category

		if err := db.Find(&Categories).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "error": "db error"})
			return
		}

		childrenMap := make(map[uint][]uint)
		for _, cat := range Categories {
			if cat.ParentID != nil {
				parent := *cat.ParentID
				childrenMap[parent] = append(childrenMap[parent], cat.ID)
			}
		}

		toVisit := []uint{categoryId}
		visited := make(map[uint]bool)
		var allIDs []uint

		for len(toVisit) > 0 {

			n := len(toVisit) - 1
			current := toVisit[n]
			toVisit = toVisit[:n]

			if visited[current] {
				continue
			}
			visited[current] = true
			allIDs = append(allIDs, current)

			for _, childID := range childrenMap[current] {
				if !visited[childID] {
					toVisit = append(toVisit, childID)
				}
			}
		}

		var products []models.Product
		if err := db.Where("category_id IN ?", allIDs).Find(&products).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "error": "db error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data":   products,
		})

	}
}

// filter product by brand (query based) ?brand=
func FilterProductByBrand(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		brand := c.Query("brand")

		if brand == "" {
			c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": "choose a category"})
			return
		}

		search := fmt.Sprintf("%%%s%%", brand)

		products := []models.Product{}

		if err := db.Where("brand ILIKE ?", search).Find(&products).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "error": "db error"})
			return
		}

		if len(products) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"status": "failed", "error": "brand not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "success", "data": products})
	}
}

//filter by price (query based) ?min_price=xx&,max_price=xx (tested)

func FilterProductByPrice(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		priceMin := c.Query("price_min")
		pricrMax := c.Query("price_max")

		if priceMin == "" && pricrMax == "" {
			c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": "invalid price range"})
			return
		}

		var products []models.Product

		if err := db.Where("price >= ? AND price <= ?", priceMin, pricrMax).Find(&products).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "error": "db error"})
			return
		}

		if len(products) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"status": "failed", "error": "invalid price range"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "success", "data": products})
	}
}
