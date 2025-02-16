package routes

import (
	"smart_campus/internal/api/handlers"
	"smart_campus/internal/api/middleware"
	"smart_campus/internal/services"

	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(
	router *gin.RouterGroup,
	authService *services.AuthService,
	deviceService *services.DeviceService,
) {
	// Create handlers
	authHandler := handlers.NewAuthHandler(authService, deviceService)

	// Auth routes
	auth := router.Group("/auth")
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/register", authHandler.Register)
		auth.POST("/verify-otp", authHandler.VerifyOTP)
		auth.POST("/forgot-password", authHandler.ForgotPassword)
		auth.POST("/reset-password", authHandler.ResetPassword)
		auth.POST("/refresh-token", authHandler.RefreshToken)

		// Protected device routes
		device := auth.Group("/device")
		device.Use(middleware.JWTAuth(authService))
		{
			device.POST("/bind", authHandler.BindDevice)
			device.POST("/verify", authHandler.VerifyDevice)
			device.POST("/unbind", authHandler.UnbindDevice)
		}
	}
}
