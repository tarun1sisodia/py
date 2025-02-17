package routes

import (
	"github.com/gin-gonic/gin"
	"smart-attendance/controllers"
	"smart-attendance/middleware"
)

// SessionRoutes defines the session management routes.
func SessionRoutes(router *gin.Engine) {
	sessionGroup := router.Group("/sessions")
	sessionGroup.Use(middleware.AuthMiddleware()) // Apply authentication middleware

	{
		sessionGroup.POST("/start", controllers.StartSession)
		sessionGroup.GET("/active", controllers.GetActiveSession)
		sessionGroup.PATCH("/end/:id", controllers.EndSession)
	}
}
