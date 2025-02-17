package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	firebase "firebase.google.com/go/v4"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"

	"smart_campus_backend/config"
	"smart_campus_backend/internal/api"
	"smart_campus_backend/internal/db"
	"smart_campus_backend/internal/services"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Error loading config: %v\n", err)
	}

	// Set Gin mode
	gin.SetMode(cfg.Server.Mode)

	// Initialize Firebase Admin SDK
	firebaseApp, err := initializeFirebase(cfg.Firebase.CredentialsPath)
	if err != nil {
		log.Fatalf("Error initializing Firebase: %v\n", err)
	}

	// Initialize database
	db, err := db.NewDatabase(cfg.Database)
	if err != nil {
		log.Fatalf("Error connecting to database: %v\n", err)
	}
	defer db.Close()

	// Initialize repositories
	repos := initializeRepositories(db)

	// Initialize services
	services := initializeServices(repos, firebaseApp)

	// Initialize router and API
	router := gin.Default()
	api.InitializeAPI(router, services, cfg)

	// Start server
	srv := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on port %s\n", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting server: %v\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}

func initializeFirebase(credentialsPath string) (*firebase.App, error) {
	opt := option.WithCredentialsFile(credentialsPath)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, err
	}
	return app, nil
}

func initializeRepositories(db *db.Database) *services.Repositories {
	return &services.Repositories{
		UserRepo: db.NewUserRepository(),
		/*SessionRepo:    db.NewSessionRepository(),
		CourseRepo:     db.NewCourseRepository(),
		DeviceRepo:     db.NewDeviceRepository(),
		AttendanceRepo: db.NewAttendanceRepository(),*/
	}
}

func initializeServices(repos *services.Repositories, firebaseApp *firebase.App) *services.Services {
	auth, err := firebaseApp.Auth(context.Background())
	if err != nil {
		log.Fatalf("Error getting Firebase Auth client: %v\n", err)
	}

	return &services.Services{
		AuthService: services.NewAuthService(repos.UserRepo, auth),
		/*SessionService:    services.NewSessionService(repos.SessionRepo, repos.AttendanceRepo),
		CourseService:     services.NewCourseService(repos.CourseRepo),
		DeviceService:     services.NewDeviceService(repos.DeviceRepo),
		AttendanceService: services.NewAttendanceService(repos.AttendanceRepo),*/
	}
}
