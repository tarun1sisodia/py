package routes

import (
	"github.com/gin-gonic/gin"
)

// SetupRoutes defines all API routes.
func SetupRoutes(router *gin.Engine) {
	AuthRoutes(router)
	SessionRoutes(router)
	AttendanceRoutes(router)
	SyncRoutes(router) // Add sync routes
	// Define other routes here
}
