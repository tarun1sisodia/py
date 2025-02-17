package controllers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupSessionRouter() *gin.Engine {
	router := gin.Default()
	return router
}

func TestStartSession(t *testing.T) {
	router := setupSessionRouter()
	router.POST("/sessions/start", StartSession)

	jsonStr := `{"TeacherID": "teacher123", "SubjectID": "subject123", "AcademicYear": "1st Year", "CountdownDuration": "30s", "WifiSSID": "test", "WifiBSSID": "test", "LocationLat": 1.0,  "LocationLong": 1.0}`
	req, _ := http.NewRequest("POST", "/sessions/start", strings.NewReader(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "Session created successfully")
}

// Add more tests for other session endpoints
