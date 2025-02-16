package repositories

import (
	"context"
	"time"
)

// AttendanceRepository defines the interface for attendance-related database operations
type AttendanceRepository interface {
	// Create creates a new attendance record
	Create(ctx context.Context, record *AttendanceRecord) error

	// FindByID finds an attendance record by ID
	FindByID(ctx context.Context, id string) (*AttendanceRecord, error)

	// ExistsForSessionAndStudent checks if a student has already marked attendance for a session
	ExistsForSessionAndStudent(ctx context.Context, sessionID, studentID string) (bool, error)

	// Update updates an attendance record
	Update(ctx context.Context, record *AttendanceRecord) error

	// List lists attendance records with optional filters
	List(ctx context.Context, filters AttendanceFilters) ([]*AttendanceRecord, error)
}

// AttendanceFilters defines filters for listing attendance records
type AttendanceFilters struct {
	SessionID          string
	StudentID          string
	CourseID           string
	StartDate          time.Time
	EndDate            time.Time
	VerificationStatus string
}

// AttendanceRecord represents an attendance record entity
type AttendanceRecord struct {
	ID                 string
	SessionID          string
	StudentID          string
	MarkedAt           time.Time
	WifiSSID           string
	WifiBSSID          string
	LocationLatitude   float64
	LocationLongitude  float64
	DeviceID           string
	VerificationStatus string
	RejectionReason    string
	CreatedAt          time.Time
	UpdatedAt          time.Time
}
