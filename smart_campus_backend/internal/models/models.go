package models

import "time"

// User represents a user in the system
type User struct {
	ID           string    `json:"id" db:"id"`
	Role         string    `json:"role" db:"role"`
	FullName     string    `json:"full_name" db:"full_name"`
	Email        string    `json:"email" db:"email"`
	Phone        string    `json:"phone" db:"phone"`
	PasswordHash string    `json:"-" db:"password_hash"`
	FirebaseUID  string    `json:"firebase_uid" db:"firebase_uid"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// Course represents a course in the system
type Course struct {
	ID          string    `json:"id" db:"id"`
	TeacherID   string    `json:"teacher_id" db:"teacher_id"`
	Code        string    `json:"code" db:"code"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Schedule    string    `json:"schedule" db:"schedule"`
	Room        string    `json:"room" db:"room"`
	Semester    string    `json:"semester" db:"semester"`
	Year        int       `json:"year" db:"year"`
	Status      string    `json:"status" db:"status"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// Session represents an attendance session
type Session struct {
	ID          string    `json:"id" db:"id"`
	CourseID    string    `json:"course_id" db:"course_id"`
	TeacherID   string    `json:"teacher_id" db:"teacher_id"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	StartTime   time.Time `json:"start_time" db:"start_time"`
	EndTime     time.Time `json:"end_time" db:"end_time"`
	Location    string    `json:"location" db:"location"`
	Status      string    `json:"status" db:"status"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// Device represents a registered device
type Device struct {
	ID          string    `json:"id" db:"id"`
	UserID      string    `json:"user_id" db:"user_id"`
	DeviceID    string    `json:"device_id" db:"device_id"`
	DeviceName  string    `json:"device_name" db:"device_name"`
	DeviceModel string    `json:"device_model" db:"device_model"`
	DeviceType  string    `json:"device_type" db:"device_type"`
	OSVersion   string    `json:"os_version" db:"os_version"`
	AppVersion  string    `json:"app_version" db:"app_version"`
	LastSeen    time.Time `json:"last_seen" db:"last_seen"`
	IsVerified  bool      `json:"is_verified" db:"is_verified"`
	Status      string    `json:"status" db:"status"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// Attendance represents an attendance record
type Attendance struct {
	ID        string    `json:"id" db:"id"`
	SessionID string    `json:"session_id" db:"session_id"`
	StudentID string    `json:"student_id" db:"student_id"`
	DeviceID  string    `json:"device_id" db:"device_id"`
	Status    string    `json:"status" db:"status"`
	Latitude  float64   `json:"latitude" db:"latitude"`
	Longitude float64   `json:"longitude" db:"longitude"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// AttendanceStats represents attendance statistics
type AttendanceStats struct {
	TotalSessions     int     `json:"total_sessions"`
	AttendedSessions  int     `json:"attended_sessions"`
	AbsentSessions    int     `json:"absent_sessions"`
	AttendanceRate    float64 `json:"attendance_rate"`
	TotalStudents     int     `json:"total_students"`
	PresentStudents   int     `json:"present_students"`
	AbsentStudents    int     `json:"absent_students"`
	ParticipationRate float64 `json:"participation_rate"`
}
