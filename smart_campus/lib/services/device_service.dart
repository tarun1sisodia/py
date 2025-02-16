import 'package:injectable/injectable.dart';
import 'package:smart_campus/data/datasources/local/database_helper.dart';
import 'package:smart_campus/core/error/app_error.dart';

@lazySingleton
class DeviceService {
  final DatabaseHelper _databaseHelper;

  DeviceService(this._databaseHelper);

  /// Checks if a device is registered to a student
  Future<bool> isDeviceRegistered({
    required String studentId,
    required String deviceId,
  }) async {
    try {
      final results = await _databaseHelper.query(
        'device_bindings',
        where: 'student_id = ? AND device_id = ? AND is_active = ?',
        whereArgs: [studentId, deviceId, 1],
      );

      return results.isNotEmpty;
    } catch (e) {
      throw AppError('Failed to check device registration: $e');
    }
  }

  /// Checks if a device is authorized (not blacklisted)
  Future<bool> isDeviceAuthorized(String deviceId) async {
    try {
      final results = await _databaseHelper.query(
        'device_bindings',
        where: 'device_id = ? AND is_blacklisted = ?',
        whereArgs: [deviceId, 0],
      );

      return results.isNotEmpty;
    } catch (e) {
      throw AppError('Failed to check device authorization: $e');
    }
  }

  /// Registers a new device for a student
  Future<void> registerDevice({
    required String studentId,
    required String deviceId,
    required String deviceName,
    required String deviceModel,
  }) async {
    try {
      // Check if device is already registered
      final existingDevice = await _databaseHelper.query(
        'device_bindings',
        where: 'student_id = ? AND device_id = ?',
        whereArgs: [studentId, deviceId],
      );

      if (existingDevice.isNotEmpty) {
        // Update existing device if inactive
        if (existingDevice.first['is_active'] == 0) {
          await _databaseHelper.update(
            'device_bindings',
            {
              'is_active': 1,
              'device_name': deviceName,
              'device_model': deviceModel,
              'updated_at': DateTime.now().toIso8601String(),
            },
            where: 'student_id = ? AND device_id = ?',
            whereArgs: [studentId, deviceId],
          );
        } else {
          throw AppError('Device is already registered');
        }
      } else {
        // Insert new device
        await _databaseHelper.insert('device_bindings', {
          'student_id': studentId,
          'device_id': deviceId,
          'device_name': deviceName,
          'device_model': deviceModel,
          'is_active': 1,
          'is_blacklisted': 0,
          'created_at': DateTime.now().toIso8601String(),
          'updated_at': DateTime.now().toIso8601String(),
        });
      }
    } catch (e) {
      throw AppError('Failed to register device: $e');
    }
  }

  /// Deactivates a device for a student
  Future<void> deactivateDevice({
    required String studentId,
    required String deviceId,
  }) async {
    try {
      final count = await _databaseHelper.update(
        'device_bindings',
        {
          'is_active': 0,
          'updated_at': DateTime.now().toIso8601String(),
        },
        where: 'student_id = ? AND device_id = ?',
        whereArgs: [studentId, deviceId],
      );

      if (count == 0) {
        throw AppError('Device not found');
      }
    } catch (e) {
      throw AppError('Failed to deactivate device: $e');
    }
  }

  /// Blacklists a device
  Future<void> blacklistDevice(String deviceId) async {
    try {
      final count = await _databaseHelper.update(
        'device_bindings',
        {
          'is_blacklisted': 1,
          'updated_at': DateTime.now().toIso8601String(),
        },
        where: 'device_id = ?',
        whereArgs: [deviceId],
      );

      if (count == 0) {
        throw AppError('Device not found');
      }
    } catch (e) {
      throw AppError('Failed to blacklist device: $e');
    }
  }

  /// Removes a device from the blacklist
  Future<void> removeFromBlacklist(String deviceId) async {
    try {
      final count = await _databaseHelper.update(
        'device_bindings',
        {
          'is_blacklisted': 0,
          'updated_at': DateTime.now().toIso8601String(),
        },
        where: 'device_id = ?',
        whereArgs: [deviceId],
      );

      if (count == 0) {
        throw AppError('Device not found');
      }
    } catch (e) {
      throw AppError('Failed to remove device from blacklist: $e');
    }
  }
}
