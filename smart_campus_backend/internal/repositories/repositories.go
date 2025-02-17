package repositories

import (
	"context"
	"smart_campus_backend/internal/models"
)

// UserRepository defines the interface for user-related database operations
type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	FindByID(ctx context.Context, id string) (*models.User, error)
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindByPhone(ctx context.Context, phone string) (*models.User, error)
	FindByRole(ctx context.Context, role string) ([]*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id string) error
}

// CourseRepository defines the interface for course-related database operations
type CourseRepository interface {
	Create(ctx context.Context, course *models.Course) error
	FindByID(ctx context.Context, id string) (*models.Course, error)
	FindByTeacherID(ctx context.Context, teacherID string) ([]*models.Course, error)
	FindBySemesterAndYear(ctx context.Context, semester string, year string) ([]*models.Course, error)
	FindByStatus(ctx context.Context, status string) ([]*models.Course, error)
	FindAll(ctx context.Context) ([]*models.Course, error)
	Update(ctx context.Context, course *models.Course) error
	Delete(ctx context.Context, id string) error
	EnrollStudent(ctx context.Context, courseID string, studentID string) error
	UnenrollStudent(ctx context.Context, courseID string, studentID string) error
	FindEnrolledStudents(ctx context.Context, courseID string) ([]*models.User, error)
}

// SessionRepository defines the interface for session-related database operations
type SessionRepository interface {
	Create(ctx context.Context, session *models.Session) error
	FindByID(ctx context.Context, id string) (*models.Session, error)
	FindByCourseID(ctx context.Context, courseID string) ([]*models.Session, error)
	FindByTeacherID(ctx context.Context, teacherID string) ([]*models.Session, error)
	FindByStatus(ctx context.Context, status string) ([]*models.Session, error)
	FindAll(ctx context.Context) ([]*models.Session, error)
	Update(ctx context.Context, session *models.Session) error
	Delete(ctx context.Context, id string) error
}

// DeviceRepository defines the interface for device-related database operations
type DeviceRepository interface {
	Create(ctx context.Context, device *models.Device) error
	FindByID(ctx context.Context, id string) (*models.Device, error)
	FindByUserID(ctx context.Context, userID string) (*models.Device, error)
	FindByDeviceID(ctx context.Context, deviceID string) (*models.Device, error)
	FindByStatus(ctx context.Context, status string) ([]*models.Device, error)
	FindByVerificationStatus(ctx context.Context, isVerified bool) ([]*models.Device, error)
	FindAll(ctx context.Context) ([]*models.Device, error)
	Update(ctx context.Context, device *models.Device) error
	Delete(ctx context.Context, id string) error
}

// AttendanceRepository defines the interface for attendance-related database operations
type AttendanceRepository interface {
	Create(ctx context.Context, attendance *models.Attendance) error
	FindByID(ctx context.Context, id string) (*models.Attendance, error)
	FindBySessionID(ctx context.Context, sessionID string) ([]*models.Attendance, error)
	FindByStudentID(ctx context.Context, studentID string) ([]*models.Attendance, error)
	FindByStatus(ctx context.Context, status string) ([]*models.Attendance, error)
	FindAll(ctx context.Context) ([]*models.Attendance, error)
	Update(ctx context.Context, attendance *models.Attendance) error
	Delete(ctx context.Context, id string) error
	GetSessionStats(ctx context.Context, sessionID string) (*models.AttendanceStats, error)
	GetStudentStats(ctx context.Context, studentID string) (*models.AttendanceStats, error)
}
