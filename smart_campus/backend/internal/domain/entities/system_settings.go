package entities

import (
	"encoding/json"
	"time"
)

// SystemSettingKey represents a key for system settings
type SystemSettingKey string

const (
	// SystemSettingKeyOTPExpiration represents the OTP expiration time in minutes
	SystemSettingKeyOTPExpiration SystemSettingKey = "otp_expiration"
	// SystemSettingKeyMaxLoginAttempts represents the maximum number of login attempts
	SystemSettingKeyMaxLoginAttempts SystemSettingKey = "max_login_attempts"
	// SystemSettingKeyLocationRadius represents the acceptable radius in meters for location verification
	SystemSettingKeyLocationRadius SystemSettingKey = "location_radius"
	// SystemSettingKeySessionDuration represents the session duration in minutes
	SystemSettingKeySessionDuration SystemSettingKey = "session_duration"
	// SystemSettingKeyRefreshTokenDuration represents the refresh token duration in days
	SystemSettingKeyRefreshTokenDuration SystemSettingKey = "refresh_token_duration"
	// SystemSettingKeyMaxDeviceBindings represents the maximum number of device bindings per user
	SystemSettingKeyMaxDeviceBindings SystemSettingKey = "max_device_bindings"
)

// SystemSetting represents a system setting
type SystemSetting struct {
	Key         SystemSettingKey `json:"key"`
	Value       string           `json:"value"`
	Description string           `json:"description"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}

// NewSystemSetting creates a new system setting
func NewSystemSetting(key SystemSettingKey, value interface{}, description string) (*SystemSetting, error) {
	valueStr, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	return &SystemSetting{
		Key:         key,
		Value:       string(valueStr),
		Description: description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

// GetIntValue returns the setting value as an integer
func (s *SystemSetting) GetIntValue() (int, error) {
	var value int
	err := json.Unmarshal([]byte(s.Value), &value)
	return value, err
}

// GetFloatValue returns the setting value as a float64
func (s *SystemSetting) GetFloatValue() (float64, error) {
	var value float64
	err := json.Unmarshal([]byte(s.Value), &value)
	return value, err
}

// GetBoolValue returns the setting value as a boolean
func (s *SystemSetting) GetBoolValue() (bool, error) {
	var value bool
	err := json.Unmarshal([]byte(s.Value), &value)
	return value, err
}

// GetStringValue returns the setting value as a string
func (s *SystemSetting) GetStringValue() (string, error) {
	var value string
	err := json.Unmarshal([]byte(s.Value), &value)
	return value, err
}

// UpdateValue updates the setting value
func (s *SystemSetting) UpdateValue(value interface{}) error {
	valueStr, err := json.Marshal(value)
	if err != nil {
		return err
	}

	s.Value = string(valueStr)
	s.UpdatedAt = time.Now()
	return nil
}
