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

// DepartmentRepository implements the repositories.DepartmentRepository interface
type DepartmentRepository struct {
	conn *Connection
}

// NewDepartmentRepository creates a new MySQL department repository
func NewDepartmentRepository(conn *Connection) repositories.DepartmentRepository {
	return &DepartmentRepository{conn: conn}
}

// Create creates a new department
func (r *DepartmentRepository) Create(ctx context.Context, department *entities.Department) error {
	query := `
		INSERT INTO departments (
			id, name, code, faculty_id, head_id,
			description, status, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.conn.DB().ExecContext(ctx, query,
		department.ID, department.Name, department.Code,
		department.FacultyID, department.HeadID, department.Description,
		department.Status, department.CreatedAt, department.UpdatedAt,
	)

	if err != nil {
		if isDuplicateKeyError(err) {
			return domain.ErrAlreadyExists
		}
		return fmt.Errorf("error creating department: %v", err)
	}

	return nil
}

// GetByID retrieves a department by ID
func (r *DepartmentRepository) GetByID(ctx context.Context, id string) (*entities.Department, error) {
	query := `
		SELECT d.id, d.name, d.code, d.faculty_id, d.head_id,
		       d.description, d.status, d.created_at, d.updated_at,
		       f.name as faculty_name,
		       CONCAT(u.first_name, ' ', u.last_name) as head_name
		FROM departments d
		LEFT JOIN faculties f ON d.faculty_id = f.id
		LEFT JOIN users u ON d.head_id = u.id
		WHERE d.id = ?
	`

	department := &entities.Department{}
	err := r.conn.DB().QueryRowContext(ctx, query, id).Scan(
		&department.ID, &department.Name, &department.Code,
		&department.FacultyID, &department.HeadID, &department.Description,
		&department.Status, &department.CreatedAt, &department.UpdatedAt,
		&department.FacultyName, &department.HeadName,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("error getting department: %v", err)
	}

	return department, nil
}

// GetByCode retrieves a department by code
func (r *DepartmentRepository) GetByCode(ctx context.Context, code string) (*entities.Department, error) {
	query := `
		SELECT d.id, d.name, d.code, d.faculty_id, d.head_id,
		       d.description, d.status, d.created_at, d.updated_at,
		       f.name as faculty_name,
		       CONCAT(u.first_name, ' ', u.last_name) as head_name
		FROM departments d
		LEFT JOIN faculties f ON d.faculty_id = f.id
		LEFT JOIN users u ON d.head_id = u.id
		WHERE d.code = ?
	`

	department := &entities.Department{}
	err := r.conn.DB().QueryRowContext(ctx, query, code).Scan(
		&department.ID, &department.Name, &department.Code,
		&department.FacultyID, &department.HeadID, &department.Description,
		&department.Status, &department.CreatedAt, &department.UpdatedAt,
		&department.FacultyName, &department.HeadName,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("error getting department by code: %v", err)
	}

	return department, nil
}

// Update updates an existing department
func (r *DepartmentRepository) Update(ctx context.Context, department *entities.Department) error {
	query := `
		UPDATE departments SET
			name = ?, code = ?, faculty_id = ?, head_id = ?,
			description = ?, status = ?, updated_at = ?
		WHERE id = ?
	`

	result, err := r.conn.DB().ExecContext(ctx, query,
		department.Name, department.Code, department.FacultyID,
		department.HeadID, department.Description, department.Status,
		department.UpdatedAt, department.ID,
	)

	if err != nil {
		return fmt.Errorf("error updating department: %v", err)
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

// Delete deletes a department by ID
func (r *DepartmentRepository) Delete(ctx context.Context, id string) error {
	query := "DELETE FROM departments WHERE id = ?"

	result, err := r.conn.DB().ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error deleting department: %v", err)
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

// List retrieves departments with optional filters
func (r *DepartmentRepository) List(ctx context.Context, filters map[string]interface{}, page, limit int) ([]*entities.Department, int, error) {
	var conditions []string
	var args []interface{}

	// Build WHERE clause based on filters
	for key, value := range filters {
		conditions = append(conditions, fmt.Sprintf("d.%s = ?", key))
		args = append(args, value)
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Build query
	query := `
		SELECT d.id, d.name, d.code, d.faculty_id, d.head_id,
		       d.description, d.status, d.created_at, d.updated_at,
		       f.name as faculty_name,
		       CONCAT(u.first_name, ' ', u.last_name) as head_name
		FROM departments d
		LEFT JOIN faculties f ON d.faculty_id = f.id
		LEFT JOIN users u ON d.head_id = u.id
	`

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY d.name LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	// Execute query
	rows, err := r.conn.DB().QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("error listing departments: %v", err)
	}
	defer rows.Close()

	var departments []*entities.Department
	for rows.Next() {
		department := &entities.Department{}
		err := rows.Scan(
			&department.ID, &department.Name, &department.Code,
			&department.FacultyID, &department.HeadID, &department.Description,
			&department.Status, &department.CreatedAt, &department.UpdatedAt,
			&department.FacultyName, &department.HeadName,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("error scanning department row: %v", err)
		}
		departments = append(departments, department)
	}

	// Get total count
	countQuery := "SELECT COUNT(*) FROM departments d"
	if len(conditions) > 0 {
		countQuery += " WHERE " + strings.Join(conditions, " AND ")
	}

	var total int
	err = r.conn.DB().QueryRowContext(ctx, countQuery, args[:len(args)-2]...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("error getting total count: %v", err)
	}

	return departments, total, nil
}

// GetByFaculty retrieves departments by faculty
func (r *DepartmentRepository) GetByFaculty(ctx context.Context, facultyID string) ([]*entities.Department, error) {
	query := `
		SELECT d.id, d.name, d.code, d.faculty_id, d.head_id,
		       d.description, d.status, d.created_at, d.updated_at,
		       f.name as faculty_name,
		       CONCAT(u.first_name, ' ', u.last_name) as head_name
		FROM departments d
		LEFT JOIN faculties f ON d.faculty_id = f.id
		LEFT JOIN users u ON d.head_id = u.id
		WHERE d.faculty_id = ?
		ORDER BY d.name
	`

	rows, err := r.conn.DB().QueryContext(ctx, query, facultyID)
	if err != nil {
		return nil, fmt.Errorf("error getting faculty departments: %v", err)
	}
	defer rows.Close()

	var departments []*entities.Department
	for rows.Next() {
		department := &entities.Department{}
		err := rows.Scan(
			&department.ID, &department.Name, &department.Code,
			&department.FacultyID, &department.HeadID, &department.Description,
			&department.Status, &department.CreatedAt, &department.UpdatedAt,
			&department.FacultyName, &department.HeadName,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning department row: %v", err)
		}
		departments = append(departments, department)
	}

	return departments, nil
}

// UpdateStatus updates a department's status
func (r *DepartmentRepository) UpdateStatus(ctx context.Context, departmentID string, status entities.DepartmentStatus) error {
	query := `
		UPDATE departments
		SET status = ?, updated_at = ?
		WHERE id = ?
	`

	result, err := r.conn.DB().ExecContext(ctx, query, status, time.Now(), departmentID)
	if err != nil {
		return fmt.Errorf("error updating department status: %v", err)
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

// Search searches for departments by name or code
func (r *DepartmentRepository) Search(ctx context.Context, query string, limit int) ([]*entities.Department, error) {
	sqlQuery := `
		SELECT d.id, d.name, d.code, d.faculty_id, d.head_id,
		       d.description, d.status, d.created_at, d.updated_at,
		       f.name as faculty_name,
		       CONCAT(u.first_name, ' ', u.last_name) as head_name
		FROM departments d
		LEFT JOIN faculties f ON d.faculty_id = f.id
		LEFT JOIN users u ON d.head_id = u.id
		WHERE d.name LIKE ? OR d.code LIKE ?
		ORDER BY d.name
		LIMIT ?
	`

	searchPattern := "%" + query + "%"
	rows, err := r.conn.DB().QueryContext(ctx, sqlQuery,
		searchPattern, searchPattern, limit)
	if err != nil {
		return nil, fmt.Errorf("error searching departments: %v", err)
	}
	defer rows.Close()

	var departments []*entities.Department
	for rows.Next() {
		department := &entities.Department{}
		err := rows.Scan(
			&department.ID, &department.Name, &department.Code,
			&department.FacultyID, &department.HeadID, &department.Description,
			&department.Status, &department.CreatedAt, &department.UpdatedAt,
			&department.FacultyName, &department.HeadName,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning department row: %v", err)
		}
		departments = append(departments, department)
	}

	return departments, nil
}
