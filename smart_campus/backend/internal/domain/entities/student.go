package entities

import "time"

// Student represents a student in the system
type Student struct {
	BaseEntity
	UserID         string        `json:"user_id"`
	DepartmentID   string        `json:"department_id"`
	StudentNumber  string        `json:"student_number"`
	Year           int           `json:"year"`
	Status         StudentStatus `json:"status"`
	EnrolledAt     time.Time     `json:"enrolled_at"`
	FirstName      string        `json:"first_name"`
	LastName       string        `json:"last_name"`
	Email          string        `json:"email"`
	DepartmentName string        `json:"department_name"`
}

// FullName returns the student's full name
func (s *Student) FullName() string {
	return s.FirstName + " " + s.LastName
}

// IsActive returns true if the student's status is active
func (s *Student) IsActive() bool {
	return s.Status == StudentStatusActive
}

// IsInactive returns true if the student's status is inactive
func (s *Student) IsInactive() bool {
	return s.Status == StudentStatusInactive
}

// IsGraduated returns true if the student's status is graduated
func (s *Student) IsGraduated() bool {
	return s.Status == StudentStatusGraduated
}

// IsSuspended returns true if the student's status is suspended
func (s *Student) IsSuspended() bool {
	return s.Status == StudentStatusSuspended
}

// IsWithdrawn returns true if the student's status is withdrawn
func (s *Student) IsWithdrawn() bool {
	return s.Status == StudentStatusWithdrawn
}

// HasStatus returns true if the student has the specified status
func (s *Student) HasStatus(status StudentStatus) bool {
	return s.Status == status
}

// UpdateStatus updates the student's status
func (s *Student) UpdateStatus(status StudentStatus) {
	s.Status = status
	s.UpdatedAt = time.Now()
}

// UpdateYear updates the student's year
func (s *Student) UpdateYear(year int) {
	s.Year = year
	s.UpdatedAt = time.Now()
}

// UpdateDepartment updates the student's department
func (s *Student) UpdateDepartment(departmentID, departmentName string) {
	s.DepartmentID = departmentID
	s.DepartmentName = departmentName
	s.UpdatedAt = time.Now()
}

// IsEligibleForGraduation returns true if the student is eligible for graduation
func (s *Student) IsEligibleForGraduation() bool {
	// This is a placeholder implementation
	// The actual logic should be based on the institution's requirements
	return s.Status == StudentStatusActive && s.Year >= 4
}

// ToPublic returns a public view of the student
func (s *Student) ToPublic() map[string]interface{} {
	return map[string]interface{}{
		"id":              s.ID,
		"user_id":         s.UserID,
		"department_id":   s.DepartmentID,
		"student_number":  s.StudentNumber,
		"year":            s.Year,
		"status":          s.Status,
		"enrolled_at":     s.EnrolledAt,
		"first_name":      s.FirstName,
		"last_name":       s.LastName,
		"email":           s.Email,
		"department_name": s.DepartmentName,
		"created_at":      s.CreatedAt,
		"updated_at":      s.UpdatedAt,
	}
}
