package services

import (
	"time"

	"smart-attendance/config"
	"smart-attendance/models"
)

// CreateAttendanceRecord creates a new attendance record.
func CreateAttendanceRecord(record *models.AttendanceRecord) error {
	_, err := config.DB.Exec("INSERT INTO attendance_records (id, session_id, student_id, marked_at, verification_method, device_info, location_lat, location_long, wifi_ssid, wifi_bssid, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		record.ID, record.SessionID, record.StudentID, record.MarkedAt, record.VerificationMethod, record.DeviceInfo, record.LocationLat, record.LocationLong, record.WifiSSID, record.WifiBSSID, record.CreatedAt)
	return err
}

// GetAttendanceRecord retrieves an attendance record for a student in a session.
func GetAttendanceRecord(studentID, sessionID string) (*models.AttendanceRecord, error) {
	record := &models.AttendanceRecord{}
	err := config.DB.QueryRow("SELECT id, session_id, student_id, marked_at, verification_method, device_info, location_lat, location_long, wifi_ssid, wifi_bssid, created_at FROM attendance_records WHERE student_id = ? AND session_id = ?", studentID, sessionID).Scan(
		&record.ID, &record.SessionID, &record.StudentID, &record.MarkedAt, &record.VerificationMethod, &record.DeviceInfo, &record.LocationLat, &record.LocationLong, &record.WifiSSID, &record.WifiBSSID, &record.CreatedAt)
	if err != nil {
		return nil, err
	}
	return record, nil
}

// HasStudentMarkedAttendance checks if a student has already marked attendance for a session on a given day.
func HasStudentMarkedAttendance(studentID, sessionID string, date time.Time) (bool, error) {
	var count int
	err := config.DB.QueryRow("SELECT COUNT(*) FROM attendance_records WHERE student_id = ? AND session_id = ? AND DATE(marked_at) = DATE(?)", studentID, sessionID, date).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
