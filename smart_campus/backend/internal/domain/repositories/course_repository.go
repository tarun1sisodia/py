package repositories

import (
	"context"
)

// CourseRepository defines the interface for course-related database operations
type CourseRepository interface {
	// Create creates a new course
	Create(ctx context.Context, course *Course) error

	// FindByID finds a course by ID
	FindByID(ctx context.Context, id string) (*Course, error)

	// FindByTeacher finds courses assigned to a teacher
	FindByTeacher(ctx context.Context, teacherID string, activeOnly bool) ([]*Course, error)

	// FindByDepartmentAndYear finds courses for a department and year
	FindByDepartmentAndYear(ctx context.Context, department string, yearOfStudy int) ([]*Course, error)

	// AssignTeacher assigns a teacher to a course
	AssignTeacher(ctx context.Context, assignment *TeacherCourseAssignment) error

	// Update updates a course
	Update(ctx context.Context, course *Course) error
}

// Course represents a course entity
type Course struct {
	ID          string
	CourseCode  string
	CourseName  string
	Department  string
	YearOfStudy int
	Semester    int
	CreatedAt   string
	UpdatedAt   string
}

// TeacherCourseAssignment represents a teacher-course assignment
type TeacherCourseAssignment struct {
	ID           string
	TeacherID    string
	CourseID     string
	AcademicYear string
	IsActive     bool
	CreatedAt    string
	UpdatedAt    string
}
