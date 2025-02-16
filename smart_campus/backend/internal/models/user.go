package models

import (
	"time"
)

type Role string

const (
	RoleStudent Role = "student"
	RoleTeacher Role = "teacher"
	RoleAdmin   Role = "admin"
)

type User struct {
	ID               string    `json:"id"`
	Role             Role      `json:"role"`
	Email            string    `json:"email"`
	PasswordHash     string    `json:"-"`
	FullName         string    `json:"full_name"`
	EnrollmentNumber *string   `json:"enrollment_number,omitempty"`
	EmployeeID       *string   `json:"employee_id,omitempty"`
	Department       *string   `json:"department,omitempty"`
	YearOfStudy      *int      `json:"year_of_study,omitempty"`
	DeviceID         *string   `json:"device_id,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

func (u *User) IsStudent() bool {
	return u.Role == RoleStudent
}

func (u *User) IsTeacher() bool {
	return u.Role == RoleTeacher
}

func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

func (u *User) HasDeviceRegistered() bool {
	return u.DeviceID != nil && *u.DeviceID != ""
}

type UserRepository interface {
	Create(user *User) error
	GetByID(id string) (*User, error)
	GetByEmail(email string) (*User, error)
	Update(user *User) error
	Delete(id string) error
	List(offset, limit int) ([]*User, error)
	UpdatePassword(id string, passwordHash string) error
	UpdateDeviceID(id string, deviceID *string) error
}
