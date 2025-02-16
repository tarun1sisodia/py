package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"smart_campus/internal/database"
	"smart_campus/internal/models"

	"github.com/google/uuid"
)

type MySQLCourseRepository struct {
	db *database.MySQLDB
}

func NewMySQLCourseRepository(db *database.MySQLDB) models.CourseRepository {
	return &MySQLCourseRepository{db: db}
}

func (r *MySQLCourseRepository) Create(course *models.Course) error {
	if course.ID == "" {
		course.ID = uuid.New().String()
	}
	now := time.Now()
	course.CreatedAt = now
	course.UpdatedAt = now

	query := `
		INSERT INTO courses (
			id, course_code, course_name, department,
			year_of_study, semester, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.Exec(query,
		course.ID, course.CourseCode, course.CourseName,
		course.Department, course.YearOfStudy, course.Semester,
		course.CreatedAt, course.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("error creating course: %v", err)
	}

	return nil
}

func (r *MySQLCourseRepository) GetByID(id string) (*models.Course, error) {
	course := &models.Course{}
	query := `
		SELECT id, course_code, course_name, department,
		year_of_study, semester, created_at, updated_at
		FROM courses WHERE id = ?
	`

	err := r.db.QueryRow(query, id).Scan(
		&course.ID, &course.CourseCode, &course.CourseName,
		&course.Department, &course.YearOfStudy, &course.Semester,
		&course.CreatedAt, &course.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error getting course by ID: %v", err)
	}

	return course, nil
}

func (r *MySQLCourseRepository) GetByCode(code string) (*models.Course, error) {
	course := &models.Course{}
	query := `
		SELECT id, course_code, course_name, department,
		year_of_study, semester, created_at, updated_at
		FROM courses WHERE course_code = ?
	`

	err := r.db.QueryRow(query, code).Scan(
		&course.ID, &course.CourseCode, &course.CourseName,
		&course.Department, &course.YearOfStudy, &course.Semester,
		&course.CreatedAt, &course.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error getting course by code: %v", err)
	}

	return course, nil
}

func (r *MySQLCourseRepository) Update(course *models.Course) error {
	course.UpdatedAt = time.Now()

	query := `
		UPDATE courses SET
			course_code = ?, course_name = ?, department = ?,
			year_of_study = ?, semester = ?, updated_at = ?
		WHERE id = ?
	`

	result, err := r.db.Exec(query,
		course.CourseCode, course.CourseName, course.Department,
		course.YearOfStudy, course.Semester, course.UpdatedAt,
		course.ID,
	)

	if err != nil {
		return fmt.Errorf("error updating course: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("course not found")
	}

	return nil
}

func (r *MySQLCourseRepository) Delete(id string) error {
	result, err := r.db.Exec("DELETE FROM courses WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("error deleting course: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("course not found")
	}

	return nil
}

func (r *MySQLCourseRepository) List(offset, limit int) ([]*models.Course, error) {
	query := `
		SELECT id, course_code, course_name, department,
		year_of_study, semester, created_at, updated_at
		FROM courses LIMIT ? OFFSET ?
	`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error listing courses: %v", err)
	}
	defer rows.Close()

	var courses []*models.Course
	for rows.Next() {
		course := &models.Course{}
		err := rows.Scan(
			&course.ID, &course.CourseCode, &course.CourseName,
			&course.Department, &course.YearOfStudy, &course.Semester,
			&course.CreatedAt, &course.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning course row: %v", err)
		}
		courses = append(courses, course)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating course rows: %v", err)
	}

	return courses, nil
}

func (r *MySQLCourseRepository) GetByDepartment(department string) ([]*models.Course, error) {
	query := `
		SELECT id, course_code, course_name, department,
		year_of_study, semester, created_at, updated_at
		FROM courses WHERE department = ?
	`

	rows, err := r.db.Query(query, department)
	if err != nil {
		return nil, fmt.Errorf("error getting courses by department: %v", err)
	}
	defer rows.Close()

	var courses []*models.Course
	for rows.Next() {
		course := &models.Course{}
		err := rows.Scan(
			&course.ID, &course.CourseCode, &course.CourseName,
			&course.Department, &course.YearOfStudy, &course.Semester,
			&course.CreatedAt, &course.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning course row: %v", err)
		}
		courses = append(courses, course)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating course rows: %v", err)
	}

	return courses, nil
}

func (r *MySQLCourseRepository) GetByTeacher(teacherID string) ([]*models.Course, error) {
	query := `
		SELECT c.id, c.course_code, c.course_name, c.department,
		c.year_of_study, c.semester, c.created_at, c.updated_at
		FROM courses c
		INNER JOIN teacher_course_assignments tca ON c.id = tca.course_id
		WHERE tca.teacher_id = ? AND tca.is_active = true
	`

	rows, err := r.db.Query(query, teacherID)
	if err != nil {
		return nil, fmt.Errorf("error getting courses by teacher: %v", err)
	}
	defer rows.Close()

	var courses []*models.Course
	for rows.Next() {
		course := &models.Course{}
		err := rows.Scan(
			&course.ID, &course.CourseCode, &course.CourseName,
			&course.Department, &course.YearOfStudy, &course.Semester,
			&course.CreatedAt, &course.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning course row: %v", err)
		}
		courses = append(courses, course)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating course rows: %v", err)
	}

	return courses, nil
}

func (r *MySQLCourseRepository) AssignTeacher(assignment *models.TeacherCourseAssignment) error {
	if assignment.ID == "" {
		assignment.ID = uuid.New().String()
	}
	now := time.Now()
	assignment.CreatedAt = now
	assignment.UpdatedAt = now

	query := `
		INSERT INTO teacher_course_assignments (
			id, teacher_id, course_id, academic_year,
			is_active, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.Exec(query,
		assignment.ID, assignment.TeacherID, assignment.CourseID,
		assignment.AcademicYear, assignment.IsActive,
		assignment.CreatedAt, assignment.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("error assigning teacher to course: %v", err)
	}

	return nil
}

func (r *MySQLCourseRepository) UnassignTeacher(teacherID, courseID string) error {
	query := `
		UPDATE teacher_course_assignments
		SET is_active = false, updated_at = ?
		WHERE teacher_id = ? AND course_id = ?
	`

	result, err := r.db.Exec(query, time.Now(), teacherID, courseID)
	if err != nil {
		return fmt.Errorf("error unassigning teacher from course: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("assignment not found")
	}

	return nil
}

func (r *MySQLCourseRepository) GetTeacherAssignments(teacherID string) ([]*models.TeacherCourseAssignment, error) {
	query := `
		SELECT id, teacher_id, course_id, academic_year,
		is_active, created_at, updated_at
		FROM teacher_course_assignments
		WHERE teacher_id = ?
	`

	rows, err := r.db.Query(query, teacherID)
	if err != nil {
		return nil, fmt.Errorf("error getting teacher assignments: %v", err)
	}
	defer rows.Close()

	var assignments []*models.TeacherCourseAssignment
	for rows.Next() {
		assignment := &models.TeacherCourseAssignment{}
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

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating assignment rows: %v", err)
	}

	return assignments, nil
}
