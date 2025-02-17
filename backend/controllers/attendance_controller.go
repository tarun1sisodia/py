package controllers

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"smart-attendance/models"
	"smart-attendance/services"
)

// MarkAttendance handles the recording of student attendance.
func MarkAttendance(c *gin.Context) {
	var req models.AttendanceRecord
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate that an active session exists
	session, err := services.GetSessionByID(req.SessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve session"})
		return
	}

	if session == nil || session.Status != "active" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No active session found"})
		return
	}

	// Check that the student's Academic Year matches the session's Academic Year
	student, err := services.GetUserByID(req.StudentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve student"})
		return
	}

	if student == nil || student.AcademicYear != session.AcademicYear {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid academic year"})
		return
	}

	// Validate GPS coordinates against campus geofence
	campusLat, err := strconv.ParseFloat(os.Getenv("CAMPUS_LAT"), 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid campus latitude"})
		return
	}

	campusLong, err := strconv.ParseFloat(os.Getenv("CAMPUS_LONG"), 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid campus longitude"})
		return
	}

	radius, err := strconv.ParseFloat(os.Getenv("CAMPUS_RADIUS"), 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid campus radius"})
		return
	}

	if !services.IsWithinRadius(req.LocationLat, req.LocationLong, campusLat, campusLong, radius) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Student is not within campus bounds"})
		return
	}

	// Check Wi-Fi SSID/BSSID against the institution's credentials
	wifiSSID := os.Getenv("WIFI_SSID")
	wifiBSSID := os.Getenv("WIFI_BSSID")

	if req.WifiSSID != wifiSSID || req.WifiBSSID != wifiBSSID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Wi-Fi credentials"})
		return
	}

	// Verify device binding and developer option status
	// Placeholder implementation:
	isDeviceBound := true  // Replace with actual device binding check
	isDeveloperModeEnabled := false // Replace with actual developer mode check

	if !isDeviceBound {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Device is not bound"})
		return
	}

	if isDeveloperModeEnabled {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Developer mode is enabled"})
		return
	}
	log.Println("Device binding check passed.")
	log.Println("Developer mode check passed.")

	// Ensure the student has not already marked attendance for the subject on the same day
	hasMarked, err := services.HasStudentMarkedAttendance(req.StudentID, req.SessionID, time.Now())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check attendance record"})
		return
	}

	if hasMarked {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Attendance already marked for this subject today"})
		return
	}

	// Create a new attendance record
	newRecord := models.AttendanceRecord{
		ID:             uuid.New().String(),
		SessionID:      req.SessionID,
		StudentID:      req.StudentID,
		MarkedAt:       time.Now(),
		VerificationMethod: req.VerificationMethod,
		DeviceInfo:     req.DeviceInfo,
		LocationLat:    req.LocationLat,
		LocationLong:   req.LocationLong,
		WifiSSID:       req.WifiSSID,
		WifiBSSID:      req.WifiBSSID,
		CreatedAt:      time.Now(),
	}

	// Save the attendance record to the database
	err = services.CreateAttendanceRecord(&newRecord)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create attendance record"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Attendance marked successfully"})
}

// GetAttendanceStatus retrieves the attendance status for a student in a session.
func GetAttendanceStatus(c *gin.Context) {
	// Get the student ID and session ID from the request
	studentID := c.Query("student_id")
	sessionID := c.Query("session_id")

	// Retrieve the attendance record from the database
	record, err := services.GetAttendanceRecord(studentID, sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve attendance record"})
		return
	}

	if record == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Attendance not marked"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"attendance": record})
}
