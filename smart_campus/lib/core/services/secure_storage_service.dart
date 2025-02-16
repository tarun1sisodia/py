import 'dart:convert';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import 'package:injectable/injectable.dart';
import 'package:smart_campus/core/errors/app_error.dart';
import 'package:smart_campus/core/services/logger_service.dart';

@lazySingleton
class SecureStorageService {
  final FlutterSecureStorage _storage;
  final LoggerService _logger;

  SecureStorageService(this._logger)
      : _storage = const FlutterSecureStorage(
          aOptions: AndroidOptions(
            encryptedSharedPreferences: true,
          ),
        );

  Future<void> write({
    required String key,
    required String value,
  }) async {
    try {
      await _storage.write(key: key, value: value);
    } catch (e) {
      _logger.error('Error writing to secure storage', e);
      throw AppError('Failed to write to secure storage: $e');
    }
  }

  Future<String?> read(String key) async {
    try {
      return await _storage.read(key: key);
    } catch (e) {
      _logger.error('Error reading from secure storage', e);
      throw AppError('Failed to read from secure storage: $e');
    }
  }

  Future<void> delete(String key) async {
    try {
      await _storage.delete(key: key);
    } catch (e) {
      _logger.error('Error deleting from secure storage', e);
      throw AppError('Failed to delete from secure storage: $e');
    }
  }

  Future<void> deleteAll() async {
    try {
      await _storage.deleteAll();
    } catch (e) {
      _logger.error('Error deleting all from secure storage', e);
      throw AppError('Failed to delete all from secure storage: $e');
    }
  }

  Future<Map<String, String>> readAll() async {
    try {
      return await _storage.readAll();
    } catch (e) {
      _logger.error('Error reading all from secure storage', e);
      throw AppError('Failed to read all from secure storage: $e');
    }
  }

  Future<bool> containsKey(String key) async {
    try {
      return await _storage.containsKey(key: key);
    } catch (e) {
      _logger.error('Error checking key in secure storage', e);
      throw AppError('Failed to check key in secure storage: $e');
    }
  }

  Future<void> writeObject({
    required String key,
    required Map<String, dynamic> value,
  }) async {
    try {
      final jsonString = jsonEncode(value);
      await write(key: key, value: jsonString);
    } catch (e) {
      _logger.error('Error writing object to secure storage', e);
      throw AppError('Failed to write object to secure storage: $e');
    }
  }

  Future<Map<String, dynamic>?> readObject(String key) async {
    try {
      final jsonString = await read(key);
      if (jsonString == null) return null;
      return jsonDecode(jsonString) as Map<String, dynamic>;
    } catch (e) {
      _logger.error('Error reading object from secure storage', e);
      throw AppError('Failed to read object from secure storage: $e');
    }
  }
}
