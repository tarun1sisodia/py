package models

import "time"

// User represents a user in the system.
type User struct {
	ID           string    `json:"id"`
	Role         string    `json:"role"`
	FullName     string    `json:"full_name" validate:"required"`
	Username     string    `json:"username" validate:"required"`
	RollNumber   string    `json:"roll_number"`
	Email        string    `json:"email" validate:"required,email"`
	Course       string    `json:"course"`
	AcademicYear string    `json:"academic_year"`
	Phone        string    `json:"phone" validate:"required"`
	HighestDegree string    `json:"highest_degree"`
	Experience   string    `json:"experience"`
	PasswordHash string    `json:"password_hash" validate:"required"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	FirebaseUID  string    `json:"firebase_uid"` // Link to Firebase user
}
