package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/junaid9001/spectr_backend/models"
	"github.com/junaid9001/spectr_backend/utils"
	"gorm.io/gorm"
)

// create a payment portal by orderId
func CreatePayment(db *gorm.DB) gin.HandlerFunc {
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

		var Order models.Order

		if err := db.Where("id=? AND user_id=?", orderId, userId).First(&Order).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"status": "failed", "error": err.Error()})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "error": "db error"})
			return
		}

		if Order.PaymentStatus != "pending" {
			c.JSON(400, gin.H{"error": "order not pending payment"})
			return
		}

		amount := Order.TotalAmount

		payment := models.Payment{
			OrderID:       orderId,
			Amount:        amount,
			PaymentStatus: "pending",
		}

		if err := db.Create(&payment).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"payment_id": payment.ID, "amount": amount})
	}
}

// confirm payment by payment id
func ConfirmPayment(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		paymentId, err := utils.StringToUint(c.Param("payment_id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": "invalid id format"})
			return
		}

		var payment models.Payment

		if err := db.First(&payment, paymentId).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"status": "failed", "error": err.Error()})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "error": "db error"})
			return
		}

		if payment.PaymentStatus != "pending" {
			c.JSON(200, gin.H{"error": "payment already completed"})
			return
		}

		orderId := payment.OrderID
		var order models.Order

		if err := db.Preload("OrderItems").First(&order, orderId).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"status": "failed", "error": err.Error()})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "error": "db error"})
			return
		}

		if order.TotalAmount != payment.Amount {
			c.JSON(400, gin.H{"error": "amount mismatch"})
			return
		}

		if err := db.Transaction(func(tx *gorm.DB) error {

			//reload order
			db.Preload("OrderItems").First(&order, orderId)

			if err := tx.Model(&payment).Update("payment_status", "completed").Error; err != nil {
				return err
			}

			if err := tx.Model(&order).Update("payment_status", "completed").Error; err != nil {
				return err
			}

			totalProducts := len(order.OrderItems)
			totalAmount := payment.Amount

			if err := tx.Model(&models.AppStats{}).Where("id=?", 1).Updates(map[string]interface{}{
				"total_sales":         gorm.Expr("total_sales + ?", 1),
				"total_products_sold": gorm.Expr("total_products_sold + ?", totalProducts),
				"total_revenue":       gorm.Expr("total_revenue + ?", totalAmount),
			}).Error; err != nil {
				return err
			}

			return nil
		}); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "error": "db error"})
			return
		}

		c.JSON(200, gin.H{"status": "paid", "order_id": order.ID})
	}
}
