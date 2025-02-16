package models

import (
	"time"

	"smart_attendance_backend/utils"

	"gorm.io/gorm"
)

type SessionStatus string

const (
	SessionStatusActive   SessionStatus = "active"
	SessionStatusComplete SessionStatus = "completed"
	SessionStatusCanceled SessionStatus = "cancelled"
)

type AttendanceSession struct {
	ID                string         `gorm:"type:varchar(36);primary_key" json:"id"`
	TeacherID         string         `gorm:"type:varchar(36);not null" json:"teacher_id"`
	Teacher           User           `gorm:"foreignKey:TeacherID" json:"teacher"`
	SubjectID         string         `gorm:"type:varchar(36);not null" json:"subject_id"`
	AcademicYear      string         `gorm:"type:varchar(50);not null" json:"academic_year"`
	StartTime         time.Time      `gorm:"not null" json:"start_time"`
	EndTime           time.Time      `gorm:"not null" json:"end_time"`
	CountdownDuration string         `gorm:"type:enum('30s','1m','3m');not null" json:"countdown_duration"`
	Status            SessionStatus  `gorm:"type:enum('active','completed','cancelled');default:'active'" json:"status"`
	WifiSSID          string         `gorm:"type:varchar(100);not null" json:"wifi_ssid"`
	WifiBSSID         string         `gorm:"type:varchar(100);not null" json:"wifi_bssid"`
	LocationLat       float64        `gorm:"type:decimal(10,8);not null" json:"location_lat"`
	LocationLong      float64        `gorm:"type:decimal(11,8);not null" json:"location_long"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
}

func (s *AttendanceSession) BeforeCreate(tx *gorm.DB) error {
	if s.ID == "" {
		s.ID = utils.GenerateUUID()
	}
	return nil
}

func (s *AttendanceSession) IsActive() bool {
	return s.Status == SessionStatusActive
}

func (s *AttendanceSession) Complete() {
	s.Status = SessionStatusComplete
}

func (s *AttendanceSession) Cancel() {
	s.Status = SessionStatusCanceled
}
