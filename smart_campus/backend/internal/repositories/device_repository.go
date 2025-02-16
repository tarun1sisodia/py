package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"smart_campus/internal/database"
	"smart_campus/internal/models"

	"github.com/google/uuid"
)

type MySQLDeviceRepository struct {
	db *database.MySQLDB
}

func NewMySQLDeviceRepository(db *database.MySQLDB) models.DeviceRepository {
	return &MySQLDeviceRepository{db: db}
}

func (r *MySQLDeviceRepository) Create(device *models.DeviceBinding) error {
	if device.ID == "" {
		device.ID = uuid.New().String()
	}
	now := time.Now()
	device.CreatedAt = now
	device.UpdatedAt = now
	device.BoundAt = now
	device.LastUsedAt = &now

	query := `
		INSERT INTO device_bindings (
			id, user_id, device_id, device_name, device_model,
			is_active, is_blacklisted, bound_at, last_used_at,
			created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.Exec(query,
		device.ID, device.UserID, device.DeviceID,
		device.DeviceName, device.DeviceModel,
		device.IsActive, device.IsBlacklisted,
		device.BoundAt, device.LastUsedAt,
		device.CreatedAt, device.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("error creating device binding: %v", err)
	}

	return nil
}

func (r *MySQLDeviceRepository) GetByID(id string) (*models.DeviceBinding, error) {
	device := &models.DeviceBinding{}
	query := `
		SELECT id, user_id, device_id, device_name, device_model,
		is_active, is_blacklisted, bound_at, last_used_at,
		created_at, updated_at
		FROM device_bindings WHERE id = ?
	`

	err := r.db.QueryRow(query, id).Scan(
		&device.ID, &device.UserID, &device.DeviceID,
		&device.DeviceName, &device.DeviceModel,
		&device.IsActive, &device.IsBlacklisted,
		&device.BoundAt, &device.LastUsedAt,
		&device.CreatedAt, &device.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error getting device binding by ID: %v", err)
	}

	return device, nil
}

func (r *MySQLDeviceRepository) GetByDeviceID(deviceID string) (*models.DeviceBinding, error) {
	device := &models.DeviceBinding{}
	query := `
		SELECT id, user_id, device_id, device_name, device_model,
		is_active, is_blacklisted, bound_at, last_used_at,
		created_at, updated_at
		FROM device_bindings WHERE device_id = ?
	`

	err := r.db.QueryRow(query, deviceID).Scan(
		&device.ID, &device.UserID, &device.DeviceID,
		&device.DeviceName, &device.DeviceModel,
		&device.IsActive, &device.IsBlacklisted,
		&device.BoundAt, &device.LastUsedAt,
		&device.CreatedAt, &device.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error getting device binding by device ID: %v", err)
	}

	return device, nil
}

func (r *MySQLDeviceRepository) Update(device *models.DeviceBinding) error {
	now := time.Now()
	device.UpdatedAt = now

	query := `
		UPDATE device_bindings SET
			user_id = ?, device_id = ?, device_name = ?,
			device_model = ?, is_active = ?, is_blacklisted = ?,
			bound_at = ?, last_used_at = ?, updated_at = ?
		WHERE id = ?
	`

	result, err := r.db.Exec(query,
		device.UserID, device.DeviceID,
		device.DeviceName, device.DeviceModel,
		device.IsActive, device.IsBlacklisted,
		device.BoundAt, device.LastUsedAt,
		device.UpdatedAt, device.ID,
	)

	if err != nil {
		return fmt.Errorf("error updating device binding: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("device binding not found")
	}

	return nil
}

func (r *MySQLDeviceRepository) Delete(id string) error {
	result, err := r.db.Exec("DELETE FROM device_bindings WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("error deleting device binding: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("device binding not found")
	}

	return nil
}

func (r *MySQLDeviceRepository) List(offset, limit int) ([]*models.DeviceBinding, error) {
	query := `
		SELECT id, user_id, device_id, device_name, device_model,
		is_active, is_blacklisted, bound_at, last_used_at,
		created_at, updated_at
		FROM device_bindings
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error listing device bindings: %v", err)
	}
	defer rows.Close()

	var devices []*models.DeviceBinding
	for rows.Next() {
		device := &models.DeviceBinding{}
		err := rows.Scan(
			&device.ID, &device.UserID, &device.DeviceID,
			&device.DeviceName, &device.DeviceModel,
			&device.IsActive, &device.IsBlacklisted,
			&device.BoundAt, &device.LastUsedAt,
			&device.CreatedAt, &device.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning device binding row: %v", err)
		}
		devices = append(devices, device)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating device binding rows: %v", err)
	}

	return devices, nil
}

func (r *MySQLDeviceRepository) GetByUserID(userID string) ([]*models.DeviceBinding, error) {
	query := `
		SELECT id, user_id, device_id, device_name, device_model,
		is_active, is_blacklisted, bound_at, last_used_at,
		created_at, updated_at
		FROM device_bindings
		WHERE user_id = ?
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("error getting device bindings by user: %v", err)
	}
	defer rows.Close()

	var devices []*models.DeviceBinding
	for rows.Next() {
		device := &models.DeviceBinding{}
		err := rows.Scan(
			&device.ID, &device.UserID, &device.DeviceID,
			&device.DeviceName, &device.DeviceModel,
			&device.IsActive, &device.IsBlacklisted,
			&device.BoundAt, &device.LastUsedAt,
			&device.CreatedAt, &device.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning device binding row: %v", err)
		}
		devices = append(devices, device)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating device binding rows: %v", err)
	}

	return devices, nil
}

func (r *MySQLDeviceRepository) Deactivate(id string) error {
	query := `
		UPDATE device_bindings
		SET is_active = ?, updated_at = ?
		WHERE id = ?
	`

	result, err := r.db.Exec(query, false, time.Now(), id)
	if err != nil {
		return fmt.Errorf("error deactivating device binding: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("device binding not found")
	}

	return nil
}

func (r *MySQLDeviceRepository) Blacklist(id string) error {
	query := `
		UPDATE device_bindings
		SET is_blacklisted = ?, is_active = ?, updated_at = ?
		WHERE id = ?
	`

	result, err := r.db.Exec(query, true, false, time.Now(), id)
	if err != nil {
		return fmt.Errorf("error blacklisting device binding: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("device binding not found")
	}

	return nil
}

func (r *MySQLDeviceRepository) UpdateLastUsed(id string) error {
	query := `
		UPDATE device_bindings
		SET last_used_at = ?, updated_at = ?
		WHERE id = ?
	`

	now := time.Now()
	result, err := r.db.Exec(query, &now, now, id)
	if err != nil {
		return fmt.Errorf("error updating device binding last used time: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("device binding not found")
	}

	return nil
}

func (r *MySQLDeviceRepository) IsDeviceRegistered(deviceID string, userID string) (bool, error) {
	var exists bool
	query := `
		SELECT EXISTS(
			SELECT 1 FROM device_bindings
			WHERE device_id = ? AND user_id = ? AND is_active = true AND is_blacklisted = false
		)
	`

	err := r.db.QueryRow(query, deviceID, userID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error checking device registration: %v", err)
	}

	return exists, nil
}

func (r *MySQLDeviceRepository) IsDeviceBlacklisted(deviceID string) (bool, error) {
	var isBlacklisted bool
	query := `
		SELECT EXISTS(
			SELECT 1 FROM device_bindings
			WHERE device_id = ? AND is_blacklisted = true
		)
	`

	err := r.db.QueryRow(query, deviceID).Scan(&isBlacklisted)
	if err != nil {
		return false, fmt.Errorf("error checking device blacklist status: %v", err)
	}

	return isBlacklisted, nil
}

func (r *MySQLDeviceRepository) RemoveFromBlacklist(deviceID string) error {
	query := `
		UPDATE device_bindings
		SET is_blacklisted = false, updated_at = ?
		WHERE device_id = ? AND is_blacklisted = true
	`

	now := time.Now()
	result, err := r.db.Exec(query, now, deviceID)
	if err != nil {
		return fmt.Errorf("error removing device from blacklist: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("device not found or not blacklisted")
	}

	return nil
}
