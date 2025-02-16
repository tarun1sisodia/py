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

// CourseRepository implements the repositories.CourseRepository interface
type CourseRepository struct {
	conn *Connection
}

// NewCourseRepository creates a new MySQL course repository
func NewCourseRepository(conn *Connection) repositories.CourseRepository {
	return &CourseRepository{conn: conn}
}

// Create creates a new course
func (r *CourseRepository) Create(ctx context.Context, course *entities.Course) error {
	query := `
		INSERT INTO courses (
			id, course_code, course_name, department,
			year_of_study, semester, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.conn.DB().ExecContext(ctx, query,
		course.ID, course.CourseCode, course.CourseName, course.Department,
		course.YearOfStudy, course.Semester, course.CreatedAt, course.UpdatedAt,
	)

	if err != nil {
		if isDuplicateKeyError(err) {
			return domain.ErrAlreadyExists
		}
		return fmt.Errorf("error creating course: %v", err)
	}

	return nil
}

// GetByID retrieves a course by ID
func (r *CourseRepository) GetByID(ctx context.Context, id string) (*entities.Course, error) {
	query := `
		SELECT id, course_code, course_name, department,
		       year_of_study, semester, created_at, updated_at
		FROM courses WHERE id = ?
	`

	course := &entities.Course{}
	err := r.conn.DB().QueryRowContext(ctx, query, id).Scan(
		&course.ID, &course.CourseCode, &course.CourseName, &course.Department,
		&course.YearOfStudy, &course.Semester, &course.CreatedAt, &course.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("error getting course: %v", err)
	}

	return course, nil
}

// GetByCourseCode retrieves a course by course code
func (r *CourseRepository) GetByCourseCode(ctx context.Context, courseCode string) (*entities.Course, error) {
	query := `
		SELECT id, course_code, course_name, department,
		       year_of_study, semester, created_at, updated_at
		FROM courses WHERE course_code = ?
	`

	course := &entities.Course{}
	err := r.conn.DB().QueryRowContext(ctx, query, courseCode).Scan(
		&course.ID, &course.CourseCode, &course.CourseName, &course.Department,
		&course.YearOfStudy, &course.Semester, &course.CreatedAt, &course.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("error getting course by code: %v", err)
	}

	return course, nil
}

// Update updates an existing course
func (r *CourseRepository) Update(ctx context.Context, course *entities.Course) error {
	query := `
		UPDATE courses SET
			course_name = ?, department = ?, year_of_study = ?,
			semester = ?, updated_at = ?
		WHERE id = ?
	`

	result, err := r.conn.DB().ExecContext(ctx, query,
		course.CourseName, course.Department, course.YearOfStudy,
		course.Semester, time.Now(), course.ID,
	)

	if err != nil {
		return fmt.Errorf("error updating course: %v", err)
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

// Delete deletes a course by ID
func (r *CourseRepository) Delete(ctx context.Context, id string) error {
	query := "DELETE FROM courses WHERE id = ?"

	result, err := r.conn.DB().ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error deleting course: %v", err)
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

// List retrieves courses with optional filters
func (r *CourseRepository) List(ctx context.Context, filters map[string]interface{}, page, limit int) ([]*entities.Course, int, error) {
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
		SELECT id, course_code, course_name, department,
		       year_of_study, semester, created_at, updated_at
		FROM courses
	`

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += " LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	// Execute query
	rows, err := r.conn.DB().QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("error listing courses: %v", err)
	}
	defer rows.Close()

	var courses []*entities.Course
	for rows.Next() {
		course := &entities.Course{}
		err := rows.Scan(
			&course.ID, &course.CourseCode, &course.CourseName, &course.Department,
			&course.YearOfStudy, &course.Semester, &course.CreatedAt, &course.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("error scanning course row: %v", err)
		}
		courses = append(courses, course)
	}

	// Get total count
	countQuery := "SELECT COUNT(*) FROM courses"
	if len(conditions) > 0 {
		countQuery += " WHERE " + strings.Join(conditions, " AND ")
	}

	var total int
	err = r.conn.DB().QueryRowContext(ctx, countQuery, args[:len(args)-2]...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("error getting total count: %v", err)
	}

	return courses, total, nil
}

// GetByDepartment retrieves courses by department
func (r *CourseRepository) GetByDepartment(ctx context.Context, department string) ([]*entities.Course, error) {
	query := `
		SELECT id, course_code, course_name, department,
		       year_of_study, semester, created_at, updated_at
		FROM courses
		WHERE department = ?
		ORDER BY year_of_study, semester
	`

	rows, err := r.conn.DB().QueryContext(ctx, query, department)
	if err != nil {
		return nil, fmt.Errorf("error getting courses by department: %v", err)
	}
	defer rows.Close()

	var courses []*entities.Course
	for rows.Next() {
		course := &entities.Course{}
		err := rows.Scan(
			&course.ID, &course.CourseCode, &course.CourseName, &course.Department,
			&course.YearOfStudy, &course.Semester, &course.CreatedAt, &course.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning course row: %v", err)
		}
		courses = append(courses, course)
	}

	return courses, nil
}

// GetByDepartmentAndYear retrieves courses by department and year
func (r *CourseRepository) GetByDepartmentAndYear(ctx context.Context, department string, year int) ([]*entities.Course, error) {
	query := `
		SELECT id, course_code, course_name, department,
		       year_of_study, semester, created_at, updated_at
		FROM courses
		WHERE department = ? AND year_of_study = ?
		ORDER BY semester
	`

	rows, err := r.conn.DB().QueryContext(ctx, query, department, year)
	if err != nil {
		return nil, fmt.Errorf("error getting courses by department and year: %v", err)
	}
	defer rows.Close()

	var courses []*entities.Course
	for rows.Next() {
		course := &entities.Course{}
		err := rows.Scan(
			&course.ID, &course.CourseCode, &course.CourseName, &course.Department,
			&course.YearOfStudy, &course.Semester, &course.CreatedAt, &course.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning course row: %v", err)
		}
		courses = append(courses, course)
	}

	return courses, nil
}

// GetByTeacher retrieves courses assigned to a teacher
func (r *CourseRepository) GetByTeacher(ctx context.Context, teacherID string) ([]*entities.Course, error) {
	query := `
		SELECT c.id, c.course_code, c.course_name, c.department,
		       c.year_of_study, c.semester, c.created_at, c.updated_at
		FROM courses c
		INNER JOIN teacher_course_assignments tca ON c.id = tca.course_id
		WHERE tca.teacher_id = ? AND tca.is_active = true
		ORDER BY c.year_of_study, c.semester
	`

	rows, err := r.conn.DB().QueryContext(ctx, query, teacherID)
	if err != nil {
		return nil, fmt.Errorf("error getting courses by teacher: %v", err)
	}
	defer rows.Close()

	var courses []*entities.Course
	for rows.Next() {
		course := &entities.Course{}
		err := rows.Scan(
			&course.ID, &course.CourseCode, &course.CourseName, &course.Department,
			&course.YearOfStudy, &course.Semester, &course.CreatedAt, &course.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning course row: %v", err)
		}
		courses = append(courses, course)
	}

	return courses, nil
}

// AssignTeacher assigns a teacher to a course
func (r *CourseRepository) AssignTeacher(ctx context.Context, courseID, teacherID, academicYear string) error {
	assignment := entities.NewTeacherCourseAssignment(teacherID, courseID, academicYear)

	query := `
		INSERT INTO teacher_course_assignments (
			id, teacher_id, course_id, academic_year,
			is_active, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.conn.DB().ExecContext(ctx, query,
		assignment.ID, assignment.TeacherID, assignment.CourseID,
		assignment.AcademicYear, assignment.IsActive,
		assignment.CreatedAt, assignment.UpdatedAt,
	)

	if err != nil {
		if isDuplicateKeyError(err) {
			return domain.ErrAlreadyExists
		}
		return fmt.Errorf("error assigning teacher to course: %v", err)
	}

	return nil
}

// UnassignTeacher removes a teacher from a course
func (r *CourseRepository) UnassignTeacher(ctx context.Context, courseID, teacherID, academicYear string) error {
	query := `
		UPDATE teacher_course_assignments
		SET is_active = false, updated_at = ?
		WHERE course_id = ? AND teacher_id = ? AND academic_year = ?
	`

	result, err := r.conn.DB().ExecContext(ctx, query,
		time.Now(), courseID, teacherID, academicYear,
	)

	if err != nil {
		return fmt.Errorf("error unassigning teacher from course: %v", err)
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

// GetTeacherAssignments retrieves all teacher course assignments
func (r *CourseRepository) GetTeacherAssignments(ctx context.Context, filters map[string]interface{}) ([]*entities.TeacherCourseAssignment, error) {
	var conditions []string
	var args []interface{}

	// Build WHERE clause based on filters
	for key, value := range filters {
		conditions = append(conditions, fmt.Sprintf("%s = ?", key))
		args = append(args, value)
	}

	query := `
		SELECT id, teacher_id, course_id, academic_year,
		       is_active, created_at, updated_at
		FROM teacher_course_assignments
	`

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY academic_year DESC, created_at DESC"

	rows, err := r.conn.DB().QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error getting teacher assignments: %v", err)
	}
	defer rows.Close()

	var assignments []*entities.TeacherCourseAssignment
	for rows.Next() {
		assignment := &entities.TeacherCourseAssignment{}
		err := rows.Scan(
			&assignment.ID, &assignment.TeacherID, &assignment.CourseID,
			&assignment.AcademicYear, &assignment.IsActive,
			&assignment.CreatedAt, &assignment.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning assignment row: %v", err)
		}
		assignments = append(assignments, assignment)
	}

	return assignments, nil
}
