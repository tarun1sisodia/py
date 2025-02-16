package entities

type OTPPurpose string

const (
	OTPPurposeRegistration      OTPPurpose = "registration"
	OTPPurposePasswordReset     OTPPurpose = "password_reset"
	OTPPurposeDeviceBinding     OTPPurpose = "device_binding"
	OTPPurposePhoneVerification OTPPurpose = "phone_verification"
)
