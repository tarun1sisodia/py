package services

import (
	"context"
	"fmt"
	"smart_campus/internal/models"
	"time"
)

type AttendanceSession struct {
	ID                string
	TeacherID         string
	CourseID          string
	SessionDate       time.Time
	StartTime         time.Time
	EndTime           time.Time
	WifiSSID          string
	WifiBSSID         string
	LocationLatitude  float64
	LocationLongitude float64
	LocationRadius    int
	Status            string
}

type SessionService struct {
	sessionRepo        models.SessionRepository
	attendanceRepo     models.AttendanceRepository
	firebaseAuth       *FirebaseAuthService
	deviceVerification *DeviceVerificationService
}

func NewSessionService(
	sessionRepo models.SessionRepository,
	attendanceRepo models.AttendanceRepository,
	firebaseAuth *FirebaseAuthService,
	deviceVerification *DeviceVerificationService,
) *SessionService {
	return &SessionService{
		sessionRepo:        sessionRepo,
		attendanceRepo:     attendanceRepo,
		firebaseAuth:       firebaseAuth,
		deviceVerification: deviceVerification,
	}
}

// CreateSession creates a new attendance session
func (s *SessionService) CreateSession(ctx context.Context, session *AttendanceSession) error {
	// Verify teacher through Firebase
	user, err := s.firebaseAuth.GetUser(ctx, session.TeacherID)
	if err != nil {
		return fmt.Errorf("error verifying teacher: %v", err)
	}

	// Verify teacher role through custom claims
	claims, err := s.firebaseAuth.VerifyIDToken(ctx, user.CustomClaims["role"].(string))
	if err != nil || claims.Claims["role"] != "teacher" {
		return fmt.Errorf("unauthorized: user is not a teacher")
	}

	// Validate session parameters
	if err := s.validateSessionParams(session); err != nil {
		return err
	}

	// Create session
	return s.sessionRepo.Create(ctx, session)
}

// GetActiveSession gets the currently active session for a course
func (s *SessionService) GetActiveSession(ctx context.Context, courseID string) (*AttendanceSession, error) {
	return s.sessionRepo.FindActiveByCourse(ctx, courseID)
}

// EndSession ends an active attendance session
func (s *SessionService) EndSession(ctx context.Context, sessionID string) error {
	session, err := s.sessionRepo.FindByID(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("error finding session: %v", err)
	}

	session.Status = "completed"
	session.EndTime = time.Now()

	return s.sessionRepo.Update(ctx, session)
}

// ValidateAttendanceWindow checks if attendance can be marked for a session
func (s *SessionService) ValidateAttendanceWindow(ctx context.Context, sessionID string) error {
	session, err := s.sessionRepo.FindByID(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("error finding session: %v", err)
	}

	if session.Status != "active" {
		return fmt.Errorf("session is not active")
	}

	now := time.Now()
	if now.Before(session.StartTime) || now.After(session.EndTime) {
		return fmt.Errorf("attendance window is closed")
	}

	return nil
}

func (s *SessionService) validateSessionParams(session *AttendanceSession) error {
	// Validate time window
	if session.EndTime.Before(session.StartTime) {
		return fmt.Errorf("end time cannot be before start time")
	}

	// Validate location
	if session.LocationRadius <= 0 {
		return fmt.Errorf("invalid location radius")
	}

	// Validate WiFi credentials
	if session.WifiSSID == "" || session.WifiBSSID == "" {
		return fmt.Errorf("WiFi credentials are required")
	}

	return nil
}

func (s *SessionService) GetSessionById(id string) (*models.Session, error) {
	return s.sessionRepo.GetByID(id)
}

func (s *SessionService) CancelSession(id string) error {
	return s.sessionRepo.CancelSession(id)
}

func (s *SessionService) MarkAttendance(record *models.AttendanceRecord) error {
	// Get session details
	session, err := s.sessionRepo.GetByID(record.SessionID)
	if err != nil {
		return fmt.Errorf("error getting session: %v", err)
	}

	// Check if session is active
	if session.Status != models.SessionStatusActive {
		return fmt.Errorf("session is not active")
	}

	// Check if attendance is within session time
	now := time.Now()
	if now.Before(session.StartTime) || now.After(session.EndTime) {
		return fmt.Errorf("attendance can only be marked during session time")
	}

	// Create attendance record
	return s.attendanceRepo.Create(record)
}

func (s *SessionService) VerifyAttendance(id string) error {
	return s.attendanceRepo.VerifyAttendance(id)
}

func (s *SessionService) RejectAttendance(id string, reason string) error {
	return s.attendanceRepo.RejectAttendance(id, reason)
}

func (s *SessionService) GetSessionAttendance(sessionID string) ([]*models.AttendanceRecord, error) {
	return s.attendanceRepo.GetBySession(sessionID)
}

func (s *SessionService) GetStudentAttendance(studentID string) ([]*models.AttendanceRecord, error) {
	return s.attendanceRepo.GetByStudent(studentID)
}

func (s *SessionService) GetAttendanceStatistics(
	studentID string,
	courseID string,
	startDate time.Time,
	endDate time.Time,
) (*models.AttendanceStatistics, error) {
	return s.attendanceRepo.GetAttendanceStatistics(studentID, courseID, startDate, endDate)
}

func (s *SessionService) GetSessionsByTeacher(teacherID string) ([]*models.Session, error) {
	return s.sessionRepo.GetByTeacher(teacherID)
}

func (s *SessionService) GetSessionsByCourse(courseID string) ([]*models.Session, error) {
	return s.sessionRepo.GetByCourse(courseID)
}

func (s *SessionService) UpdateSession(session *models.Session) error {
	return s.sessionRepo.Update(session)
}
