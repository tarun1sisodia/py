package models

import (
	"encoding/json"
	"time"

	"smart_attendance_backend/utils"

	"gorm.io/gorm"
)

type VerificationMethod string

const (
	VerificationMethodLocation VerificationMethod = "location"
	VerificationMethodWifi     VerificationMethod = "wifi"
	VerificationMethodBoth     VerificationMethod = "both"
)

type AttendanceRecord struct {
	ID                 string             `gorm:"type:varchar(36);primary_key" json:"id"`
	SessionID          string             `gorm:"type:varchar(36);not null" json:"session_id"`
	Session            AttendanceSession  `gorm:"foreignKey:SessionID" json:"session"`
	StudentID          string             `gorm:"type:varchar(36);not null" json:"student_id"`
	Student            User               `gorm:"foreignKey:StudentID" json:"student"`
	MarkedAt           time.Time          `gorm:"not null" json:"marked_at"`
	VerificationMethod VerificationMethod `gorm:"type:enum('location','wifi','both');not null" json:"verification_method"`
	DeviceInfo         json.RawMessage    `gorm:"type:json" json:"device_info"`
	LocationLat        *float64           `gorm:"type:decimal(10,8)" json:"location_lat,omitempty"`
	LocationLong       *float64           `gorm:"type:decimal(11,8)" json:"location_long,omitempty"`
	WifiSSID           *string            `gorm:"type:varchar(100)" json:"wifi_ssid,omitempty"`
	WifiBSSID          *string            `gorm:"type:varchar(100)" json:"wifi_bssid,omitempty"`
	CreatedAt          time.Time          `json:"created_at"`
	UpdatedAt          time.Time          `json:"updated_at"`
	DeletedAt          gorm.DeletedAt     `gorm:"index" json:"-"`
}

func (a *AttendanceRecord) BeforeCreate(tx *gorm.DB) error {
	if a.ID == "" {
		a.ID = utils.GenerateUUID()
	}
	if a.MarkedAt.IsZero() {
		a.MarkedAt = time.Now()
	}
	return nil
}

// DeviceInfoStruct represents the structure of device information
type DeviceInfoStruct struct {
	DeviceID         string `json:"device_id"`
	DeviceModel      string `json:"device_model"`
	DeviceBrand      string `json:"device_brand"`
	OSVersion        string `json:"os_version"`
	DeveloperMode    bool   `json:"developer_mode"`
	DeviceRegistered bool   `json:"device_registered"`
}

func (a *AttendanceRecord) SetDeviceInfo(info DeviceInfoStruct) error {
	data, err := json.Marshal(info)
	if err != nil {
		return err
	}
	a.DeviceInfo = data
	return nil
}

func (a *AttendanceRecord) GetDeviceInfo() (*DeviceInfoStruct, error) {
	if a.DeviceInfo == nil {
		return nil, nil
	}
	var info DeviceInfoStruct
	if err := json.Unmarshal(a.DeviceInfo, &info); err != nil {
		return nil, err
	}
	return &info, nil
}
