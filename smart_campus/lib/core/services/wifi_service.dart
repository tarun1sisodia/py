import 'package:injectable/injectable.dart';
import 'package:network_info_plus/network_info_plus.dart';
import 'package:smart_campus/core/errors/app_error.dart';
import 'package:smart_campus/core/services/logger_service.dart';
import 'dart:async';

@lazySingleton
class WifiService {
  final NetworkInfo _networkInfo;
  final LoggerService _logger;

  WifiService(this._logger) : _networkInfo = NetworkInfo();

  Future<bool> isWifiEnabled() async {
    try {
      final wifiName = await _networkInfo.getWifiName();
      return wifiName != null;
    } catch (e) {
      _logger.error('Error checking WiFi status', e);
      throw AppError.wifiService('Error checking WiFi status: $e');
    }
  }

  Future<String?> getWifiName() async {
    try {
      final wifiName = await _networkInfo.getWifiName();
      _logger.info('WiFi name: $wifiName');
      // Remove quotes if present
      return wifiName?.replaceAll('"', '');
    } catch (e) {
      _logger.error('Error getting WiFi name', e);
      throw AppError.wifiService('Error getting WiFi name: $e');
    }
  }

  Future<String?> getWifiBSSID() async {
    try {
      final bssid = await _networkInfo.getWifiBSSID();
      _logger.info('WiFi BSSID: $bssid');
      return bssid;
    } catch (e) {
      _logger.error('Error getting WiFi BSSID', e);
      throw AppError.wifiService('Error getting WiFi BSSID: $e');
    }
  }

  Future<String?> getWifiIP() async {
    try {
      final ip = await _networkInfo.getWifiIP();
      _logger.info('WiFi IP: $ip');
      return ip;
    } catch (e) {
      _logger.error('Error getting WiFi IP', e);
      throw AppError.wifiService('Error getting WiFi IP: $e');
    }
  }

  Future<bool> isConnectedToWifi({
    String? requiredSSID,
    String? requiredBSSID,
  }) async {
    try {
      final isEnabled = await isWifiEnabled();
      if (!isEnabled) {
        return false;
      }

      if (requiredSSID != null || requiredBSSID != null) {
        final currentSSID = await getWifiName();
        final currentBSSID = await getWifiBSSID();

        if (requiredSSID != null && currentSSID != requiredSSID) {
          _logger.info(
            'WiFi SSID mismatch. Required: $requiredSSID, Current: $currentSSID',
          );
          return false;
        }

        if (requiredBSSID != null && currentBSSID != requiredBSSID) {
          _logger.info(
            'WiFi BSSID mismatch. Required: $requiredBSSID, Current: $currentBSSID',
          );
          return false;
        }
      }

      return true;
    } catch (e) {
      _logger.error('Error checking WiFi connection', e);
      throw AppError.wifiService('Error checking WiFi connection: $e');
    }
  }

  Future<Map<String, String?>> getWifiInfo() async {
    try {
      final ssid = await getWifiName();
      final bssid = await getWifiBSSID();
      final ip = await getWifiIP();

      return {
        'ssid': ssid,
        'bssid': bssid,
        'ip': ip,
      };
    } catch (e) {
      _logger.error('Error getting WiFi info', e);
      throw AppError.wifiService('Error getting WiFi info: $e');
    }
  }

  /// Verifies if the user's WiFi connection is valid for attendance marking
  Future<bool> verifyAttendanceWifi({
    required String sessionSSID,
    required String sessionBSSID,
    Duration timeout = const Duration(seconds: 10),
  }) async {
    try {
      final isEnabled = await isWifiEnabled();
      if (!isEnabled) {
        _logger.warning('WiFi is not enabled');
        return false;
      }

      // Start a timeout timer
      final completer = Completer<bool>();
      Timer(timeout, () {
        if (!completer.isCompleted) {
          completer.complete(false);
          _logger.warning(
              'WiFi verification timed out after ${timeout.inSeconds} seconds');
        }
      });

      // Verify WiFi connection
      final isValid = await isConnectedToWifi(
        requiredSSID: sessionSSID,
        requiredBSSID: sessionBSSID,
      ).timeout(
        timeout,
        onTimeout: () {
          _logger.warning('WiFi verification timed out');
          return false;
        },
      );

      if (!isValid && !completer.isCompleted) {
        _logger.warning('WiFi verification failed - Invalid network');
        completer.complete(false);
      } else if (!completer.isCompleted) {
        completer.complete(true);
      }

      return await completer.future;
    } catch (e) {
      _logger.error('Error verifying attendance WiFi', e);
      throw AppError.wifiService('Error verifying attendance WiFi: $e');
    }
  }

  /// Starts monitoring WiFi connection for attendance session
  Stream<bool> monitorAttendanceWifi({
    required String sessionSSID,
    required String sessionBSSID,
    Duration checkInterval = const Duration(seconds: 5),
  }) {
    return Stream.periodic(checkInterval).asyncMap((_) async {
      try {
        final isValid = await isConnectedToWifi(
          requiredSSID: sessionSSID,
          requiredBSSID: sessionBSSID,
        );

        if (!isValid) {
          _logger.warning(
            'WiFi connection changed. No longer connected to required network. '
            'SSID: $sessionSSID, BSSID: $sessionBSSID',
          );
        }

        return isValid;
      } catch (e) {
        _logger.error('Error monitoring attendance WiFi', e);
        return false;
      }
    });
  }
}
