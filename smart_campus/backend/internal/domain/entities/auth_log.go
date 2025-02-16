package entities

import (
	"time"

	"github.com/google/uuid"
)

// AuthLog represents an authentication log entry in the system
type AuthLog struct {
	BaseEntity
	UserID      string        `json:"user_id"`
	Type        AuthLogType   `json:"type"`
	Status      AuthLogStatus `json:"status"`
	DeviceID    string        `json:"device_id"`
	IPAddress   string        `json:"ip_address"`
	UserAgent   string        `json:"user_agent"`
	Location    string        `json:"location"`
	Description string        `json:"description"`
}

// NewAuthLog creates a new authentication log entry
func NewAuthLog(userID string, logType AuthLogType, status AuthLogStatus, deviceID, ipAddress, userAgent, location, description string) *AuthLog {
	now := time.Now()
	return &AuthLog{
		BaseEntity: BaseEntity{
			ID:        uuid.New().String(),
			CreatedAt: now,
			UpdatedAt: now,
		},
		UserID:      userID,
		Type:        logType,
		Status:      status,
		DeviceID:    deviceID,
		IPAddress:   ipAddress,
		UserAgent:   userAgent,
		Location:    location,
		Description: description,
	}
}

// IsSuccess returns true if the log status is success
func (l *AuthLog) IsSuccess() bool {
	return l.Status == AuthLogStatusSuccess
}

// IsFailed returns true if the log status is failed
func (l *AuthLog) IsFailed() bool {
	return l.Status == AuthLogStatusFailed
}

// IsLogin returns true if the log type is login
func (l *AuthLog) IsLogin() bool {
	return l.Type == AuthLogTypeLogin
}

// IsLogout returns true if the log type is logout
func (l *AuthLog) IsLogout() bool {
	return l.Type == AuthLogTypeLogout
}

// IsPasswordChange returns true if the log type is password change
func (l *AuthLog) IsPasswordChange() bool {
	return l.Type == AuthLogTypePasswordChange
}

// IsPasswordReset returns true if the log type is password reset
func (l *AuthLog) IsPasswordReset() bool {
	return l.Type == AuthLogTypePasswordReset
}

// IsDeviceBind returns true if the log type is device bind
func (l *AuthLog) IsDeviceBind() bool {
	return l.Type == AuthLogTypeDeviceBind
}

// IsDeviceUnbind returns true if the log type is device unbind
func (l *AuthLog) IsDeviceUnbind() bool {
	return l.Type == AuthLogTypeDeviceUnbind
}

// HasType returns true if the log has the specified type
func (l *AuthLog) HasType(logType AuthLogType) bool {
	return l.Type == logType
}

// HasStatus returns true if the log has the specified status
func (l *AuthLog) HasStatus(status AuthLogStatus) bool {
	return l.Status == status
}

// ToPublic returns a public view of the auth log
func (l *AuthLog) ToPublic() map[string]interface{} {
	return map[string]interface{}{
		"id":          l.ID,
		"user_id":     l.UserID,
		"type":        l.Type,
		"status":      l.Status,
		"device_id":   l.DeviceID,
		"ip_address":  l.IPAddress,
		"user_agent":  l.UserAgent,
		"location":    l.Location,
		"description": l.Description,
		"created_at":  l.CreatedAt,
		"updated_at":  l.UpdatedAt,
	}
}

// FormatDescription returns a formatted description of the auth log
func (l *AuthLog) FormatDescription() string {
	return l.Description
}

// GetEventTime returns the time of the authentication event
func (l *AuthLog) GetEventTime() time.Time {
	return l.CreatedAt
}
