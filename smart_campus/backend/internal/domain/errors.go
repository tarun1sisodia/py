package domain

import "errors"

// Common domain errors
var (
	// ErrNotFound is returned when a requested entity is not found
	ErrNotFound = errors.New("entity not found")

	// ErrDuplicate is returned when trying to create a duplicate entity
	ErrDuplicate = errors.New("entity already exists")

	// ErrInvalidInput is returned when the input data is invalid
	ErrInvalidInput = errors.New("invalid input data")

	// ErrUnauthorized is returned when the user is not authorized to perform an action
	ErrUnauthorized = errors.New("unauthorized")

	// ErrInvalidCredentials is returned when the provided credentials are invalid
	ErrInvalidCredentials = errors.New("invalid credentials")

	// ErrDeviceNotBound is returned when trying to use an unbound device
	ErrDeviceNotBound = errors.New("device not bound")

	// ErrInvalidOTP is returned when the provided OTP is invalid
	ErrInvalidOTP = errors.New("invalid OTP")

	// ErrExpiredOTP is returned when the OTP has expired
	ErrExpiredOTP = errors.New("OTP has expired")

	// ErrAlreadyExists is returned when attempting to create a resource that already exists
	ErrAlreadyExists = errors.New("resource already exists")

	// ErrForbidden is returned when the user is forbidden from performing an action
	ErrForbidden = errors.New("forbidden")

	// ErrAccountLocked is returned when the user account is locked
	ErrAccountLocked = errors.New("account locked")

	// ErrAccountDisabled is returned when the user account is disabled
	ErrAccountDisabled = errors.New("account disabled")

	// ErrSessionExpired is returned when the user session has expired
	ErrSessionExpired = errors.New("session expired")

	// ErrInvalidToken is returned when the provided token is invalid
	ErrInvalidToken = errors.New("invalid token")

	// ErrTokenExpired is returned when the provided token has expired
	ErrTokenExpired = errors.New("token expired")

	// ErrDeviceNotRegistered is returned when the device is not registered
	ErrDeviceNotRegistered = errors.New("device not registered")

	// ErrDeviceBlocked is returned when the device is blocked
	ErrDeviceBlocked = errors.New("device blocked")

	// ErrTooManyDevices is returned when the user has too many registered devices
	ErrTooManyDevices = errors.New("too many devices")

	// ErrInvalidLocation is returned when the location is invalid or out of range
	ErrInvalidLocation = errors.New("invalid location")

	// ErrInvalidWiFi is returned when the WiFi network is invalid or not allowed
	ErrInvalidWiFi = errors.New("invalid wifi network")

	// ErrSessionNotActive is returned when the attendance session is not active
	ErrSessionNotActive = errors.New("session not active")

	// ErrSessionClosed is returned when the attendance session is closed
	ErrSessionClosed = errors.New("session closed")

	// ErrAttendanceAlreadyMarked is returned when attendance is already marked
	ErrAttendanceAlreadyMarked = errors.New("attendance already marked")

	// ErrStudentNotEnrolled is returned when the student is not enrolled in the course
	ErrStudentNotEnrolled = errors.New("student not enrolled")

	// ErrTeacherNotAssigned is returned when the teacher is not assigned to the course
	ErrTeacherNotAssigned = errors.New("teacher not assigned")

	// ErrDepartmentNotActive is returned when the department is not active
	ErrDepartmentNotActive = errors.New("department not active")

	// ErrFacultyNotActive is returned when the faculty is not active
	ErrFacultyNotActive = errors.New("faculty not active")

	// ErrInvalidStatus is returned when the status is invalid
	ErrInvalidStatus = errors.New("invalid status")

	// ErrInvalidRole is returned when the role is invalid
	ErrInvalidRole = errors.New("invalid role")

	// ErrInvalidDate is returned when the date is invalid
	ErrInvalidDate = errors.New("invalid date")

	// ErrInvalidTime is returned when the time is invalid
	ErrInvalidTime = errors.New("invalid time")

	// ErrInvalidDuration is returned when the duration is invalid
	ErrInvalidDuration = errors.New("invalid duration")

	// ErrInvalidLimit is returned when the limit is invalid
	ErrInvalidLimit = errors.New("invalid limit")

	// ErrInvalidPage is returned when the page number is invalid
	ErrInvalidPage = errors.New("invalid page")

	// ErrInvalidSort is returned when the sort parameter is invalid
	ErrInvalidSort = errors.New("invalid sort")

	// ErrInvalidFilter is returned when the filter parameter is invalid
	ErrInvalidFilter = errors.New("invalid filter")

	// ErrInvalidSearch is returned when the search parameter is invalid
	ErrInvalidSearch = errors.New("invalid search")

	// ErrInvalidFormat is returned when the format is invalid
	ErrInvalidFormat = errors.New("invalid format")

	// ErrInvalidSize is returned when the size is invalid
	ErrInvalidSize = errors.New("invalid size")

	// ErrInvalidType is returned when the type is invalid
	ErrInvalidType = errors.New("invalid type")

	// ErrInvalidValue is returned when the value is invalid
	ErrInvalidValue = errors.New("invalid value")

	// ErrInvalidState is returned when the state is invalid
	ErrInvalidState = errors.New("invalid state")

	// ErrInvalidOperation is returned when an operation is invalid for the current state
	ErrInvalidOperation = errors.New("invalid operation for current state")

	// ErrInternalError is returned when an internal error occurs
	ErrInternalError = errors.New("internal error")

	// ErrDatabaseError is returned when a database error occurs
	ErrDatabaseError = errors.New("database error")

	// ErrNetworkError is returned when a network error occurs
	ErrNetworkError = errors.New("network error")

	// ErrTimeout is returned when an operation times out
	ErrTimeout = errors.New("timeout")

	// ErrCancelled is returned when an operation is cancelled
	ErrCancelled = errors.New("cancelled")

	// ErrBusy is returned when the system is busy
	ErrBusy = errors.New("system busy")

	// ErrUnavailable is returned when a service is unavailable
	ErrUnavailable = errors.New("service unavailable")

	// ErrMaintenance is returned when the system is under maintenance
	ErrMaintenance = errors.New("system under maintenance")
)

// Error wraps a domain error with additional context
type Error struct {
	// Original is the original error
	Original error
	// Code is the error code
	Code string
	// Message is a human-readable error message
	Message string
	// Details contains additional error details
	Details map[string]interface{}
}

// Error returns the error message
func (e *Error) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return e.Original.Error()
}

// Unwrap returns the original error
func (e *Error) Unwrap() error {
	return e.Original
}

// NewError creates a new domain error
func NewError(err error, code string, message string) *Error {
	return &Error{
		Original: err,
		Code:     code,
		Message:  message,
		Details:  make(map[string]interface{}),
	}
}

// WithDetails adds details to the error
func (e *Error) WithDetails(details map[string]interface{}) *Error {
	e.Details = details
	return e
}

// Is reports whether the error matches the target error
func (e *Error) Is(target error) bool {
	return errors.Is(e.Original, target)
}

// As attempts to convert the error to the target type
func (e *Error) As(target interface{}) bool {
	return errors.As(e.Original, target)
}
