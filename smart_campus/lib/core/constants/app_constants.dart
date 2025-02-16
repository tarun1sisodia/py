class AppConstants {
  // App Info
  static const String appName = 'Smart Campus';
  static const String appVersion = '1.0.0';

  // Shared Preferences Keys
  static const String prefKeyAuthToken = 'auth_token';
  static const String prefKeyRefreshToken = 'refresh_token';
  static const String prefKeyUser = 'user';
  static const String prefKeyDeviceId = 'device_id';
  static const String prefKeyLastSync = 'last_sync';

  // API Endpoints
  static const String baseUrl = 'https://api.smartcampus.com';
  static const String apiVersion = 'v1';

  // API Paths
  static const String pathLogin = '/auth/login';
  static const String pathRegister = '/auth/register';
  static const String pathRefreshToken = '/auth/refresh';
  static const String pathLogout = '/auth/logout';
  static const String pathForgotPassword = '/auth/forgot-password';
  static const String pathResetPassword = '/auth/reset-password';
  static const String pathVerifyOTP = '/auth/verify-otp';
  static const String pathChangePassword = '/auth/change-password';

  // Session Paths
  static const String pathCreateSession = '/sessions/create';
  static const String pathActiveSessions = '/sessions/active';
  static const String pathSessionDetails = '/sessions/details';
  static const String pathMarkAttendance = '/attendance/mark';
  static const String pathAttendanceHistory = '/attendance/history';

  // Timeouts
  static const int connectionTimeout = 30000; // 30 seconds
  static const int receiveTimeout = 30000; // 30 seconds

  // Validation
  static const int minPasswordLength = 8;
  static const int otpLength = 6;
  static const int maxLoginAttempts = 3;

  // Location
  static const int locationTimeout = 10000; // 10 seconds
  static const int locationInterval = 1000; // 1 second
  static const double defaultLocationAccuracy = 50.0; // meters
  static const int defaultLocationRadius = 100; // meters

  // Session
  static const int defaultSessionDuration = 300; // 5 minutes
  static const int minimumSessionDuration = 60; // 1 minute
  static const int maximumSessionDuration = 3600; // 1 hour

  // UI Constants
  static const double defaultPadding = 16.0;
  static const double defaultMargin = 16.0;
  static const double defaultRadius = 8.0;
  static const double defaultElevation = 2.0;

  // Animation Durations
  static const Duration shortAnimationDuration = Duration(milliseconds: 200);
  static const Duration mediumAnimationDuration = Duration(milliseconds: 500);
  static const Duration longAnimationDuration = Duration(milliseconds: 800);

  // Error Messages
  static const String errorNoInternet = 'No internet connection';
  static const String errorTimeout = 'Request timeout';
  static const String errorUnauthorized = 'Unauthorized access';
  static const String errorServer = 'Server error';
  static const String errorUnknown = 'Unknown error occurred';
  static const String errorInvalidCredentials = 'Invalid credentials';
  static const String errorWeakPassword = 'Password is too weak';
  static const String errorInvalidEmail = 'Invalid email address';
  static const String errorInvalidOTP = 'Invalid OTP';
  static const String errorLocationPermission = 'Location permission denied';
  static const String errorLocationDisabled = 'Location services are disabled';
  static const String errorWifiDisabled = 'WiFi is disabled';
  static const String errorDeviceNotBound = 'Device not bound';
}
