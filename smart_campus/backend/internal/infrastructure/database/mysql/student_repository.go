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

// StudentRepository implements the repositories.StudentRepository interface
type StudentRepository struct {
	conn *Connection
}

// NewStudentRepository creates a new MySQL student repository
func NewStudentRepository(conn *Connection) repositories.StudentRepository {
	return &StudentRepository{conn: conn}
}

// Create creates a new student
func (r *StudentRepository) Create(ctx context.Context, student *entities.Student) error {
	query := `
		INSERT INTO students (
			id, user_id, student_number, department_id, year,
			status, enrolled_at, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.conn.DB().ExecContext(ctx, query,
		student.ID, student.UserID, student.StudentNumber,
		student.DepartmentID, student.Year, student.Status,
		student.EnrolledAt, student.CreatedAt, student.UpdatedAt,
	)

	if err != nil {
		if isDuplicateKeyError(err) {
			return domain.ErrAlreadyExists
		}
		return fmt.Errorf("error creating student: %v", err)
	}

	return nil
}

// GetByID retrieves a student by ID
func (r *StudentRepository) GetByID(ctx context.Context, id string) (*entities.Student, error) {
	query := `
		SELECT s.id, s.user_id, s.student_number, s.department_id,
		       s.year, s.status, s.enrolled_at, s.created_at, s.updated_at,
		       u.first_name, u.last_name, u.email, d.name as department_name
		FROM students s
		JOIN users u ON s.user_id = u.id
		JOIN departments d ON s.department_id = d.id
		WHERE s.id = ?
	`

	student := &entities.Student{}
	err := r.conn.DB().QueryRowContext(ctx, query, id).Scan(
		&student.ID, &student.UserID, &student.StudentNumber,
		&student.DepartmentID, &student.Year, &student.Status,
		&student.EnrolledAt, &student.CreatedAt, &student.UpdatedAt,
		&student.FirstName, &student.LastName, &student.Email,
		&student.DepartmentName,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("error getting student: %v", err)
	}

	return student, nil
}

// GetByUserID retrieves a student by user ID
func (r *StudentRepository) GetByUserID(ctx context.Context, userID string) (*entities.Student, error) {
	query := `
		SELECT s.id, s.user_id, s.student_number, s.department_id,
		       s.year, s.status, s.enrolled_at, s.created_at, s.updated_at,
		       u.first_name, u.last_name, u.email, d.name as department_name
		FROM students s
		JOIN users u ON s.user_id = u.id
		JOIN departments d ON s.department_id = d.id
		WHERE s.user_id = ?
	`

	student := &entities.Student{}
	err := r.conn.DB().QueryRowContext(ctx, query, userID).Scan(
		&student.ID, &student.UserID, &student.StudentNumber,
		&student.DepartmentID, &student.Year, &student.Status,
		&student.EnrolledAt, &student.CreatedAt, &student.UpdatedAt,
		&student.FirstName, &student.LastName, &student.Email,
		&student.DepartmentName,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("error getting student by user ID: %v", err)
	}

	return student, nil
}

// GetByStudentNumber retrieves a student by student number
func (r *StudentRepository) GetByStudentNumber(ctx context.Context, studentNumber string) (*entities.Student, error) {
	query := `
		SELECT s.id, s.user_id, s.student_number, s.department_id,
		       s.year, s.status, s.enrolled_at, s.created_at, s.updated_at,
		       u.first_name, u.last_name, u.email, d.name as department_name
		FROM students s
		JOIN users u ON s.user_id = u.id
		JOIN departments d ON s.department_id = d.id
		WHERE s.student_number = ?
	`

	student := &entities.Student{}
	err := r.conn.DB().QueryRowContext(ctx, query, studentNumber).Scan(
		&student.ID, &student.UserID, &student.StudentNumber,
		&student.DepartmentID, &student.Year, &student.Status,
		&student.EnrolledAt, &student.CreatedAt, &student.UpdatedAt,
		&student.FirstName, &student.LastName, &student.Email,
		&student.DepartmentName,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("error getting student by student number: %v", err)
	}

	return student, nil
}

// Update updates an existing student
func (r *StudentRepository) Update(ctx context.Context, student *entities.Student) error {
	query := `
		UPDATE students SET
			department_id = ?, year = ?, status = ?,
			enrolled_at = ?, updated_at = ?
		WHERE id = ?
	`

	result, err := r.conn.DB().ExecContext(ctx, query,
		student.DepartmentID, student.Year, student.Status,
		student.EnrolledAt, student.UpdatedAt, student.ID,
	)

	if err != nil {
		return fmt.Errorf("error updating student: %v", err)
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

// Delete deletes a student by ID
func (r *StudentRepository) Delete(ctx context.Context, id string) error {
	query := "DELETE FROM students WHERE id = ?"

	result, err := r.conn.DB().ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error deleting student: %v", err)
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

// List retrieves students with optional filters
func (r *StudentRepository) List(ctx context.Context, filters map[string]interface{}, page, limit int) ([]*entities.Student, int, error) {
	var conditions []string
	var args []interface{}

	// Build WHERE clause based on filters
	for key, value := range filters {
		conditions = append(conditions, fmt.Sprintf("s.%s = ?", key))
		args = append(args, value)
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Build query
	query := `
		SELECT s.id, s.user_id, s.student_number, s.department_id,
		       s.year, s.status, s.enrolled_at, s.created_at, s.updated_at,
		       u.first_name, u.last_name, u.email, d.name as department_name
		FROM students s
		JOIN users u ON s.user_id = u.id
		JOIN departments d ON s.department_id = d.id
	`

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY s.student_number LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	// Execute query
	rows, err := r.conn.DB().QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("error listing students: %v", err)
	}
	defer rows.Close()

	var students []*entities.Student
	for rows.Next() {
		student := &entities.Student{}
		err := rows.Scan(
			&student.ID, &student.UserID, &student.StudentNumber,
			&student.DepartmentID, &student.Year, &student.Status,
			&student.EnrolledAt, &student.CreatedAt, &student.UpdatedAt,
			&student.FirstName, &student.LastName, &student.Email,
			&student.DepartmentName,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("error scanning student row: %v", err)
		}
		students = append(students, student)
	}

	// Get total count
	countQuery := `
		SELECT COUNT(*)
		FROM students s
	`
	if len(conditions) > 0 {
		countQuery += " WHERE " + strings.Join(conditions, " AND ")
	}

	var total int
	err = r.conn.DB().QueryRowContext(ctx, countQuery, args[:len(args)-2]...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("error getting total count: %v", err)
	}

	return students, total, nil
}

// GetByDepartment retrieves students by department
func (r *StudentRepository) GetByDepartment(ctx context.Context, departmentID string) ([]*entities.Student, error) {
	query := `
		SELECT s.id, s.user_id, s.student_number, s.department_id,
		       s.year, s.status, s.enrolled_at, s.created_at, s.updated_at,
		       u.first_name, u.last_name, u.email, d.name as department_name
		FROM students s
		JOIN users u ON s.user_id = u.id
		JOIN departments d ON s.department_id = d.id
		WHERE s.department_id = ?
		ORDER BY s.student_number
	`

	rows, err := r.conn.DB().QueryContext(ctx, query, departmentID)
	if err != nil {
		return nil, fmt.Errorf("error getting department students: %v", err)
	}
	defer rows.Close()

	var students []*entities.Student
	for rows.Next() {
		student := &entities.Student{}
		err := rows.Scan(
			&student.ID, &student.UserID, &student.StudentNumber,
			&student.DepartmentID, &student.Year, &student.Status,
			&student.EnrolledAt, &student.CreatedAt, &student.UpdatedAt,
			&student.FirstName, &student.LastName, &student.Email,
			&student.DepartmentName,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning student row: %v", err)
		}
		students = append(students, student)
	}

	return students, nil
}

// GetByDepartmentAndYear retrieves students by department and year
func (r *StudentRepository) GetByDepartmentAndYear(ctx context.Context, departmentID string, year int) ([]*entities.Student, error) {
	query := `
		SELECT s.id, s.user_id, s.student_number, s.department_id,
		       s.year, s.status, s.enrolled_at, s.created_at, s.updated_at,
		       u.first_name, u.last_name, u.email, d.name as department_name
		FROM students s
		JOIN users u ON s.user_id = u.id
		JOIN departments d ON s.department_id = d.id
		WHERE s.department_id = ? AND s.year = ?
		ORDER BY s.student_number
	`

	rows, err := r.conn.DB().QueryContext(ctx, query, departmentID, year)
	if err != nil {
		return nil, fmt.Errorf("error getting department year students: %v", err)
	}
	defer rows.Close()

	var students []*entities.Student
	for rows.Next() {
		student := &entities.Student{}
		err := rows.Scan(
			&student.ID, &student.UserID, &student.StudentNumber,
			&student.DepartmentID, &student.Year, &student.Status,
			&student.EnrolledAt, &student.CreatedAt, &student.UpdatedAt,
			&student.FirstName, &student.LastName, &student.Email,
			&student.DepartmentName,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning student row: %v", err)
		}
		students = append(students, student)
	}

	return students, nil
}

// GetByCourse retrieves students enrolled in a course
func (r *StudentRepository) GetByCourse(ctx context.Context, courseID string) ([]*entities.Student, error) {
	query := `
		SELECT s.id, s.user_id, s.student_number, s.department_id,
		       s.year, s.status, s.enrolled_at, s.created_at, s.updated_at,
		       u.first_name, u.last_name, u.email, d.name as department_name
		FROM students s
		JOIN users u ON s.user_id = u.id
		JOIN departments d ON s.department_id = d.id
		JOIN course_enrollments ce ON s.id = ce.student_id
		WHERE ce.course_id = ?
		ORDER BY s.student_number
	`

	rows, err := r.conn.DB().QueryContext(ctx, query, courseID)
	if err != nil {
		return nil, fmt.Errorf("error getting course students: %v", err)
	}
	defer rows.Close()

	var students []*entities.Student
	for rows.Next() {
		student := &entities.Student{}
		err := rows.Scan(
			&student.ID, &student.UserID, &student.StudentNumber,
			&student.DepartmentID, &student.Year, &student.Status,
			&student.EnrolledAt, &student.CreatedAt, &student.UpdatedAt,
			&student.FirstName, &student.LastName, &student.Email,
			&student.DepartmentName,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning student row: %v", err)
		}
		students = append(students, student)
	}

	return students, nil
}

// UpdateStatus updates a student's status
func (r *StudentRepository) UpdateStatus(ctx context.Context, studentID string, status entities.StudentStatus) error {
	query := `
		UPDATE students
		SET status = ?, updated_at = ?
		WHERE id = ?
	`

	result, err := r.conn.DB().ExecContext(ctx, query, status, time.Now(), studentID)
	if err != nil {
		return fmt.Errorf("error updating student status: %v", err)
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

// Search searches for students by name or student number
func (r *StudentRepository) Search(ctx context.Context, query string, limit int) ([]*entities.Student, error) {
	sqlQuery := `
		SELECT s.id, s.user_id, s.student_number, s.department_id,
		       s.year, s.status, s.enrolled_at, s.created_at, s.updated_at,
		       u.first_name, u.last_name, u.email, d.name as department_name
		FROM students s
		JOIN users u ON s.user_id = u.id
		JOIN departments d ON s.department_id = d.id
		WHERE s.student_number LIKE ? OR
		      u.first_name LIKE ? OR
		      u.last_name LIKE ?
		ORDER BY s.student_number
		LIMIT ?
	`

	searchPattern := "%" + query + "%"
	rows, err := r.conn.DB().QueryContext(ctx, sqlQuery,
		searchPattern, searchPattern, searchPattern, limit)
	if err != nil {
		return nil, fmt.Errorf("error searching students: %v", err)
	}
	defer rows.Close()

	var students []*entities.Student
	for rows.Next() {
		student := &entities.Student{}
		err := rows.Scan(
			&student.ID, &student.UserID, &student.StudentNumber,
			&student.DepartmentID, &student.Year, &student.Status,
			&student.EnrolledAt, &student.CreatedAt, &student.UpdatedAt,
			&student.FirstName, &student.LastName, &student.Email,
			&student.DepartmentName,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning student row: %v", err)
		}
		students = append(students, student)
	}

	return students, nil
}
