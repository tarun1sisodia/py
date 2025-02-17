package routes

import (
	"github.com/gin-gonic/gin"
	"smart-attendance/controllers"
	"smart-attendance/middleware"
)

// SyncRoutes defines the synchronization routes.
func SyncRoutes(router *gin.Engine) {
	syncGroup := router.Group("/sync")
	syncGroup.Use(middleware.AuthMiddleware()) // Apply authentication middleware

	{
		syncGroup.POST("/attendance", controllers.SyncAttendanceData)
	}
}
