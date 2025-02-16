package entities

import "time"

// Faculty represents a faculty in the system
type Faculty struct {
	BaseEntity
	Name        string        `json:"name"`
	Code        string        `json:"code"`
	DeanID      string        `json:"dean_id"`
	Description string        `json:"description"`
	Status      FacultyStatus `json:"status"`
	DeanName    string        `json:"dean_name"`
}

// IsActive returns true if the faculty's status is active
func (f *Faculty) IsActive() bool {
	return f.Status == FacultyStatusActive
}

// IsInactive returns true if the faculty's status is inactive
func (f *Faculty) IsInactive() bool {
	return f.Status == FacultyStatusInactive
}

// IsMerged returns true if the faculty's status is merged
func (f *Faculty) IsMerged() bool {
	return f.Status == FacultyStatusMerged
}

// HasStatus returns true if the faculty has the specified status
func (f *Faculty) HasStatus(status FacultyStatus) bool {
	return f.Status == status
}

// UpdateStatus updates the faculty's status
func (f *Faculty) UpdateStatus(status FacultyStatus) {
	f.Status = status
	f.UpdatedAt = time.Now()
}

// UpdateDean updates the faculty's dean
func (f *Faculty) UpdateDean(deanID, deanName string) {
	f.DeanID = deanID
	f.DeanName = deanName
	f.UpdatedAt = time.Now()
}

// UpdateInfo updates the faculty's basic information
func (f *Faculty) UpdateInfo(name, code, description string) {
	f.Name = name
	f.Code = code
	f.Description = description
	f.UpdatedAt = time.Now()
}

// ToPublic returns a public view of the faculty
func (f *Faculty) ToPublic() map[string]interface{} {
	return map[string]interface{}{
		"id":          f.ID,
		"name":        f.Name,
		"code":        f.Code,
		"dean_id":     f.DeanID,
		"description": f.Description,
		"status":      f.Status,
		"dean_name":   f.DeanName,
		"created_at":  f.CreatedAt,
		"updated_at":  f.UpdatedAt,
	}
}
