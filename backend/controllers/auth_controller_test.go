package controllers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *gin.Engine {
	router := gin.Default()
	return router
}

func TestRegisterTeacher(t *testing.T) {
	router := setupRouter()
	router.POST("/auth/register/teacher", RegisterTeacher)

	jsonStr := `{"FullName": "Test Teacher", "Username": "testteacher", "Email": "test@example.com", "Phone": "1234567890", "HighestDegree": "PhD", "Experience": "5 years", "PasswordHash": "password"}`
	req, _ := http.NewRequest("POST", "/auth/register/teacher", strings.NewReader(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "Teacher registered successfully")
}

func TestRegisterStudent(t *testing.T) {
	router := setupRouter()
	router.POST("/auth/register/student", RegisterStudent)

	jsonStr := `{"FullName": "Test Student", "RollNumber": "123", "Course": "BCA", "AcademicYear": "1st Year", "Phone": "1234567890", "PasswordHash": "password"}`
	req, _ := http.NewRequest("POST", "/auth/register/student", strings.NewReader(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "Student registered successfully")
}

// Add more tests for other auth endpoints
