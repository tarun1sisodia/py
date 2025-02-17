package controllers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupAttendanceRouter() *gin.Engine {
	router := gin.Default()
	return router
}

func TestMarkAttendance(t *testing.T) {
	router := setupAttendanceRouter()
	router.POST("/attendance/mark", MarkAttendance)

	jsonStr := `{"SessionID": "session123", "StudentID": "student123", "VerificationMethod": "location", "DeviceInfo": "test", "LocationLat": 1.0,  "LocationLong": 1.0, "WifiSSID": "test", "WifiBSSID": "test"}`
	req, _ := http.NewRequest("POST", "/attendance/mark", strings.NewReader(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "Attendance marked successfully")
}

// Add more tests for other attendance endpoints
