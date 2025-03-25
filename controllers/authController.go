package controllers

import (
	"fmt"
	"net/http"

	"github.com/Yuossof/messaging-app-server/models"

	"github.com/gin-gonic/gin"
	// "gorm.io/gorm"
	"github.com/Yuossof/messaging-app-server/database"
	"github.com/Yuossof/messaging-app-server/utils"
)

func Register(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "success": false})
		return
	}

	if err := user.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "success": false})
		return
	}

	var existingUser models.User
	if err := database.DB.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "user already exists", "success": false})
		return
	}

	hashedPassword, err := utils.HashPassword(user.Password)
	user.Password = hashedPassword
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "hashing password error", "success": false})
		return
	}

	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "register error", "success": false})
		return
	}

	token, err := utils.GenerateToken(user.ID.String())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "auth error", "success": false})
		return
	}

	c.SetCookie("token", token, 60*60*24, "/", "localhost", false, true)

	c.JSON(http.StatusCreated, gin.H{
		"message": "account created successfully",
		"token":   token,
		"id":      user.ID,
	})

}

//---------------------------------------------------------------//

func Login(c *gin.Context) {
	var user models.User
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "success": false})
		return
	}

	if err := database.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		fmt.Println("Register Error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "user not found", "success": false})
		return
	}

	isCheckPassword := utils.CheckPassword(user.Password, input.Password)
	if !isCheckPassword {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email or password", "success": false})
		return
	}

	token, err := utils.GenerateToken(user.ID.String())

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "auth error", "success": false})
		return
	}

	c.SetCookie("token", token, 60*60*24, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "login success", "token": token, "id": user.ID})
}

func VerifyToken(c *gin.Context) {
	tokenString, err := c.Cookie("token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized", "success": false})
		return
	}

	_, claims, err := utils.VerifyToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token", "success": false})
		return
	}

	userID := claims["user_id"].(string)

	c.JSON(http.StatusOK, gin.H{"message": "Token Verified", "success": true, "user_id": userID})
}
