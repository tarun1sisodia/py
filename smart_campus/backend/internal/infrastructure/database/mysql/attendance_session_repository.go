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

// AttendanceSessionRepository implements the repositories.AttendanceSessionRepository interface
type AttendanceSessionRepository struct {
	conn *Connection
}

// NewAttendanceSessionRepository creates a new MySQL attendance session repository
func NewAttendanceSessionRepository(conn *Connection) repositories.AttendanceSessionRepository {
	return &AttendanceSessionRepository{conn: conn}
}

// Create creates a new attendance session
func (r *AttendanceSessionRepository) Create(ctx context.Context, session *entities.AttendanceSession) error {
	query := `
		INSERT INTO attendance_sessions (
			id, teacher_id, course_id, session_date, start_time,
			end_time, wifi_ssid, wifi_bssid, location_latitude,
			location_longitude, location_radius, status,
			created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.conn.DB().ExecContext(ctx, query,
		session.ID, session.TeacherID, session.CourseID,
		session.StartTime.Format("2006-01-02"), session.StartTime,
		session.EndTime, session.WiFiSSID, session.WiFiBSSID,
		session.LocationLat, session.LocationLong, 50, // Default radius of 50 meters
		session.Status, session.CreatedAt, session.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("error creating attendance session: %v", err)
	}

	return nil
}

// GetByID retrieves an attendance session by ID
func (r *AttendanceSessionRepository) GetByID(ctx context.Context, id string) (*entities.AttendanceSession, error) {
	query := `
		SELECT id, teacher_id, course_id, session_date, start_time,
		       end_time, wifi_ssid, wifi_bssid, location_latitude,
		       location_longitude, location_radius, status,
		       created_at, updated_at
		FROM attendance_sessions WHERE id = ?
	`

	session := &entities.AttendanceSession{}
	err := r.conn.DB().QueryRowContext(ctx, query, id).Scan(
		&session.ID, &session.TeacherID, &session.CourseID,
		&session.StartTime, &session.StartTime,
		&session.EndTime, &session.WiFiSSID, &session.WiFiBSSID,
		&session.LocationLat, &session.LocationLong, new(int),
		&session.Status, &session.CreatedAt, &session.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("error getting attendance session: %v", err)
	}

	return session, nil
}

// Update updates an existing attendance session
func (r *AttendanceSessionRepository) Update(ctx context.Context, session *entities.AttendanceSession) error {
	query := `
		UPDATE attendance_sessions SET
			end_time = ?, wifi_ssid = ?, wifi_bssid = ?,
			location_latitude = ?, location_longitude = ?,
			status = ?, updated_at = ?
		WHERE id = ?
	`

	result, err := r.conn.DB().ExecContext(ctx, query,
		session.EndTime, session.WiFiSSID, session.WiFiBSSID,
		session.LocationLat, session.LocationLong,
		session.Status, time.Now(), session.ID,
	)

	if err != nil {
		return fmt.Errorf("error updating attendance session: %v", err)
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

// Delete deletes an attendance session by ID
func (r *AttendanceSessionRepository) Delete(ctx context.Context, id string) error {
	query := "DELETE FROM attendance_sessions WHERE id = ?"

	result, err := r.conn.DB().ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error deleting attendance session: %v", err)
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

// List retrieves attendance sessions with optional filters
func (r *AttendanceSessionRepository) List(ctx context.Context, filters map[string]interface{}, page, limit int) ([]*entities.AttendanceSession, int, error) {
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
		SELECT id, teacher_id, course_id, session_date, start_time,
		       end_time, wifi_ssid, wifi_bssid, location_latitude,
		       location_longitude, location_radius, status,
		       created_at, updated_at
		FROM attendance_sessions
	`

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY session_date DESC, start_time DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	// Execute query
	rows, err := r.conn.DB().QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("error listing attendance sessions: %v", err)
	}
	defer rows.Close()

	var sessions []*entities.AttendanceSession
	for rows.Next() {
		session := &entities.AttendanceSession{}
		err := rows.Scan(
			&session.ID, &session.TeacherID, &session.CourseID,
			&session.StartTime, &session.StartTime,
			&session.EndTime, &session.WiFiSSID, &session.WiFiBSSID,
			&session.LocationLat, &session.LocationLong, new(int),
			&session.Status, &session.CreatedAt, &session.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("error scanning session row: %v", err)
		}
		sessions = append(sessions, session)
	}

	// Get total count
	countQuery := "SELECT COUNT(*) FROM attendance_sessions"
	if len(conditions) > 0 {
		countQuery += " WHERE " + strings.Join(conditions, " AND ")
	}

	var total int
	err = r.conn.DB().QueryRowContext(ctx, countQuery, args[:len(args)-2]...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("error getting total count: %v", err)
	}

	return sessions, total, nil
}

// GetActiveSessions retrieves all active sessions
func (r *AttendanceSessionRepository) GetActiveSessions(ctx context.Context) ([]*entities.AttendanceSession, error) {
	query := `
		SELECT id, teacher_id, course_id, session_date, start_time,
		       end_time, wifi_ssid, wifi_bssid, location_latitude,
		       location_longitude, location_radius, status,
		       created_at, updated_at
		FROM attendance_sessions
		WHERE status = 'active' AND end_time > ?
		ORDER BY session_date, start_time
	`

	rows, err := r.conn.DB().QueryContext(ctx, query, time.Now())
	if err != nil {
		return nil, fmt.Errorf("error getting active sessions: %v", err)
	}
	defer rows.Close()

	var sessions []*entities.AttendanceSession
	for rows.Next() {
		session := &entities.AttendanceSession{}
		err := rows.Scan(
			&session.ID, &session.TeacherID, &session.CourseID,
			&session.StartTime, &session.StartTime,
			&session.EndTime, &session.WiFiSSID, &session.WiFiBSSID,
			&session.LocationLat, &session.LocationLong, new(int),
			&session.Status, &session.CreatedAt, &session.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning session row: %v", err)
		}
		sessions = append(sessions, session)
	}

	return sessions, nil
}

// GetByTeacher retrieves sessions created by a teacher
func (r *AttendanceSessionRepository) GetByTeacher(ctx context.Context, teacherID string) ([]*entities.AttendanceSession, error) {
	query := `
		SELECT id, teacher_id, course_id, session_date, start_time,
		       end_time, wifi_ssid, wifi_bssid, location_latitude,
		       location_longitude, location_radius, status,
		       created_at, updated_at
		FROM attendance_sessions
		WHERE teacher_id = ?
		ORDER BY session_date DESC, start_time DESC
	`

	rows, err := r.conn.DB().QueryContext(ctx, query, teacherID)
	if err != nil {
		return nil, fmt.Errorf("error getting sessions by teacher: %v", err)
	}
	defer rows.Close()

	var sessions []*entities.AttendanceSession
	for rows.Next() {
		session := &entities.AttendanceSession{}
		err := rows.Scan(
			&session.ID, &session.TeacherID, &session.CourseID,
			&session.StartTime, &session.StartTime,
			&session.EndTime, &session.WiFiSSID, &session.WiFiBSSID,
			&session.LocationLat, &session.LocationLong, new(int),
			&session.Status, &session.CreatedAt, &session.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning session row: %v", err)
		}
		sessions = append(sessions, session)
	}

	return sessions, nil
}

// GetByCourse retrieves sessions for a course
func (r *AttendanceSessionRepository) GetByCourse(ctx context.Context, courseID string) ([]*entities.AttendanceSession, error) {
	query := `
		SELECT id, teacher_id, course_id, session_date, start_time,
		       end_time, wifi_ssid, wifi_bssid, location_latitude,
		       location_longitude, location_radius, status,
		       created_at, updated_at
		FROM attendance_sessions
		WHERE course_id = ?
		ORDER BY session_date DESC, start_time DESC
	`

	rows, err := r.conn.DB().QueryContext(ctx, query, courseID)
	if err != nil {
		return nil, fmt.Errorf("error getting sessions by course: %v", err)
	}
	defer rows.Close()

	var sessions []*entities.AttendanceSession
	for rows.Next() {
		session := &entities.AttendanceSession{}
		err := rows.Scan(
			&session.ID, &session.TeacherID, &session.CourseID,
			&session.StartTime, &session.StartTime,
			&session.EndTime, &session.WiFiSSID, &session.WiFiBSSID,
			&session.LocationLat, &session.LocationLong, new(int),
			&session.Status, &session.CreatedAt, &session.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning session row: %v", err)
		}
		sessions = append(sessions, session)
	}

	return sessions, nil
}

// GetByDateRange retrieves sessions within a date range
func (r *AttendanceSessionRepository) GetByDateRange(ctx context.Context, startDate, endDate time.Time) ([]*entities.AttendanceSession, error) {
	query := `
		SELECT id, teacher_id, course_id, session_date, start_time,
		       end_time, wifi_ssid, wifi_bssid, location_latitude,
		       location_longitude, location_radius, status,
		       created_at, updated_at
		FROM attendance_sessions
		WHERE session_date BETWEEN ? AND ?
		ORDER BY session_date, start_time
	`

	rows, err := r.conn.DB().QueryContext(ctx, query,
		startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	if err != nil {
		return nil, fmt.Errorf("error getting sessions by date range: %v", err)
	}
	defer rows.Close()

	var sessions []*entities.AttendanceSession
	for rows.Next() {
		session := &entities.AttendanceSession{}
		err := rows.Scan(
			&session.ID, &session.TeacherID, &session.CourseID,
			&session.StartTime, &session.StartTime,
			&session.EndTime, &session.WiFiSSID, &session.WiFiBSSID,
			&session.LocationLat, &session.LocationLong, new(int),
			&session.Status, &session.CreatedAt, &session.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning session row: %v", err)
		}
		sessions = append(sessions, session)
	}

	return sessions, nil
}

// GetActiveSessionsByStudent retrieves active sessions available for a student
func (r *AttendanceSessionRepository) GetActiveSessionsByStudent(ctx context.Context, studentID string) ([]*entities.AttendanceSession, error) {
	query := `
		SELECT DISTINCT as.id, as.teacher_id, as.course_id, as.session_date,
		       as.start_time, as.end_time, as.wifi_ssid, as.wifi_bssid,
		       as.location_latitude, as.location_longitude, as.location_radius,
		       as.status, as.created_at, as.updated_at
		FROM attendance_sessions as
		INNER JOIN courses c ON as.course_id = c.id
		INNER JOIN users u ON u.department = c.department
		WHERE u.id = ? AND u.year_of_study = c.year_of_study
		AND as.status = 'active' AND as.end_time > ?
		AND NOT EXISTS (
			SELECT 1 FROM attendance_records ar
			WHERE ar.session_id = as.id AND ar.student_id = ?
		)
		ORDER BY as.session_date, as.start_time
	`

	rows, err := r.conn.DB().QueryContext(ctx, query, studentID, time.Now(), studentID)
	if err != nil {
		return nil, fmt.Errorf("error getting active sessions for student: %v", err)
	}
	defer rows.Close()

	var sessions []*entities.AttendanceSession
	for rows.Next() {
		session := &entities.AttendanceSession{}
		err := rows.Scan(
			&session.ID, &session.TeacherID, &session.CourseID,
			&session.StartTime, &session.StartTime,
			&session.EndTime, &session.WiFiSSID, &session.WiFiBSSID,
			&session.LocationLat, &session.LocationLong, new(int),
			&session.Status, &session.CreatedAt, &session.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning session row: %v", err)
		}
		sessions = append(sessions, session)
	}

	return sessions, nil
}

// CompleteSession marks a session as completed
func (r *AttendanceSessionRepository) CompleteSession(ctx context.Context, sessionID string) error {
	query := `
		UPDATE attendance_sessions
		SET status = 'completed', updated_at = ?
		WHERE id = ? AND status = 'active'
	`

	result, err := r.conn.DB().ExecContext(ctx, query, time.Now(), sessionID)
	if err != nil {
		return fmt.Errorf("error completing session: %v", err)
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

// CancelSession marks a session as cancelled
func (r *AttendanceSessionRepository) CancelSession(ctx context.Context, sessionID string) error {
	query := `
		UPDATE attendance_sessions
		SET status = 'cancelled', updated_at = ?
		WHERE id = ? AND status = 'active'
	`

	result, err := r.conn.DB().ExecContext(ctx, query, time.Now(), sessionID)
	if err != nil {
		return fmt.Errorf("error cancelling session: %v", err)
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

// GetSessionStatistics retrieves attendance statistics for a session
func (r *AttendanceSessionRepository) GetSessionStatistics(ctx context.Context, sessionID string) (map[string]interface{}, error) {
	query := `
		SELECT 
			COUNT(DISTINCT ar.student_id) as total_marked,
			SUM(CASE WHEN ar.status = 'present' THEN 1 ELSE 0 END) as present_count,
			SUM(CASE WHEN ar.status = 'late' THEN 1 ELSE 0 END) as late_count,
			(
				SELECT COUNT(DISTINCT u.id)
				FROM users u
				INNER JOIN courses c ON c.department = u.department
				INNER JOIN attendance_sessions as2 ON as2.course_id = c.id
				WHERE u.role = 'student'
				AND u.year_of_study = c.year_of_study
				AND as2.id = ?
			) as total_students
		FROM attendance_records ar
		WHERE ar.session_id = ?
	`

	var stats struct {
		TotalMarked   int
		PresentCount  int
		LateCount     int
		TotalStudents int
	}

	err := r.conn.DB().QueryRowContext(ctx, query, sessionID, sessionID).Scan(
		&stats.TotalMarked, &stats.PresentCount,
		&stats.LateCount, &stats.TotalStudents,
	)
	if err != nil {
		return nil, fmt.Errorf("error getting session statistics: %v", err)
	}

	return map[string]interface{}{
		"total_marked":    stats.TotalMarked,
		"present_count":   stats.PresentCount,
		"late_count":      stats.LateCount,
		"absent_count":    stats.TotalStudents - stats.TotalMarked,
		"total_students":  stats.TotalStudents,
		"attendance_rate": float64(stats.PresentCount) / float64(stats.TotalStudents) * 100,
	}, nil
}
