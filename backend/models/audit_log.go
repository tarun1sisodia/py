package models

import (
	"encoding/json"
	"time"

	"smart_attendance_backend/utils"

	"gorm.io/gorm"
)

type AuditLog struct {
	ID        string          `gorm:"type:varchar(36);primary_key" json:"id"`
	UserID    string          `gorm:"type:varchar(36);not null" json:"user_id"`
	User      User            `gorm:"foreignKey:UserID" json:"user"`
	Action    string          `gorm:"type:varchar(100);not null" json:"action"`
	Details   json.RawMessage `gorm:"type:json" json:"details"`
	IPAddress string          `gorm:"type:varchar(45)" json:"ip_address"`
	CreatedAt time.Time       `json:"created_at"`
}

func (a *AuditLog) BeforeCreate(tx *gorm.DB) error {
	if a.ID == "" {
		a.ID = utils.GenerateUUID()
	}
	return nil
}

func (a *AuditLog) SetDetails(details interface{}) error {
	data, err := json.Marshal(details)
	if err != nil {
		return err
	}
	a.Details = data
	return nil
}

func (a *AuditLog) GetDetails(out interface{}) error {
	if a.Details == nil {
		return nil
	}
	return json.Unmarshal(a.Details, out)
}
