package models

import (
	"time"
)

type DeviceBinding struct {
	ID            string     `json:"id"`
	UserID        string     `json:"user_id"`
	DeviceID      string     `json:"device_id"`
	DeviceName    string     `json:"device_name"`
	DeviceModel   string     `json:"device_model"`
	IsActive      bool       `json:"is_active"`
	IsBlacklisted bool       `json:"is_blacklisted"`
	BoundAt       time.Time  `json:"bound_at"`
	LastUsedAt    *time.Time `json:"last_used_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type DeviceRepository interface {
	Create(binding *DeviceBinding) error
	GetByID(id string) (*DeviceBinding, error)
	GetByUserID(userID string) ([]*DeviceBinding, error)
	GetByDeviceID(deviceID string) (*DeviceBinding, error)
	Update(binding *DeviceBinding) error
	Delete(id string) error
	Deactivate(id string) error
	Blacklist(deviceID string) error
	RemoveFromBlacklist(deviceID string) error
	IsDeviceRegistered(userID, deviceID string) (bool, error)
	IsDeviceBlacklisted(deviceID string) (bool, error)
	UpdateLastUsed(id string) error
}
