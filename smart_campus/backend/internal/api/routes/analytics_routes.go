package routes

import (
	"smart_campus/internal/api/handlers"
	"smart_campus/internal/domain/repositories"
	"smart_campus/internal/services"

	"github.com/gin-gonic/gin"
)

// RegisterAnalyticsRoutes registers the analytics routes
func RegisterAnalyticsRoutes(
	router *gin.RouterGroup,
	attendanceRepo repositories.AttendanceRecordRepository,
	courseRepo repositories.CourseRepository,
	studentRepo repositories.StudentRepository,
	teacherRepo repositories.TeacherRepository,
	authService *services.AuthService,
) {
	analyticsHandler := handlers.NewAnalyticsHandler(
		attendanceRepo,
		courseRepo,
		studentRepo,
		teacherRepo,
		authService,
	)
	analyticsHandler.RegisterRoutes(router)
}
