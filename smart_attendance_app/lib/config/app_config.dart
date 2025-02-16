class AppConfig {
  static const String appName = 'Smart Attendance';
  static const String apiBaseUrl =
      'http://localhost:8080'; // Change this for production

  // API Endpoints
  static const String teacherRegister = '/auth/register/teacher';
  static const String studentRegister = '/auth/register/student';
  static const String verifyOtp = '/auth/verify-otp';
  static const String teacherLogin = '/auth/login/teacher';
  static const String studentLogin = '/auth/login/student';
  static const String resetPassword = '/auth/reset-password';

  // Session Endpoints
  static const String startSession = '/sessions/start';
  static const String activeSession = '/sessions/active';
  static const String endSession = '/sessions/end';

  // Attendance Endpoints
  static const String markAttendance = '/attendance/mark';
  static const String attendanceStatus = '/attendance/status';

  // Security Endpoints
  static const String checkDeveloperMode = '/security/check-developer-mode';
  static const String checkDeviceBinding = '/security/check-device-binding';
  static const String reportFraud = '/security/report-fraud';

  // Timeouts
  static const int connectionTimeout = 30000; // 30 seconds
  static const int receiveTimeout = 30000; // 30 seconds

  // Location Settings
  static const double defaultLatitude =
      0.0; // Set your institution's default latitude
  static const double defaultLongitude =
      0.0; // Set your institution's default longitude
  static const int locationUpdateInterval = 10000; // 10 seconds

  // Cache Settings
  static const int cacheValidityDuration = 3600; // 1 hour in seconds

  // Session Settings
  static const List<int> sessionDurations = [
    30,
    60,
    180,
  ]; // Duration options in seconds
}
