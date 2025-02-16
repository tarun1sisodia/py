package repositories

import (
	"context"
	"time"
)

// SessionRepository defines the interface for session-related database operations
type SessionRepository interface {
	// Create creates a new session
	Create(ctx context.Context, session *AttendanceSession) error

	// FindByID finds a session by ID
	FindByID(ctx context.Context, id string) (*AttendanceSession, error)

	// FindActiveByCourse finds active sessions for a course
	FindActiveByCourse(ctx context.Context, courseID string) (*AttendanceSession, error)

	// Update updates a session
	Update(ctx context.Context, session *AttendanceSession) error

	// List lists all sessions with optional filters
	List(ctx context.Context, filters SessionFilters) ([]*AttendanceSession, error)
}

// SessionFilters defines filters for listing sessions
type SessionFilters struct {
	TeacherID string
	CourseID  string
	Status    string
	StartDate time.Time
	EndDate   time.Time
}

// AttendanceSession represents a session entity
type AttendanceSession struct {
	ID                string
	TeacherID         string
	CourseID          string
	SessionDate       time.Time
	StartTime         time.Time
	EndTime           time.Time
	WifiSSID          string
	WifiBSSID         string
	LocationLatitude  float64
	LocationLongitude float64
	LocationRadius    int
	Status            string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
