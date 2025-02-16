package entities

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	// ErrInvalidOperation is returned when an operation is invalid for the current state
	ErrInvalidOperation = errors.New("invalid operation for current state")
	// ErrInvalidEmail is returned when the email is invalid
	ErrInvalidEmail = errors.New("invalid email format")
	// ErrInvalidName is returned when the name is invalid
	ErrInvalidName = errors.New("invalid name format")
	// ErrInvalidDepartment is returned when the department is invalid
	ErrInvalidDepartment = errors.New("invalid department")
	// ErrInvalidRole is returned when the role is invalid
	ErrInvalidRole = errors.New("invalid role")
	// ErrInvalidYear is returned when the year of study is invalid
	ErrInvalidYear = errors.New("invalid year of study")
	// ErrInvalidID is returned when the ID is invalid
	ErrInvalidID = errors.New("invalid ID format")
	// ErrEmptyField is returned when a required field is empty
	ErrEmptyField = errors.New("required field is empty")
)

var (
	emailRegex      = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	nameRegex       = regexp.MustCompile(`^[a-zA-Z\s-']{2,50}$`)
	enrollmentRegex = regexp.MustCompile(`^[A-Z0-9-]{5,20}$`)
	employeeIDRegex = regexp.MustCompile(`^[A-Z0-9-]{5,20}$`)
)

// User represents a user in the system
type User struct {
	BaseEntity
	Role             UserRole  `json:"role"`
	Email            string    `json:"email"`
	PasswordHash     string    `json:"-"`
	Name             string    `json:"full_name"`
	EnrollmentNumber string    `json:"enrollment_number,omitempty"`
	EmployeeID       string    `json:"employee_id,omitempty"`
	Department       string    `json:"department"`
	YearOfStudy      *int      `json:"year_of_study,omitempty"`
	DeviceID         string    `json:"device_id,omitempty"`
	LastLogin        time.Time `json:"last_login,omitempty"`
	Active           bool      `json:"is_active"`
}

// NewUser creates a new user instance
func NewUser(role UserRole, email, fullName, department string) *User {
	now := time.Now()
	return &User{
		BaseEntity: BaseEntity{
			ID:        uuid.New().String(),
			CreatedAt: now,
			UpdatedAt: now,
		},
		Role:       role,
		Email:      strings.ToLower(strings.TrimSpace(email)),
		Name:       strings.TrimSpace(fullName),
		Department: strings.TrimSpace(department),
		Active:     true,
	}
}

// Validate validates the user entity
func (u *User) Validate() error {
	// Check required fields
	if u.ID == "" {
		return fmt.Errorf("%w: ID", ErrEmptyField)
	}
	if u.Email == "" {
		return fmt.Errorf("%w: email", ErrEmptyField)
	}
	if u.Name == "" {
		return fmt.Errorf("%w: name", ErrEmptyField)
	}
	if u.Department == "" {
		return fmt.Errorf("%w: department", ErrEmptyField)
	}

	// Validate email format
	if !emailRegex.MatchString(u.Email) {
		return ErrInvalidEmail
	}

	// Validate name format
	if !nameRegex.MatchString(u.Name) {
		return ErrInvalidName
	}

	// Validate role
	if !u.IsValidRole() {
		return ErrInvalidRole
	}

	// Role-specific validation
	if u.IsStudent() {
		if u.YearOfStudy != nil {
			if *u.YearOfStudy < 1 || *u.YearOfStudy > 6 {
				return ErrInvalidYear
			}
		}
		if u.EnrollmentNumber != "" && !enrollmentRegex.MatchString(u.EnrollmentNumber) {
			return fmt.Errorf("invalid enrollment number format")
		}
	}

	if u.IsTeacher() && u.EmployeeID != "" && !employeeIDRegex.MatchString(u.EmployeeID) {
		return fmt.Errorf("invalid employee ID format")
	}

	return nil
}

// IsValidRole checks if the user's role is valid
func (u *User) IsValidRole() bool {
	switch u.Role {
	case UserRoleAdmin, UserRoleTeacher, UserRoleStudent, UserRoleStaff:
		return true
	default:
		return false
	}
}

// IsStudent checks if the user is a student
func (u *User) IsStudent() bool {
	return u.Role == UserRoleStudent
}

// IsTeacher checks if the user is a teacher
func (u *User) IsTeacher() bool {
	return u.Role == UserRoleTeacher
}

// IsAdmin checks if the user is an admin
func (u *User) IsAdmin() bool {
	return u.Role == UserRoleAdmin
}

// IsStaff checks if the user is a staff member
func (u *User) IsStaff() bool {
	return u.Role == UserRoleStaff
}

// SetPassword sets the hashed password for the user
func (u *User) SetPassword(hash string) {
	u.PasswordHash = hash
	u.UpdatedAt = time.Now()
}

// SetDeviceID sets the device ID for the user
func (u *User) SetDeviceID(deviceID string) {
	u.DeviceID = deviceID
	u.UpdatedAt = time.Now()
}

// SetEnrollmentNumber sets the enrollment number for student users
func (u *User) SetEnrollmentNumber(number string) error {
	if !u.IsStudent() {
		return ErrInvalidOperation
	}
	if !enrollmentRegex.MatchString(number) {
		return fmt.Errorf("invalid enrollment number format")
	}
	u.EnrollmentNumber = number
	u.UpdatedAt = time.Now()
	return nil
}

// SetEmployeeID sets the employee ID for teacher users
func (u *User) SetEmployeeID(id string) error {
	if !u.IsTeacher() {
		return ErrInvalidOperation
	}
	if !employeeIDRegex.MatchString(id) {
		return fmt.Errorf("invalid employee ID format")
	}
	u.EmployeeID = id
	u.UpdatedAt = time.Now()
	return nil
}

// SetYearOfStudy sets the year of study for student users
func (u *User) SetYearOfStudy(year int) error {
	if !u.IsStudent() {
		return ErrInvalidOperation
	}
	if year < 1 || year > 6 {
		return ErrInvalidYear
	}
	u.YearOfStudy = &year
	u.UpdatedAt = time.Now()
	return nil
}

// UpdateLastLogin updates the user's last login timestamp
func (u *User) UpdateLastLogin() {
	u.LastLogin = time.Now()
	u.UpdatedAt = time.Now()
}

// IsActive returns true if the user's status is active
func (u *User) IsActive() bool {
	return u.Active
}

// Activate activates the user
func (u *User) Activate() {
	u.Active = true
	u.UpdatedAt = time.Now()
}

// Deactivate deactivates the user
func (u *User) Deactivate() {
	u.Active = false
	u.UpdatedAt = time.Now()
}

// ToPublic returns a public view of the user
func (u *User) ToPublic() map[string]interface{} {
	data := map[string]interface{}{
		"id":         u.ID,
		"full_name":  u.Name,
		"email":      u.Email,
		"role":       u.Role,
		"department": u.Department,
		"is_active":  u.Active,
		"created_at": u.CreatedAt,
		"updated_at": u.UpdatedAt,
	}

	if !u.LastLogin.IsZero() {
		data["last_login"] = u.LastLogin
	}

	if u.IsStudent() {
		data["enrollment_number"] = u.EnrollmentNumber
		if u.YearOfStudy != nil {
			data["year_of_study"] = *u.YearOfStudy
		}
	}

	if u.IsTeacher() {
		data["employee_id"] = u.EmployeeID
	}

	return data
}

// Clone creates a deep copy of the user
func (u *User) Clone() *User {
	clone := *u
	if u.YearOfStudy != nil {
		year := *u.YearOfStudy
		clone.YearOfStudy = &year
	}
	return &clone
}

// Sanitize removes sensitive information from the user
func (u *User) Sanitize() {
	u.PasswordHash = ""
	u.DeviceID = ""
}
