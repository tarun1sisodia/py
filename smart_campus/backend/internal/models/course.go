package models

import (
	"time"
)

type Course struct {
	ID          string    `json:"id"`
	CourseCode  string    `json:"course_code"`
	CourseName  string    `json:"course_name"`
	Department  string    `json:"department"`
	YearOfStudy int       `json:"year_of_study"`
	Semester    int       `json:"semester"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type TeacherCourseAssignment struct {
	ID           string    `json:"id"`
	TeacherID    string    `json:"teacher_id"`
	CourseID     string    `json:"course_id"`
	AcademicYear string    `json:"academic_year"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type CourseRepository interface {
	Create(course *Course) error
	GetByID(id string) (*Course, error)
	GetByCode(code string) (*Course, error)
	Update(course *Course) error
	Delete(id string) error
	List(offset, limit int) ([]*Course, error)
	GetByDepartment(department string) ([]*Course, error)
	GetByTeacher(teacherID string) ([]*Course, error)
	AssignTeacher(assignment *TeacherCourseAssignment) error
	UnassignTeacher(teacherID, courseID string) error
	GetTeacherAssignments(teacherID string) ([]*TeacherCourseAssignment, error)
}
