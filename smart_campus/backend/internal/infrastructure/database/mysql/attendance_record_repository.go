package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"smart_campus/internal/domain"
	"smart_campus/internal/domain/entities"
	"smart_campus/internal/domain/repositories"
)

// AttendanceRecordRepository implements the repositories.AttendanceRecordRepository interface
type AttendanceRecordRepository struct {
	conn *Connection
}

// NewAttendanceRecordRepository creates a new MySQL attendance record repository
func NewAttendanceRecordRepository(conn *Connection) repositories.AttendanceRecordRepository {
	return &AttendanceRecordRepository{conn: conn}
}

// Create creates a new attendance record
func (r *AttendanceRecordRepository) Create(ctx context.Context, record *entities.AttendanceRecord) error {
	query := `
		INSERT INTO attendance_records (
			id, session_id, student_id, status, location_lat,
			location_long, wifi_ssid, wifi_bssid, device_id,
			verification_log, marked_at, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.conn.DB().ExecContext(ctx, query,
		record.ID, record.SessionID, record.StudentID, record.Status,
		record.LocationLat, record.LocationLong, record.WiFiSSID,
		record.WiFiBSSID, record.DeviceID, record.VerificationLog,
		record.MarkedAt, record.CreatedAt, record.UpdatedAt,
	)

	if err != nil {
		if isDuplicateKeyError(err) {
			return domain.ErrAlreadyExists
		}
		return fmt.Errorf("error creating attendance record: %v", err)
	}

	return nil
}

// GetByID retrieves an attendance record by ID
func (r *AttendanceRecordRepository) GetByID(ctx context.Context, id string) (*entities.AttendanceRecord, error) {
	query := `
		SELECT id, session_id, student_id, status, location_lat,
		       location_long, wifi_ssid, wifi_bssid, device_id,
		       verification_log, marked_at, created_at, updated_at
		FROM attendance_records WHERE id = ?
	`

	record := &entities.AttendanceRecord{}
	err := r.conn.DB().QueryRowContext(ctx, query, id).Scan(
		&record.ID, &record.SessionID, &record.StudentID, &record.Status,
		&record.LocationLat, &record.LocationLong, &record.WiFiSSID,
		&record.WiFiBSSID, &record.DeviceID, &record.VerificationLog,
		&record.MarkedAt, &record.CreatedAt, &record.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("error getting attendance record: %v", err)
	}

	return record, nil
}

// Update updates an existing attendance record
func (r *AttendanceRecordRepository) Update(ctx context.Context, record *entities.AttendanceRecord) error {
	query := `
		UPDATE attendance_records SET
			status = ?, location_lat = ?, location_long = ?,
			wifi_ssid = ?, wifi_bssid = ?, device_id = ?,
			verification_log = ?, marked_at = ?, updated_at = ?
		WHERE id = ?
	`

	result, err := r.conn.DB().ExecContext(ctx, query,
		record.Status, record.LocationLat, record.LocationLong,
		record.WiFiSSID, record.WiFiBSSID, record.DeviceID,
		record.VerificationLog, record.MarkedAt, record.UpdatedAt,
		record.ID,
	)

	if err != nil {
		return fmt.Errorf("error updating attendance record: %v", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting affected rows: %v", err)
	}

	if rows == 0 {
		return domain.ErrNotFound
	}

	return nil
}

// Delete deletes an attendance record by ID
func (r *AttendanceRecordRepository) Delete(ctx context.Context, id string) error {
	query := "DELETE FROM attendance_records WHERE id = ?"

	result, err := r.conn.DB().ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error deleting attendance record: %v", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting affected rows: %v", err)
	}

	if rows == 0 {
		return domain.ErrNotFound
	}

	return nil
}

// List retrieves attendance records with optional filters
func (r *AttendanceRecordRepository) List(ctx context.Context, filters map[string]interface{}, page, limit int) ([]*entities.AttendanceRecord, int, error) {
	var conditions []string
	var args []interface{}

	// Build WHERE clause based on filters
	for key, value := range filters {
		conditions = append(conditions, fmt.Sprintf("%s = ?", key))
		args = append(args, value)
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Build query
	query := `
		SELECT id, session_id, student_id, status, location_lat,
		       location_long, wifi_ssid, wifi_bssid, device_id,
		       verification_log, marked_at, created_at, updated_at
		FROM attendance_records
	`

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY marked_at DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	// Execute query
	rows, err := r.conn.DB().QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("error listing attendance records: %v", err)
	}
	defer rows.Close()

	var records []*entities.AttendanceRecord
	for rows.Next() {
		record := &entities.AttendanceRecord{}
		err := rows.Scan(
			&record.ID, &record.SessionID, &record.StudentID, &record.Status,
			&record.LocationLat, &record.LocationLong, &record.WiFiSSID,
			&record.WiFiBSSID, &record.DeviceID, &record.VerificationLog,
			&record.MarkedAt, &record.CreatedAt, &record.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("error scanning attendance record row: %v", err)
		}
		records = append(records, record)
	}

	// Get total count
	countQuery := "SELECT COUNT(*) FROM attendance_records"
	if len(conditions) > 0 {
		countQuery += " WHERE " + strings.Join(conditions, " AND ")
	}

	var total int
	err = r.conn.DB().QueryRowContext(ctx, countQuery, args[:len(args)-2]...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("error getting total count: %v", err)
	}

	return records, total, nil
}

// GetBySession retrieves all attendance records for a session
func (r *AttendanceRecordRepository) GetBySession(ctx context.Context, sessionID string) ([]*entities.AttendanceRecord, error) {
	query := `
		SELECT id, session_id, student_id, status, location_lat,
		       location_long, wifi_ssid, wifi_bssid, device_id,
		       verification_log, marked_at, created_at, updated_at
		FROM attendance_records
		WHERE session_id = ?
		ORDER BY marked_at DESC
	`

	rows, err := r.conn.DB().QueryContext(ctx, query, sessionID)
	if err != nil {
		return nil, fmt.Errorf("error getting session attendance records: %v", err)
	}
	defer rows.Close()

	var records []*entities.AttendanceRecord
	for rows.Next() {
		record := &entities.AttendanceRecord{}
		err := rows.Scan(
			&record.ID, &record.SessionID, &record.StudentID, &record.Status,
			&record.LocationLat, &record.LocationLong, &record.WiFiSSID,
			&record.WiFiBSSID, &record.DeviceID, &record.VerificationLog,
			&record.MarkedAt, &record.CreatedAt, &record.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning attendance record row: %v", err)
		}
		records = append(records, record)
	}

	return records, nil
}

// GetByStudent retrieves attendance records for a student
func (r *AttendanceRecordRepository) GetByStudent(ctx context.Context, studentID string) ([]*entities.AttendanceRecord, error) {
	query := `
		SELECT id, session_id, student_id, status, location_lat,
		       location_long, wifi_ssid, wifi_bssid, device_id,
		       verification_log, marked_at, created_at, updated_at
		FROM attendance_records
		WHERE student_id = ?
		ORDER BY marked_at DESC
	`

	rows, err := r.conn.DB().QueryContext(ctx, query, studentID)
	if err != nil {
		return nil, fmt.Errorf("error getting student attendance records: %v", err)
	}
	defer rows.Close()

	var records []*entities.AttendanceRecord
	for rows.Next() {
		record := &entities.AttendanceRecord{}
		err := rows.Scan(
			&record.ID, &record.SessionID, &record.StudentID, &record.Status,
			&record.LocationLat, &record.LocationLong, &record.WiFiSSID,
			&record.WiFiBSSID, &record.DeviceID, &record.VerificationLog,
			&record.MarkedAt, &record.CreatedAt, &record.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning attendance record row: %v", err)
		}
		records = append(records, record)
	}

	return records, nil
}

// GetByStudentAndDateRange retrieves attendance records for a student within a date range
func (r *AttendanceRecordRepository) GetByStudentAndDateRange(ctx context.Context, studentID string, startDate, endDate time.Time) ([]*entities.AttendanceRecord, error) {
	query := `
		SELECT ar.id, ar.session_id, ar.student_id, ar.status, ar.location_lat,
		       ar.location_long, ar.wifi_ssid, ar.wifi_bssid, ar.device_id,
		       ar.verification_log, ar.marked_at, ar.created_at, ar.updated_at
		FROM attendance_records ar
		JOIN attendance_sessions s ON ar.session_id = s.id
		WHERE ar.student_id = ? AND s.session_date BETWEEN ? AND ?
		ORDER BY ar.marked_at DESC
	`

	rows, err := r.conn.DB().QueryContext(ctx, query, studentID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("error getting student attendance records by date range: %v", err)
	}
	defer rows.Close()

	var records []*entities.AttendanceRecord
	for rows.Next() {
		record := &entities.AttendanceRecord{}
		err := rows.Scan(
			&record.ID, &record.SessionID, &record.StudentID, &record.Status,
			&record.LocationLat, &record.LocationLong, &record.WiFiSSID,
			&record.WiFiBSSID, &record.DeviceID, &record.VerificationLog,
			&record.MarkedAt, &record.CreatedAt, &record.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning attendance record row: %v", err)
		}
		records = append(records, record)
	}

	return records, nil
}

// GetByStudentAndCourse retrieves attendance records for a student in a course
func (r *AttendanceRecordRepository) GetByStudentAndCourse(ctx context.Context, studentID, courseID string) ([]*entities.AttendanceRecord, error) {
	query := `
		SELECT ar.id, ar.session_id, ar.student_id, ar.status, ar.location_lat,
		       ar.location_long, ar.wifi_ssid, ar.wifi_bssid, ar.device_id,
		       ar.verification_log, ar.marked_at, ar.created_at, ar.updated_at
		FROM attendance_records ar
		JOIN attendance_sessions s ON ar.session_id = s.id
		WHERE ar.student_id = ? AND s.course_id = ?
		ORDER BY ar.marked_at DESC
	`

	rows, err := r.conn.DB().QueryContext(ctx, query, studentID, courseID)
	if err != nil {
		return nil, fmt.Errorf("error getting student course attendance records: %v", err)
	}
	defer rows.Close()

	var records []*entities.AttendanceRecord
	for rows.Next() {
		record := &entities.AttendanceRecord{}
		err := rows.Scan(
			&record.ID, &record.SessionID, &record.StudentID, &record.Status,
			&record.LocationLat, &record.LocationLong, &record.WiFiSSID,
			&record.WiFiBSSID, &record.DeviceID, &record.VerificationLog,
			&record.MarkedAt, &record.CreatedAt, &record.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning attendance record row: %v", err)
		}
		records = append(records, record)
	}

	return records, nil
}

// GetBySessionAndStudent retrieves an attendance record for a specific session and student
func (r *AttendanceRecordRepository) GetBySessionAndStudent(ctx context.Context, sessionID, studentID string) (*entities.AttendanceRecord, error) {
	query := `
		SELECT id, session_id, student_id, status, location_lat,
		       location_long, wifi_ssid, wifi_bssid, device_id,
		       verification_log, marked_at, created_at, updated_at
		FROM attendance_records
		WHERE session_id = ? AND student_id = ?
	`

	record := &entities.AttendanceRecord{}
	err := r.conn.DB().QueryRowContext(ctx, query, sessionID, studentID).Scan(
		&record.ID, &record.SessionID, &record.StudentID, &record.Status,
		&record.LocationLat, &record.LocationLong, &record.WiFiSSID,
		&record.WiFiBSSID, &record.DeviceID, &record.VerificationLog,
		&record.MarkedAt, &record.CreatedAt, &record.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("error getting session student attendance record: %v", err)
	}

	return record, nil
}

// UpdateVerificationStatus updates the verification status of an attendance record
func (r *AttendanceRecordRepository) UpdateVerificationStatus(ctx context.Context, recordID string, status entities.AttendanceStatus) error {
	query := `
		UPDATE attendance_records
		SET status = ?, updated_at = ?
		WHERE id = ?
	`

	result, err := r.conn.DB().ExecContext(ctx, query, status, time.Now(), recordID)
	if err != nil {
		return fmt.Errorf("error updating verification status: %v", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting affected rows: %v", err)
	}

	if rows == 0 {
		return domain.ErrNotFound
	}

	return nil
}

// AddVerificationLog adds a verification log entry to an attendance record
func (r *AttendanceRecordRepository) AddVerificationLog(ctx context.Context, recordID, log string) error {
	query := `
		UPDATE attendance_records
		SET verification_log = CONCAT(COALESCE(verification_log, ''), ?, '\n'),
		    updated_at = ?
		WHERE id = ?
	`

	result, err := r.conn.DB().ExecContext(ctx, query, log, time.Now(), recordID)
	if err != nil {
		return fmt.Errorf("error adding verification log: %v", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting affected rows: %v", err)
	}

	if rows == 0 {
		return domain.ErrNotFound
	}

	return nil
}

// GetStudentAttendanceStats retrieves attendance statistics for a student
func (r *AttendanceRecordRepository) GetStudentAttendanceStats(ctx context.Context, studentID string, courseID *string) (map[string]interface{}, error) {
	var query string
	var args []interface{}

	if courseID != nil {
		query = `
			SELECT
				COUNT(*) as total_sessions,
				SUM(CASE WHEN ar.status = 'present' THEN 1 ELSE 0 END) as present_count,
				SUM(CASE WHEN ar.status = 'late' THEN 1 ELSE 0 END) as late_count,
				SUM(CASE WHEN ar.status = 'absent' THEN 1 ELSE 0 END) as absent_count
			FROM attendance_records ar
			JOIN attendance_sessions s ON ar.session_id = s.id
			WHERE ar.student_id = ? AND s.course_id = ?
		`
		args = append(args, studentID, *courseID)
	} else {
		query = `
			SELECT
				COUNT(*) as total_sessions,
				SUM(CASE WHEN status = 'present' THEN 1 ELSE 0 END) as present_count,
				SUM(CASE WHEN status = 'late' THEN 1 ELSE 0 END) as late_count,
				SUM(CASE WHEN status = 'absent' THEN 1 ELSE 0 END) as absent_count
			FROM attendance_records
			WHERE student_id = ?
		`
		args = append(args, studentID)
	}

	var totalSessions, presentCount, lateCount, absentCount int
	err := r.conn.DB().QueryRowContext(ctx, query, args...).Scan(
		&totalSessions, &presentCount, &lateCount, &absentCount,
	)
	if err != nil {
		return nil, fmt.Errorf("error getting attendance statistics: %v", err)
	}

	attendanceRate := float64(0)
	if totalSessions > 0 {
		attendanceRate = float64(presentCount+lateCount) / float64(totalSessions) * 100
	}

	return map[string]interface{}{
		"total_sessions":  totalSessions,
		"present_count":   presentCount,
		"late_count":      lateCount,
		"absent_count":    absentCount,
		"attendance_rate": attendanceRate,
	}, nil
}
