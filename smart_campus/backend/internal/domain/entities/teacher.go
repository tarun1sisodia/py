package entities

import "time"

// Teacher represents a teacher in the system
type Teacher struct {
	BaseEntity
	UserID         string        `json:"user_id"`
	DepartmentID   string        `json:"department_id"`
	EmployeeID     string        `json:"employee_id"`
	Position       string        `json:"position"`
	Status         TeacherStatus `json:"status"`
	JoinedAt       time.Time     `json:"joined_at"`
	FirstName      string        `json:"first_name"`
	LastName       string        `json:"last_name"`
	Email          string        `json:"email"`
	DepartmentName string        `json:"department_name"`
}

// FullName returns the teacher's full name
func (t *Teacher) FullName() string {
	return t.FirstName + " " + t.LastName
}

// IsActive returns true if the teacher's status is active
func (t *Teacher) IsActive() bool {
	return t.Status == TeacherStatusActive
}

// IsInactive returns true if the teacher's status is inactive
func (t *Teacher) IsInactive() bool {
	return t.Status == TeacherStatusInactive
}

// IsOnLeave returns true if the teacher's status is on leave
func (t *Teacher) IsOnLeave() bool {
	return t.Status == TeacherStatusOnLeave
}

// IsTerminated returns true if the teacher's status is terminated
func (t *Teacher) IsTerminated() bool {
	return t.Status == TeacherStatusTerminated
}

// HasStatus returns true if the teacher has the specified status
func (t *Teacher) HasStatus(status TeacherStatus) bool {
	return t.Status == status
}

// UpdateStatus updates the teacher's status
func (t *Teacher) UpdateStatus(status TeacherStatus) {
	t.Status = status
	t.UpdatedAt = time.Now()
}

// UpdatePosition updates the teacher's position
func (t *Teacher) UpdatePosition(position string) {
	t.Position = position
	t.UpdatedAt = time.Now()
}

// UpdateDepartment updates the teacher's department
func (t *Teacher) UpdateDepartment(departmentID, departmentName string) {
	t.DepartmentID = departmentID
	t.DepartmentName = departmentName
	t.UpdatedAt = time.Now()
}

// ToPublic returns a public view of the teacher
func (t *Teacher) ToPublic() map[string]interface{} {
	return map[string]interface{}{
		"id":              t.ID,
		"user_id":         t.UserID,
		"department_id":   t.DepartmentID,
		"employee_id":     t.EmployeeID,
		"position":        t.Position,
		"status":          t.Status,
		"joined_at":       t.JoinedAt,
		"first_name":      t.FirstName,
		"last_name":       t.LastName,
		"email":           t.Email,
		"department_name": t.DepartmentName,
		"created_at":      t.CreatedAt,
		"updated_at":      t.UpdatedAt,
	}
}
