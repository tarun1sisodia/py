package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"smart_campus/internal/api/requests"
	"smart_campus/internal/api/response"
	"smart_campus/internal/services"
)

// DeviceHandler handles device-related requests
type DeviceHandler struct {
	deviceService *services.DeviceService
}

// NewDeviceHandler creates a new device handler
func NewDeviceHandler(deviceService *services.DeviceService) *DeviceHandler {
	return &DeviceHandler{
		deviceService: deviceService,
	}
}

// BindDevice binds a device to a user
func (h *DeviceHandler) BindDevice(c *gin.Context) {
	var req requests.BindDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body")
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

// VerifyDevice verifies a device binding
func (h *DeviceHandler) VerifyDevice(c *gin.Context) {
	var req struct {
		DeviceID string `json:"device_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body")
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

// UnbindDevice removes a device binding
func (h *DeviceHandler) UnbindDevice(c *gin.Context) {
	var req struct {
		DeviceID string `json:"device_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body")
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

// ListDevices returns a list of devices bound to a user
func (h *DeviceHandler) ListDevices(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		response.Error(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	devices, err := h.deviceService.ListDevices(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Devices retrieved successfully", devices)
}

// GetDevice returns details of a specific device
func (h *DeviceHandler) GetDevice(c *gin.Context) {
	deviceID := c.Param("id")
	if deviceID == "" {
		response.Error(c, http.StatusBadRequest, "Device ID is required")
		return
	}

	userID := c.GetString("userID")
	if userID == "" {
		response.Error(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	device, err := h.deviceService.GetDevice(userID, deviceID)
	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Device retrieved successfully", device)
}
