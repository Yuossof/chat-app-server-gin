package routes

import (
	"github.com/Yuossof/messaging-app-server/controllers"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	api := r.Group("/api")

	// get userID
	api.GET("/getUserId", controllers.GetUserID)

	// auth
	{
		api.POST("/register", controllers.Register)
		api.POST("/login", controllers.Login)
		api.GET("/verify", controllers.VerifyToken)
	}

	// conversation
	{
		api.POST("/conv/create-direct", controllers.CreateDirectConversation)
		api.GET("/conv/sp-direct", controllers.GetSdirectConversations)
		api.GET("/conv/sp-direct/:id", controllers.GetSdirectCinversation)
		api.GET("/conv/all-direct", controllers.GetAlldirectConversations)
	}

	// users
	{
		api.GET("/users/all", controllers.GetUsers)
		api.GET("/users/search-users", controllers.SearchUser)
	}

	r.GET("/ws", controllers.WebSocketHandler)
}
