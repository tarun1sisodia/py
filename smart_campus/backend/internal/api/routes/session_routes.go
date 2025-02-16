package routes

import (
	"smart_campus/internal/api/handlers"
	"smart_campus/internal/api/middleware"
	"smart_campus/internal/services"

	"github.com/gin-gonic/gin"
)

func RegisterSessionRoutes(
	router *gin.RouterGroup,
	sessionService *services.SessionService,
	authService *services.AuthService,
) {
	// Create handlers
	sessionHandler := handlers.NewSessionHandler(sessionService)

	// Session routes
	sessions := router.Group("/sessions")
	sessions.Use(middleware.JWTAuth(authService))
	{
		// Teacher routes
		teacher := sessions.Group("")
		teacher.Use(middleware.RequireRole("teacher"))
		{
			teacher.POST("", sessionHandler.CreateSession)
			teacher.PUT("/:id/end", sessionHandler.EndSession)
			teacher.PUT("/:id/cancel", sessionHandler.CancelSession)
			teacher.GET("/:id/attendance", sessionHandler.GetSessionAttendance)
			teacher.PUT("/attendance/:id/verify", sessionHandler.VerifyAttendance)
			teacher.PUT("/attendance/:id/reject", sessionHandler.RejectAttendance)
		}

		// Student routes
		student := sessions.Group("")
		student.Use(middleware.RequireRole("student"))
		{
			student.GET("/active", sessionHandler.GetActiveSessions)
			student.POST("/:id/attendance", sessionHandler.MarkAttendance)
			student.GET("/attendance/history", sessionHandler.GetAttendanceHistory)
			student.GET("/attendance/statistics", sessionHandler.GetAttendanceStatistics)
		}

		// Common routes
		sessions.GET("", sessionHandler.GetSessions)
		sessions.GET("/:id", sessionHandler.GetSessionById)
	}
}
