package models

import "time"

// VerificationStatus represents the status of an attendance verification
type VerificationStatus string

const (
	VerificationStatusPending  VerificationStatus = "pending"
	VerificationStatusVerified VerificationStatus = "verified"
	VerificationStatusRejected VerificationStatus = "rejected"
)

// AttendanceRecord represents a student's attendance record for a session
type AttendanceRecord struct {
	ID                 string             `json:"id"`
	SessionID          string             `json:"session_id"`
	StudentID          string             `json:"student_id"`
	MarkedAt           time.Time          `json:"marked_at"`
	WifiSSID           *string            `json:"wifi_ssid,omitempty"`
	WifiBSSID          *string            `json:"wifi_bssid,omitempty"`
	LocationLatitude   *float64           `json:"location_latitude,omitempty"`
	LocationLongitude  *float64           `json:"location_longitude,omitempty"`
	DeviceID           string             `json:"device_id"`
	VerificationStatus VerificationStatus `json:"verification_status"`
	RejectionReason    *string            `json:"rejection_reason,omitempty"`
	CreatedAt          time.Time          `json:"created_at"`
	UpdatedAt          time.Time          `json:"updated_at"`
}

// AttendanceStatistics represents attendance statistics for a student in a course
type AttendanceStatistics struct {
	TotalSessions        int     `json:"total_sessions"`
	VerifiedSessions     int     `json:"verified_sessions"`
	RejectedSessions     int     `json:"rejected_sessions"`
	PendingSessions      int     `json:"pending_sessions"`
	AttendancePercentage float64 `json:"attendance_percentage"`
}

// AttendanceRepository defines the interface for attendance record operations
type AttendanceRepository interface {
	Create(record *AttendanceRecord) error
	GetByID(id string) (*AttendanceRecord, error)
	Update(record *AttendanceRecord) error
	Delete(id string) error
	List(offset, limit int) ([]*AttendanceRecord, error)
	GetBySession(sessionID string) ([]*AttendanceRecord, error)
	GetByStudent(studentID string) ([]*AttendanceRecord, error)
	VerifyAttendance(id string) error
	RejectAttendance(id string, reason string) error
	GetAttendanceStatistics(studentID string, courseID string, startDate time.Time, endDate time.Time) (*AttendanceStatistics, error)
}
