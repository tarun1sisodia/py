package models

import "time"

// Session represents an attendance session.
type Session struct {
	ID             string    `json:"id"`
	TeacherID      string    `json:"teacher_id"`
	SubjectID      string    `json:"subject_id"`
	AcademicYear   string    `json:"academic_year"`
	StartTime      time.Time `json:"start_time"`
	EndTime        time.Time `json:"end_time"`
	CountdownDuration string    `json:"countdown_duration"`
	Status         string    `json:"status"`
	WifiSSID       string    `json:"wifi_ssid"`
	WifiBSSID      string    `json:"wifi_bssid"`
	LocationLat    float64   `json:"location_lat"`
	LocationLong   float64   `json:"location_long"`
	CreatedAt      time.Time `json:"created_at"`
}
