package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetUserId(c *gin.Context) (uint, bool) {

	uID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": "failed",
			"error":  "unauthenticated",
		})
		return 0, false
	}

	userId, ok := uID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "failed",
			"error":  "invalid user id",
		})
		return 0, false
	}
	return userId, true

}
