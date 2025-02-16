package routes

import (
	"smart_campus/internal/api/handlers"
	"smart_campus/internal/api/middleware"
	"smart_campus/internal/services"

	"github.com/gin-gonic/gin"
)

// RegisterDeviceRoutes registers the device routes
func RegisterDeviceRoutes(
	router *gin.RouterGroup,
	deviceService *services.DeviceService,
	authService *services.AuthService,
) {
	// Create handlers
	deviceHandler := handlers.NewDeviceHandler(deviceService)

	// Device routes
	device := router.Group("/devices")
	device.Use(middleware.JWTAuth(authService))
	{
		device.POST("/bind", deviceHandler.BindDevice)
		device.POST("/verify", deviceHandler.VerifyDevice)
		device.POST("/unbind", deviceHandler.UnbindDevice)
		device.GET("/list", deviceHandler.ListDevices)
		device.GET("/:id", deviceHandler.GetDevice)
	}
}
