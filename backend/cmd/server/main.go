package main

import (
	"log"
	"os"
	"time"

	"smart_attendance_backend/config"
	"smart_attendance_backend/models"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	if err := config.Load(); err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database
	if err := models.InitDB(config.GetDSN()); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Set Gin mode
	if os.Getenv("GIN_MODE") != "debug" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize router
	router := gin.Default()

	// Configure CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     config.AppConfig.Server.AllowOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Initialize routes
	initializeRoutes(router)

	// Start server
	port := config.AppConfig.Server.Port
	log.Printf("Server starting on port %s...\n", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func initializeRoutes(router *gin.Engine) {
	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	// API version group
	v1 := router.Group("/api/v1")
	{
		// Auth routes
		auth := v1.Group("/auth")
		{
			auth.POST("/register/teacher", nil) // TODO: Implement handler
			auth.POST("/register/student", nil) // TODO: Implement handler
			auth.POST("/verify-otp", nil)       // TODO: Implement handler
			auth.POST("/login/teacher", nil)    // TODO: Implement handler
			auth.POST("/login/student", nil)    // TODO: Implement handler
			auth.POST("/reset-password", nil)   // TODO: Implement handler
		}

		// Session routes
		sessions := v1.Group("/sessions")
		{
			sessions.POST("/start", nil) // TODO: Implement handler
			sessions.GET("/active", nil) // TODO: Implement handler
			sessions.PATCH("/end", nil)  // TODO: Implement handler
		}

		// Attendance routes
		attendance := v1.Group("/attendance")
		{
			attendance.POST("/mark", nil)  // TODO: Implement handler
			attendance.GET("/status", nil) // TODO: Implement handler
		}

		// Security routes
		security := v1.Group("/security")
		{
			security.GET("/check-developer-mode", nil) // TODO: Implement handler
			security.GET("/check-device-binding", nil) // TODO: Implement handler
			security.POST("/report-fraud", nil)        // TODO: Implement handler
		}
	}
}
