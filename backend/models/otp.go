package models

import (
	"time"

	"smart_attendance_backend/utils"

	"gorm.io/gorm"
)

type OTPVerification struct {
	ID        string         `gorm:"type:varchar(36);primary_key" json:"id"`
	UserID    string         `gorm:"type:varchar(36);not null" json:"user_id"`
	User      User           `gorm:"foreignKey:UserID" json:"user"`
	OTPCode   string         `gorm:"type:varchar(10);not null" json:"-"`
	ExpiresAt time.Time      `gorm:"not null" json:"expires_at"`
	Verified  bool           `gorm:"default:false" json:"verified"`
	Attempts  int            `gorm:"default:0" json:"attempts"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (o *OTPVerification) BeforeCreate(tx *gorm.DB) error {
	if o.ID == "" {
		o.ID = utils.GenerateUUID()
	}
	return nil
}

func (o *OTPVerification) IsExpired() bool {
	return time.Now().After(o.ExpiresAt)
}

func (o *OTPVerification) IncrementAttempts() {
	o.Attempts++
}

func (o *OTPVerification) HasExceededMaxAttempts(maxAttempts int) bool {
	return o.Attempts >= maxAttempts
}

func (o *OTPVerification) MarkAsVerified() {
	o.Verified = true
}
