package entities

import (
	"time"

	"github.com/google/uuid"
)

// Course represents a course in the system
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

// NewCourse creates a new course instance
func NewCourse(courseCode, courseName, department string, yearOfStudy, semester int) *Course {
	now := time.Now()
	return &Course{
		ID:          uuid.New().String(),
		CourseCode:  courseCode,
		CourseName:  courseName,
		Department:  department,
		YearOfStudy: yearOfStudy,
		Semester:    semester,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// Update updates the course details
func (c *Course) Update(courseName, department string, yearOfStudy, semester int) {
	c.CourseName = courseName
	c.Department = department
	c.YearOfStudy = yearOfStudy
	c.Semester = semester
	c.UpdatedAt = time.Now()
}

// IsValidForYear checks if the course is valid for a given year of study
func (c *Course) IsValidForYear(year int) bool {
	return c.YearOfStudy == year
}

// IsValidForDepartment checks if the course is valid for a given department
func (c *Course) IsValidForDepartment(department string) bool {
	return c.Department == department
}
