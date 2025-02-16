package entities

import (
	"time"
)

// DeviceBinding represents a device binding in the system
type DeviceBinding struct {
	ID            string    `json:"id"`
	UserID        string    `json:"user_id"`
	DeviceID      string    `json:"device_id"`
	DeviceName    string    `json:"device_name"`
	DeviceModel   string    `json:"device_model"`
	IsActive      bool      `json:"is_active"`
	IsBlacklisted bool      `json:"is_blacklisted"`
	BoundAt       time.Time `json:"bound_at"`
	LastUsedAt    time.Time `json:"last_used_at,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// Validate validates the device binding entity
func (d *DeviceBinding) Validate() error {
	// Add validation logic here
	return nil
}

// IsValid returns true if the device binding is valid for use
func (d *DeviceBinding) IsValid() bool {
	return d.IsActive && !d.IsBlacklisted
}

// UpdateLastUsed updates the last used timestamp
func (d *DeviceBinding) UpdateLastUsed() {
	d.LastUsedAt = time.Now()
	d.UpdatedAt = time.Now()
}

// Activate activates the device binding
func (d *DeviceBinding) Activate() {
	d.IsActive = true
	d.UpdatedAt = time.Now()
}

// Deactivate deactivates the device binding
func (d *DeviceBinding) Deactivate() {
	d.IsActive = false
	d.UpdatedAt = time.Now()
}

// Blacklist blacklists the device binding
func (d *DeviceBinding) Blacklist() {
	d.IsBlacklisted = true
	d.IsActive = false
	d.UpdatedAt = time.Now()
}

// NewDeviceBinding creates a new device binding instance
func NewDeviceBinding(userID, deviceID, deviceName, deviceModel string) *DeviceBinding {
	now := time.Now()
	return &DeviceBinding{
		UserID:      userID,
		DeviceID:    deviceID,
		DeviceName:  deviceName,
		DeviceModel: deviceModel,
		IsActive:    false,
		BoundAt:     now,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}
