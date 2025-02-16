package models

import (
	"time"
)

type SessionStatus string

const (
	SessionStatusActive    SessionStatus = "active"
	SessionStatusCompleted SessionStatus = "completed"
	SessionStatusCancelled SessionStatus = "cancelled"
)

type Session struct {
	ID                string        `json:"id"`
	TeacherID         string        `json:"teacher_id"`
	CourseID          string        `json:"course_id"`
	SessionDate       time.Time     `json:"session_date"`
	StartTime         time.Time     `json:"start_time"`
	EndTime           time.Time     `json:"end_time"`
	WifiSSID          *string       `json:"wifi_ssid,omitempty"`
	WifiBSSID         *string       `json:"wifi_bssid,omitempty"`
	LocationLatitude  *float64      `json:"location_latitude,omitempty"`
	LocationLongitude *float64      `json:"location_longitude,omitempty"`
	LocationRadius    *int          `json:"location_radius,omitempty"`
	Status            SessionStatus `json:"status"`
	CreatedAt         time.Time     `json:"created_at"`
	UpdatedAt         time.Time     `json:"updated_at"`
}

type SessionRepository interface {
	Create(session *Session) error
	GetByID(id string) (*Session, error)
	Update(session *Session) error
	Delete(id string) error
	List(offset, limit int) ([]*Session, error)
	GetActiveSessions() ([]*Session, error)
	GetByTeacher(teacherID string) ([]*Session, error)
	GetByCourse(courseID string) ([]*Session, error)
	EndSession(id string) error
	CancelSession(id string) error
}
