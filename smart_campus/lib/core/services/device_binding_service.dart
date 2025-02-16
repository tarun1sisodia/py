import 'dart:async';
import 'package:device_info_plus/device_info_plus.dart';
import 'package:injectable/injectable.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:uuid/uuid.dart';
import 'package:smart_campus/core/errors/app_error.dart';
import 'package:smart_campus/core/services/logger_service.dart';

@lazySingleton
class DeviceBindingService {
  final DeviceInfoPlugin _deviceInfo;
  final SharedPreferences _prefs;
  final LoggerService _logger;
  static const _deviceIdKey = 'device_binding_id';
  static const _deviceBindingKey = 'device_binding_data';

  DeviceBindingService(
    this._prefs,
    this._logger,
  ) : _deviceInfo = DeviceInfoPlugin();

  /// Get the unique device identifier
  Future<String> getDeviceId() async {
    try {
      // Try to get existing device ID
      String? deviceId = _prefs.getString(_deviceIdKey);

      if (deviceId == null) {
        // Generate new device ID if not exists
        deviceId = const Uuid().v4();
        await _prefs.setString(_deviceIdKey, deviceId);
      }

      return deviceId;
    } catch (e) {
      _logger.error('Error getting device ID', e);
      throw AppError.deviceBinding('Failed to get device ID: $e');
    }
  }

  /// Get detailed device information
  Future<Map<String, dynamic>> getDeviceInfo() async {
    try {
      if (await isDeviceDeveloperModeEnabled()) {
        throw AppError.deviceBinding('Developer mode is enabled');
      }

      final androidInfo = await _deviceInfo.androidInfo;

      return {
        'manufacturer': androidInfo.manufacturer,
        'model': androidInfo.model,
        'androidId': androidInfo.id,
        'brand': androidInfo.brand,
        'device': androidInfo.device,
        'hardware': androidInfo.hardware,
        'isPhysicalDevice': androidInfo.isPhysicalDevice,
        'fingerprint': androidInfo.fingerprint,
        'securityPatch': androidInfo.version.securityPatch,
        'sdkInt': androidInfo.version.sdkInt,
        'release': androidInfo.version.release,
      };
    } catch (e) {
      _logger.error('Error getting device info', e);
      throw AppError.deviceBinding('Failed to get device info: $e');
    }
  }

  /// Check if device is in developer mode
  Future<bool> isDeviceDeveloperModeEnabled() async {
    try {
      final androidInfo = await _deviceInfo.androidInfo;
      return androidInfo.isPhysicalDevice == false ||
          androidInfo.version.sdkInt >= 30; // Android 11 or higher
    } catch (e) {
      _logger.error('Error checking developer mode', e);
      return true; // Fail safe: assume developer mode is enabled
    }
  }

  /// Bind device to user
  Future<void> bindDevice(String userId) async {
    try {
      if (await isDeviceDeveloperModeEnabled()) {
        throw AppError.deviceBinding(
            'Cannot bind device with developer mode enabled');
      }

      final deviceId = await getDeviceId();
      final deviceInfo = await getDeviceInfo();

      final bindingData = {
        'userId': userId,
        'deviceId': deviceId,
        'bindingTime': DateTime.now().toIso8601String(),
        'deviceInfo': deviceInfo,
      };

      await _prefs.setString(_deviceBindingKey, deviceId);
      _logger.info('Device bound successfully', bindingData);
    } catch (e) {
      _logger.error('Error binding device', e);
      throw AppError.deviceBinding('Failed to bind device: $e');
    }
  }

  /// Verify if device is bound to user
  Future<bool> verifyDeviceBinding(String userId) async {
    try {
      if (await isDeviceDeveloperModeEnabled()) {
        _logger.warning('Device verification failed: Developer mode enabled');
        return false;
      }

      final boundDeviceId = _prefs.getString(_deviceBindingKey);
      if (boundDeviceId == null) {
        _logger.warning('Device verification failed: No device binding found');
        return false;
      }

      final currentDeviceId = await getDeviceId();
      if (boundDeviceId != currentDeviceId) {
        _logger.warning(
          'Device verification failed: Device ID mismatch',
          {'bound': boundDeviceId, 'current': currentDeviceId},
        );
        return false;
      }

      return true;
    } catch (e) {
      _logger.error('Error verifying device binding', e);
      return false;
    }
  }

  /// Remove device binding
  Future<void> unbindDevice() async {
    try {
      await _prefs.remove(_deviceBindingKey);
      _logger.info('Device unbound successfully');
    } catch (e) {
      _logger.error('Error unbinding device', e);
      throw AppError.deviceBinding('Failed to unbind device: $e');
    }
  }

  /// Get binding status
  Future<Map<String, dynamic>> getBindingStatus() async {
    try {
      final deviceId = await getDeviceId();
      final boundDeviceId = _prefs.getString(_deviceBindingKey);
      final isDeveloperMode = await isDeviceDeveloperModeEnabled();

      return {
        'isBound': boundDeviceId != null && boundDeviceId == deviceId,
        'deviceId': deviceId,
        'isDeveloperMode': isDeveloperMode,
        'bindingTime': boundDeviceId != null
            ? _prefs.getString('${_deviceBindingKey}_time')
            : null,
      };
    } catch (e) {
      _logger.error('Error getting binding status', e);
      throw AppError.deviceBinding('Failed to get binding status: $e');
    }
  }
}
