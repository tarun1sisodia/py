class AppError implements Exception {
  final String message;
  final String? code;
  final dynamic details;

  AppError(this.message, {this.code, this.details});

  @override
  String toString() =>
      'AppError: $message${code != null ? ' (Code: $code)' : ''}';

  factory AppError.network(String message) => AppError(
        message,
        code: 'NETWORK_ERROR',
      );

  factory AppError.server(String message) => AppError(
        message,
        code: 'SERVER_ERROR',
      );

  factory AppError.unauthorized() => AppError(
        'Unauthorized access',
        code: 'UNAUTHORIZED',
      );

  factory AppError.validation(String message) => AppError(
        message,
        code: 'VALIDATION_ERROR',
      );

  factory AppError.deviceBinding(String message) => AppError(
        message,
        code: 'DEVICE_BINDING_ERROR',
      );

  factory AppError.locationService(String message) => AppError(
        message,
        code: 'LOCATION_SERVICE_ERROR',
      );

  factory AppError.wifiService(String message) => AppError(
        message,
        code: 'WIFI_SERVICE_ERROR',
      );

  factory AppError.attendanceVerification(String message) => AppError(
        message,
        code: 'ATTENDANCE_VERIFICATION_ERROR',
      );
}

class AttendanceVerificationError extends AppError {
  AttendanceVerificationError(String message) : super(message);
}
