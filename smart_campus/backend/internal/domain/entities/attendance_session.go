package entities

import (
	"time"

	"github.com/google/uuid"
)

// AttendanceSession represents an attendance session in the system
type AttendanceSession struct {
	BaseEntity
	CourseID      string        `json:"course_id"`
	TeacherID     string        `json:"teacher_id"`
	SessionDate   time.Time     `json:"session_date"`
	StartTime     time.Time     `json:"start_time"`
	EndTime       time.Time     `json:"end_time"`
	LocationLat   float64       `json:"location_lat"`
	LocationLong  float64       `json:"location_long"`
	LocationName  string        `json:"location_name"`
	WiFiSSID      string        `json:"wifi_ssid"`
	WiFiBSSID     string        `json:"wifi_bssid"`
	Status        SessionStatus `json:"status"`
	Description   string        `json:"description"`
	CourseName    string        `json:"course_name"`
	TeacherName   string        `json:"teacher_name"`
	TotalStudents int           `json:"total_students"`
	PresentCount  int           `json:"present_count"`
	LateCount     int           `json:"late_count"`
	AbsentCount   int           `json:"absent_count"`
}

// NewAttendanceSession creates a new attendance session
func NewAttendanceSession(courseID, teacherID string, startTime, endTime time.Time, lat, long float64, ssid, bssid string) *AttendanceSession {
	now := time.Now()
	return &AttendanceSession{
		BaseEntity: BaseEntity{
			ID:        uuid.New().String(),
			CreatedAt: now,
			UpdatedAt: now,
		},
		CourseID:     courseID,
		TeacherID:    teacherID,
		StartTime:    startTime,
		EndTime:      endTime,
		LocationLat:  lat,
		LocationLong: long,
		WiFiSSID:     ssid,
		WiFiBSSID:    bssid,
		Status:       SessionStatusActive,
	}
}

// IsScheduled returns true if the session's status is scheduled
func (s *AttendanceSession) IsScheduled() bool {
	return s.Status == SessionStatusScheduled
}

// IsActive returns true if the session's status is active
func (s *AttendanceSession) IsActive() bool {
	return s.Status == SessionStatusActive
}

// IsComplete returns true if the session's status is complete
func (s *AttendanceSession) IsComplete() bool {
	return s.Status == SessionStatusComplete
}

// IsCancelled returns true if the session's status is cancelled
func (s *AttendanceSession) IsCancelled() bool {
	return s.Status == SessionStatusCancelled
}

// HasStatus returns true if the session has the specified status
func (s *AttendanceSession) HasStatus(status SessionStatus) bool {
	return s.Status == status
}

// UpdateStatus updates the session's status
func (s *AttendanceSession) UpdateStatus(status SessionStatus) {
	s.Status = status
	s.UpdatedAt = time.Now()
}

// UpdateLocation updates the session's location
func (s *AttendanceSession) UpdateLocation(lat, long float64, name string) {
	s.LocationLat = lat
	s.LocationLong = long
	s.LocationName = name
	s.UpdatedAt = time.Now()
}

// UpdateWiFi updates the session's WiFi information
func (s *AttendanceSession) UpdateWiFi(ssid, bssid string) {
	s.WiFiSSID = ssid
	s.WiFiBSSID = bssid
	s.UpdatedAt = time.Now()
}

// UpdateTime updates the session's time information
func (s *AttendanceSession) UpdateTime(sessionDate, startTime, endTime time.Time) {
	s.SessionDate = sessionDate
	s.StartTime = startTime
	s.EndTime = endTime
	s.UpdatedAt = time.Now()
}

// UpdateCounts updates the session's attendance counts
func (s *AttendanceSession) UpdateCounts(totalStudents, presentCount, lateCount, absentCount int) {
	s.TotalStudents = totalStudents
	s.PresentCount = presentCount
	s.LateCount = lateCount
	s.AbsentCount = absentCount
	s.UpdatedAt = time.Now()
}

// GetAttendanceRate returns the attendance rate for the session
func (s *AttendanceSession) GetAttendanceRate() float64 {
	if s.TotalStudents == 0 {
		return 0
	}
	return float64(s.PresentCount+s.LateCount) / float64(s.TotalStudents) * 100
}

// IsInProgress returns true if the session is currently in progress
func (s *AttendanceSession) IsInProgress() bool {
	now := time.Now()
	return s.Status == SessionStatusActive &&
		now.After(s.StartTime) && now.Before(s.EndTime)
}

// CanStart returns true if the session can be started
func (s *AttendanceSession) CanStart() bool {
	now := time.Now()
	return s.Status == SessionStatusScheduled &&
		now.After(s.StartTime.Add(-15*time.Minute)) &&
		now.Before(s.EndTime)
}

// CanComplete returns true if the session can be completed
func (s *AttendanceSession) CanComplete() bool {
	return s.Status == SessionStatusActive &&
		time.Now().After(s.EndTime)
}

// ToPublic returns a public view of the attendance session
func (s *AttendanceSession) ToPublic() map[string]interface{} {
	return map[string]interface{}{
		"id":             s.ID,
		"course_id":      s.CourseID,
		"teacher_id":     s.TeacherID,
		"session_date":   s.SessionDate,
		"start_time":     s.StartTime,
		"end_time":       s.EndTime,
		"location_lat":   s.LocationLat,
		"location_long":  s.LocationLong,
		"location_name":  s.LocationName,
		"wifi_ssid":      s.WiFiSSID,
		"wifi_bssid":     s.WiFiBSSID,
		"status":         s.Status,
		"description":    s.Description,
		"course_name":    s.CourseName,
		"teacher_name":   s.TeacherName,
		"total_students": s.TotalStudents,
		"present_count":  s.PresentCount,
		"late_count":     s.LateCount,
		"absent_count":   s.AbsentCount,
		"created_at":     s.CreatedAt,
		"updated_at":     s.UpdatedAt,
	}
}

// ValidateLocation checks if the provided location is within acceptable range
func (s *AttendanceSession) ValidateLocation(lat, long float64) bool {
	// TODO: Implement location validation logic
	// This should check if the provided coordinates are within an acceptable range
	// of the session's location
	return true
}

// ValidateWiFi checks if the provided WiFi details match the session
func (s *AttendanceSession) ValidateWiFi(ssid, bssid string) bool {
	return s.WiFiSSID == ssid && s.WiFiBSSID == bssid
}
