package repositories

import (
	"context"
	"smart_campus/internal/domain/entities"
)

type DeviceBindingRepository interface {
	FindByUserAndDevice(ctx context.Context, userID, deviceID string) (*entities.DeviceBinding, error)
	IsDeviceRegistered(userID, deviceID string) (bool, error)
	IsDeviceBlacklisted(deviceID string) (bool, error)
	Create(ctx context.Context, binding *entities.DeviceBinding) error
	Update(ctx context.Context, binding *entities.DeviceBinding) error
	Delete(ctx context.Context, id string) error
}
