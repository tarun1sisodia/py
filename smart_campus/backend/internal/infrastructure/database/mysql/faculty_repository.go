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

// FacultyRepository implements the repositories.FacultyRepository interface
type FacultyRepository struct {
	conn *Connection
}

// NewFacultyRepository creates a new MySQL faculty repository
func NewFacultyRepository(conn *Connection) repositories.FacultyRepository {
	return &FacultyRepository{conn: conn}
}

// Create creates a new faculty
func (r *FacultyRepository) Create(ctx context.Context, faculty *entities.Faculty) error {
	query := `
		INSERT INTO faculties (
			id, name, code, dean_id, description,
			status, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.conn.DB().ExecContext(ctx, query,
		faculty.ID, faculty.Name, faculty.Code,
		faculty.DeanID, faculty.Description, faculty.Status,
		faculty.CreatedAt, faculty.UpdatedAt,
	)

	if err != nil {
		if isDuplicateKeyError(err) {
			return domain.ErrAlreadyExists
		}
		return fmt.Errorf("error creating faculty: %v", err)
	}

	return nil
}

// GetByID retrieves a faculty by ID
func (r *FacultyRepository) GetByID(ctx context.Context, id string) (*entities.Faculty, error) {
	query := `
		SELECT f.id, f.name, f.code, f.dean_id, f.description,
		       f.status, f.created_at, f.updated_at,
		       CONCAT(u.first_name, ' ', u.last_name) as dean_name
		FROM faculties f
		LEFT JOIN users u ON f.dean_id = u.id
		WHERE f.id = ?
	`

	faculty := &entities.Faculty{}
	err := r.conn.DB().QueryRowContext(ctx, query, id).Scan(
		&faculty.ID, &faculty.Name, &faculty.Code,
		&faculty.DeanID, &faculty.Description, &faculty.Status,
		&faculty.CreatedAt, &faculty.UpdatedAt, &faculty.DeanName,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("error getting faculty: %v", err)
	}

	return faculty, nil
}

// GetByCode retrieves a faculty by code
func (r *FacultyRepository) GetByCode(ctx context.Context, code string) (*entities.Faculty, error) {
	query := `
		SELECT f.id, f.name, f.code, f.dean_id, f.description,
		       f.status, f.created_at, f.updated_at,
		       CONCAT(u.first_name, ' ', u.last_name) as dean_name
		FROM faculties f
		LEFT JOIN users u ON f.dean_id = u.id
		WHERE f.code = ?
	`

	faculty := &entities.Faculty{}
	err := r.conn.DB().QueryRowContext(ctx, query, code).Scan(
		&faculty.ID, &faculty.Name, &faculty.Code,
		&faculty.DeanID, &faculty.Description, &faculty.Status,
		&faculty.CreatedAt, &faculty.UpdatedAt, &faculty.DeanName,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("error getting faculty by code: %v", err)
	}

	return faculty, nil
}

// Update updates an existing faculty
func (r *FacultyRepository) Update(ctx context.Context, faculty *entities.Faculty) error {
	query := `
		UPDATE faculties SET
			name = ?, code = ?, dean_id = ?,
			description = ?, status = ?, updated_at = ?
		WHERE id = ?
	`

	result, err := r.conn.DB().ExecContext(ctx, query,
		faculty.Name, faculty.Code, faculty.DeanID,
		faculty.Description, faculty.Status, faculty.UpdatedAt,
		faculty.ID,
	)

	if err != nil {
		return fmt.Errorf("error updating faculty: %v", err)
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

// Delete deletes a faculty by ID
func (r *FacultyRepository) Delete(ctx context.Context, id string) error {
	query := "DELETE FROM faculties WHERE id = ?"

	result, err := r.conn.DB().ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error deleting faculty: %v", err)
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

// List retrieves faculties with optional filters
func (r *FacultyRepository) List(ctx context.Context, filters map[string]interface{}, page, limit int) ([]*entities.Faculty, int, error) {
	var conditions []string
	var args []interface{}

	// Build WHERE clause based on filters
	for key, value := range filters {
		conditions = append(conditions, fmt.Sprintf("f.%s = ?", key))
		args = append(args, value)
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Build query
	query := `
		SELECT f.id, f.name, f.code, f.dean_id, f.description,
		       f.status, f.created_at, f.updated_at,
		       CONCAT(u.first_name, ' ', u.last_name) as dean_name
		FROM faculties f
		LEFT JOIN users u ON f.dean_id = u.id
	`

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY f.name LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	// Execute query
	rows, err := r.conn.DB().QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("error listing faculties: %v", err)
	}
	defer rows.Close()

	var faculties []*entities.Faculty
	for rows.Next() {
		faculty := &entities.Faculty{}
		err := rows.Scan(
			&faculty.ID, &faculty.Name, &faculty.Code,
			&faculty.DeanID, &faculty.Description, &faculty.Status,
			&faculty.CreatedAt, &faculty.UpdatedAt, &faculty.DeanName,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("error scanning faculty row: %v", err)
		}
		faculties = append(faculties, faculty)
	}

	// Get total count
	countQuery := "SELECT COUNT(*) FROM faculties f"
	if len(conditions) > 0 {
		countQuery += " WHERE " + strings.Join(conditions, " AND ")
	}

	var total int
	err = r.conn.DB().QueryRowContext(ctx, countQuery, args[:len(args)-2]...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("error getting total count: %v", err)
	}

	return faculties, total, nil
}

// UpdateStatus updates a faculty's status
func (r *FacultyRepository) UpdateStatus(ctx context.Context, facultyID string, status entities.FacultyStatus) error {
	query := `
		UPDATE faculties
		SET status = ?, updated_at = ?
		WHERE id = ?
	`

	result, err := r.conn.DB().ExecContext(ctx, query, status, time.Now(), facultyID)
	if err != nil {
		return fmt.Errorf("error updating faculty status: %v", err)
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

// Search searches for faculties by name or code
func (r *FacultyRepository) Search(ctx context.Context, query string, limit int) ([]*entities.Faculty, error) {
	sqlQuery := `
		SELECT f.id, f.name, f.code, f.dean_id, f.description,
		       f.status, f.created_at, f.updated_at,
		       CONCAT(u.first_name, ' ', u.last_name) as dean_name
		FROM faculties f
		LEFT JOIN users u ON f.dean_id = u.id
		WHERE f.name LIKE ? OR f.code LIKE ?
		ORDER BY f.name
		LIMIT ?
	`

	searchPattern := "%" + query + "%"
	rows, err := r.conn.DB().QueryContext(ctx, sqlQuery,
		searchPattern, searchPattern, limit)
	if err != nil {
		return nil, fmt.Errorf("error searching faculties: %v", err)
	}
	defer rows.Close()

	var faculties []*entities.Faculty
	for rows.Next() {
		faculty := &entities.Faculty{}
		err := rows.Scan(
			&faculty.ID, &faculty.Name, &faculty.Code,
			&faculty.DeanID, &faculty.Description, &faculty.Status,
			&faculty.CreatedAt, &faculty.UpdatedAt, &faculty.DeanName,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning faculty row: %v", err)
		}
		faculties = append(faculties, faculty)
	}

	return faculties, nil
}

// GetDepartmentCount gets the number of departments in a faculty
func (r *FacultyRepository) GetDepartmentCount(ctx context.Context, facultyID string) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM departments
		WHERE faculty_id = ?
	`

	var count int
	err := r.conn.DB().QueryRowContext(ctx, query, facultyID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("error getting department count: %v", err)
	}

	return count, nil
}
