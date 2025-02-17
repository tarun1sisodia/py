package services

import (
	"smart-attendance/config"
	"smart-attendance/models"
)

// CreateSession creates a new session.
func CreateSession(session *models.Session) error {
	_, err := config.DB.Exec("INSERT INTO attendance_sessions (id, teacher_id, subject_id, academic_year, start_time, end_time, countdown_duration, status, wifi_ssid, wifi_bssid, location_lat, location_long, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		session.ID, session.TeacherID, session.SubjectID, session.AcademicYear, session.StartTime, session.EndTime, session.CountdownDuration, session.Status, session.WifiSSID, session.WifiBSSID, session.LocationLat, session.LocationLong, session.CreatedAt)
	return err
}

// GetSessionByID gets a session by ID.
func GetSessionByID(id string) (*models.Session, error) {
	session := &models.Session{}
	err := config.DB.QueryRow("SELECT id, teacher_id, subject_id, academic_year, start_time, end_time, countdown_duration, status, wifi_ssid, wifi_bssid, location_lat, location_long, created_at FROM attendance_sessions WHERE id = ?", id).Scan(
		&session.ID, &session.TeacherID, &session.SubjectID, &session.AcademicYear, &session.StartTime, &session.EndTime, &session.CountdownDuration, &session.Status, &session.WifiSSID, &session.WifiBSSID, &session.LocationLat, &session.LocationLong, &session.CreatedAt)
	if err != nil {
		return nil, err
	}
	return session, nil
}

// GetActiveSessionByTeacherID gets the active session for a teacher.
func GetActiveSessionByTeacherID(teacherID string) (*models.Session, error) {
	session := &models.Session{}
	err := config.DB.QueryRow("SELECT id, teacher_id, subject_id, academic_year, start_time, end_time, countdown_duration, status, wifi_ssid, wifi_bssid, location_lat, location_long, created_at FROM attendance_sessions WHERE teacher_id = ? AND status = 'active'", teacherID).Scan(
		&session.ID, &session.TeacherID, &session.SubjectID, &session.AcademicYear, &session.StartTime, &session.EndTime, &session.CountdownDuration, &session.Status, &session.WifiSSID, &session.WifiBSSID, &session.LocationLat, &session.LocationLong, &session.CreatedAt)
	if err != nil {
		return nil, err
	}
	return session, nil
}

// UpdateSession updates a session.
func UpdateSession(session *models.Session) error {
	_, err := config.DB.Exec("UPDATE attendance_sessions SET teacher_id = ?, subject_id = ?, academic_year = ?, start_time = ?, end_time = ?, countdown_duration = ?, status = ?, wifi_ssid = ?, wifi_bssid = ?, location_lat = ?, location_long = ?, created_at = ? WHERE id = ?",
		session.TeacherID, session.SubjectID, session.AcademicYear, session.StartTime, session.EndTime, session.CountdownDuration, session.Status, session.WifiSSID, session.WifiBSSID, session.LocationLat, session.LocationLong, session.CreatedAt, session.ID)
	return err
}
