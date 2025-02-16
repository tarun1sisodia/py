package models

import (
	"time"

	"smart_attendance_backend/utils"

	"gorm.io/gorm"
)

type UserRole string

const (
	RoleTeacher UserRole = "teacher"
	RoleStudent UserRole = "student"
)

type User struct {
	ID            string         `gorm:"type:varchar(36);primary_key" json:"id"`
	Role          UserRole       `gorm:"type:enum('teacher','student');not null" json:"role"`
	FullName      string         `gorm:"type:varchar(100);not null" json:"full_name"`
	Username      *string        `gorm:"type:varchar(100);unique" json:"username,omitempty"`
	RollNumber    *string        `gorm:"type:varchar(100);unique" json:"roll_number,omitempty"`
	Email         *string        `gorm:"type:varchar(255)" json:"email,omitempty"`
	Course        *string        `gorm:"type:varchar(100)" json:"course,omitempty"`
	AcademicYear  *string        `gorm:"type:varchar(50)" json:"academic_year,omitempty"`
	Phone         string         `gorm:"type:varchar(20);not null" json:"phone"`
	HighestDegree *string        `gorm:"type:varchar(100)" json:"highest_degree,omitempty"`
	Experience    *string        `gorm:"type:varchar(50)" json:"experience,omitempty"`
	PasswordHash  string         `gorm:"type:varchar(255);not null" json:"-"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = utils.GenerateUUID()
	}
	return nil
}

func (u *User) IsTeacher() bool {
	return u.Role == RoleTeacher
}

func (u *User) IsStudent() bool {
	return u.Role == RoleStudent
}
