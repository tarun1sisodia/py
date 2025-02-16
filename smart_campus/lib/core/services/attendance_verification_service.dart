import 'dart:async';
import 'package:injectable/injectable.dart';
import 'package:rxdart/rxdart.dart';
import 'package:smart_campus/core/errors/app_error.dart';
import 'package:smart_campus/core/services/location_service.dart';
import 'package:smart_campus/core/services/logger_service.dart';
import 'package:smart_campus/core/services/wifi_service.dart';
import 'package:smart_campus/core/services/device_binding_service.dart';

@lazySingleton
class AttendanceVerificationService {
  final LocationService _locationService;
  final WifiService _wifiService;
  final DeviceBindingService _deviceBindingService;
  final LoggerService _logger;

  AttendanceVerificationService(
    this._locationService,
    this._wifiService,
    this._deviceBindingService,
    this._logger,
  );

  /// Verifies location, WiFi, and device binding for attendance marking
  Future<bool> verifyAttendance({
    required String userId,
    required double sessionLatitude,
    required double sessionLongitude,
    required int allowedRadius,
    required String sessionSSID,
    required String sessionBSSID,
    Duration verificationTimeout = const Duration(seconds: 15),
  }) async {
    try {
      // First check device binding
      final isDeviceBound =
          await _deviceBindingService.verifyDeviceBinding(userId);
      if (!isDeviceBound) {
        _logger.warning('Attendance verification failed: Device not bound');
        return false;
      }

      final locationFuture = _locationService.verifyAttendanceLocation(
        sessionLatitude: sessionLatitude,
        sessionLongitude: sessionLongitude,
        allowedRadius: allowedRadius,
        timeLimit: verificationTimeout,
      );

      final wifiFuture = _wifiService.verifyAttendanceWifi(
        sessionSSID: sessionSSID,
        sessionBSSID: sessionBSSID,
        timeout: verificationTimeout,
      );

      final results = await Future.wait([locationFuture, wifiFuture]);
      final isLocationValid = results[0];
      final isWifiValid = results[1];

      if (!isLocationValid) {
        _logger.warning('Attendance verification failed: Invalid location');
      }
      if (!isWifiValid) {
        _logger
            .warning('Attendance verification failed: Invalid WiFi connection');
      }

      return isLocationValid && isWifiValid;
    } catch (e) {
      _logger.error('Error during attendance verification', e);
      throw AppError.attendanceVerification('Verification failed: $e');
    }
  }

  /// Monitors location, WiFi, and device binding continuously for attendance session
  Stream<bool> monitorAttendance({
    required String userId,
    required double sessionLatitude,
    required double sessionLongitude,
    required int allowedRadius,
    required String sessionSSID,
    required String sessionBSSID,
  }) {
    final locationStream = _locationService.monitorAttendanceLocation(
      sessionLatitude: sessionLatitude,
      sessionLongitude: sessionLongitude,
      allowedRadius: allowedRadius,
    );

    final wifiStream = _wifiService.monitorAttendanceWifi(
      sessionSSID: sessionSSID,
      sessionBSSID: sessionBSSID,
    );

    // Create a periodic stream for device binding verification
    final deviceBindingStream = Stream.periodic(
      const Duration(seconds: 30),
      (_) => _deviceBindingService.verifyDeviceBinding(userId),
    ).asyncMap((future) => future);

    return Rx.combineLatest3(
      locationStream,
      wifiStream,
      deviceBindingStream,
      (bool isLocationValid, bool isWifiValid, bool isDeviceBound) {
        final isValid = isLocationValid && isWifiValid && isDeviceBound;

        if (!isValid) {
          _logger.warning(
            'Attendance validation failed. '
            'Location valid: $isLocationValid, '
            'WiFi valid: $isWifiValid, '
            'Device bound: $isDeviceBound',
          );
        }

        return isValid;
      },
    );
  }

  /// Stops all monitoring streams
  void dispose() {
    // Add any cleanup logic if needed
  }
}
