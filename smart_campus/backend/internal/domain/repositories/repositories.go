package repositories

import (
	"context"
	"errors"
	"time"

	"smart_campus/internal/domain/entities"
)

var (
	// ErrNotFound is returned when a requested entity is not found
	ErrNotFound = errors.New("entity not found")

	// ErrDuplicate is returned when trying to create a duplicate entity
	ErrDuplicate = errors.New("entity already exists")

	// ErrInvalidInput is returned when the input data is invalid
	ErrInvalidInput = errors.New("invalid input data")

	// ErrDatabase is returned when a database error occurs
	ErrDatabase = errors.New("database error")
)

// UserRepository defines the interface for user data access
type UserRepository interface {
	Create(ctx context.Context, user *entities.User) error
	Update(ctx context.Context, user *entities.User) error
	Delete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (*entities.User, error)
	FindByEmail(ctx context.Context, email string) (*entities.User, error)
	FindByRole(ctx context.Context, role entities.UserRole) ([]*entities.User, error)
	List(ctx context.Context, offset, limit int) ([]*entities.User, error)
	Count(ctx context.Context) (int64, error)
	GetByDeviceID(ctx context.Context, deviceID string) (*entities.User, error)
	UpdatePassword(ctx context.Context, userID, passwordHash string) error
	UpdateDeviceID(ctx context.Context, userID, deviceID string) error
	GetStudentsByDepartmentAndYear(ctx context.Context, department string, year int) ([]*entities.User, error)
	GetTeachersByDepartment(ctx context.Context, department string) ([]*entities.User, error)
}

// TeacherRepository defines the interface for teacher data access
type TeacherRepository interface {
	Create(ctx context.Context, teacher *entities.Teacher) error
	GetByID(ctx context.Context, id string) (*entities.Teacher, error)
	GetByUserID(ctx context.Context, userID string) (*entities.Teacher, error)
	GetByEmployeeID(ctx context.Context, employeeID string) (*entities.Teacher, error)
	Update(ctx context.Context, teacher *entities.Teacher) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, filters map[string]interface{}, page, limit int) ([]*entities.Teacher, int, error)
	GetByDepartment(ctx context.Context, departmentID string) ([]*entities.Teacher, error)
	UpdateStatus(ctx context.Context, teacherID string, status entities.TeacherStatus) error
	Search(ctx context.Context, query string, limit int) ([]*entities.Teacher, error)
	GetCourseCount(ctx context.Context, teacherID string) (int, error)
}

// StudentRepository defines the interface for student data access
type StudentRepository interface {
	Create(ctx context.Context, student *entities.Student) error
	GetByID(ctx context.Context, id string) (*entities.Student, error)
	GetByUserID(ctx context.Context, userID string) (*entities.Student, error)
	GetByStudentNumber(ctx context.Context, studentNumber string) (*entities.Student, error)
	Update(ctx context.Context, student *entities.Student) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, filters map[string]interface{}, page, limit int) ([]*entities.Student, int, error)
	GetByDepartment(ctx context.Context, departmentID string) ([]*entities.Student, error)
	GetByDepartmentAndYear(ctx context.Context, departmentID string, year int) ([]*entities.Student, error)
	GetByCourse(ctx context.Context, courseID string) ([]*entities.Student, error)
	UpdateStatus(ctx context.Context, studentID string, status entities.StudentStatus) error
	Search(ctx context.Context, query string, limit int) ([]*entities.Student, error)
}

// DepartmentRepository defines the interface for department data access
type DepartmentRepository interface {
	Create(ctx context.Context, department *entities.Department) error
	GetByID(ctx context.Context, id string) (*entities.Department, error)
	GetByCode(ctx context.Context, code string) (*entities.Department, error)
	Update(ctx context.Context, department *entities.Department) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, filters map[string]interface{}, page, limit int) ([]*entities.Department, int, error)
	GetByFaculty(ctx context.Context, facultyID string) ([]*entities.Department, error)
	UpdateStatus(ctx context.Context, departmentID string, status entities.DepartmentStatus) error
	Search(ctx context.Context, query string, limit int) ([]*entities.Department, error)
}

// FacultyRepository defines the interface for faculty data access
type FacultyRepository interface {
	Create(ctx context.Context, faculty *entities.Faculty) error
	GetByID(ctx context.Context, id string) (*entities.Faculty, error)
	GetByCode(ctx context.Context, code string) (*entities.Faculty, error)
	Update(ctx context.Context, faculty *entities.Faculty) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, filters map[string]interface{}, page, limit int) ([]*entities.Faculty, int, error)
	UpdateStatus(ctx context.Context, facultyID string, status entities.FacultyStatus) error
	Search(ctx context.Context, query string, limit int) ([]*entities.Faculty, error)
	GetDepartmentCount(ctx context.Context, facultyID string) (int, error)
}

// AuthLogRepository defines the interface for authentication log data access
type AuthLogRepository interface {
	Create(ctx context.Context, log *entities.AuthLog) error
	GetByID(ctx context.Context, id string) (*entities.AuthLog, error)
	List(ctx context.Context, filters map[string]interface{}, page, limit int) ([]*entities.AuthLog, int, error)
	GetByUser(ctx context.Context, userID string) ([]*entities.AuthLog, error)
	GetByUserAndType(ctx context.Context, userID string, logType entities.AuthLogType) ([]*entities.AuthLog, error)
	GetByDateRange(ctx context.Context, startDate, endDate time.Time) ([]*entities.AuthLog, error)
	GetByUserAndDateRange(ctx context.Context, userID string, startDate, endDate time.Time) ([]*entities.AuthLog, error)
	GetFailedAttempts(ctx context.Context, userID string, since time.Time) ([]*entities.AuthLog, error)
	GetLastSuccessfulLogin(ctx context.Context, userID string) (*entities.AuthLog, error)
	CountFailedAttempts(ctx context.Context, userID string, since time.Time) (int, error)
	DeleteOldLogs(ctx context.Context, before time.Time) error
}

// DeviceBindingRepository defines the interface for device binding data access
type DeviceBindingRepository interface {
	Create(ctx context.Context, binding *entities.DeviceBinding) error
	Update(ctx context.Context, binding *entities.DeviceBinding) error
	Delete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (*entities.DeviceBinding, error)
	FindByUserAndDeviceID(ctx context.Context, userID, deviceID string) (*entities.DeviceBinding, error)
	FindByUser(ctx context.Context, userID string) ([]*entities.DeviceBinding, error)
	FindActiveByUser(ctx context.Context, userID string) ([]*entities.DeviceBinding, error)
	List(ctx context.Context, offset, limit int) ([]*entities.DeviceBinding, error)
	Count(ctx context.Context) (int64, error)
}

// AttendanceSessionRepository defines the interface for attendance session data access
type AttendanceSessionRepository interface {
	Create(ctx context.Context, session *entities.AttendanceSession) error
	GetByID(ctx context.Context, id string) (*entities.AttendanceSession, error)
	Update(ctx context.Context, session *entities.AttendanceSession) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, filters map[string]interface{}, page, limit int) ([]*entities.AttendanceSession, int, error)
	GetActiveSessions(ctx context.Context) ([]*entities.AttendanceSession, error)
	GetByTeacher(ctx context.Context, teacherID string) ([]*entities.AttendanceSession, error)
	GetByCourse(ctx context.Context, courseID string) ([]*entities.AttendanceSession, error)
	GetByDateRange(ctx context.Context, startDate, endDate time.Time) ([]*entities.AttendanceSession, error)
	GetActiveSessionsByStudent(ctx context.Context, studentID string) ([]*entities.AttendanceSession, error)
	CompleteSession(ctx context.Context, id string) error
	CancelSession(ctx context.Context, id string) error
	GetSessionStatistics(ctx context.Context, id string) (map[string]interface{}, error)
}

// AttendanceRecordRepository defines the interface for attendance record data access
type AttendanceRecordRepository interface {
	Create(ctx context.Context, record *entities.AttendanceRecord) error
	GetByID(ctx context.Context, id string) (*entities.AttendanceRecord, error)
	Update(ctx context.Context, record *entities.AttendanceRecord) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, filters map[string]interface{}, page, limit int) ([]*entities.AttendanceRecord, int, error)
	GetBySession(ctx context.Context, sessionID string) ([]*entities.AttendanceRecord, error)
	GetByStudent(ctx context.Context, studentID string) ([]*entities.AttendanceRecord, error)
	GetByStudentAndDateRange(ctx context.Context, studentID string, startDate, endDate time.Time) ([]*entities.AttendanceRecord, error)
	GetByStudentAndCourse(ctx context.Context, studentID, courseID string) ([]*entities.AttendanceRecord, error)
	GetBySessionAndStudent(ctx context.Context, sessionID, studentID string) (*entities.AttendanceRecord, error)
	UpdateVerificationStatus(ctx context.Context, recordID string, status entities.AttendanceStatus) error
	AddVerificationLog(ctx context.Context, recordID, log string) error
	GetStudentAttendanceStats(ctx context.Context, studentID string, courseID *string) (map[string]interface{}, error)
}

// OTPVerificationRepository defines the interface for OTP verification data access
type OTPVerificationRepository interface {
	Create(ctx context.Context, verification *entities.OTPVerification) error
	Update(ctx context.Context, verification *entities.OTPVerification) error
	Delete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (*entities.OTPVerification, error)
	FindByCode(ctx context.Context, code string) (*entities.OTPVerification, error)
	FindLatestByUserID(ctx context.Context, userID string) (*entities.OTPVerification, error)
	FindByUserAndPurpose(ctx context.Context, userID string, purpose entities.OTPPurpose) ([]*entities.OTPVerification, error)
	List(ctx context.Context, offset, limit int) ([]*entities.OTPVerification, error)
	Count(ctx context.Context) (int64, error)
	DeleteExpired(ctx context.Context) error
}

// CourseRepository defines the interface for course data access
type CourseRepository interface {
	Create(ctx context.Context, course *entities.Course) error
	GetByID(ctx context.Context, id string) (*entities.Course, error)
	GetByCourseCode(ctx context.Context, courseCode string) (*entities.Course, error)
	Update(ctx context.Context, course *entities.Course) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, filters map[string]interface{}, page, limit int) ([]*entities.Course, int, error)
	GetByDepartment(ctx context.Context, department string) ([]*entities.Course, error)
	GetByDepartmentAndYear(ctx context.Context, department string, year int) ([]*entities.Course, error)
	GetByTeacher(ctx context.Context, teacherID string) ([]*entities.Course, error)
	AssignTeacher(ctx context.Context, courseID, teacherID, academicYear string) error
	UnassignTeacher(ctx context.Context, courseID, teacherID, academicYear string) error
	GetTeacherAssignments(ctx context.Context, filters map[string]interface{}) ([]*entities.TeacherCourseAssignment, error)
}

// SystemSettingsRepository defines the interface for system settings data access
type SystemSettingsRepository interface {
	Create(ctx context.Context, setting *entities.SystemSetting) error
	GetByID(ctx context.Context, id string) (*entities.SystemSetting, error)
	GetByKey(ctx context.Context, key entities.SystemSettingKey) (*entities.SystemSetting, error)
	Update(ctx context.Context, setting *entities.SystemSetting) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]*entities.SystemSetting, error)
	GetMultipleByKeys(ctx context.Context, keys []entities.SystemSettingKey) ([]*entities.SystemSetting, error)
	UpdateValue(ctx context.Context, key entities.SystemSettingKey, value interface{}) error
	GetIntValue(ctx context.Context, key entities.SystemSettingKey) (int, error)
	GetFloatValue(ctx context.Context, key entities.SystemSettingKey) (float64, error)
	GetBoolValue(ctx context.Context, key entities.SystemSettingKey) (bool, error)
	GetStringValue(ctx context.Context, key entities.SystemSettingKey) (string, error)
	Exists(ctx context.Context, key entities.SystemSettingKey) (bool, error)
}
