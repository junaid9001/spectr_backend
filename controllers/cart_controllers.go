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

// add products to cart (user)
func AddProductToCart(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		//get user id also safe accertion
		userId, ok := utils.GetUserId(c)
		if !ok {
			return
		}
		var input struct {
			ProductId uint `json:"product_id" binding:"required,gt=0"`
			Quantity  int  `json:"quantity" binding:"required,gt=0"`
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

		if product.StockQuantity < input.Quantity {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "failed",
				"error":  "not enough stock is available",
			})
			return
		}

		//check if product already exists /if then increase the quantity
		var item models.CartItem
		err := db.Where("product_id=? AND user_id=?", input.ProductId, userId).First(&item).Error

		if err == nil {
			item.Quantity = item.Quantity + input.Quantity
			item.UnitPrice = product.Price
			item.TotalPrice = float64(item.Quantity) * product.Price

			if err := db.Save(&item).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "error": "could not update cart item"})
				return
			}

			c.JSON(http.StatusCreated, gin.H{"status": "success"})
			return

		}

		//if first time

		if errors.Is(err, gorm.ErrRecordNotFound) {

			userCartItem := models.CartItem{
				UserId:     userId,
				ProductId:  product.ID,
				Quantity:   input.Quantity,
				UnitPrice:  product.Price,
				TotalPrice: float64(input.Quantity) * product.Price,
				Product:    product,
			}

			result := db.Create(&userCartItem)
			if result.Error != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "error": result.Error.Error()})
				return
			}

			c.JSON(http.StatusCreated, gin.H{"status": "success"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "error": err.Error()})

	}
}

// user's cart

func GetUserCart(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		//get user if also safe accertion
		userId, ok := utils.GetUserId(c)
		if !ok {
			return
		}
		var products []models.CartItem

		if err := db.Preload("Product").Where("user_id=?", userId).Find(&products).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "error": err.Error()})
			return
		}

		if len(products) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"status": "failed", "error": "your cart is empty"})
			return
		}

		var subTotal float64

		for _, prod := range products {
			subTotal += prod.TotalPrice
		}
		c.JSON(http.StatusOK, gin.H{
			"status":   "success",
			"count":    len(products),
			"subtotal": subTotal,
			"data":     products,
		})

	}
}

//update quantity of product in cart by cart id

func UpdateQuantityInCartByID(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		idstr := c.Param("id")

		cartItemId, err := utils.StringToUint(idstr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "failed",
				"error":  "invalid id format",
			})
			return
		}

		userId, ok := utils.GetUserId(c)
		if !ok {
			return //responds is auto
		}

		var input struct {
			Delta int `json:"delta" binding:"required,ne=0"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": err.Error()})
			return
		}

		usercartItem := models.CartItem{}

		if err := db.Where("id=? AND user_id=?", cartItemId, userId).First(&usercartItem).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"status": "failed", "error": err.Error()})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "error": err.Error()})
			return
		}

		productId := usercartItem.ProductId

		var product models.Product

		if err := db.Where("id=?", productId).First(&product).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"status": "failed", "error": err.Error()})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "error": err.Error()})
			return
		}

		newQuantity := usercartItem.Quantity + input.Delta

		if newQuantity <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": "quantity should be > 1"})
			return
		}

		if product.StockQuantity < newQuantity {
			c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": "not enough stock is available "})
			return
		}

		unitPrice := product.Price
		totalPrice := product.Price * float64(newQuantity)

		updates := map[string]interface{}{
			"quantity":    newQuantity,
			"unit_price":  unitPrice,
			"total_price": totalPrice,
		}

		if err := db.Model(&usercartItem).Updates(updates).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "success", "data": usercartItem})

	}
}

//delete cartitem by id

func DeleteCartItemByID(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idstr := c.Param("id")

		cartItemId, err := utils.StringToUint(idstr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "failed",
				"error":  "invalid id format",
			})
			return
		}

		userId, ok := utils.GetUserId(c)
		if !ok {
			return
		}

		result := db.Where("user_id=? AND id=?", userId, cartItemId).Delete(&models.CartItem{})

		if result.Error != nil {

			c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "error": result.Error.Error()})
			return
		}

		if result.RowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"status": "failed", "error": "cart item not found"})
			return
		}

		c.Status(http.StatusNoContent)

	}
}
