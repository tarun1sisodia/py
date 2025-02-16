package api

import (
	"context"
	"fmt"
	"net/http"
	"smart_campus/internal/api/middleware"
	"smart_campus/internal/api/routes"
	"smart_campus/internal/config"
	"smart_campus/internal/domain/repositories"
	"smart_campus/internal/services"

	"github.com/gin-gonic/gin"
)

type Server struct {
	router         *gin.Engine
	httpServer     *http.Server
	config         *config.Config
	authService    *services.AuthService
	deviceService  *services.DeviceService
	sessionService *services.SessionService
	attendanceRepo repositories.AttendanceRecordRepository
	courseRepo     repositories.CourseRepository
	studentRepo    repositories.StudentRepository
	teacherRepo    repositories.TeacherRepository
}

func NewServer(
	config *config.Config,
	authService *services.AuthService,
	deviceService *services.DeviceService,
	sessionService *services.SessionService,
	attendanceRepo repositories.AttendanceRecordRepository,
	courseRepo repositories.CourseRepository,
	studentRepo repositories.StudentRepository,
	teacherRepo repositories.TeacherRepository,
) *Server {
	server := &Server{
		router:         gin.Default(),
		config:         config,
		authService:    authService,
		deviceService:  deviceService,
		sessionService: sessionService,
		attendanceRepo: attendanceRepo,
		courseRepo:     courseRepo,
		studentRepo:    studentRepo,
		teacherRepo:    teacherRepo,
	}

	// Setup CORS
	server.router.Use(middleware.CORS())

	// Setup recovery middleware
	server.router.Use(middleware.Recovery())

	// Setup request logging
	server.router.Use(middleware.RequestLogger())

	// Initialize routes
	v1 := server.router.Group("/api/v1")
	routes.RegisterAuthRoutes(v1, authService, deviceService)
	routes.RegisterSessionRoutes(v1, sessionService, authService)
	routes.RegisterDeviceRoutes(v1, deviceService, authService)
	routes.RegisterAnalyticsRoutes(v1, attendanceRepo, courseRepo, studentRepo, teacherRepo, authService)

	// Setup HTTP server
	server.httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%s", config.Server.Port),
		Handler: server.router,
	}

	return server
}

func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
