import 'dart:async';
import 'package:injectable/injectable.dart';
import 'package:smart_campus/core/errors/app_error.dart';
import 'package:smart_campus/core/services/api_client.dart';
import 'package:smart_campus/core/services/device_binding_service.dart';
import 'package:smart_campus/core/services/logger_service.dart';
import 'package:smart_campus/core/services/secure_storage_service.dart';
import 'package:smart_campus/data/models/user_model.dart';

@lazySingleton
class AuthService {
  static const _tokenKey = 'auth_token';
  static const _refreshTokenKey = 'refresh_token';
  static const _userKey = 'current_user';

  final ApiClient _apiClient;
  final DeviceBindingService _deviceBindingService;
  final LoggerService _logger;
  final SecureStorageService _secureStorage;
  final _authStateController = StreamController<AuthState>.broadcast();

  AuthService(
    this._apiClient,
    this._deviceBindingService,
    this._logger,
    this._secureStorage,
  );

  Stream<AuthState> get authStateChanges => _authStateController.stream;

  Future<UserModel> login({
    required String email,
    required String password,
  }) async {
    try {
      final response = await _apiClient.post<Map<String, dynamic>>(
        '/auth/login',
        data: {
          'email': email,
          'password': password,
          'device_info': await _deviceBindingService.getDeviceInfo(),
        },
      );

      final token = response['token'] as String;
      final refreshToken = response['refresh_token'] as String;
      final user = UserModel.fromJson(response['user'] as Map<String, dynamic>);

      await Future.wait([
        _secureStorage.write(key: _tokenKey, value: token),
        _secureStorage.write(key: _refreshTokenKey, value: refreshToken),
        _secureStorage.writeObject(key: _userKey, value: user.toJson()),
        _deviceBindingService.bindDevice(user.id),
      ]);

      _authStateController.add(AuthState.authenticated);
      return user;
    } catch (e) {
      _logger.error('Login failed', e);
      throw AppError.unauthorized();
    }
  }

  Future<void> logout() async {
    try {
      final token = await getToken();
      if (token != null) {
        await _apiClient.post<void>(
          '/auth/logout',
          data: {'token': token},
        );
      }
    } catch (e) {
      _logger.error('Error during logout', e);
    } finally {
      await Future.wait([
        _secureStorage.delete(_tokenKey),
        _secureStorage.delete(_refreshTokenKey),
        _secureStorage.delete(_userKey),
        _deviceBindingService.unbindDevice(),
      ]);
      _authStateController.add(AuthState.unauthenticated);
    }
  }

  Future<String?> getToken() async {
    return _secureStorage.read(_tokenKey);
  }

  Future<String?> refreshToken() async {
    try {
      final currentRefreshToken = await _secureStorage.read(_refreshTokenKey);
      if (currentRefreshToken == null) {
        throw AppError.unauthorized();
      }

      final response = await _apiClient.post<Map<String, dynamic>>(
        '/auth/refresh',
        data: {'refresh_token': currentRefreshToken},
      );

      final newToken = response['token'] as String;
      final newRefreshToken = response['refresh_token'] as String;

      await Future.wait([
        _secureStorage.write(key: _tokenKey, value: newToken),
        _secureStorage.write(key: _refreshTokenKey, value: newRefreshToken),
      ]);

      return newToken;
    } catch (e) {
      _logger.error('Token refresh failed', e);
      await logout();
      throw AppError.unauthorized();
    }
  }

  Future<UserModel?> getCurrentUser() async {
    try {
      final userData = await _secureStorage.readObject(_userKey);
      if (userData == null) return null;

      return UserModel.fromJson(userData);
    } catch (e) {
      _logger.error('Error getting current user', e);
      return null;
    }
  }

  Future<bool> isAuthenticated() async {
    final token = await getToken();
    final user = await getCurrentUser();
    return token != null && user != null;
  }

  Future<void> sendPasswordResetEmail(String email) async {
    try {
      await _apiClient.post<void>(
        '/auth/password/reset',
        data: {'email': email},
      );
    } catch (e) {
      _logger.error('Error sending password reset email', e);
      rethrow;
    }
  }

  Future<void> resetPassword({
    required String token,
    required String newPassword,
  }) async {
    try {
      await _apiClient.post<void>(
        '/auth/password/reset/confirm',
        data: {
          'token': token,
          'password': newPassword,
        },
      );
    } catch (e) {
      _logger.error('Error resetting password', e);
      rethrow;
    }
  }

  Future<void> changePassword({
    required String currentPassword,
    required String newPassword,
  }) async {
    try {
      await _apiClient.post<void>(
        '/auth/password/change',
        data: {
          'current_password': currentPassword,
          'new_password': newPassword,
        },
      );
    } catch (e) {
      _logger.error('Error changing password', e);
      rethrow;
    }
  }

  void dispose() {
    _authStateController.close();
  }
}

enum AuthState {
  initial,
  authenticated,
  unauthenticated,
}
