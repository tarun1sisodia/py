package models

import "time"

// AttendanceRecord represents a student's attendance record.
type AttendanceRecord struct {
	ID             string    `json:"id"`
	SessionID      string    `json:"session_id"`
	StudentID      string    `json:"student_id"`
	MarkedAt       time.Time `json:"marked_at"`
	VerificationMethod string    `json:"verification_method"`
	DeviceInfo     string    `json:"device_info"`
	LocationLat    float64   `json:"location_lat"`
	LocationLong   float64   `json:"location_long"`
	WifiSSID       string    `json:"wifi_ssid"`
	WifiBSSID      string    `json:"wifi_bssid"`
	CreatedAt      time.Time `json:"created_at"`
}
