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

// DeviceBindingRepository implements the repositories.DeviceBindingRepository interface
type DeviceBindingRepository struct {
	conn *Connection
}

// NewDeviceBindingRepository creates a new MySQL device binding repository
func NewDeviceBindingRepository(conn *Connection) repositories.DeviceBindingRepository {
	return &DeviceBindingRepository{conn: conn}
}

// Create creates a new device binding
func (r *DeviceBindingRepository) Create(ctx context.Context, binding *entities.DeviceBinding) error {
	query := `
		INSERT INTO device_bindings (
			id, user_id, device_id, device_name, device_model,
			device_os, device_version, status, last_used_at,
			created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.conn.DB().ExecContext(ctx, query,
		binding.ID, binding.UserID, binding.DeviceID, binding.DeviceName,
		binding.DeviceModel, binding.DeviceOS, binding.DeviceVersion,
		binding.Status, binding.LastUsedAt, binding.CreatedAt, binding.UpdatedAt,
	)

	if err != nil {
		if isDuplicateKeyError(err) {
			return domain.ErrAlreadyExists
		}
		return fmt.Errorf("error creating device binding: %v", err)
	}

	return nil
}

// GetByID retrieves a device binding by ID
func (r *DeviceBindingRepository) GetByID(ctx context.Context, id string) (*entities.DeviceBinding, error) {
	query := `
		SELECT id, user_id, device_id, device_name, device_model,
		       device_os, device_version, status, last_used_at,
		       created_at, updated_at
		FROM device_bindings WHERE id = ?
	`

	binding := &entities.DeviceBinding{}
	err := r.conn.DB().QueryRowContext(ctx, query, id).Scan(
		&binding.ID, &binding.UserID, &binding.DeviceID, &binding.DeviceName,
		&binding.DeviceModel, &binding.DeviceOS, &binding.DeviceVersion,
		&binding.Status, &binding.LastUsedAt, &binding.CreatedAt, &binding.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("error getting device binding: %v", err)
	}

	return binding, nil
}

// Update updates an existing device binding
func (r *DeviceBindingRepository) Update(ctx context.Context, binding *entities.DeviceBinding) error {
	query := `
		UPDATE device_bindings SET
			device_name = ?, device_model = ?, device_os = ?,
			device_version = ?, status = ?, last_used_at = ?,
			updated_at = ?
		WHERE id = ?
	`

	result, err := r.conn.DB().ExecContext(ctx, query,
		binding.DeviceName, binding.DeviceModel, binding.DeviceOS,
		binding.DeviceVersion, binding.Status, binding.LastUsedAt,
		binding.UpdatedAt, binding.ID,
	)

	if err != nil {
		return fmt.Errorf("error updating device binding: %v", err)
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

// Delete deletes a device binding by ID
func (r *DeviceBindingRepository) Delete(ctx context.Context, id string) error {
	query := "DELETE FROM device_bindings WHERE id = ?"

	result, err := r.conn.DB().ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error deleting device binding: %v", err)
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

// List retrieves device bindings with optional filters
func (r *DeviceBindingRepository) List(ctx context.Context, filters map[string]interface{}, page, limit int) ([]*entities.DeviceBinding, int, error) {
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
		SELECT id, user_id, device_id, device_name, device_model,
		       device_os, device_version, status, last_used_at,
		       created_at, updated_at
		FROM device_bindings
	`

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY last_used_at DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	// Execute query
	rows, err := r.conn.DB().QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("error listing device bindings: %v", err)
	}
	defer rows.Close()

	var bindings []*entities.DeviceBinding
	for rows.Next() {
		binding := &entities.DeviceBinding{}
		err := rows.Scan(
			&binding.ID, &binding.UserID, &binding.DeviceID, &binding.DeviceName,
			&binding.DeviceModel, &binding.DeviceOS, &binding.DeviceVersion,
			&binding.Status, &binding.LastUsedAt, &binding.CreatedAt, &binding.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("error scanning device binding row: %v", err)
		}
		bindings = append(bindings, binding)
	}

	// Get total count
	countQuery := "SELECT COUNT(*) FROM device_bindings"
	if len(conditions) > 0 {
		countQuery += " WHERE " + strings.Join(conditions, " AND ")
	}

	var total int
	err = r.conn.DB().QueryRowContext(ctx, countQuery, args[:len(args)-2]...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("error getting total count: %v", err)
	}

	return bindings, total, nil
}

// GetByUser retrieves all device bindings for a user
func (r *DeviceBindingRepository) GetByUser(ctx context.Context, userID string) ([]*entities.DeviceBinding, error) {
	query := `
		SELECT id, user_id, device_id, device_name, device_model,
		       device_os, device_version, status, last_used_at,
		       created_at, updated_at
		FROM device_bindings
		WHERE user_id = ?
		ORDER BY last_used_at DESC
	`

	rows, err := r.conn.DB().QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("error getting user device bindings: %v", err)
	}
	defer rows.Close()

	var bindings []*entities.DeviceBinding
	for rows.Next() {
		binding := &entities.DeviceBinding{}
		err := rows.Scan(
			&binding.ID, &binding.UserID, &binding.DeviceID, &binding.DeviceName,
			&binding.DeviceModel, &binding.DeviceOS, &binding.DeviceVersion,
			&binding.Status, &binding.LastUsedAt, &binding.CreatedAt, &binding.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning device binding row: %v", err)
		}
		bindings = append(bindings, binding)
	}

	return bindings, nil
}

// GetByDeviceID retrieves a device binding by device ID
func (r *DeviceBindingRepository) GetByDeviceID(ctx context.Context, deviceID string) (*entities.DeviceBinding, error) {
	query := `
		SELECT id, user_id, device_id, device_name, device_model,
		       device_os, device_version, status, last_used_at,
		       created_at, updated_at
		FROM device_bindings
		WHERE device_id = ? AND status = 'active'
	`

	binding := &entities.DeviceBinding{}
	err := r.conn.DB().QueryRowContext(ctx, query, deviceID).Scan(
		&binding.ID, &binding.UserID, &binding.DeviceID, &binding.DeviceName,
		&binding.DeviceModel, &binding.DeviceOS, &binding.DeviceVersion,
		&binding.Status, &binding.LastUsedAt, &binding.CreatedAt, &binding.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("error getting device binding by device ID: %v", err)
	}

	return binding, nil
}

// GetByUserAndDevice retrieves a device binding by user ID and device ID
func (r *DeviceBindingRepository) GetByUserAndDevice(ctx context.Context, userID, deviceID string) (*entities.DeviceBinding, error) {
	query := `
		SELECT id, user_id, device_id, device_name, device_model,
		       device_os, device_version, status, last_used_at,
		       created_at, updated_at
		FROM device_bindings
		WHERE user_id = ? AND device_id = ?
	`

	binding := &entities.DeviceBinding{}
	err := r.conn.DB().QueryRowContext(ctx, query, userID, deviceID).Scan(
		&binding.ID, &binding.UserID, &binding.DeviceID, &binding.DeviceName,
		&binding.DeviceModel, &binding.DeviceOS, &binding.DeviceVersion,
		&binding.Status, &binding.LastUsedAt, &binding.CreatedAt, &binding.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("error getting device binding by user and device: %v", err)
	}

	return binding, nil
}

// GetActiveBindings retrieves all active device bindings for a user
func (r *DeviceBindingRepository) GetActiveBindings(ctx context.Context, userID string) ([]*entities.DeviceBinding, error) {
	query := `
		SELECT id, user_id, device_id, device_name, device_model,
		       device_os, device_version, status, last_used_at,
		       created_at, updated_at
		FROM device_bindings
		WHERE user_id = ? AND status = 'active'
		ORDER BY last_used_at DESC
	`

	rows, err := r.conn.DB().QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("error getting active device bindings: %v", err)
	}
	defer rows.Close()

	var bindings []*entities.DeviceBinding
	for rows.Next() {
		binding := &entities.DeviceBinding{}
		err := rows.Scan(
			&binding.ID, &binding.UserID, &binding.DeviceID, &binding.DeviceName,
			&binding.DeviceModel, &binding.DeviceOS, &binding.DeviceVersion,
			&binding.Status, &binding.LastUsedAt, &binding.CreatedAt, &binding.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning device binding row: %v", err)
		}
		bindings = append(bindings, binding)
	}

	return bindings, nil
}

// RevokeBinding revokes a device binding
func (r *DeviceBindingRepository) RevokeBinding(ctx context.Context, id string) error {
	query := `
		UPDATE device_bindings
		SET status = 'revoked', updated_at = ?
		WHERE id = ?
	`

	result, err := r.conn.DB().ExecContext(ctx, query, time.Now(), id)
	if err != nil {
		return fmt.Errorf("error revoking device binding: %v", err)
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

// RevokeAllUserBindings revokes all device bindings for a user
func (r *DeviceBindingRepository) RevokeAllUserBindings(ctx context.Context, userID string) error {
	query := `
		UPDATE device_bindings
		SET status = 'revoked', updated_at = ?
		WHERE user_id = ? AND status = 'active'
	`

	_, err := r.conn.DB().ExecContext(ctx, query, time.Now(), userID)
	if err != nil {
		return fmt.Errorf("error revoking all user device bindings: %v", err)
	}

	return nil
}

// UpdateLastUsed updates the last used timestamp for a device binding
func (r *DeviceBindingRepository) UpdateLastUsed(ctx context.Context, id string) error {
	query := `
		UPDATE device_bindings
		SET last_used_at = ?, updated_at = ?
		WHERE id = ?
	`

	now := time.Now()
	result, err := r.conn.DB().ExecContext(ctx, query, now, now, id)
	if err != nil {
		return fmt.Errorf("error updating last used timestamp: %v", err)
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

// CountActiveBindings counts the number of active device bindings for a user
func (r *DeviceBindingRepository) CountActiveBindings(ctx context.Context, userID string) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM device_bindings
		WHERE user_id = ? AND status = 'active'
	`

	var count int
	err := r.conn.DB().QueryRowContext(ctx, query, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("error counting active device bindings: %v", err)
	}

	return count, nil
}
