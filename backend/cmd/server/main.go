package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"smart-attendance/config"
	"smart-attendance/middleware"
	"smart-attendance/routes"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Connect to the database
	config.ConnectDB()
	defer config.CloseDB()

	// Initialize Gin router
	router := gin.Default()

	// Add error handling middleware
	router.Use(middleware.ErrorHandler())

	// Add rate limiting middleware
	rateLimit := os.Getenv("RATE_LIMIT")
	rateBurst := os.Getenv("RATE_BURST")

	limit, err := time.ParseDuration(rateLimit)

	if err != nil {
		log.Println("RATE_LIMIT environment variable not set, using default value")
		limit = time.Minute
	}

	burst, err := strconv.Atoi(rateBurst)

	if err != nil {
		log.Println("RATE_BURST environment variable not set, using default value")
		burst = 100
	}

	router.Use(middleware.RateLimiter(rate.Every(limit), burst))

	// Setup routes
	routes.SetupRoutes(router)

	// Define a simple health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// Start the server
	serverAddr := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("Starting server on %s", serverAddr)
	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
