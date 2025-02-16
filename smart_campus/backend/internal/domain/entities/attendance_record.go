package entities

import (
	"time"

	"github.com/google/uuid"
)

// AttendanceRecord represents an attendance record in the system
type AttendanceRecord struct {
	BaseEntity
	SessionID       string           `json:"session_id"`
	StudentID       string           `json:"student_id"`
	Status          AttendanceStatus `json:"status"`
	LocationLat     float64          `json:"location_lat"`
	LocationLong    float64          `json:"location_long"`
	WiFiSSID        string           `json:"wifi_ssid"`
	WiFiBSSID       string           `json:"wifi_bssid"`
	DeviceID        string           `json:"device_id"`
	VerificationLog string           `json:"verification_log"`
	MarkedAt        time.Time        `json:"marked_at"`
}

// IsPresent returns true if the record's status is present
func (r *AttendanceRecord) IsPresent() bool {
	return r.Status == AttendanceStatusPresent
}

// IsLate returns true if the record's status is late
func (r *AttendanceRecord) IsLate() bool {
	return r.Status == AttendanceStatusLate
}

// IsAbsent returns true if the record's status is absent
func (r *AttendanceRecord) IsAbsent() bool {
	return r.Status == AttendanceStatusAbsent
}

// IsExcused returns true if the record's status is excused
func (r *AttendanceRecord) IsExcused() bool {
	return r.Status == AttendanceStatusExcused
}

// IsPending returns true if the record's status is pending
func (r *AttendanceRecord) IsPending() bool {
	return r.Status == AttendanceStatusPending
}

// IsVerified returns true if the record's status is verified
func (r *AttendanceRecord) IsVerified() bool {
	return r.Status == AttendanceStatusVerified
}

// IsRejected returns true if the record's status is rejected
func (r *AttendanceRecord) IsRejected() bool {
	return r.Status == AttendanceStatusRejected
}

// HasStatus returns true if the record has the specified status
func (r *AttendanceRecord) HasStatus(status AttendanceStatus) bool {
	return r.Status == status
}

// UpdateStatus updates the record's status
func (r *AttendanceRecord) UpdateStatus(status AttendanceStatus) {
	r.Status = status
	r.UpdatedAt = time.Now()
}

// UpdateLocation updates the record's location
func (r *AttendanceRecord) UpdateLocation(lat, long float64) {
	r.LocationLat = lat
	r.LocationLong = long
	r.UpdatedAt = time.Now()
}

// UpdateWiFi updates the record's WiFi information
func (r *AttendanceRecord) UpdateWiFi(ssid, bssid string) {
	r.WiFiSSID = ssid
	r.WiFiBSSID = bssid
	r.UpdatedAt = time.Now()
}

// AddVerificationLog adds a log entry to the verification log
func (r *AttendanceRecord) AddVerificationLog(log string) {
	if r.VerificationLog != "" {
		r.VerificationLog += "\n"
	}
	r.VerificationLog += log
	r.UpdatedAt = time.Now()
}

// ToPublic returns a public view of the attendance record
func (r *AttendanceRecord) ToPublic() map[string]interface{} {
	return map[string]interface{}{
		"id":               r.ID,
		"session_id":       r.SessionID,
		"student_id":       r.StudentID,
		"status":           r.Status,
		"location_lat":     r.LocationLat,
		"location_long":    r.LocationLong,
		"wifi_ssid":        r.WiFiSSID,
		"wifi_bssid":       r.WiFiBSSID,
		"device_id":        r.DeviceID,
		"verification_log": r.VerificationLog,
		"marked_at":        r.MarkedAt,
		"created_at":       r.CreatedAt,
		"updated_at":       r.UpdatedAt,
	}
}

// NewAttendanceRecord creates a new attendance record
func NewAttendanceRecord(sessionID, studentID string, status AttendanceStatus, lat, long float64, ssid, bssid, deviceID string) *AttendanceRecord {
	now := time.Now()
	return &AttendanceRecord{
		ID:           uuid.New().String(),
		SessionID:    sessionID,
		StudentID:    studentID,
		Status:       status,
		LocationLat:  lat,
		LocationLong: long,
		WiFiSSID:     ssid,
		WiFiBSSID:    bssid,
		DeviceID:     deviceID,
		MarkedAt:     now,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// IsMarkedByValidDevice checks if the attendance was marked by a valid device
func (r *AttendanceRecord) IsMarkedByValidDevice(validDeviceID string) bool {
	return r.DeviceID == validDeviceID
}
