package entities

import (
	"time"

	"github.com/google/uuid"
)

// TeacherCourseAssignment represents an assignment of a teacher to a course
type TeacherCourseAssignment struct {
	ID           string    `json:"id"`
	TeacherID    string    `json:"teacher_id"`
	CourseID     string    `json:"course_id"`
	AcademicYear string    `json:"academic_year"` // Format: 2023-2024
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// NewTeacherCourseAssignment creates a new teacher course assignment
func NewTeacherCourseAssignment(teacherID, courseID, academicYear string) *TeacherCourseAssignment {
	now := time.Now()
	return &TeacherCourseAssignment{
		ID:           uuid.New().String(),
		TeacherID:    teacherID,
		CourseID:     courseID,
		AcademicYear: academicYear,
		IsActive:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// Deactivate marks the assignment as inactive
func (t *TeacherCourseAssignment) Deactivate() {
	t.IsActive = false
	t.UpdatedAt = time.Now()
}

// Activate marks the assignment as active
func (t *TeacherCourseAssignment) Activate() {
	t.IsActive = true
	t.UpdatedAt = time.Now()
}

// IsActiveForAcademicYear checks if the assignment is active for the given academic year
func (t *TeacherCourseAssignment) IsActiveForAcademicYear(academicYear string) bool {
	return t.IsActive && t.AcademicYear == academicYear
}
