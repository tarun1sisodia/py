package entities

import "time"

// Department represents a department in the system
type Department struct {
	BaseEntity
	Name        string           `json:"name"`
	Code        string           `json:"code"`
	FacultyID   string           `json:"faculty_id"`
	HeadID      string           `json:"head_id"`
	Description string           `json:"description"`
	Status      DepartmentStatus `json:"status"`
	FacultyName string           `json:"faculty_name"`
	HeadName    string           `json:"head_name"`
}

// IsActive returns true if the department's status is active
func (d *Department) IsActive() bool {
	return d.Status == DepartmentStatusActive
}

// IsInactive returns true if the department's status is inactive
func (d *Department) IsInactive() bool {
	return d.Status == DepartmentStatusInactive
}

// IsMerged returns true if the department's status is merged
func (d *Department) IsMerged() bool {
	return d.Status == DepartmentStatusMerged
}

// HasStatus returns true if the department has the specified status
func (d *Department) HasStatus(status DepartmentStatus) bool {
	return d.Status == status
}

// UpdateStatus updates the department's status
func (d *Department) UpdateStatus(status DepartmentStatus) {
	d.Status = status
	d.UpdatedAt = time.Now()
}

// UpdateHead updates the department's head
func (d *Department) UpdateHead(headID, headName string) {
	d.HeadID = headID
	d.HeadName = headName
	d.UpdatedAt = time.Now()
}

// UpdateFaculty updates the department's faculty
func (d *Department) UpdateFaculty(facultyID, facultyName string) {
	d.FacultyID = facultyID
	d.FacultyName = facultyName
	d.UpdatedAt = time.Now()
}

// UpdateInfo updates the department's basic information
func (d *Department) UpdateInfo(name, code, description string) {
	d.Name = name
	d.Code = code
	d.Description = description
	d.UpdatedAt = time.Now()
}

// ToPublic returns a public view of the department
func (d *Department) ToPublic() map[string]interface{} {
	return map[string]interface{}{
		"id":           d.ID,
		"name":         d.Name,
		"code":         d.Code,
		"faculty_id":   d.FacultyID,
		"head_id":      d.HeadID,
		"description":  d.Description,
		"status":       d.Status,
		"faculty_name": d.FacultyName,
		"head_name":    d.HeadName,
		"created_at":   d.CreatedAt,
		"updated_at":   d.UpdatedAt,
	}
}
