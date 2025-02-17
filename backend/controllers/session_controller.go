package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"smart-attendance/models"
	"smart-attendance/services"
)

// StartSession handles the creation of a new attendance session.
func StartSession(c *gin.Context) {
	var req models.Session
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate teacher credentials, academic year, and subject
	// (Add your validation logic here)

	// Create a new session
	newSession := models.Session{
		ID:              uuid.New().String(),
		TeacherID:       req.TeacherID,
		SubjectID:       req.SubjectID,
		AcademicYear:    req.AcademicYear,
		StartTime:       time.Now(),
		EndTime:         time.Now().Add(time.Minute * 30), // Default end time (adjust as needed)
		CountdownDuration: req.CountdownDuration,
		Status:          "active",
		WifiSSID:        req.WifiSSID,
		WifiBSSID:       req.WifiBSSID,
		LocationLat:     req.LocationLat,
		LocationLong:    req.LocationLong,
		CreatedAt:       time.Now(),
	}

	// Save the session to the database
	err := services.CreateSession(&newSession)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Session created successfully", "session_id": newSession.ID})
}

// GetActiveSession retrieves the details of an active session.
func GetActiveSession(c *gin.Context) {
	// Get the teacher ID from the request (e.g., from the JWT token)
	teacherID := c.GetString("teacher_id") // Assuming you have a middleware to extract teacher ID

	// Retrieve the active session from the database
	session, err := services.GetActiveSessionByTeacherID(teacherID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve active session"})
		return
	}

	if session == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "No active session found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"session": session})
}

// EndSession handles the ending of an active session.
func EndSession(c *gin.Context) {
	sessionID := c.Param("id")

	// Retrieve the session from the database
	session, err := services.GetSessionByID(sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve session"})
		return
	}

	if session == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Session not found"})
		return
	}

	// Update the session status to "completed"
	session.Status = "completed"
	err = services.UpdateSession(session)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Session ended successfully"})
}
