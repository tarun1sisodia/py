package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"smart_campus/internal/api"
	"smart_campus/internal/config"
	"smart_campus/internal/database"
	"smart_campus/internal/repositories"
	"smart_campus/internal/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Set Gin mode
	gin.SetMode(cfg.Server.Mode)

	// Initialize database
	db, err := database.NewMySQLDB(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Initialize repositories
	userRepo := repositories.NewMySQLUserRepository(db)
	courseRepo := repositories.NewMySQLCourseRepository(db)
	sessionRepo := repositories.NewMySQLSessionRepository(db)
	deviceRepo := repositories.NewMySQLDeviceRepository(db)
	attendanceRepo := repositories.NewMySQLAttendanceRepository(db)

	// Initialize services
	authService := services.NewAuthService(cfg, userRepo, deviceRepo)
	deviceService := services.NewDeviceService(deviceRepo)
	sessionService := services.NewSessionService(sessionRepo, attendanceRepo)

	// Initialize server
	server := api.NewServer(cfg, authService, deviceService, sessionService)

	// Server run context
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	// Listen for syscall signals for process to interrupt/quit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		// Shutdown signal with grace period of 30 seconds
		shutdownCtx, _ := context.WithTimeout(serverCtx, 30*time.Second)

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		// Trigger graceful shutdown
		err := server.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal(err)
		}
		serverStopCtx()
	}()

	// Run the server
	fmt.Printf("Server is running on port %s in %s mode\n", cfg.Server.Port, cfg.Server.Mode)
	err = server.Start()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

	// Wait for server context to be stopped
	<-serverCtx.Done()
}
