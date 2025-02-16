package services

import (
	"context"
	"fmt"
	"smart_campus/internal/domain/entities"
	"smart_campus/internal/domain/repositories"
)

// DeviceService handles device-related operations
type DeviceService struct {
	deviceRepo repositories.DeviceBindingRepository
}

// NewDeviceService creates a new device service
func NewDeviceService(deviceRepo repositories.DeviceBindingRepository) *DeviceService {
	return &DeviceService{
		deviceRepo: deviceRepo,
	}
}

// BindDevice binds a device to a user
func (s *DeviceService) BindDevice(userID, deviceID, deviceName, deviceModel string) error {
	binding := entities.NewDeviceBinding(userID, deviceID, deviceName, deviceModel)
	return s.deviceRepo.Create(context.Background(), binding)
}

// VerifyDevice verifies if a device is bound to a user
func (s *DeviceService) VerifyDevice(userID, deviceID string) (bool, error) {
	binding, err := s.deviceRepo.FindByUserAndDeviceID(context.Background(), userID, deviceID)
	if err != nil {
		if err == repositories.ErrNotFound {
			return false, nil
		}
		return false, fmt.Errorf("error verifying device: %w", err)
	}

	return binding.IsValid(), nil
}

// UnbindDevice removes a device binding
func (s *DeviceService) UnbindDevice(userID, deviceID string) error {
	binding, err := s.deviceRepo.FindByUserAndDeviceID(context.Background(), userID, deviceID)
	if err != nil {
		return fmt.Errorf("error finding device binding: %w", err)
	}

	binding.Deactivate()
	return s.deviceRepo.Update(context.Background(), binding)
}

// ListDevices returns a list of devices bound to a user
func (s *DeviceService) ListDevices(userID string) ([]*entities.DeviceBinding, error) {
	return s.deviceRepo.FindByUser(context.Background(), userID)
}

// GetDevice returns details of a specific device
func (s *DeviceService) GetDevice(userID, deviceID string) (*entities.DeviceBinding, error) {
	binding, err := s.deviceRepo.FindByUserAndDeviceID(context.Background(), userID, deviceID)
	if err != nil {
		return nil, fmt.Errorf("error finding device: %w", err)
	}
	return binding, nil
}

// UpdateLastUsed updates the last used timestamp of a device
func (s *DeviceService) UpdateLastUsed(userID, deviceID string) error {
	binding, err := s.deviceRepo.FindByUserAndDeviceID(context.Background(), userID, deviceID)
	if err != nil {
		return fmt.Errorf("error finding device: %w", err)
	}

	binding.UpdateLastUsed()
	return s.deviceRepo.Update(context.Background(), binding)
}
