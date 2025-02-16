package mysql

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"smart_campus/internal/domain"
	"smart_campus/internal/domain/entities"
	"smart_campus/internal/domain/repositories"
)

// SystemSettingsRepository implements the repositories.SystemSettingsRepository interface
type SystemSettingsRepository struct {
	conn *Connection
}

// NewSystemSettingsRepository creates a new MySQL system settings repository
func NewSystemSettingsRepository(conn *Connection) repositories.SystemSettingsRepository {
	return &SystemSettingsRepository{conn: conn}
}

// Create creates a new system setting
func (r *SystemSettingsRepository) Create(ctx context.Context, setting *entities.SystemSetting) error {
	query := `
		INSERT INTO system_settings (
			key, value, description, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?)
	`

	_, err := r.conn.DB().ExecContext(ctx, query,
		setting.Key, setting.Value, setting.Description,
		setting.CreatedAt, setting.UpdatedAt,
	)

	if err != nil {
		if isDuplicateKeyError(err) {
			return domain.ErrAlreadyExists
		}
		return fmt.Errorf("error creating system setting: %v", err)
	}

	return nil
}

// GetByID retrieves a system setting by ID
func (r *SystemSettingsRepository) GetByID(ctx context.Context, id string) (*entities.SystemSetting, error) {
	query := `
		SELECT key, value, description, created_at, updated_at
		FROM system_settings WHERE id = ?
	`

	setting := &entities.SystemSetting{}
	err := r.conn.DB().QueryRowContext(ctx, query, id).Scan(
		&setting.Key, &setting.Value, &setting.Description,
		&setting.CreatedAt, &setting.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("error getting system setting: %v", err)
	}

	return setting, nil
}

// GetByKey retrieves a system setting by key
func (r *SystemSettingsRepository) GetByKey(ctx context.Context, key entities.SystemSettingKey) (*entities.SystemSetting, error) {
	query := `
		SELECT key, value, description, created_at, updated_at
		FROM system_settings WHERE key = ?
	`

	setting := &entities.SystemSetting{}
	err := r.conn.DB().QueryRowContext(ctx, query, key).Scan(
		&setting.Key, &setting.Value, &setting.Description,
		&setting.CreatedAt, &setting.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("error getting system setting by key: %v", err)
	}

	return setting, nil
}

// Update updates an existing system setting
func (r *SystemSettingsRepository) Update(ctx context.Context, setting *entities.SystemSetting) error {
	query := `
		UPDATE system_settings
		SET value = ?, description = ?, updated_at = ?
		WHERE key = ?
	`

	result, err := r.conn.DB().ExecContext(ctx, query,
		setting.Value, setting.Description,
		time.Now(), setting.Key,
	)

	if err != nil {
		return fmt.Errorf("error updating system setting: %v", err)
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

// Delete deletes a system setting by ID
func (r *SystemSettingsRepository) Delete(ctx context.Context, id string) error {
	query := "DELETE FROM system_settings WHERE id = ?"

	result, err := r.conn.DB().ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error deleting system setting: %v", err)
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

// List retrieves all system settings
func (r *SystemSettingsRepository) List(ctx context.Context) ([]*entities.SystemSetting, error) {
	query := `
		SELECT key, value, description, created_at, updated_at
		FROM system_settings
		ORDER BY key
	`

	rows, err := r.conn.DB().QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error listing system settings: %v", err)
	}
	defer rows.Close()

	var settings []*entities.SystemSetting
	for rows.Next() {
		setting := &entities.SystemSetting{}
		err := rows.Scan(
			&setting.Key, &setting.Value, &setting.Description,
			&setting.CreatedAt, &setting.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning system setting row: %v", err)
		}
		settings = append(settings, setting)
	}

	return settings, nil
}

// GetMultipleByKeys retrieves multiple system settings by their keys
func (r *SystemSettingsRepository) GetMultipleByKeys(ctx context.Context, keys []entities.SystemSettingKey) ([]*entities.SystemSetting, error) {
	if len(keys) == 0 {
		return []*entities.SystemSetting{}, nil
	}

	// Create placeholders for SQL IN clause
	placeholders := make([]string, len(keys))
	args := make([]interface{}, len(keys))
	for i, key := range keys {
		placeholders[i] = "?"
		args[i] = key
	}

	query := fmt.Sprintf(`
		SELECT key, value, description, created_at, updated_at
		FROM system_settings
		WHERE key IN (%s)
		ORDER BY key
	`, strings.Join(placeholders, ","))

	rows, err := r.conn.DB().QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error getting multiple system settings: %v", err)
	}
	defer rows.Close()

	var settings []*entities.SystemSetting
	for rows.Next() {
		setting := &entities.SystemSetting{}
		err := rows.Scan(
			&setting.Key, &setting.Value, &setting.Description,
			&setting.CreatedAt, &setting.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning system setting row: %v", err)
		}
		settings = append(settings, setting)
	}

	return settings, nil
}

// UpdateValue updates the value of a system setting by key
func (r *SystemSettingsRepository) UpdateValue(ctx context.Context, key entities.SystemSettingKey, value interface{}) error {
	valueJSON, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("error marshaling value: %v", err)
	}

	query := `
		UPDATE system_settings
		SET value = ?, updated_at = ?
		WHERE key = ?
	`

	result, err := r.conn.DB().ExecContext(ctx, query, string(valueJSON), time.Now(), key)
	if err != nil {
		return fmt.Errorf("error updating system setting value: %v", err)
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

// GetIntValue retrieves an integer value for a system setting
func (r *SystemSettingsRepository) GetIntValue(ctx context.Context, key entities.SystemSettingKey) (int, error) {
	setting, err := r.GetByKey(ctx, key)
	if err != nil {
		return 0, err
	}

	var value int
	err = json.Unmarshal([]byte(setting.Value), &value)
	if err != nil {
		return 0, fmt.Errorf("error unmarshaling int value: %v", err)
	}

	return value, nil
}

// GetFloatValue retrieves a float value for a system setting
func (r *SystemSettingsRepository) GetFloatValue(ctx context.Context, key entities.SystemSettingKey) (float64, error) {
	setting, err := r.GetByKey(ctx, key)
	if err != nil {
		return 0, err
	}

	var value float64
	err = json.Unmarshal([]byte(setting.Value), &value)
	if err != nil {
		return 0, fmt.Errorf("error unmarshaling float value: %v", err)
	}

	return value, nil
}

// GetBoolValue retrieves a boolean value for a system setting
func (r *SystemSettingsRepository) GetBoolValue(ctx context.Context, key entities.SystemSettingKey) (bool, error) {
	setting, err := r.GetByKey(ctx, key)
	if err != nil {
		return false, err
	}

	var value bool
	err = json.Unmarshal([]byte(setting.Value), &value)
	if err != nil {
		return false, fmt.Errorf("error unmarshaling bool value: %v", err)
	}

	return value, nil
}

// GetStringValue retrieves a string value for a system setting
func (r *SystemSettingsRepository) GetStringValue(ctx context.Context, key entities.SystemSettingKey) (string, error) {
	setting, err := r.GetByKey(ctx, key)
	if err != nil {
		return "", err
	}

	var value string
	err = json.Unmarshal([]byte(setting.Value), &value)
	if err != nil {
		return "", fmt.Errorf("error unmarshaling string value: %v", err)
	}

	return value, nil
}

// Exists checks if a system setting exists
func (r *SystemSettingsRepository) Exists(ctx context.Context, key entities.SystemSettingKey) (bool, error) {
	query := "SELECT COUNT(*) FROM system_settings WHERE key = ?"

	var count int
	err := r.conn.DB().QueryRowContext(ctx, query, key).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("error checking system setting existence: %v", err)
	}

	return count > 0, nil
}

// Helper function to check for duplicate key errors
func isDuplicateKeyError(err error) bool {
	return strings.Contains(err.Error(), "Duplicate entry") || strings.Contains(err.Error(), "unique constraint")
}
