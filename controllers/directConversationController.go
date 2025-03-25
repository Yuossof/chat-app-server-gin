package controllers

import (
	"net/http"

	"github.com/Yuossof/messaging-app-server/database"
	"github.com/Yuossof/messaging-app-server/models"
	"github.com/Yuossof/messaging-app-server/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// --------------CreateDirectConversation---------------//
func CreateDirectConversation(c *gin.Context) {
	tokenString, err := c.Cookie("token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token", "success": false})
		return
	}

	_, claims, err := utils.VerifyToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "missing token", "success": false})
		return
	}

	userIDStr := claims["user_id"].(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID format for userID"})
		return
	}

	user2IDStr := c.Query("anotherUser")
	user2ID, err := uuid.Parse(user2IDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID format for userID2"})
		return
	}

	var existingConv models.DirectConversation
	result := database.DB.Where(
		"(user1_id = ? AND user2_id = ?) OR (user1_id = ? AND user2_id = ?)",
		userID, user2ID, user2ID, userID,
	).First(&existingConv)

	if result.RowsAffected > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "conversation already exists", "success": false})
		return
	}

	conv := models.DirectConversation{
		User1ID: userID,
		User2ID: user2ID,
	}

	if err := database.DB.Create(&conv).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save conversation", "success": false})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "conversation created", "conversationID": conv.ID, "success": true})
}

// --------------GetSPdirectConversations---------------//
func GetSdirectConversations(c *gin.Context) {
	var directConversations []models.DirectConversation

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

	userID_TK := claims["user_id"].(string)
	result := database.DB.
		Preload("User1").
		Preload("User2").
		Where("user1_id = ? OR user2_id = ?", userID_TK, userID_TK).
		Find(&directConversations)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error(), "success": false})
		return
	}

	c.JSON(http.StatusOK, gin.H{"directConvs": directConversations, "success": true, "userID": userID_TK})
}

func GetSdirectCinversation(c *gin.Context) {
	var directConversation models.DirectConversation
	convID := c.Param("id")
	result := database.DB.Preload("Messages").Where("id = ?", convID).First(&directConversation)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error(), "success": false})
		return
	}

	c.JSON(http.StatusOK, gin.H{"conversation": directConversation, "success": true})
}

// --------------GetAlldirectConversations---------------//
func GetAlldirectConversations(c *gin.Context) {
	var directConversations []models.DirectConversation

	result := database.DB.Find(&directConversations)
	if result.Error != nil {
		c.JSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"directConvs": directConversations, "success": true})
}
