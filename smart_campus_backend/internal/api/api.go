package api

import (
	"smart_campus_backend/config"
	"smart_campus_backend/internal/middleware"
	"smart_campus_backend/internal/services"
	"time"

	"github.com/gin-gonic/gin"
)

func InitializeAPI(router *gin.Engine, services *services.Services, cfg *config.Config) {
	// Add middleware
	router.Use(gin.Recovery())
	router.Use(gin.Logger())
	router.Use(middleware.CORS(cfg.Server.AllowedOrigins))
	router.Use(middleware.RateLimit(10, time.Minute))

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Auth routes
		authRoutes := v1.Group("/auth")
		{
			authRoutes.POST("/register", services.AuthService.Register)
			authRoutes.POST("/login", services.AuthService.Login)
			authRoutes.POST("/verify-phone", services.AuthService.VerifyPhone)

			// Protected auth routes
			protected := authRoutes.Group("")
			protected.Use(middleware.JWT(cfg.JWT.Secret))
			{
				protected.POST("/refresh-token", services.AuthService.RefreshToken)
				protected.POST("/logout", services.AuthService.Logout)
			}
		}

		// Protected routes
		// protected := v1.Group("")
		// protected.Use(middleware.JWT(cfg.JWT.Secret))
		// {
		// 	// User routes
		// 	userRoutes := protected.Group("/users")
		// 	{
		// 		userRoutes.GET("/me", services.AuthService.GetCurrentUser)
		// 		userRoutes.PUT("/me", services.AuthService.UpdateProfile)
		// 	}

		// 	// Course routes
		// 	courseRoutes := protected.Group("/courses")
		// 	{
		// 		courseRoutes.POST("", services.CourseService.CreateCourse)
		// 		courseRoutes.GET("", services.CourseService.ListCourses)
		// 		courseRoutes.GET("/:id", services.CourseService.GetCourse)
		// 		courseRoutes.PUT("/:id", services.CourseService.UpdateCourse)
		// 	}

		// 	// Session routes
		// 	sessionRoutes := protected.Group("/sessions")
		// 	{
		// 		sessionRoutes.POST("", services.SessionService.CreateSession)
		// 		sessionRoutes.GET("", services.SessionService.ListSessions)
		// 		sessionRoutes.GET("/:id", services.SessionService.GetSession)
		// 		sessionRoutes.PUT("/:id", services.SessionService.UpdateSession)
		// 		sessionRoutes.POST("/:id/end", services.SessionService.EndSession)
		// 	}

		// 	// Attendance routes
		// 	attendanceRoutes := protected.Group("/attendance")
		// 	{
		// 		attendanceRoutes.POST("", services.AttendanceService.MarkAttendance)
		// 		attendanceRoutes.GET("/sessions/:id", services.AttendanceService.GetSessionAttendance)
		// 		attendanceRoutes.GET("/me", services.AttendanceService.GetMyAttendance)
		// 	}

		// 	// Device routes
		// 	deviceRoutes := protected.Group("/devices")
		// 	{
		// 		deviceRoutes.POST("/bind", services.DeviceService.BindDevice)
		// 		deviceRoutes.POST("/verify", services.DeviceService.VerifyDevice)
		// 		deviceRoutes.POST("/unbind", services.DeviceService.UnbindDevice)
		// 	}
		// }
	}
}
