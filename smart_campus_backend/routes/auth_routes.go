package routes

import (
	"smart_campus_backend/handlers"
	"smart_campus_backend/services"

	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(router *gin.RouterGroup, authService *services.AuthService) {
	handler := handlers.NewAuthHandler(authService)

	auth := router.Group("/auth")
	{
		auth.POST("/register", handler.Register)
		auth.POST("/login", handler.Login)
		auth.POST("/verify-phone", handler.VerifyPhone)
	}
}
