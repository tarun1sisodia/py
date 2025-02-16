package services

import (
	"context"
	"fmt"
	"time"
)

type Location struct {
	Latitude  float64
	Longitude float64
}

type WifiInfo struct {
	SSID  string
	BSSID string
}

type VerificationParams struct {
	UserID      string
	DeviceID    string
	Location    *Location
	WifiInfo    *WifiInfo
	FirebaseUID string
}

type VerificationResult struct {
	IsValid   bool
	Timestamp time.Time
	Reason    string
}

type DeviceVerificationService struct {
	deviceRepo   DeviceBindingRepository
	firebaseAuth *FirebaseAuthService
}

func NewDeviceVerificationService(
	deviceRepo DeviceBindingRepository,
	firebaseAuth *FirebaseAuthService,
) *DeviceVerificationService {
	return &DeviceVerificationService{
		deviceRepo:   deviceRepo,
		firebaseAuth: firebaseAuth,
	}
}

// VerifyDeviceForAttendance checks if a device is valid for attendance
func (s *DeviceVerificationService) VerifyDeviceForAttendance(ctx context.Context, params VerificationParams) (*VerificationResult, error) {
	// Verify Firebase user
	user, err := s.firebaseAuth.GetUser(ctx, params.FirebaseUID)
	if err != nil {
		return nil, fmt.Errorf("error verifying Firebase user: %v", err)
	}

	// Check if phone is verified
	if !user.PhoneNumber || !user.PhoneNumberVerified {
		return &VerificationResult{
			IsValid:   false,
			Timestamp: time.Now(),
			Reason:    "phone number not verified",
		}, nil
	}

	// Check device binding
	isRegistered, err := s.deviceRepo.IsDeviceRegistered(params.UserID, params.DeviceID)
	if err != nil {
		return nil, fmt.Errorf("error checking device registration: %v", err)
	}
	if !isRegistered {
		return &VerificationResult{
			IsValid:   false,
			Timestamp: time.Now(),
			Reason:    "device not registered to user",
		}, nil
	}

	// Check if device is blacklisted
	isBlacklisted, err := s.deviceRepo.IsDeviceBlacklisted(params.DeviceID)
	if err != nil {
		return nil, fmt.Errorf("error checking device blacklist: %v", err)
	}
	if isBlacklisted {
		return &VerificationResult{
			IsValid:   false,
			Timestamp: time.Now(),
			Reason:    "device is blacklisted",
		}, nil
	}

	// Verify location if provided
	if params.Location != nil {
		if !s.isLocationValid(params.Location) {
			return &VerificationResult{
				IsValid:   false,
				Timestamp: time.Now(),
				Reason:    "invalid location",
			}, nil
		}
	}

	// Verify WiFi if provided
	if params.WifiInfo != nil {
		if !s.isWifiValid(params.WifiInfo) {
			return &VerificationResult{
				IsValid:   false,
				Timestamp: time.Now(),
				Reason:    "invalid wifi network",
			}, nil
		}
	}

	return &VerificationResult{
		IsValid:   true,
		Timestamp: time.Now(),
	}, nil
}

func (s *DeviceVerificationService) isLocationValid(loc *Location) bool {
	// TODO: Implement location validation logic
	// Check if location is within campus bounds
	// Check for location spoofing
	return true
}

func (s *DeviceVerificationService) isWifiValid(wifi *WifiInfo) bool {
	// TODO: Implement WiFi validation logic
	// Check if WiFi network is authorized
	// Check for WiFi spoofing
	return true
}
