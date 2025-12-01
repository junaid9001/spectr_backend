package controllers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/junaid9001/spectr_backend/models"
	"github.com/junaid9001/spectr_backend/utils"
	"gorm.io/gorm"
)

//done postman

// add new category
func AddCategory(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		var input struct {
			CategoryName string `json:"category_name" binding:"required"`
			ParentID     *uint  `json:"parent_id"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": err.Error()})
			return
		}

		//check if parent category exists
		if input.ParentID != nil {
			var parentCat models.Category

			if err := db.First(&parentCat, *input.ParentID).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					c.JSON(http.StatusBadRequest, gin.H{
						"status": "failed",
						"error":  "parent category not found",
					})
					return
				}

				c.JSON(http.StatusInternalServerError, gin.H{
					"status": "failed",
					"error":  "db error",
				})
				return
			}
		}

		category := models.Category{
			CategoryName: input.CategoryName,
			ParentID:     input.ParentID,
		}

		result := db.Create(&category)

		if result.Error != nil {
			if strings.Contains(result.Error.Error(), "duplicate") || strings.Contains(result.Error.Error(), "UNIQUE") {

				c.JSON(http.StatusBadRequest, gin.H{
					"status": "failed",
					"error":  "category name already exists",
				})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "error": "db error"})
			return
		}

		c.JSON(http.StatusOK,
			gin.H{"status": "success", "message": "category " + input.CategoryName + " created", "category_id": category.ID})

	}
}

//delete category by catgID

func DeleteCategoryByID(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		CategoryID, err := utils.StringToUint(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": "invalid id format"})
			return
		}

		var category models.Category

		if err := db.First(&category, CategoryID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": "category not found"})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "error": "db error"})
			return
		}

		if err := db.Delete(&category).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "error": "db error"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "success", "message": "category deleted"})
	}
}

//all categories

func AllCategories(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		var Categories []models.Category

		if err := db.Find(&Categories).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "error": "db error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "success", "data": Categories})
	}
}

type JsTreeNode struct {
	ID     string `json:"id"`
	Parent string `json:"parent"`
	Text   string `json:"text"`
	IsLeaf bool   `json:"is_leaf"`
}

func CategoryTree(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var categories []models.Category

		if err := db.Find(&categories).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "error": "db error"})
			return
		}

		//category id and its child count  no child == leaf (im showing prods by leaf only)
		childCount := make(map[uint]int)
		for _, cat := range categories {
			if cat.ParentID != nil {
				childCount[*cat.ParentID]++
			}
		}
		nodes := make([]JsTreeNode, 0, len(categories))

		for _, cat := range categories {
			node := JsTreeNode{
				ID:   fmt.Sprintf("%d", cat.ID),
				Text: cat.CategoryName,
			}

			if cat.ParentID == nil {
				node.Parent = "#"
			} else {
				node.Parent = fmt.Sprintf("%d", *cat.ParentID)
			}

			// leaf if no children
			node.IsLeaf = childCount[cat.ID] == 0

			nodes = append(nodes, node)
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data":   nodes,
		})

	}
}

func EditCategoryByID(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		categoryId, err := utils.StringToUint(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": "invalid id format"})
			return
		}

		var input struct {
			CategoryName *string `json:"category_name"`
			ParentID     *uint   `json:"parent_id."`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": err.Error()})
			return
		}

		var Category models.Category

		if err := db.First(&Category, categoryId).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": "category not found"})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "error": "db error"})
			return
		}

		if input.CategoryName != nil {
			Category.CategoryName = *input.CategoryName
		}

		if input.ParentID != nil {
			Category.ParentID = input.ParentID
		}

		db.Save(&Category)

		c.JSON(http.StatusOK, gin.H{"status": "success", "data": Category})
	}
}
