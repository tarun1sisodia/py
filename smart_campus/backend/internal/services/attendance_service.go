package services

import (
	"context"
	"fmt"
	"time"
)

type AttendanceRecord struct {
	ID                 string
	SessionID          string
	StudentID          string
	MarkedAt           time.Time
	WifiSSID           string
	WifiBSSID          string
	LocationLatitude   float64
	LocationLongitude  float64
	DeviceID           string
	VerificationStatus string
	RejectionReason    string
}

type AttendanceService struct {
	attendanceRepo     AttendanceRepository
	sessionService     *SessionService
	deviceVerification *DeviceVerificationService
	firebaseAuth       *FirebaseAuthService
}

func NewAttendanceService(
	attendanceRepo AttendanceRepository,
	sessionService *SessionService,
	deviceVerification *DeviceVerificationService,
	firebaseAuth *FirebaseAuthService,
) *AttendanceService {
	return &AttendanceService{
		attendanceRepo:     attendanceRepo,
		sessionService:     sessionService,
		deviceVerification: deviceVerification,
		firebaseAuth:       firebaseAuth,
	}
}

// MarkAttendance marks attendance for a student in a session
func (s *AttendanceService) MarkAttendance(ctx context.Context, record *AttendanceRecord) error {
	// Verify session is active
	if err := s.sessionService.ValidateAttendanceWindow(ctx, record.SessionID); err != nil {
		return fmt.Errorf("invalid session: %v", err)
	}

	// Verify student through Firebase
	user, err := s.firebaseAuth.GetUser(ctx, record.StudentID)
	if err != nil {
		return fmt.Errorf("error verifying student: %v", err)
	}

	// Verify student role
	claims, err := s.firebaseAuth.VerifyIDToken(ctx, user.CustomClaims["role"].(string))
	if err != nil || claims.Claims["role"] != "student" {
		return fmt.Errorf("unauthorized: user is not a student")
	}

	// Verify device
	verificationResult, err := s.deviceVerification.VerifyDeviceForAttendance(ctx, VerificationParams{
		UserID:   record.StudentID,
		DeviceID: record.DeviceID,
		Location: &Location{
			Latitude:  record.LocationLatitude,
			Longitude: record.LocationLongitude,
		},
		WifiInfo: &WifiInfo{
			SSID:  record.WifiSSID,
			BSSID: record.WifiBSSID,
		},
		FirebaseUID: record.StudentID,
	})
	if err != nil {
		return fmt.Errorf("device verification failed: %v", err)
	}

	if !verificationResult.IsValid {
		record.VerificationStatus = "rejected"
		record.RejectionReason = verificationResult.Reason
	} else {
		record.VerificationStatus = "verified"
	}

	// Check for duplicate attendance
	exists, err := s.attendanceRepo.ExistsForSessionAndStudent(ctx, record.SessionID, record.StudentID)
	if err != nil {
		return fmt.Errorf("error checking duplicate attendance: %v", err)
	}
	if exists {
		return fmt.Errorf("attendance already marked for this session")
	}

	// Save attendance record
	return s.attendanceRepo.Create(ctx, record)
}

// GetAttendanceReport generates an attendance report
func (s *AttendanceService) GetAttendanceReport(ctx context.Context, filters AttendanceFilters) ([]*AttendanceRecord, error) {
	return s.attendanceRepo.List(ctx, filters)
}

// UpdateAttendanceStatus updates the verification status of an attendance record
func (s *AttendanceService) UpdateAttendanceStatus(ctx context.Context, recordID string, status string, reason string) error {
	record, err := s.attendanceRepo.FindByID(ctx, recordID)
	if err != nil {
		return fmt.Errorf("error finding attendance record: %v", err)
	}

	record.VerificationStatus = status
	record.RejectionReason = reason

	return s.attendanceRepo.Update(ctx, record)
}
