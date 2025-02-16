package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}
}

func main() {
	// Set Gin mode
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize router
	router := gin.Default()

	// Add middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Initialize routes
	initializeRoutes(router)

	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func initializeRoutes(router *gin.Engine) {
	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Auth routes
		auth := v1.Group("/auth")
		{
			auth.POST("/login", nil)    // TODO: Implement login handler
			auth.POST("/register", nil) // TODO: Implement register handler
			auth.POST("/refresh", nil)  // TODO: Implement token refresh handler
		}

		// Protected routes
		protected := v1.Group("/")
		{
			// TODO: Add authentication middleware

			// User routes
			users := protected.Group("/users")
			{
				users.GET("/", nil)       // TODO: Implement get users handler
				users.GET("/:id", nil)    // TODO: Implement get user handler
				users.PUT("/:id", nil)    // TODO: Implement update user handler
				users.DELETE("/:id", nil) // TODO: Implement delete user handler
			}

			// Course routes
			courses := protected.Group("/courses")
			{
				courses.POST("/", nil)      // TODO: Implement create course handler
				courses.GET("/", nil)       // TODO: Implement get courses handler
				courses.GET("/:id", nil)    // TODO: Implement get course handler
				courses.PUT("/:id", nil)    // TODO: Implement update course handler
				courses.DELETE("/:id", nil) // TODO: Implement delete course handler
			}

			// Session routes
			sessions := protected.Group("/sessions")
			{
				sessions.POST("/", nil)      // TODO: Implement create session handler
				sessions.GET("/", nil)       // TODO: Implement get sessions handler
				sessions.GET("/:id", nil)    // TODO: Implement get session handler
				sessions.PUT("/:id", nil)    // TODO: Implement update session handler
				sessions.DELETE("/:id", nil) // TODO: Implement delete session handler

				// Attendance routes
				sessions.POST("/:id/attendance", nil) // TODO: Implement mark attendance handler
				sessions.GET("/:id/attendance", nil)  // TODO: Implement get attendance handler
			}
		}
	}
}
