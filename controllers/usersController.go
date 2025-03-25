package controllers

import (
	"net/http"

	"github.com/Yuossof/messaging-app-server/database"
	"github.com/Yuossof/messaging-app-server/models"
	"github.com/gin-gonic/gin"
)

func GetUsers(c *gin.Context) {
	var users []models.User

	if err := database.DB.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch users"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}

func SearchUser(c *gin.Context) {
	var users []models.User
	searchKey := c.Query("searchKey")

	if err := database.DB.Where("username LIKE ?", "%"+searchKey+"%").Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	if len(users) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "no users found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}
