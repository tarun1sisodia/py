import 'package:injectable/injectable.dart';
import 'package:sentry_flutter/sentry_flutter.dart';
import 'package:smart_campus/core/services/logger_service.dart';

@lazySingleton
class ErrorReportingService {
  final LoggerService _logger;

  ErrorReportingService(this._logger);

  Future<void> reportError(
    dynamic error,
    StackTrace stackTrace, {
    String? hint,
    Map<String, dynamic>? extras,
    SentryLevel level = SentryLevel.error,
  }) async {
    try {
      _logger.error(error.toString(), error, stackTrace);

      await Sentry.captureException(
        error,
        stackTrace: stackTrace,
        withScope: (scope) {
          scope.level = level;
          if (extras != null) {
            for (final entry in extras.entries) {
              scope.setExtra(entry.key, entry.value);
            }
          }
          if (hint != null) {
            scope.setTag('hint', hint);
          }
        },
      );
    } catch (e) {
      _logger.error('Failed to report error to Sentry', e);
    }
  }

  Future<void> addBreadcrumb({
    required String message,
    String? category,
    Map<String, dynamic>? data,
    SentryLevel level = SentryLevel.info,
  }) async {
    try {
      final breadcrumb = Breadcrumb(
        message: message,
        category: category,
        data: data,
        level: level,
        timestamp: DateTime.now(),
      );

      await Sentry.addBreadcrumb(breadcrumb);
    } catch (e) {
      _logger.error('Failed to add breadcrumb', e);
    }
  }

  Future<void> setUser({
    required String id,
    String? email,
    String? username,
    Map<String, dynamic>? data,
  }) async {
    try {
      await Sentry.configureScope((scope) {
        scope.setUser(SentryUser(
          id: id,
          email: email,
          username: username,
          data: data,
        ));
      });
    } catch (e) {
      _logger.error('Failed to set user context', e);
    }
  }

  Future<void> clearUser() async {
    try {
      await Sentry.configureScope((scope) => scope.setUser(null));
    } catch (e) {
      _logger.error('Failed to clear user context', e);
    }
  }

  Future<void> setTag(String key, String value) async {
    try {
      await Sentry.configureScope((scope) => scope.setTag(key, value));
    } catch (e) {
      _logger.error('Failed to set tag', e);
    }
  }

  Future<void> setExtra(String key, dynamic value) async {
    try {
      await Sentry.configureScope((scope) => scope.setExtra(key, value));
    } catch (e) {
      _logger.error('Failed to set extra', e);
    }
  }
}
