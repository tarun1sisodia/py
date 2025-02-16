package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"smart_campus/internal/api/middleware"
	"smart_campus/internal/api/response"
	"smart_campus/internal/domain/entities"
	"smart_campus/internal/domain/repositories"
	"smart_campus/internal/services"
)

// AnalyticsHandler handles analytics-related requests
type AnalyticsHandler struct {
	attendanceRepo repositories.AttendanceRecordRepository
	courseRepo     repositories.CourseRepository
	studentRepo    repositories.StudentRepository
	teacherRepo    repositories.TeacherRepository
	authService    *services.AuthService
}

// NewAnalyticsHandler creates a new analytics handler
func NewAnalyticsHandler(
	attendanceRepo repositories.AttendanceRecordRepository,
	courseRepo repositories.CourseRepository,
	studentRepo repositories.StudentRepository,
	teacherRepo repositories.TeacherRepository,
	authService *services.AuthService,
) *AnalyticsHandler {
	return &AnalyticsHandler{
		attendanceRepo: attendanceRepo,
		courseRepo:     courseRepo,
		studentRepo:    studentRepo,
		teacherRepo:    teacherRepo,
		authService:    authService,
	}
}

// RegisterRoutes registers the analytics routes
func (h *AnalyticsHandler) RegisterRoutes(router *gin.RouterGroup) {
	analytics := router.Group("/analytics")
	analytics.Use(middleware.JWTAuth(h.authService))
	{
		analytics.GET("/attendance", h.GetAttendanceAnalytics)
		analytics.GET("/courses", h.GetCourseAnalytics)
		analytics.GET("/students", h.GetStudentAnalytics)
	}

	reports := router.Group("/reports")
	reports.Use(middleware.JWTAuth(h.authService))
	{
		reports.GET("/attendance", h.GetAttendanceReport)
		reports.POST("/export", h.ExportReport)
	}
}

// GetAttendanceAnalytics returns attendance analytics
func (h *AnalyticsHandler) GetAttendanceAnalytics(c *gin.Context) {
	// Get query parameters
	courseID := c.Query("course_id")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	departmentID := c.Query("department_id")

	// Parse dates
	var start, end time.Time
	var err error
	if startDate != "" {
		start, err = time.Parse("2006-01-02", startDate)
		if err != nil {
			response.Error(c, http.StatusBadRequest, "Invalid start date format")
			return
		}
	}
	if endDate != "" {
		end, err = time.Parse("2006-01-02", endDate)
		if err != nil {
			response.Error(c, http.StatusBadRequest, "Invalid end date format")
			return
		}
	}

	// Build filters
	filters := make(map[string]interface{})
	if courseID != "" {
		filters["course_id"] = courseID
	}
	if departmentID != "" {
		filters["department_id"] = departmentID
	}
	if !start.IsZero() {
		filters["start_date"] = start
	}
	if !end.IsZero() {
		filters["end_date"] = end
	}

	// Get attendance records
	records, err := h.attendanceRepo.List(c.Request.Context(), filters, 0, 0)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Error getting attendance records")
		return
	}

	// Calculate analytics
	analytics := calculateAttendanceAnalytics(records)
	response.Success(c, http.StatusOK, "Attendance analytics retrieved successfully", analytics)
}

// GetCourseAnalytics returns course analytics
func (h *AnalyticsHandler) GetCourseAnalytics(c *gin.Context) {
	departmentID := c.Query("department_id")
	academicYear := c.Query("academic_year")

	filters := make(map[string]interface{})
	if departmentID != "" {
		filters["department_id"] = departmentID
	}
	if academicYear != "" {
		filters["academic_year"] = academicYear
	}

	courses, err := h.courseRepo.List(c.Request.Context(), filters, 0, 0)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Error getting course analytics")
		return
	}

	analytics := calculateCourseAnalytics(courses)
	response.Success(c, http.StatusOK, "Course analytics retrieved successfully", analytics)
}

// GetStudentAnalytics returns student analytics
func (h *AnalyticsHandler) GetStudentAnalytics(c *gin.Context) {
	departmentID := c.Query("department_id")
	yearOfStudy := c.Query("year_of_study")

	filters := make(map[string]interface{})
	if departmentID != "" {
		filters["department_id"] = departmentID
	}
	if yearOfStudy != "" {
		year, err := strconv.Atoi(yearOfStudy)
		if err != nil {
			response.Error(c, http.StatusBadRequest, "Invalid year of study")
			return
		}
		filters["year_of_study"] = year
	}

	students, err := h.studentRepo.List(c.Request.Context(), filters, 0, 0)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Error getting student analytics")
		return
	}

	analytics := calculateStudentAnalytics(students)
	response.Success(c, http.StatusOK, "Student analytics retrieved successfully", analytics)
}

// GetAttendanceReport returns a detailed attendance report
func (h *AnalyticsHandler) GetAttendanceReport(c *gin.Context) {
	courseID := c.Query("course_id")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	studentID := c.Query("student_id")

	// Parse dates
	var start, end time.Time
	var err error
	if startDate != "" {
		start, err = time.Parse("2006-01-02", startDate)
		if err != nil {
			response.Error(c, http.StatusBadRequest, "Invalid start date format")
			return
		}
	}
	if endDate != "" {
		end, err = time.Parse("2006-01-02", endDate)
		if err != nil {
			response.Error(c, http.StatusBadRequest, "Invalid end date format")
			return
		}
	}

	// Get attendance records
	var records []*entities.AttendanceRecord
	if studentID != "" {
		records, err = h.attendanceRepo.GetByStudentAndDateRange(c.Request.Context(), studentID, start, end)
	} else if courseID != "" {
		records, err = h.attendanceRepo.GetBySession(c.Request.Context(), courseID)
	} else {
		response.Error(c, http.StatusBadRequest, "Either student_id or course_id is required")
		return
	}

	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Error getting attendance report")
		return
	}

	report := generateAttendanceReport(records)
	response.Success(c, http.StatusOK, "Attendance report generated successfully", report)
}

// ExportReport exports analytics data in various formats
func (h *AnalyticsHandler) ExportReport(c *gin.Context) {
	var req struct {
		Type    string                 `json:"type" binding:"required,oneof=attendance courses students"`
		Format  string                 `json:"format" binding:"required,oneof=csv excel pdf"`
		Filters map[string]interface{} `json:"filters"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	var (
		data interface{}
		err  error
	)

	ctx := c.Request.Context()

	switch req.Type {
	case "attendance":
		data, err = h.generateAttendanceExport(ctx, req.Filters)
	case "courses":
		data, err = h.generateCourseExport(ctx, req.Filters)
	case "students":
		data, err = h.generateStudentExport(ctx, req.Filters)
	default:
		response.Error(c, http.StatusBadRequest, "Invalid export type")
		return
	}

	if err != nil {
		response.Error(c, http.StatusInternalServerError, fmt.Sprintf("Error generating export: %v", err))
		return
	}

	var exportData []byte
	switch req.Format {
	case "csv":
		exportData, err = exportToCSV(data)
	case "excel":
		exportData, err = exportToExcel(data)
	case "pdf":
		exportData, err = exportToPDF(data)
	}

	if err != nil {
		response.Error(c, http.StatusInternalServerError, fmt.Sprintf("Error exporting data: %v", err))
		return
	}

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=report_%s.%s", req.Type, req.Format))
	c.Data(http.StatusOK, getContentType(req.Format), exportData)
}

// Helper functions

func calculateAttendanceAnalytics(records []*entities.AttendanceRecord) map[string]interface{} {
	var (
		totalSessions   int
		totalPresent    int
		totalLate       int
		totalAbsent     int
		attendanceRate  float64
		averageLateness float64
	)

	// Calculate metrics
	for _, record := range records {
		totalSessions++
		switch record.Status {
		case entities.AttendanceStatusPresent:
			totalPresent++
		case entities.AttendanceStatusLate:
			totalLate++
		case entities.AttendanceStatusAbsent:
			totalAbsent++
		}
	}

	if totalSessions > 0 {
		attendanceRate = float64(totalPresent+totalLate) / float64(totalSessions) * 100
		averageLateness = float64(totalLate) / float64(totalSessions) * 100
	}

	return map[string]interface{}{
		"total_sessions":    totalSessions,
		"total_present":     totalPresent,
		"total_late":        totalLate,
		"total_absent":      totalAbsent,
		"attendance_rate":   attendanceRate,
		"average_lateness":  averageLateness,
		"attendance_trend":  calculateAttendanceTrend(records),
		"monthly_breakdown": calculateMonthlyBreakdown(records),
	}
}

func calculateAttendanceTrend(records []*entities.AttendanceRecord) []map[string]interface{} {
	// Implementation for attendance trend calculation
	return nil
}

func calculateMonthlyBreakdown(records []*entities.AttendanceRecord) []map[string]interface{} {
	// Implementation for monthly breakdown calculation
	return nil
}

func calculateCourseAnalytics(courses []*entities.Course) map[string]interface{} {
	// Implementation for course analytics calculation
	return nil
}

func calculateStudentAnalytics(students []*entities.Student) map[string]interface{} {
	// Implementation for student analytics calculation
	return nil
}

func generateAttendanceReport(records []*entities.AttendanceRecord) map[string]interface{} {
	// Implementation for attendance report generation
	return nil
}

func (h *AnalyticsHandler) generateAttendanceExport(ctx context.Context, filters map[string]interface{}) (interface{}, error) {
	// Implementation for attendance export generation
	return nil, nil
}

func (h *AnalyticsHandler) generateCourseExport(ctx context.Context, filters map[string]interface{}) (interface{}, error) {
	// Implementation for course export generation
	return nil, nil
}

func (h *AnalyticsHandler) generateStudentExport(ctx context.Context, filters map[string]interface{}) (interface{}, error) {
	// Implementation for student export generation
	return nil, nil
}

func exportToCSV(data interface{}) ([]byte, error) {
	// Implementation for CSV export
	return nil, nil
}

func exportToExcel(data interface{}) ([]byte, error) {
	// Implementation for Excel export
	return nil, nil
}

func exportToPDF(data interface{}) ([]byte, error) {
	// Implementation for PDF export
	return nil, nil
}

func getContentType(format string) string {
	switch format {
	case "csv":
		return "text/csv"
	case "excel":
		return "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	case "pdf":
		return "application/pdf"
	default:
		return "application/octet-stream"
	}
}
