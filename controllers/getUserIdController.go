package controllers

import (
	"net/http"

	"github.com/Yuossof/messaging-app-server/utils"
	"github.com/gin-gonic/gin"
)

func GetUserID(c *gin.Context) {
	tokenString, err := c.Cookie("token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token", "success": false})
		return
	}
	_, claims, err := utils.VerifyToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing token", "success": false})
		return
	}
	userID := claims["user_id"].(string)
	c.JSON(http.StatusOK, gin.H{"user_id": userID})
}
