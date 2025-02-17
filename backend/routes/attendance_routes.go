package routes

import (
	"github.com/gin-gonic/gin"
	"smart-attendance/controllers"
	"smart-attendance/middleware"
)

// AttendanceRoutes defines the attendance management routes.
func AttendanceRoutes(router *gin.Engine) {
	attendanceGroup := router.Group("/attendance")
	attendanceGroup.Use(middleware.AuthMiddleware()) // Apply authentication middleware

	{
		attendanceGroup.POST("/mark", controllers.MarkAttendance)
		attendanceGroup.GET("/status", controllers.GetAttendanceStatus)
	}
}
