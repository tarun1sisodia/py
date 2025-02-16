package services

import (
	"context"
	"fmt"
)

type Course struct {
	ID          string
	CourseCode  string
	CourseName  string
	Department  string
	YearOfStudy int
	Semester    int
}

type TeacherCourseAssignment struct {
	ID           string
	TeacherID    string
	CourseID     string
	AcademicYear string
	IsActive     bool
}

type CourseService struct {
	courseRepo   CourseRepository
	firebaseAuth *FirebaseAuthService
}

func NewCourseService(
	courseRepo CourseRepository,
	firebaseAuth *FirebaseAuthService,
) *CourseService {
	return &CourseService{
		courseRepo:   courseRepo,
		firebaseAuth: firebaseAuth,
	}
}

// CreateCourse creates a new course
func (s *CourseService) CreateCourse(ctx context.Context, course *Course) error {
	// Validate course data
	if err := s.validateCourse(course); err != nil {
		return err
	}

	return s.courseRepo.Create(ctx, course)
}

// AssignTeacherToCourse assigns a teacher to a course
func (s *CourseService) AssignTeacherToCourse(ctx context.Context, assignment *TeacherCourseAssignment) error {
	// Verify teacher through Firebase
	user, err := s.firebaseAuth.GetUser(ctx, assignment.TeacherID)
	if err != nil {
		return fmt.Errorf("error verifying teacher: %v", err)
	}

	// Verify teacher role
	claims, err := s.firebaseAuth.VerifyIDToken(ctx, user.CustomClaims["role"].(string))
	if err != nil || claims.Claims["role"] != "teacher" {
		return fmt.Errorf("unauthorized: user is not a teacher")
	}

	// Check if course exists
	course, err := s.courseRepo.FindByID(ctx, assignment.CourseID)
	if err != nil {
		return fmt.Errorf("error finding course: %v", err)
	}
	if course == nil {
		return fmt.Errorf("course not found")
	}

	return s.courseRepo.AssignTeacher(ctx, assignment)
}

// GetTeacherCourses gets all courses assigned to a teacher
func (s *CourseService) GetTeacherCourses(ctx context.Context, teacherID string, activeOnly bool) ([]*Course, error) {
	return s.courseRepo.FindByTeacher(ctx, teacherID, activeOnly)
}

// GetStudentCourses gets all courses for a student's year and department
func (s *CourseService) GetStudentCourses(ctx context.Context, department string, yearOfStudy int) ([]*Course, error) {
	return s.courseRepo.FindByDepartmentAndYear(ctx, department, yearOfStudy)
}

func (s *CourseService) validateCourse(course *Course) error {
	if course.CourseCode == "" {
		return fmt.Errorf("course code is required")
	}
	if course.CourseName == "" {
		return fmt.Errorf("course name is required")
	}
	if course.Department == "" {
		return fmt.Errorf("department is required")
	}
	if course.YearOfStudy < 1 || course.YearOfStudy > 5 {
		return fmt.Errorf("invalid year of study")
	}
	if course.Semester < 1 || course.Semester > 10 {
		return fmt.Errorf("invalid semester")
	}
	return nil
}
