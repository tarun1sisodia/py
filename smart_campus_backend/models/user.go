package models

import "time"

type UserRole string

const (
	RoleTeacher UserRole = "teacher"
	RoleStudent UserRole = "student"
)

type User struct {
	ID            string    `json:"id" db:"id"`
	Role          UserRole  `json:"role" db:"role"`
	FullName      string    `json:"full_name" db:"full_name"`
	Username      string    `json:"username,omitempty" db:"username"`       // For teachers
	RollNumber    string    `json:"roll_number,omitempty" db:"roll_number"` // For students
	Email         string    `json:"email,omitempty" db:"email"`
	Course        string    `json:"course,omitempty" db:"course"`               // For students
	AcademicYear  string    `json:"academic_year,omitempty" db:"academic_year"` // For students
	Phone         string    `json:"phone" db:"phone"`
	HighestDegree string    `json:"highest_degree,omitempty" db:"highest_degree"` // For teachers
	Experience    string    `json:"experience,omitempty" db:"experience"`         // For teachers
	PasswordHash  string    `json:"-" db:"password_hash"`
	FirebaseUID   string    `json:"-" db:"firebase_uid"`
	PhoneVerified bool      `json:"phone_verified" db:"phone_verified"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

type RegisterRequest struct {
	Role          UserRole `json:"role" binding:"required,oneof=teacher student"`
	FullName      string   `json:"full_name" binding:"required"`
	Username      string   `json:"username,omitempty"`      // Required for teachers
	RollNumber    string   `json:"roll_number,omitempty"`   // Required for students
	Email         string   `json:"email,omitempty"`         // Required for teachers
	Course        string   `json:"course,omitempty"`        // Required for students
	AcademicYear  string   `json:"academic_year,omitempty"` // Required for students
	Phone         string   `json:"phone" binding:"required"`
	HighestDegree string   `json:"highest_degree,omitempty"` // Required for teachers
	Experience    string   `json:"experience,omitempty"`     // Required for teachers
	Password      string   `json:"password" binding:"required,min=8"`
}

type LoginRequest struct {
	Role       UserRole `json:"role" binding:"required,oneof=teacher student"`
	Email      string   `json:"email,omitempty"`       // For teachers
	RollNumber string   `json:"roll_number,omitempty"` // For students
	Password   string   `json:"password" binding:"required"`
}

type AuthResponse struct {
	User         *User  `json:"user"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
