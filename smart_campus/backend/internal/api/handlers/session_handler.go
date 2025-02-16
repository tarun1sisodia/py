package handlers

import (
	"net/http"
	"smart_campus/internal/models"
	"smart_campus/internal/services"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SessionHandler struct {
	sessionService *services.SessionService
}

func NewSessionHandler(sessionService *services.SessionService) *SessionHandler {
	return &SessionHandler{
		sessionService: sessionService,
	}
}

type createSessionRequest struct {
	CourseID          string    `json:"course_id" binding:"required"`
	SessionDate       time.Time `json:"session_date" binding:"required"`
	StartTime         time.Time `json:"start_time" binding:"required"`
	EndTime           time.Time `json:"end_time" binding:"required"`
	LocationLatitude  *float64  `json:"location_latitude,omitempty"`
	LocationLongitude *float64  `json:"location_longitude,omitempty"`
	LocationRadius    *int      `json:"location_radius,omitempty"`
	WifiSSID          *string   `json:"wifi_ssid,omitempty"`
	WifiBSSID         *string   `json:"wifi_bssid,omitempty"`
}

func (h *SessionHandler) CreateSession(c *gin.Context) {
	var req createSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	teacherID := c.GetString("userID")
	if teacherID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	session := &models.Session{
		ID:                uuid.New().String(),
		TeacherID:         teacherID,
		CourseID:          req.CourseID,
		SessionDate:       req.SessionDate,
		StartTime:         req.StartTime,
		EndTime:           req.EndTime,
		LocationLatitude:  req.LocationLatitude,
		LocationLongitude: req.LocationLongitude,
		LocationRadius:    req.LocationRadius,
		WifiSSID:          req.WifiSSID,
		WifiBSSID:         req.WifiBSSID,
		Status:            models.SessionStatusActive,
	}

	err := h.sessionService.CreateSession(session)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, session)
}

func (h *SessionHandler) GetSessions(c *gin.Context) {
	sessions, err := h.sessionService.GetActiveSessions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, sessions)
}

func (h *SessionHandler) GetSessionById(c *gin.Context) {
	sessionID := c.Param("id")
	session, err := h.sessionService.GetSessionById(sessionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
		return
	}

	c.JSON(http.StatusOK, session)
}

func (h *SessionHandler) EndSession(c *gin.Context) {
	sessionID := c.Param("id")
	err := h.sessionService.EndSession(sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Session ended successfully"})
}

func (h *SessionHandler) CancelSession(c *gin.Context) {
	sessionID := c.Param("id")
	err := h.sessionService.CancelSession(sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Session cancelled successfully"})
}

func (h *SessionHandler) GetSessionAttendance(c *gin.Context) {
	sessionID := c.Param("id")
	records, err := h.sessionService.GetSessionAttendance(sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, records)
}

func (h *SessionHandler) MarkAttendance(c *gin.Context) {
	sessionID := c.Param("id")
	studentID := c.GetString("userID")
	if studentID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req struct {
		LocationLatitude  *float64 `json:"location_latitude,omitempty"`
		LocationLongitude *float64 `json:"location_longitude,omitempty"`
		WifiSSID          *string  `json:"wifi_ssid,omitempty"`
		WifiBSSID         *string  `json:"wifi_bssid,omitempty"`
		DeviceID          string   `json:"device_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	record := &models.AttendanceRecord{
		ID:                 uuid.New().String(),
		SessionID:          sessionID,
		StudentID:          studentID,
		MarkedAt:           time.Now(),
		LocationLatitude:   req.LocationLatitude,
		LocationLongitude:  req.LocationLongitude,
		WifiSSID:           req.WifiSSID,
		WifiBSSID:          req.WifiBSSID,
		DeviceID:           req.DeviceID,
		VerificationStatus: models.VerificationStatusPending,
	}

	err := h.sessionService.MarkAttendance(record)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, record)
}

func (h *SessionHandler) VerifyAttendance(c *gin.Context) {
	attendanceID := c.Param("id")
	err := h.sessionService.VerifyAttendance(attendanceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Attendance verified successfully"})
}

func (h *SessionHandler) RejectAttendance(c *gin.Context) {
	attendanceID := c.Param("id")
	var req struct {
		Reason string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.sessionService.RejectAttendance(attendanceID, req.Reason)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Attendance rejected successfully"})
}

func (h *SessionHandler) GetActiveSessions(c *gin.Context) {
	sessions, err := h.sessionService.GetActiveSessions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, sessions)
}

func (h *SessionHandler) GetAttendanceHistory(c *gin.Context) {
	studentID := c.GetString("userID")
	if studentID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	records, err := h.sessionService.GetStudentAttendance(studentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, records)
}

func (h *SessionHandler) GetAttendanceStatistics(c *gin.Context) {
	studentID := c.GetString("userID")
	if studentID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	courseID := c.Query("course_id")
	if courseID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Course ID is required"})
		return
	}

	startDate := time.Now().AddDate(0, -1, 0) // Last month
	endDate := time.Now()

	stats, err := h.sessionService.GetAttendanceStatistics(studentID, courseID, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}
