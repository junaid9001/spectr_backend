package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/junaid9001/spectr_backend/models"
	"github.com/junaid9001/spectr_backend/utils"
	"gorm.io/gorm"
)

//done

//place order (user)

func PlaceOrder(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			ShippingAddress string `json:"shipping_address" binding:"required"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": err.Error()})
			return
		}

		userId, ok := utils.GetUserId(c)
		if !ok {
			return //response are already send
		}

		var UserCartItems []models.CartItem

		if err := db.Where("user_id=?", userId).Preload("Product").Find(&UserCartItems).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "error": "failed to read cart"})
			return
		}

		if len(UserCartItems) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": "cart is empty"})
			return
		}

		var total float64 = 0

		orderItems := make([]models.OrderItem, 0, len(UserCartItems))

		//orderItems
		for _, val := range UserCartItems {
			unitPrice := val.UnitPrice
			tl := unitPrice * float64(val.Quantity)

			total += tl

			orderItems = append(orderItems, models.OrderItem{
				ProductID:  val.ProductId,
				UnitPrice:  unitPrice,
				Quantity:   val.Quantity,
				TotalPrice: tl,
			})
		}

		var createdOrder models.Order

		//start transaction

		err := db.Transaction(func(tx *gorm.DB) error {
			// temporarily enable debug logging (prints SQL to stdout)
			tx = tx.Debug()

			order := models.Order{
				UserID:      userId,
				TotalAmount: total,
				Address:     input.ShippingAddress,
				Status:      "pending",
				CreatedAt:   time.Now(),
			}

			if err := tx.Create(&order).Error; err != nil {
				return err
			}

			//orderid for each order item
			for i := range orderItems {
				orderItems[i].OrderID = order.ID
			}

			if err := tx.Create(&orderItems).Error; err != nil {
				return err
			}

			//decrement stock of all product
			for _, val := range orderItems {
				res := tx.Model(&models.Product{}).Where("id=? AND stock_quantity >=?", val.ProductID, val.Quantity).
					UpdateColumn("stock_quantity", gorm.Expr("stock_quantity - ?", val.Quantity))
				if res.Error != nil {
					return res.Error
				}
				if res.RowsAffected == 0 {
					return fmt.Errorf("product %d does not have enough stock", val.ProductID)
				}
			}

			if err := tx.Where("user_id=?", userId).Delete(&models.CartItem{}).Error; err != nil {
				return err
			}
			createdOrder = order
			return nil

		})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": err.Error()})
			return

		}

		//get order items for return

		if err := db.Preload("OrderItems").First(&createdOrder, createdOrder.ID).Error; err != nil {
			c.JSON(http.StatusCreated, gin.H{"status": "success", "data": createdOrder.ID})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"status": "success", "data": createdOrder})

	}
}

//get users order history (user)

func GetOrderHistory(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, ok := utils.GetUserId(c)
		if !ok {
			return
		}

		var userOrders []models.Order

		result := db.Preload("OrderItems").Where("user_id=?", userId).Find(&userOrders)

		if result.Error != nil {

			c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "error": "db error"})
			return
		}

		if len(userOrders) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"status": "failed", "error": "no orders found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "success", "data": userOrders})
	}
}

//get details of one specific order by order id (user)

func GetDetailsOfOrder(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idstr := c.Param("id")
		orderId, err := utils.StringToUint(idstr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": "invalid id format"})
			return
		}

		userId, ok := utils.GetUserId(c)
		if !ok {
			return
		}
		var order models.Order

		if err := db.Where("id=? AND user_id=?", orderId, userId).First(&order).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"status": "failed", "error": "order not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "error": "db error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "success", "data": order})
	}
}

//delete specific order by id (user)

func DeleteOrderById(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idstr := c.Param("id")

		orderId, err := utils.StringToUint(idstr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": "invalid id format"})
			return
		}

		userId, ok := utils.GetUserId(c)
		if !ok {
			return
		}

		var order models.Order

		if err := db.Where("id=? AND user_id=?", orderId, userId).First(&order).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"status": "failed", "error": "order not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "error": "db error"})
			return
		}

		//only if cancelled or delivered
		if order.Status == "pending" {
			c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": "cannot delete a pending order"})
			return
		}

		result := db.Delete(&order)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "error": "db error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "success"})
	}
}

//cancel a order by id  (user) and restock

func CancelOrderAndRestock(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idstr := c.Param("id")

		orderId, err := utils.StringToUint(idstr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": "invalid id format"})
			return
		}

		userId, ok := utils.GetUserId(c)
		if !ok {
			return
		}

		var order models.Order

		if err := db.Preload("OrderItems").Where("id=? AND user_id=?", orderId, userId).First(&order).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"status": "failed", "error": "order not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "error": "db error"})
			return
		}

		if order.Status != "pending" {
			c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": "only pending orders can be cancelled"})
			return
		}

		if err := db.Transaction(func(tx *gorm.DB) error {

			if err := tx.Model(&order).Update("status", "cancelled").Error; err != nil {
				return err
			}

			//restock each product
			for _, val := range order.OrderItems {
				res := tx.Model(&models.Product{}).Where("id=?", val.ProductID).
					UpdateColumn("stock_quantity", gorm.Expr("stock_quantity + ?", val.Quantity))

				if res.Error != nil {
					return res.Error
				}
			}

			return nil
		}); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "success"})
	}
}

//---------------**---------------------
//---------------**-------------------------
//---------------**-------------------------

//get all orders (admin)

func GetAllOrders(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var AllOrders []models.Order

		result := db.Preload("OrderItems").Find(&AllOrders)

		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "error": "db error"})
			return
		}

		if len(AllOrders) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"status": "failed", "error": "no orders found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "success", "data": AllOrders})
	}
}

//update order status by id (admin)

func UpdateOrderStatus(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		idstr := c.Param("id")

		orderId, err := utils.StringToUint(idstr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": "invalid id format"})
			return
		}

		var input struct {
			Status string `json:"status" binding:"required"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": err.Error()})
			return
		}

		statusDelivered := "delivered"
		statusCancelled := "cancelled"

		var order models.Order

		if err := db.Preload("OrderItems").First(&order, orderId).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"status": "failed", "error": "order not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "error": "db error"})
			return
		}
		//validate request

		if input.Status != statusDelivered && input.Status != statusCancelled {
			c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": "invalid status"})
			return
		}

		//already same?
		if order.Status == input.Status {
			c.JSON(http.StatusOK, gin.H{"status": "success", "data": order})
			return
		}

		if input.Status == "cancelled" && order.Status == "pending" {
			if err := db.Transaction(func(tx *gorm.DB) error {

				if err := tx.Model(&order).Update("status", "cancelled").Error; err != nil {
					return err
				}

				//restock each product
				for _, val := range order.OrderItems {
					res := tx.Model(&models.Product{}).Where("id=?", val.ProductID).
						UpdateColumn("stock_quantity", gorm.Expr("stock_quantity + ?", val.Quantity))

					if res.Error != nil {
						return res.Error
					}
				}

				return nil

			}); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "error": err.Error()})
				return
			}
			db.First(&order, orderId)
			c.JSON(http.StatusOK, gin.H{"status": "success", "data": order})
			return
		}

		result := db.Model(&order).Update("status", input.Status)

		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "error": result.Error.Error()})
			return
		}

		if result.RowsAffected == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": "order status was not updated"})
			return
		}

		db.First(&order, orderId)

		c.JSON(http.StatusOK, gin.H{"status": "success", "data": order})
	}
}
