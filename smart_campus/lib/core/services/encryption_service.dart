import 'dart:convert';
import 'package:encrypt/encrypt.dart';
import 'package:injectable/injectable.dart';
import 'package:smart_campus/core/errors/app_error.dart';
import 'package:smart_campus/core/services/logger_service.dart';
import 'package:smart_campus/core/services/secure_storage_service.dart';

@lazySingleton
class EncryptionService {
  static const _keyKey = 'encryption_key';
  static const _ivKey = 'encryption_iv';
  final SecureStorageService _secureStorage;
  final LoggerService _logger;
  Encrypter? _encrypter;
  IV? _iv;

  EncryptionService(
    this._secureStorage,
    this._logger,
  );

  Future<void> initialize() async {
    try {
      // Try to load existing key and IV
      String? keyString = await _secureStorage.read(_keyKey);
      String? ivString = await _secureStorage.read(_ivKey);

      if (keyString == null || ivString == null) {
        // Generate new key and IV if not found
        final key = Key.fromSecureRandom(32);
        _iv = IV.fromSecureRandom(16);

        // Save the new key and IV
        await Future.wait([
          _secureStorage.write(
            key: _keyKey,
            value: base64Encode(key.bytes),
          ),
          _secureStorage.write(
            key: _ivKey,
            value: base64Encode(_iv!.bytes),
          ),
        ]);

        _encrypter = Encrypter(AES(key));
      } else {
        // Use existing key and IV
        final key = Key(base64Decode(keyString));
        _iv = IV(base64Decode(ivString));
        _encrypter = Encrypter(AES(key));
      }
    } catch (e) {
      _logger.error('Error initializing encryption service', e);
      throw AppError('Failed to initialize encryption service: $e');
    }
  }

  Future<String> encrypt(String data) async {
    try {
      if (_encrypter == null || _iv == null) {
        await initialize();
      }

      final encrypted = _encrypter!.encrypt(data, iv: _iv!);
      return encrypted.base64;
    } catch (e) {
      _logger.error('Error encrypting data', e);
      throw AppError('Failed to encrypt data: $e');
    }
  }

  Future<String> decrypt(String encryptedData) async {
    try {
      if (_encrypter == null || _iv == null) {
        await initialize();
      }

      final encrypted = Encrypted.fromBase64(encryptedData);
      return _encrypter!.decrypt(encrypted, iv: _iv!);
    } catch (e) {
      _logger.error('Error decrypting data', e);
      throw AppError('Failed to decrypt data: $e');
    }
  }

  Future<String> encryptObject(Map<String, dynamic> data) async {
    try {
      final jsonString = jsonEncode(data);
      return await encrypt(jsonString);
    } catch (e) {
      _logger.error('Error encrypting object', e);
      throw AppError('Failed to encrypt object: $e');
    }
  }

  Future<Map<String, dynamic>> decryptObject(String encryptedData) async {
    try {
      final jsonString = await decrypt(encryptedData);
      return jsonDecode(jsonString) as Map<String, dynamic>;
    } catch (e) {
      _logger.error('Error decrypting object', e);
      throw AppError('Failed to decrypt object: $e');
    }
  }

  Future<List<int>> generateSecureRandomBytes(int length) async {
    try {
      final random = SecureRandom(length);
      return random.bytes;
    } catch (e) {
      _logger.error('Error generating secure random bytes', e);
      throw AppError('Failed to generate secure random bytes: $e');
    }
  }

  Future<String> generateSecureToken([int length = 32]) async {
    try {
      final bytes = await generateSecureRandomBytes(length);
      return base64Url.encode(bytes);
    } catch (e) {
      _logger.error('Error generating secure token', e);
      throw AppError('Failed to generate secure token: $e');
    }
  }

  Future<bool> verifyHash({
    required String data,
    required String hash,
  }) async {
    try {
      final dataHash = await encrypt(data);
      return hash == dataHash;
    } catch (e) {
      _logger.error('Error verifying hash', e);
      throw AppError('Failed to verify hash: $e');
    }
  }
}
