package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"smart_campus_backend/internal/models"
	"smart_campus_backend/internal/services"
	"smart_campus_backend/tests/helpers"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRegisterUser(t *testing.T) {
	// Initialize test database
	db := helpers.InitTestDB()
	defer helpers.CleanupTestDB()

	// Setup router and services
	router := gin.Default()
	authService := services.NewAuthService(db)

	// Setup test route
	router.POST("/api/auth/register", func(c *gin.Context) {
		var req models.RegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Error:   "ValidationError",
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			})
			return
		}

		user, err := authService.Register(c.Request.Context(), &req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Error:   "RegistrationError",
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			})
			return
		}

		c.JSON(http.StatusCreated, models.SuccessResponse{
			Message: "User registered successfully",
			Data:    user,
		})
	})

	// Test cases
	tests := []struct {
		name           string
		requestBody    models.RegisterRequest
		expectedStatus int
		validateResp   func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "Valid Registration",
			requestBody: models.RegisterRequest{
				Role:     "student",
				FullName: "Test Student",
				Email:    "test.student@example.com",
				Phone:    "+1234567890",
				Password: "password123",
			},
			expectedStatus: http.StatusCreated,
			validateResp: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response models.SuccessResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "User registered successfully", response.Message)

				user, ok := response.Data.(map[string]interface{})
				assert.True(t, ok)
				assert.Equal(t, "test.student@example.com", user["email"])
				assert.Equal(t, "student", user["role"])
			},
		},
		{
			name: "Invalid Role",
			requestBody: models.RegisterRequest{
				Role:     "invalid",
				FullName: "Test Student",
				Email:    "test.student@example.com",
				Phone:    "+1234567890",
				Password: "password123",
			},
			expectedStatus: http.StatusBadRequest,
			validateResp: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response models.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "ValidationError", response.Error)
			},
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			tt.validateResp(t, w)
		})
	}
}
