/// Custom error class for application-specific errors
class AppError implements Exception {
  final String message;
  final dynamic cause;
  final StackTrace? stackTrace;

  AppError(
    this.message, [
    this.cause,
    this.stackTrace,
  ]);

  @override
  String toString() {
    if (cause != null) {
      return '$message (Cause: $cause)';
    }
    return message;
  }
}
