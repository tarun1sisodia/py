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

// AuthLogRepository implements the repositories.AuthLogRepository interface
type AuthLogRepository struct {
	conn *Connection
}

// NewAuthLogRepository creates a new MySQL authentication log repository
func NewAuthLogRepository(conn *Connection) repositories.AuthLogRepository {
	return &AuthLogRepository{conn: conn}
}

// Create creates a new authentication log entry
func (r *AuthLogRepository) Create(ctx context.Context, log *entities.AuthLog) error {
	query := `
		INSERT INTO authentication_logs (
			id, user_id, type, status, device_id,
			ip_address, user_agent, location, description,
			created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.conn.DB().ExecContext(ctx, query,
		log.ID, log.UserID, log.Type, log.Status,
		log.DeviceID, log.IPAddress, log.UserAgent,
		log.Location, log.Description, log.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("error creating auth log: %v", err)
	}

	return nil
}

// GetByID retrieves an authentication log entry by ID
func (r *AuthLogRepository) GetByID(ctx context.Context, id string) (*entities.AuthLog, error) {
	query := `
		SELECT id, user_id, type, status, device_id,
		       ip_address, user_agent, location, description,
		       created_at
		FROM authentication_logs WHERE id = ?
	`

	log := &entities.AuthLog{}
	err := r.conn.DB().QueryRowContext(ctx, query, id).Scan(
		&log.ID, &log.UserID, &log.Type, &log.Status,
		&log.DeviceID, &log.IPAddress, &log.UserAgent,
		&log.Location, &log.Description, &log.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("error getting auth log: %v", err)
	}

	return log, nil
}

// List retrieves authentication logs with optional filters
func (r *AuthLogRepository) List(ctx context.Context, filters map[string]interface{}, page, limit int) ([]*entities.AuthLog, int, error) {
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
		SELECT id, user_id, type, status, device_id,
		       ip_address, user_agent, location, description,
		       created_at
		FROM authentication_logs
	`

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY created_at DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	// Execute query
	rows, err := r.conn.DB().QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("error listing auth logs: %v", err)
	}
	defer rows.Close()

	var logs []*entities.AuthLog
	for rows.Next() {
		log := &entities.AuthLog{}
		err := rows.Scan(
			&log.ID, &log.UserID, &log.Type, &log.Status,
			&log.DeviceID, &log.IPAddress, &log.UserAgent,
			&log.Location, &log.Description, &log.CreatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("error scanning auth log row: %v", err)
		}
		logs = append(logs, log)
	}

	// Get total count
	countQuery := "SELECT COUNT(*) FROM authentication_logs"
	if len(conditions) > 0 {
		countQuery += " WHERE " + strings.Join(conditions, " AND ")
	}

	var total int
	err = r.conn.DB().QueryRowContext(ctx, countQuery, args[:len(args)-2]...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("error getting total count: %v", err)
	}

	return logs, total, nil
}

// GetByUser retrieves authentication logs for a user
func (r *AuthLogRepository) GetByUser(ctx context.Context, userID string) ([]*entities.AuthLog, error) {
	query := `
		SELECT id, user_id, type, status, device_id,
		       ip_address, user_agent, location, description,
		       created_at
		FROM authentication_logs
		WHERE user_id = ?
		ORDER BY created_at DESC
	`

	rows, err := r.conn.DB().QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("error getting user auth logs: %v", err)
	}
	defer rows.Close()

	var logs []*entities.AuthLog
	for rows.Next() {
		log := &entities.AuthLog{}
		err := rows.Scan(
			&log.ID, &log.UserID, &log.Type, &log.Status,
			&log.DeviceID, &log.IPAddress, &log.UserAgent,
			&log.Location, &log.Description, &log.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning auth log row: %v", err)
		}
		logs = append(logs, log)
	}

	return logs, nil
}

// GetByUserAndType retrieves authentication logs for a user of a specific type
func (r *AuthLogRepository) GetByUserAndType(ctx context.Context, userID string, logType entities.AuthLogType) ([]*entities.AuthLog, error) {
	query := `
		SELECT id, user_id, type, status, device_id,
		       ip_address, user_agent, location, description,
		       created_at
		FROM authentication_logs
		WHERE user_id = ? AND type = ?
		ORDER BY created_at DESC
	`

	rows, err := r.conn.DB().QueryContext(ctx, query, userID, logType)
	if err != nil {
		return nil, fmt.Errorf("error getting user auth logs by type: %v", err)
	}
	defer rows.Close()

	var logs []*entities.AuthLog
	for rows.Next() {
		log := &entities.AuthLog{}
		err := rows.Scan(
			&log.ID, &log.UserID, &log.Type, &log.Status,
			&log.DeviceID, &log.IPAddress, &log.UserAgent,
			&log.Location, &log.Description, &log.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning auth log row: %v", err)
		}
		logs = append(logs, log)
	}

	return logs, nil
}

// GetByDateRange retrieves authentication logs within a date range
func (r *AuthLogRepository) GetByDateRange(ctx context.Context, startDate, endDate time.Time) ([]*entities.AuthLog, error) {
	query := `
		SELECT id, user_id, type, status, device_id,
		       ip_address, user_agent, location, description,
		       created_at
		FROM authentication_logs
		WHERE created_at BETWEEN ? AND ?
		ORDER BY created_at DESC
	`

	rows, err := r.conn.DB().QueryContext(ctx, query, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("error getting auth logs by date range: %v", err)
	}
	defer rows.Close()

	var logs []*entities.AuthLog
	for rows.Next() {
		log := &entities.AuthLog{}
		err := rows.Scan(
			&log.ID, &log.UserID, &log.Type, &log.Status,
			&log.DeviceID, &log.IPAddress, &log.UserAgent,
			&log.Location, &log.Description, &log.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning auth log row: %v", err)
		}
		logs = append(logs, log)
	}

	return logs, nil
}

// GetByUserAndDateRange retrieves authentication logs for a user within a date range
func (r *AuthLogRepository) GetByUserAndDateRange(ctx context.Context, userID string, startDate, endDate time.Time) ([]*entities.AuthLog, error) {
	query := `
		SELECT id, user_id, type, status, device_id,
		       ip_address, user_agent, location, description,
		       created_at
		FROM authentication_logs
		WHERE user_id = ? AND created_at BETWEEN ? AND ?
		ORDER BY created_at DESC
	`

	rows, err := r.conn.DB().QueryContext(ctx, query, userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("error getting user auth logs by date range: %v", err)
	}
	defer rows.Close()

	var logs []*entities.AuthLog
	for rows.Next() {
		log := &entities.AuthLog{}
		err := rows.Scan(
			&log.ID, &log.UserID, &log.Type, &log.Status,
			&log.DeviceID, &log.IPAddress, &log.UserAgent,
			&log.Location, &log.Description, &log.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning auth log row: %v", err)
		}
		logs = append(logs, log)
	}

	return logs, nil
}

// GetFailedAttempts retrieves failed authentication attempts for a user within a time window
func (r *AuthLogRepository) GetFailedAttempts(ctx context.Context, userID string, since time.Time) ([]*entities.AuthLog, error) {
	query := `
		SELECT id, user_id, type, status, device_id,
		       ip_address, user_agent, location, description,
		       created_at
		FROM authentication_logs
		WHERE user_id = ? AND status = ? AND created_at >= ?
		ORDER BY created_at DESC
	`

	rows, err := r.conn.DB().QueryContext(ctx, query, userID, entities.AuthLogStatusFailed, since)
	if err != nil {
		return nil, fmt.Errorf("error getting failed auth attempts: %v", err)
	}
	defer rows.Close()

	var logs []*entities.AuthLog
	for rows.Next() {
		log := &entities.AuthLog{}
		err := rows.Scan(
			&log.ID, &log.UserID, &log.Type, &log.Status,
			&log.DeviceID, &log.IPAddress, &log.UserAgent,
			&log.Location, &log.Description, &log.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning auth log row: %v", err)
		}
		logs = append(logs, log)
	}

	return logs, nil
}

// GetLastSuccessfulLogin retrieves the last successful login for a user
func (r *AuthLogRepository) GetLastSuccessfulLogin(ctx context.Context, userID string) (*entities.AuthLog, error) {
	query := `
		SELECT id, user_id, type, status, device_id,
		       ip_address, user_agent, location, description,
		       created_at
		FROM authentication_logs
		WHERE user_id = ? AND type = ? AND status = ?
		ORDER BY created_at DESC
		LIMIT 1
	`

	log := &entities.AuthLog{}
	err := r.conn.DB().QueryRowContext(ctx, query, userID,
		entities.AuthLogTypeLogin, entities.AuthLogStatusSuccess).Scan(
		&log.ID, &log.UserID, &log.Type, &log.Status,
		&log.DeviceID, &log.IPAddress, &log.UserAgent,
		&log.Location, &log.Description, &log.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("error getting last successful login: %v", err)
	}

	return log, nil
}

// CountFailedAttempts counts the number of failed authentication attempts for a user within a time window
func (r *AuthLogRepository) CountFailedAttempts(ctx context.Context, userID string, since time.Time) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM authentication_logs
		WHERE user_id = ? AND status = ? AND created_at >= ?
	`

	var count int
	err := r.conn.DB().QueryRowContext(ctx, query,
		userID, entities.AuthLogStatusFailed, since).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("error counting failed attempts: %v", err)
	}

	return count, nil
}

// DeleteOldLogs deletes authentication logs older than a specified date
func (r *AuthLogRepository) DeleteOldLogs(ctx context.Context, before time.Time) error {
	query := "DELETE FROM authentication_logs WHERE created_at < ?"

	_, err := r.conn.DB().ExecContext(ctx, query, before)
	if err != nil {
		return fmt.Errorf("error deleting old auth logs: %v", err)
	}

	return nil
}
