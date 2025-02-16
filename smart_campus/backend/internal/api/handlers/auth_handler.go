package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"smart_campus/internal/api/requests"
	"smart_campus/internal/api/response"
	"smart_campus/internal/domain/entities"
	"smart_campus/internal/services"
	"smart_campus/pkg/utils"
)

// AuthHandler handles authentication-related requests
type AuthHandler struct {
	authService   *services.AuthService
	deviceService *services.DeviceService
}

// NewAuthHandler creates a new authentication handler
func NewAuthHandler(
	authService *services.AuthService,
	deviceService *services.DeviceService,
) *AuthHandler {
	return &AuthHandler{
		authService:   authService,
		deviceService: deviceService,
	}
}

// Login handles user login
func (h *AuthHandler) Login(c *gin.Context) {
	var req requests.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	loginReq := &services.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
		DeviceID: req.DeviceID,
	}

	resp, err := h.authService.Login(c.Request.Context(), loginReq)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Login successful", gin.H{
		"access_token":  resp.AccessToken,
		"refresh_token": resp.RefreshToken,
		"user":          resp.User,
	})
}

// Register handles user registration
func (h *AuthHandler) Register(c *gin.Context) {
	var req requests.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	user := entities.NewUser(
		entities.UserRole(req.Role),
		req.Email,
		req.FullName,
		req.Department,
	)

	if req.YearOfStudy != nil {
		if err := user.SetYearOfStudy(*req.YearOfStudy); err != nil {
			response.Error(c, http.StatusBadRequest, err.Error())
			return
		}
	}
	if req.EnrollmentNumber != "" {
		if err := user.SetEnrollmentNumber(req.EnrollmentNumber); err != nil {
			response.Error(c, http.StatusBadRequest, err.Error())
			return
		}
	}
	if req.EmployeeID != "" {
		if err := user.SetEmployeeID(req.EmployeeID); err != nil {
			response.Error(c, http.StatusBadRequest, err.Error())
			return
		}
	}

	// Hash the password
	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to process password")
		return
	}
	user.SetPassword(passwordHash)

	// Create the user
	err = h.authService.CreateUser(c.Request.Context(), user)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "Registration successful", nil)
}

// VerifyOTP handles OTP verification
func (h *AuthHandler) VerifyOTP(c *gin.Context) {
	var req requests.VerifyOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	verifyReq := &services.VerifyDeviceRequest{
		UserID:   c.GetString("userID"),
		DeviceID: c.GetString("deviceID"),
		OTP:      req.OTP,
	}

	err := h.authService.VerifyDevice(c.Request.Context(), verifyReq)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "OTP verified successfully", nil)
}

// ForgotPassword handles password reset requests
func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	sendOTPReq := &services.SendOTPRequest{
		Email: req.Email,
	}

	err := h.authService.SendDeviceVerificationOTP(c.Request.Context(), sendOTPReq)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Password reset instructions sent", nil)
}

// ResetPassword handles password reset
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req requests.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	verifyReq := &services.VerifyDeviceRequest{
		UserID:   c.GetString("userID"),
		DeviceID: c.GetString("deviceID"),
		OTP:      req.OTP,
	}

	err := h.authService.VerifyDevice(c.Request.Context(), verifyReq)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Password reset successful", nil)
}

// RefreshToken handles token refresh
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req services.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.authService.RefreshToken(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Token refreshed successfully", gin.H{
		"access_token":  resp.AccessToken,
		"refresh_token": resp.RefreshToken,
	})
}

// BindDevice handles device binding
func (h *AuthHandler) BindDevice(c *gin.Context) {
	var req requests.BindDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	userID := c.GetString("userID")
	if userID == "" {
		response.Error(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	err := h.deviceService.BindDevice(userID, req.DeviceID, req.DeviceName, req.DeviceModel)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Device bound successfully", nil)
}

// VerifyDevice handles device verification
func (h *AuthHandler) VerifyDevice(c *gin.Context) {
	var req struct {
		DeviceID string `json:"device_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	userID := c.GetString("userID")
	if userID == "" {
		response.Error(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	isValid, err := h.deviceService.VerifyDevice(userID, req.DeviceID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Device verification successful", gin.H{
		"is_valid": isValid,
	})
}

// UnbindDevice handles device unbinding
func (h *AuthHandler) UnbindDevice(c *gin.Context) {
	var req struct {
		DeviceID string `json:"device_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	userID := c.GetString("userID")
	if userID == "" {
		response.Error(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	err := h.deviceService.UnbindDevice(userID, req.DeviceID)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Device unbound successfully", nil)
}
