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

// TeacherRepository implements the repositories.TeacherRepository interface
type TeacherRepository struct {
	conn *Connection
}

// NewTeacherRepository creates a new MySQL teacher repository
func NewTeacherRepository(conn *Connection) repositories.TeacherRepository {
	return &TeacherRepository{conn: conn}
}

// Create creates a new teacher
func (r *TeacherRepository) Create(ctx context.Context, teacher *entities.Teacher) error {
	query := `
		INSERT INTO teachers (
			id, user_id, department_id, employee_id,
			position, status, joined_at, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.conn.DB().ExecContext(ctx, query,
		teacher.ID, teacher.UserID, teacher.DepartmentID,
		teacher.EmployeeID, teacher.Position, teacher.Status,
		teacher.JoinedAt, teacher.CreatedAt, teacher.UpdatedAt,
	)

	if err != nil {
		if isDuplicateKeyError(err) {
			return domain.ErrAlreadyExists
		}
		return fmt.Errorf("error creating teacher: %v", err)
	}

	return nil
}

// GetByID retrieves a teacher by ID
func (r *TeacherRepository) GetByID(ctx context.Context, id string) (*entities.Teacher, error) {
	query := `
		SELECT t.id, t.user_id, t.department_id, t.employee_id,
		       t.position, t.status, t.joined_at, t.created_at, t.updated_at,
		       u.first_name, u.last_name, u.email,
		       d.name as department_name
		FROM teachers t
		JOIN users u ON t.user_id = u.id
		JOIN departments d ON t.department_id = d.id
		WHERE t.id = ?
	`

	teacher := &entities.Teacher{}
	err := r.conn.DB().QueryRowContext(ctx, query, id).Scan(
		&teacher.ID, &teacher.UserID, &teacher.DepartmentID,
		&teacher.EmployeeID, &teacher.Position, &teacher.Status,
		&teacher.JoinedAt, &teacher.CreatedAt, &teacher.UpdatedAt,
		&teacher.FirstName, &teacher.LastName, &teacher.Email,
		&teacher.DepartmentName,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("error getting teacher: %v", err)
	}

	return teacher, nil
}

// GetByUserID retrieves a teacher by user ID
func (r *TeacherRepository) GetByUserID(ctx context.Context, userID string) (*entities.Teacher, error) {
	query := `
		SELECT t.id, t.user_id, t.department_id, t.employee_id,
		       t.position, t.status, t.joined_at, t.created_at, t.updated_at,
		       u.first_name, u.last_name, u.email,
		       d.name as department_name
		FROM teachers t
		JOIN users u ON t.user_id = u.id
		JOIN departments d ON t.department_id = d.id
		WHERE t.user_id = ?
	`

	teacher := &entities.Teacher{}
	err := r.conn.DB().QueryRowContext(ctx, query, userID).Scan(
		&teacher.ID, &teacher.UserID, &teacher.DepartmentID,
		&teacher.EmployeeID, &teacher.Position, &teacher.Status,
		&teacher.JoinedAt, &teacher.CreatedAt, &teacher.UpdatedAt,
		&teacher.FirstName, &teacher.LastName, &teacher.Email,
		&teacher.DepartmentName,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("error getting teacher by user ID: %v", err)
	}

	return teacher, nil
}

// GetByEmployeeID retrieves a teacher by employee ID
func (r *TeacherRepository) GetByEmployeeID(ctx context.Context, employeeID string) (*entities.Teacher, error) {
	query := `
		SELECT t.id, t.user_id, t.department_id, t.employee_id,
		       t.position, t.status, t.joined_at, t.created_at, t.updated_at,
		       u.first_name, u.last_name, u.email,
		       d.name as department_name
		FROM teachers t
		JOIN users u ON t.user_id = u.id
		JOIN departments d ON t.department_id = d.id
		WHERE t.employee_id = ?
	`

	teacher := &entities.Teacher{}
	err := r.conn.DB().QueryRowContext(ctx, query, employeeID).Scan(
		&teacher.ID, &teacher.UserID, &teacher.DepartmentID,
		&teacher.EmployeeID, &teacher.Position, &teacher.Status,
		&teacher.JoinedAt, &teacher.CreatedAt, &teacher.UpdatedAt,
		&teacher.FirstName, &teacher.LastName, &teacher.Email,
		&teacher.DepartmentName,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("error getting teacher by employee ID: %v", err)
	}

	return teacher, nil
}

// Update updates an existing teacher
func (r *TeacherRepository) Update(ctx context.Context, teacher *entities.Teacher) error {
	query := `
		UPDATE teachers SET
			department_id = ?, employee_id = ?, position = ?,
			status = ?, joined_at = ?, updated_at = ?
		WHERE id = ?
	`

	result, err := r.conn.DB().ExecContext(ctx, query,
		teacher.DepartmentID, teacher.EmployeeID, teacher.Position,
		teacher.Status, teacher.JoinedAt, teacher.UpdatedAt,
		teacher.ID,
	)

	if err != nil {
		return fmt.Errorf("error updating teacher: %v", err)
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

// Delete deletes a teacher by ID
func (r *TeacherRepository) Delete(ctx context.Context, id string) error {
	query := "DELETE FROM teachers WHERE id = ?"

	result, err := r.conn.DB().ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error deleting teacher: %v", err)
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

// List retrieves teachers with optional filters
func (r *TeacherRepository) List(ctx context.Context, filters map[string]interface{}, page, limit int) ([]*entities.Teacher, int, error) {
	var conditions []string
	var args []interface{}

	// Build WHERE clause based on filters
	for key, value := range filters {
		conditions = append(conditions, fmt.Sprintf("t.%s = ?", key))
		args = append(args, value)
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Build query
	query := `
		SELECT t.id, t.user_id, t.department_id, t.employee_id,
		       t.position, t.status, t.joined_at, t.created_at, t.updated_at,
		       u.first_name, u.last_name, u.email,
		       d.name as department_name
		FROM teachers t
		JOIN users u ON t.user_id = u.id
		JOIN departments d ON t.department_id = d.id
	`

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY t.employee_id LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	// Execute query
	rows, err := r.conn.DB().QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("error listing teachers: %v", err)
	}
	defer rows.Close()

	var teachers []*entities.Teacher
	for rows.Next() {
		teacher := &entities.Teacher{}
		err := rows.Scan(
			&teacher.ID, &teacher.UserID, &teacher.DepartmentID,
			&teacher.EmployeeID, &teacher.Position, &teacher.Status,
			&teacher.JoinedAt, &teacher.CreatedAt, &teacher.UpdatedAt,
			&teacher.FirstName, &teacher.LastName, &teacher.Email,
			&teacher.DepartmentName,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("error scanning teacher row: %v", err)
		}
		teachers = append(teachers, teacher)
	}

	// Get total count
	countQuery := "SELECT COUNT(*) FROM teachers t"
	if len(conditions) > 0 {
		countQuery += " WHERE " + strings.Join(conditions, " AND ")
	}

	var total int
	err = r.conn.DB().QueryRowContext(ctx, countQuery, args[:len(args)-2]...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("error getting total count: %v", err)
	}

	return teachers, total, nil
}

// GetByDepartment retrieves teachers by department
func (r *TeacherRepository) GetByDepartment(ctx context.Context, departmentID string) ([]*entities.Teacher, error) {
	query := `
		SELECT t.id, t.user_id, t.department_id, t.employee_id,
		       t.position, t.status, t.joined_at, t.created_at, t.updated_at,
		       u.first_name, u.last_name, u.email,
		       d.name as department_name
		FROM teachers t
		JOIN users u ON t.user_id = u.id
		JOIN departments d ON t.department_id = d.id
		WHERE t.department_id = ?
		ORDER BY t.employee_id
	`

	rows, err := r.conn.DB().QueryContext(ctx, query, departmentID)
	if err != nil {
		return nil, fmt.Errorf("error getting department teachers: %v", err)
	}
	defer rows.Close()

	var teachers []*entities.Teacher
	for rows.Next() {
		teacher := &entities.Teacher{}
		err := rows.Scan(
			&teacher.ID, &teacher.UserID, &teacher.DepartmentID,
			&teacher.EmployeeID, &teacher.Position, &teacher.Status,
			&teacher.JoinedAt, &teacher.CreatedAt, &teacher.UpdatedAt,
			&teacher.FirstName, &teacher.LastName, &teacher.Email,
			&teacher.DepartmentName,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning teacher row: %v", err)
		}
		teachers = append(teachers, teacher)
	}

	return teachers, nil
}

// UpdateStatus updates a teacher's status
func (r *TeacherRepository) UpdateStatus(ctx context.Context, teacherID string, status entities.TeacherStatus) error {
	query := `
		UPDATE teachers
		SET status = ?, updated_at = ?
		WHERE id = ?
	`

	result, err := r.conn.DB().ExecContext(ctx, query, status, time.Now(), teacherID)
	if err != nil {
		return fmt.Errorf("error updating teacher status: %v", err)
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

// Search searches for teachers by name or employee ID
func (r *TeacherRepository) Search(ctx context.Context, query string, limit int) ([]*entities.Teacher, error) {
	sqlQuery := `
		SELECT t.id, t.user_id, t.department_id, t.employee_id,
		       t.position, t.status, t.joined_at, t.created_at, t.updated_at,
		       u.first_name, u.last_name, u.email,
		       d.name as department_name
		FROM teachers t
		JOIN users u ON t.user_id = u.id
		JOIN departments d ON t.department_id = d.id
		WHERE t.employee_id LIKE ? OR
		      u.first_name LIKE ? OR
		      u.last_name LIKE ?
		ORDER BY t.employee_id
		LIMIT ?
	`

	searchPattern := "%" + query + "%"
	rows, err := r.conn.DB().QueryContext(ctx, sqlQuery,
		searchPattern, searchPattern, searchPattern, limit)
	if err != nil {
		return nil, fmt.Errorf("error searching teachers: %v", err)
	}
	defer rows.Close()

	var teachers []*entities.Teacher
	for rows.Next() {
		teacher := &entities.Teacher{}
		err := rows.Scan(
			&teacher.ID, &teacher.UserID, &teacher.DepartmentID,
			&teacher.EmployeeID, &teacher.Position, &teacher.Status,
			&teacher.JoinedAt, &teacher.CreatedAt, &teacher.UpdatedAt,
			&teacher.FirstName, &teacher.LastName, &teacher.Email,
			&teacher.DepartmentName,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning teacher row: %v", err)
		}
		teachers = append(teachers, teacher)
	}

	return teachers, nil
}

// GetCourseCount gets the number of courses assigned to a teacher
func (r *TeacherRepository) GetCourseCount(ctx context.Context, teacherID string) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM course_teachers
		WHERE teacher_id = ?
	`

	var count int
	err := r.conn.DB().QueryRowContext(ctx, query, teacherID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("error getting course count: %v", err)
	}

	return count, nil
}
