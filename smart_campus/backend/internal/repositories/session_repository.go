package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"smart_campus/internal/database"
	"smart_campus/internal/models"

	"github.com/google/uuid"
)

type MySQLSessionRepository struct {
	*BaseRepository
	db *database.MySQLDB
}

func NewMySQLSessionRepository(db *database.MySQLDB) models.SessionRepository {
	return &MySQLSessionRepository{
		BaseRepository: &BaseRepository{db: db.DB},
		db:             db,
	}
}

func (r *MySQLSessionRepository) Create(session *models.Session) error {
	if session.ID == "" {
		session.ID = uuid.New().String()
	}
	now := time.Now()
	session.CreatedAt = now
	session.UpdatedAt = now

	query := `
		INSERT INTO attendance_sessions (
			id, teacher_id, course_id, session_date, start_time, end_time,
			wifi_ssid, wifi_bssid, location_latitude, location_longitude,
			location_radius, status, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.Exec(query,
		session.ID, session.TeacherID, session.CourseID,
		session.SessionDate, session.StartTime, session.EndTime,
		session.WifiSSID, session.WifiBSSID,
		session.LocationLatitude, session.LocationLongitude,
		session.LocationRadius, session.Status,
		session.CreatedAt, session.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("error creating session: %v", err)
	}

	return nil
}

func (r *MySQLSessionRepository) GetByID(id string) (*models.Session, error) {
	session := &models.Session{}
	query := `
		SELECT id, teacher_id, course_id, session_date, start_time, end_time,
		wifi_ssid, wifi_bssid, location_latitude, location_longitude,
		location_radius, status, created_at, updated_at
		FROM attendance_sessions WHERE id = ?
	`

	err := r.db.QueryRow(query, id).Scan(
		&session.ID, &session.TeacherID, &session.CourseID,
		&session.SessionDate, &session.StartTime, &session.EndTime,
		&session.WifiSSID, &session.WifiBSSID,
		&session.LocationLatitude, &session.LocationLongitude,
		&session.LocationRadius, &session.Status,
		&session.CreatedAt, &session.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error getting session by ID: %v", err)
	}

	return session, nil
}

func (r *MySQLSessionRepository) Update(session *models.Session) error {
	session.UpdatedAt = time.Now()

	query := `
		UPDATE attendance_sessions SET
			teacher_id = ?, course_id = ?, session_date = ?,
			start_time = ?, end_time = ?, wifi_ssid = ?,
			wifi_bssid = ?, location_latitude = ?, location_longitude = ?,
			location_radius = ?, status = ?, updated_at = ?
		WHERE id = ?
	`

	result, err := r.db.Exec(query,
		session.TeacherID, session.CourseID, session.SessionDate,
		session.StartTime, session.EndTime, session.WifiSSID,
		session.WifiBSSID, session.LocationLatitude, session.LocationLongitude,
		session.LocationRadius, session.Status, session.UpdatedAt,
		session.ID,
	)

	if err != nil {
		return fmt.Errorf("error updating session: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("session not found")
	}

	return nil
}

func (r *MySQLSessionRepository) Delete(id string) error {
	result, err := r.db.Exec("DELETE FROM attendance_sessions WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("error deleting session: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("session not found")
	}

	return nil
}

func (r *MySQLSessionRepository) List(offset, limit int) ([]*models.Session, error) {
	query := `
		SELECT id, teacher_id, course_id, session_date, start_time, end_time,
		wifi_ssid, wifi_bssid, location_latitude, location_longitude,
		location_radius, status, created_at, updated_at
		FROM attendance_sessions
		ORDER BY session_date DESC, start_time DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error listing sessions: %v", err)
	}
	defer rows.Close()

	var sessions []*models.Session
	for rows.Next() {
		session := &models.Session{}
		err := rows.Scan(
			&session.ID, &session.TeacherID, &session.CourseID,
			&session.SessionDate, &session.StartTime, &session.EndTime,
			&session.WifiSSID, &session.WifiBSSID,
			&session.LocationLatitude, &session.LocationLongitude,
			&session.LocationRadius, &session.Status,
			&session.CreatedAt, &session.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning session row: %v", err)
		}
		sessions = append(sessions, session)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating session rows: %v", err)
	}

	return sessions, nil
}

func (r *MySQLSessionRepository) GetActiveSessions() ([]*models.Session, error) {
	query := `
		SELECT id, teacher_id, course_id, session_date, start_time, end_time,
		wifi_ssid, wifi_bssid, location_latitude, location_longitude,
		location_radius, status, created_at, updated_at
		FROM attendance_sessions
		WHERE status = ? AND session_date = CURRENT_DATE
		AND start_time <= CURRENT_TIME AND end_time >= CURRENT_TIME
	`

	rows, err := r.db.Query(query, models.SessionStatusActive)
	if err != nil {
		return nil, fmt.Errorf("error getting active sessions: %v", err)
	}
	defer rows.Close()

	var sessions []*models.Session
	for rows.Next() {
		session := &models.Session{}
		err := rows.Scan(
			&session.ID, &session.TeacherID, &session.CourseID,
			&session.SessionDate, &session.StartTime, &session.EndTime,
			&session.WifiSSID, &session.WifiBSSID,
			&session.LocationLatitude, &session.LocationLongitude,
			&session.LocationRadius, &session.Status,
			&session.CreatedAt, &session.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning session row: %v", err)
		}
		sessions = append(sessions, session)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating session rows: %v", err)
	}

	return sessions, nil
}

func (r *MySQLSessionRepository) GetByTeacher(teacherID string) ([]*models.Session, error) {
	query := `
		SELECT id, teacher_id, course_id, session_date, start_time, end_time,
		wifi_ssid, wifi_bssid, location_latitude, location_longitude,
		location_radius, status, created_at, updated_at
		FROM attendance_sessions
		WHERE teacher_id = ?
		ORDER BY session_date DESC, start_time DESC
	`

	rows, err := r.db.Query(query, teacherID)
	if err != nil {
		return nil, fmt.Errorf("error getting sessions by teacher: %v", err)
	}
	defer rows.Close()

	var sessions []*models.Session
	for rows.Next() {
		session := &models.Session{}
		err := rows.Scan(
			&session.ID, &session.TeacherID, &session.CourseID,
			&session.SessionDate, &session.StartTime, &session.EndTime,
			&session.WifiSSID, &session.WifiBSSID,
			&session.LocationLatitude, &session.LocationLongitude,
			&session.LocationRadius, &session.Status,
			&session.CreatedAt, &session.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning session row: %v", err)
		}
		sessions = append(sessions, session)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating session rows: %v", err)
	}

	return sessions, nil
}

func (r *MySQLSessionRepository) GetByCourse(courseID string) ([]*models.Session, error) {
	query := `
		SELECT id, teacher_id, course_id, session_date, start_time, end_time,
		wifi_ssid, wifi_bssid, location_latitude, location_longitude,
		location_radius, status, created_at, updated_at
		FROM attendance_sessions
		WHERE course_id = ?
		ORDER BY session_date DESC, start_time DESC
	`

	rows, err := r.db.Query(query, courseID)
	if err != nil {
		return nil, fmt.Errorf("error getting sessions by course: %v", err)
	}
	defer rows.Close()

	var sessions []*models.Session
	for rows.Next() {
		session := &models.Session{}
		err := rows.Scan(
			&session.ID, &session.TeacherID, &session.CourseID,
			&session.SessionDate, &session.StartTime, &session.EndTime,
			&session.WifiSSID, &session.WifiBSSID,
			&session.LocationLatitude, &session.LocationLongitude,
			&session.LocationRadius, &session.Status,
			&session.CreatedAt, &session.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning session row: %v", err)
		}
		sessions = append(sessions, session)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating session rows: %v", err)
	}

	return sessions, nil
}

func (r *MySQLSessionRepository) EndSession(id string) error {
	query := `
		UPDATE attendance_sessions
		SET status = ?, updated_at = ?
		WHERE id = ? AND status = ?
	`

	result, err := r.db.Exec(query,
		models.SessionStatusCompleted,
		time.Now(),
		id,
		models.SessionStatusActive,
	)

	if err != nil {
		return fmt.Errorf("error ending session: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("session not found or already ended")
	}

	return nil
}

func (r *MySQLSessionRepository) CancelSession(id string) error {
	query := `
		UPDATE attendance_sessions
		SET status = ?, updated_at = ?
		WHERE id = ? AND status = ?
	`

	result, err := r.db.Exec(query,
		models.SessionStatusCancelled,
		time.Now(),
		id,
		models.SessionStatusActive,
	)

	if err != nil {
		return fmt.Errorf("error cancelling session: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("session not found or already ended/cancelled")
	}

	return nil
}

func (r *MySQLSessionRepository) BatchMarkAttendance(records []*models.AttendanceRecord) error {
	return r.WithTransaction(func(tx *sql.Tx) error {
		stmt, err := tx.Prepare(`
			INSERT INTO attendance_records (
				id, session_id, student_id, marked_at, 
				wifi_ssid, wifi_bssid, location_latitude, location_longitude,
				device_id, verification_status, created_at, updated_at
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
		`)
		if err != nil {
			return fmt.Errorf("error preparing statement: %v", err)
		}
		defer stmt.Close()

		for _, record := range records {
			_, err := stmt.Exec(
				record.ID,
				record.SessionID,
				record.StudentID,
				record.MarkedAt,
				record.WifiSSID,
				record.WifiBSSID,
				record.LocationLatitude,
				record.LocationLongitude,
				record.DeviceID,
				record.VerificationStatus,
			)
			if err != nil {
				return fmt.Errorf("error executing statement: %v", err)
			}
		}

		return nil
	})
}

func (r *MySQLSessionRepository) BatchUpdateSessions(sessions []*models.Session) error {
	return r.WithTransaction(func(tx *sql.Tx) error {
		stmt, err := tx.Prepare(`
			UPDATE attendance_sessions 
			SET status = ?, updated_at = NOW()
			WHERE id = ?
		`)
		if err != nil {
			return fmt.Errorf("error preparing statement: %v", err)
		}
		defer stmt.Close()

		for _, session := range sessions {
			_, err := stmt.Exec(
				session.Status,
				session.ID,
			)
			if err != nil {
				return fmt.Errorf("error executing statement: %v", err)
			}
		}

		return nil
	})
}

type MySQLAttendanceRepository struct {
	db *database.MySQLDB
}

func NewMySQLAttendanceRepository(db *database.MySQLDB) models.AttendanceRepository {
	return &MySQLAttendanceRepository{db: db}
}

func (r *MySQLAttendanceRepository) Create(record *models.AttendanceRecord) error {
	if record.ID == "" {
		record.ID = uuid.New().String()
	}
	now := time.Now()
	record.CreatedAt = now
	record.UpdatedAt = now

	query := `
		INSERT INTO attendance_records (
			id, session_id, student_id, marked_at,
			wifi_ssid, wifi_bssid, location_latitude, location_longitude,
			device_id, verification_status, rejection_reason,
			created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.Exec(query,
		record.ID, record.SessionID, record.StudentID, record.MarkedAt,
		record.WifiSSID, record.WifiBSSID,
		record.LocationLatitude, record.LocationLongitude,
		record.DeviceID, record.VerificationStatus, record.RejectionReason,
		record.CreatedAt, record.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("error creating attendance record: %v", err)
	}

	return nil
}

func (r *MySQLAttendanceRepository) GetByID(id string) (*models.AttendanceRecord, error) {
	record := &models.AttendanceRecord{}
	query := `
		SELECT id, session_id, student_id, marked_at,
		wifi_ssid, wifi_bssid, location_latitude, location_longitude,
		device_id, verification_status, rejection_reason,
		created_at, updated_at
		FROM attendance_records WHERE id = ?
	`

	err := r.db.QueryRow(query, id).Scan(
		&record.ID, &record.SessionID, &record.StudentID, &record.MarkedAt,
		&record.WifiSSID, &record.WifiBSSID,
		&record.LocationLatitude, &record.LocationLongitude,
		&record.DeviceID, &record.VerificationStatus, &record.RejectionReason,
		&record.CreatedAt, &record.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error getting attendance record by ID: %v", err)
	}

	return record, nil
}

func (r *MySQLAttendanceRepository) Update(record *models.AttendanceRecord) error {
	record.UpdatedAt = time.Now()

	query := `
		UPDATE attendance_records SET
			session_id = ?, student_id = ?, marked_at = ?,
			wifi_ssid = ?, wifi_bssid = ?,
			location_latitude = ?, location_longitude = ?,
			device_id = ?, verification_status = ?, rejection_reason = ?,
			updated_at = ?
		WHERE id = ?
	`

	result, err := r.db.Exec(query,
		record.SessionID, record.StudentID, record.MarkedAt,
		record.WifiSSID, record.WifiBSSID,
		record.LocationLatitude, record.LocationLongitude,
		record.DeviceID, record.VerificationStatus, record.RejectionReason,
		record.UpdatedAt, record.ID,
	)

	if err != nil {
		return fmt.Errorf("error updating attendance record: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("attendance record not found")
	}

	return nil
}

func (r *MySQLAttendanceRepository) Delete(id string) error {
	result, err := r.db.Exec("DELETE FROM attendance_records WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("error deleting attendance record: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("attendance record not found")
	}

	return nil
}

func (r *MySQLAttendanceRepository) List(offset, limit int) ([]*models.AttendanceRecord, error) {
	query := `
		SELECT id, session_id, student_id, marked_at,
		wifi_ssid, wifi_bssid, location_latitude, location_longitude,
		device_id, verification_status, rejection_reason,
		created_at, updated_at
		FROM attendance_records
		ORDER BY marked_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error listing attendance records: %v", err)
	}
	defer rows.Close()

	var records []*models.AttendanceRecord
	for rows.Next() {
		record := &models.AttendanceRecord{}
		err := rows.Scan(
			&record.ID, &record.SessionID, &record.StudentID, &record.MarkedAt,
			&record.WifiSSID, &record.WifiBSSID,
			&record.LocationLatitude, &record.LocationLongitude,
			&record.DeviceID, &record.VerificationStatus, &record.RejectionReason,
			&record.CreatedAt, &record.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning attendance record row: %v", err)
		}
		records = append(records, record)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating attendance record rows: %v", err)
	}

	return records, nil
}

func (r *MySQLAttendanceRepository) GetBySession(sessionID string) ([]*models.AttendanceRecord, error) {
	query := `
		SELECT id, session_id, student_id, marked_at,
		wifi_ssid, wifi_bssid, location_latitude, location_longitude,
		device_id, verification_status, rejection_reason,
		created_at, updated_at
		FROM attendance_records
		WHERE session_id = ?
		ORDER BY marked_at DESC
	`

	rows, err := r.db.Query(query, sessionID)
	if err != nil {
		return nil, fmt.Errorf("error getting attendance records by session: %v", err)
	}
	defer rows.Close()

	var records []*models.AttendanceRecord
	for rows.Next() {
		record := &models.AttendanceRecord{}
		err := rows.Scan(
			&record.ID, &record.SessionID, &record.StudentID, &record.MarkedAt,
			&record.WifiSSID, &record.WifiBSSID,
			&record.LocationLatitude, &record.LocationLongitude,
			&record.DeviceID, &record.VerificationStatus, &record.RejectionReason,
			&record.CreatedAt, &record.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning attendance record row: %v", err)
		}
		records = append(records, record)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating attendance record rows: %v", err)
	}

	return records, nil
}

func (r *MySQLAttendanceRepository) GetByStudent(studentID string) ([]*models.AttendanceRecord, error) {
	query := `
		SELECT id, session_id, student_id, marked_at,
		wifi_ssid, wifi_bssid, location_latitude, location_longitude,
		device_id, verification_status, rejection_reason,
		created_at, updated_at
		FROM attendance_records
		WHERE student_id = ?
		ORDER BY marked_at DESC
	`

	rows, err := r.db.Query(query, studentID)
	if err != nil {
		return nil, fmt.Errorf("error getting attendance records by student: %v", err)
	}
	defer rows.Close()

	var records []*models.AttendanceRecord
	for rows.Next() {
		record := &models.AttendanceRecord{}
		err := rows.Scan(
			&record.ID, &record.SessionID, &record.StudentID, &record.MarkedAt,
			&record.WifiSSID, &record.WifiBSSID,
			&record.LocationLatitude, &record.LocationLongitude,
			&record.DeviceID, &record.VerificationStatus, &record.RejectionReason,
			&record.CreatedAt, &record.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning attendance record row: %v", err)
		}
		records = append(records, record)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating attendance record rows: %v", err)
	}

	return records, nil
}

func (r *MySQLAttendanceRepository) VerifyAttendance(id string) error {
	query := `
		UPDATE attendance_records
		SET verification_status = ?, updated_at = ?
		WHERE id = ? AND verification_status = ?
	`

	result, err := r.db.Exec(query,
		models.VerificationStatusVerified,
		time.Now(),
		id,
		models.VerificationStatusPending,
	)

	if err != nil {
		return fmt.Errorf("error verifying attendance: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("attendance record not found or already verified/rejected")
	}

	return nil
}

func (r *MySQLAttendanceRepository) RejectAttendance(id string, reason string) error {
	query := `
		UPDATE attendance_records
		SET verification_status = ?, rejection_reason = ?, updated_at = ?
		WHERE id = ? AND verification_status = ?
	`

	result, err := r.db.Exec(query,
		models.VerificationStatusRejected,
		reason,
		time.Now(),
		id,
		models.VerificationStatusPending,
	)

	if err != nil {
		return fmt.Errorf("error rejecting attendance: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("attendance record not found or already verified/rejected")
	}

	return nil
}

func (r *MySQLAttendanceRepository) GetAttendanceStatistics(studentID string, courseID string, startDate time.Time, endDate time.Time) (*models.AttendanceStatistics, error) {
	query := `
		SELECT 
			COUNT(DISTINCT ar.id) as total_sessions,
			SUM(CASE WHEN ar.verification_status = ? THEN 1 ELSE 0 END) as verified_sessions,
			SUM(CASE WHEN ar.verification_status = ? THEN 1 ELSE 0 END) as rejected_sessions,
			SUM(CASE WHEN ar.verification_status = ? THEN 1 ELSE 0 END) as pending_sessions
		FROM attendance_records ar
		JOIN attendance_sessions s ON ar.session_id = s.id
		WHERE ar.student_id = ?
		AND s.course_id = ?
		AND s.session_date BETWEEN ? AND ?
	`

	stats := &models.AttendanceStatistics{}
	err := r.db.QueryRow(
		query,
		models.VerificationStatusVerified,
		models.VerificationStatusRejected,
		models.VerificationStatusPending,
		studentID,
		courseID,
		startDate.Format("2006-01-02"),
		endDate.Format("2006-01-02"),
	).Scan(
		&stats.TotalSessions,
		&stats.VerifiedSessions,
		&stats.RejectedSessions,
		&stats.PendingSessions,
	)

	if err != nil {
		return nil, fmt.Errorf("error getting attendance statistics: %v", err)
	}

	// Calculate attendance percentage
	if stats.TotalSessions > 0 {
		stats.AttendancePercentage = float64(stats.VerifiedSessions) / float64(stats.TotalSessions) * 100
	}

	return stats, nil
}
