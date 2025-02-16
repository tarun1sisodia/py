package entities

// Status represents the status of an entity
type Status string

const (
	// StatusActive represents an active entity
	StatusActive Status = "active"
	// StatusInactive represents an inactive entity
	StatusInactive Status = "inactive"
	// StatusPending represents a pending entity
	StatusPending Status = "pending"
	// StatusBlocked represents a blocked entity
	StatusBlocked Status = "blocked"
)

// UserRole represents the role of a user
type UserRole string

const (
	// UserRoleAdmin represents an admin user
	UserRoleAdmin UserRole = "admin"
	// UserRoleTeacher represents a teacher user
	UserRoleTeacher UserRole = "teacher"
	// UserRoleStudent represents a student user
	UserRoleStudent UserRole = "student"
	// UserRoleStaff represents a staff user
	UserRoleStaff UserRole = "staff"
)

// OTPPurpose represents the purpose of an OTP
type OTPPurpose string

const (
	// OTPPurposeRegistration represents OTP for registration
	OTPPurposeRegistration OTPPurpose = "registration"
	// OTPPurposePasswordReset represents OTP for password reset
	OTPPurposePasswordReset OTPPurpose = "password_reset"
	// OTPPurposeDeviceBinding represents OTP for device binding
	OTPPurposeDeviceBinding OTPPurpose = "device_binding"
	// OTPPurposeEmailVerification represents OTP for email verification
	OTPPurposeEmailVerification OTPPurpose = "email_verification"
)

// OTPStatus represents the status of an OTP verification
type OTPStatus string

const (
	// OTPStatusPending represents a pending OTP verification
	OTPStatusPending OTPStatus = "pending"
	// OTPStatusVerified represents a verified OTP
	OTPStatusVerified OTPStatus = "verified"
	// OTPStatusExpired represents an expired OTP
	OTPStatusExpired OTPStatus = "expired"
	// OTPStatusInvalid represents an invalid OTP
	OTPStatusInvalid OTPStatus = "invalid"
)

// VerificationStatus represents the status of a verification
type VerificationStatus string

const (
	// VerificationStatusPending represents a pending verification
	VerificationStatusPending VerificationStatus = "pending"
	// VerificationStatusVerified represents a verified status
	VerificationStatusVerified VerificationStatus = "verified"
	// VerificationStatusRejected represents a rejected verification
	VerificationStatusRejected VerificationStatus = "rejected"
)

// AttendanceStatus represents the status of an attendance record
type AttendanceStatus string

const (
	// AttendanceStatusPresent represents present status
	AttendanceStatusPresent AttendanceStatus = "present"
	// AttendanceStatusAbsent represents absent status
	AttendanceStatusAbsent AttendanceStatus = "absent"
	// AttendanceStatusLate represents late status
	AttendanceStatusLate AttendanceStatus = "late"
)

// SessionStatus represents the status of a session
type SessionStatus string

const (
	// SessionStatusActive represents an active session
	SessionStatusActive SessionStatus = "active"
	// SessionStatusScheduled represents a scheduled session
	SessionStatusScheduled SessionStatus = "scheduled"
	// SessionStatusComplete represents a completed session
	SessionStatusComplete SessionStatus = "complete"
	// SessionStatusCancelled represents a cancelled session
	SessionStatusCancelled SessionStatus = "cancelled"
)

// TeacherStatus represents the status of a teacher
type TeacherStatus string

const (
	// TeacherStatusActive represents an active teacher
	TeacherStatusActive TeacherStatus = "active"
	// TeacherStatusInactive represents an inactive teacher
	TeacherStatusInactive TeacherStatus = "inactive"
	// TeacherStatusOnLeave represents a teacher on leave
	TeacherStatusOnLeave TeacherStatus = "on_leave"
	// TeacherStatusTerminated represents a terminated teacher
	TeacherStatusTerminated TeacherStatus = "terminated"
)

// StudentStatus represents the status of a student
type StudentStatus string

const (
	// StudentStatusActive represents an active student
	StudentStatusActive StudentStatus = "active"
	// StudentStatusInactive represents an inactive student
	StudentStatusInactive StudentStatus = "inactive"
	// StudentStatusGraduated represents a graduated student
	StudentStatusGraduated StudentStatus = "graduated"
	// StudentStatusSuspended represents a suspended student
	StudentStatusSuspended StudentStatus = "suspended"
	// StudentStatusWithdrawn represents a withdrawn student
	StudentStatusWithdrawn StudentStatus = "withdrawn"
)

// DepartmentStatus represents the status of a department
type DepartmentStatus string

const (
	// DepartmentStatusActive represents an active department
	DepartmentStatusActive DepartmentStatus = "active"
	// DepartmentStatusInactive represents an inactive department
	DepartmentStatusInactive DepartmentStatus = "inactive"
	// DepartmentStatusMerged represents a merged department
	DepartmentStatusMerged DepartmentStatus = "merged"
)

// FacultyStatus represents the status of a faculty
type FacultyStatus string

const (
	// FacultyStatusActive represents an active faculty
	FacultyStatusActive FacultyStatus = "active"
	// FacultyStatusInactive represents an inactive faculty
	FacultyStatusInactive FacultyStatus = "inactive"
	// FacultyStatusMerged represents a merged faculty
	FacultyStatusMerged FacultyStatus = "merged"
)

// AuthLogType represents the type of authentication log
type AuthLogType string

const (
	// AuthLogTypeLogin represents a login event
	AuthLogTypeLogin AuthLogType = "login"
	// AuthLogTypeLogout represents a logout event
	AuthLogTypeLogout AuthLogType = "logout"
	// AuthLogTypeOTPRequest represents an OTP request event
	AuthLogTypeOTPRequest AuthLogType = "otp_request"
	// AuthLogTypeOTPVerify represents an OTP verification event
	AuthLogTypeOTPVerify AuthLogType = "otp_verify"
	// AuthLogTypeTokenRefresh represents a token refresh event
	AuthLogTypeTokenRefresh AuthLogType = "token_refresh"
	// AuthLogTypeDeviceBind represents a device binding event
	AuthLogTypeDeviceBind AuthLogType = "device_bind"
	// AuthLogTypeDeviceRevoke represents a device revocation event
	AuthLogTypeDeviceRevoke AuthLogType = "device_revoke"
	// AuthLogTypePasswordChange represents a password change event
	AuthLogTypePasswordChange AuthLogType = "password_change"
	// AuthLogTypePasswordReset represents a password reset event
	AuthLogTypePasswordReset AuthLogType = "password_reset"
	// AuthLogTypeDeviceUnbind represents a device unbinding event
	AuthLogTypeDeviceUnbind AuthLogType = "device_unbind"
)

// AuthLogStatus represents the status of an authentication log
type AuthLogStatus string

const (
	// AuthLogStatusSuccess represents a successful authentication event
	AuthLogStatusSuccess AuthLogStatus = "success"
	// AuthLogStatusFailed represents a failed authentication event
	AuthLogStatusFailed AuthLogStatus = "failed"
)
