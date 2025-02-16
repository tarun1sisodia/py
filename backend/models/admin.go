package models

import (
	"encoding/json"
	"time"

	"smart_attendance_backend/utils"

	"gorm.io/gorm"
)

type AdminRole string

const (
	AdminRoleSuperAdmin AdminRole = "super_admin"
	AdminRoleModerator  AdminRole = "moderator"
)

type Admin struct {
	ID           string         `gorm:"type:varchar(36);primary_key" json:"id"`
	Username     string         `gorm:"type:varchar(100);unique;not null" json:"username"`
	Email        string         `gorm:"type:varchar(255);unique;not null" json:"email"`
	PasswordHash string         `gorm:"type:varchar(255);not null" json:"-"`
	Role         AdminRole      `gorm:"type:enum('super_admin','moderator');not null" json:"role"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (a *Admin) BeforeCreate(tx *gorm.DB) error {
	if a.ID == "" {
		a.ID = utils.GenerateUUID()
	}
	return nil
}

func (a *Admin) IsSuperAdmin() bool {
	return a.Role == AdminRoleSuperAdmin
}

type AdminAuditLog struct {
	ID        string          `gorm:"type:varchar(36);primary_key" json:"id"`
	AdminID   string          `gorm:"type:varchar(36);not null" json:"admin_id"`
	Admin     Admin           `gorm:"foreignKey:AdminID" json:"admin"`
	Action    string          `gorm:"type:varchar(100);not null" json:"action"`
	Details   json.RawMessage `gorm:"type:json" json:"details"`
	IPAddress string          `gorm:"type:varchar(45)" json:"ip_address"`
	CreatedAt time.Time       `json:"created_at"`
}

func (a *AdminAuditLog) BeforeCreate(tx *gorm.DB) error {
	if a.ID == "" {
		a.ID = utils.GenerateUUID()
	}
	return nil
}

type Report struct {
	ID          string          `gorm:"type:varchar(36);primary_key" json:"id"`
	ReportType  string          `gorm:"type:enum('attendance','user_activity','security');not null" json:"report_type"`
	GeneratedBy string          `gorm:"type:varchar(36);not null" json:"generated_by"`
	Admin       Admin           `gorm:"foreignKey:GeneratedBy" json:"admin"`
	Data        json.RawMessage `gorm:"type:json" json:"data"`
	GeneratedAt time.Time       `gorm:"not null" json:"generated_at"`
}

func (r *Report) BeforeCreate(tx *gorm.DB) error {
	if r.ID == "" {
		r.ID = utils.GenerateUUID()
	}
	if r.GeneratedAt.IsZero() {
		r.GeneratedAt = time.Now()
	}
	return nil
}
