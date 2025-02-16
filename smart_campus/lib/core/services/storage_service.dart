import 'dart:convert';
import 'package:injectable/injectable.dart';
import 'package:shared_preferences/shared_preferences.dart';

@lazySingleton
class StorageService {
  static const String _tokenKey = 'auth_token';
  static const String _refreshTokenKey = 'refresh_token';
  static const String _userKey = 'current_user';
  static const String _deviceIdKey = 'device_id';
  static const String _lastSyncKey = 'last_sync';

  final SharedPreferences _prefs;

  StorageService(this._prefs);

  // Auth Token
  Future<void> saveAuthToken(String token) async {
    await _prefs.setString(_tokenKey, token);
  }

  String? getAuthToken() {
    return _prefs.getString(_tokenKey);
  }

  // Refresh Token
  Future<void> saveRefreshToken(String token) async {
    await _prefs.setString(_refreshTokenKey, token);
  }

  String? getRefreshToken() {
    return _prefs.getString(_refreshTokenKey);
  }

  // User Data
  Future<void> saveUser(Map<String, dynamic> userData) async {
    await _prefs.setString(_userKey, jsonEncode(userData));
  }

  Map<String, dynamic>? getUser() {
    final userStr = _prefs.getString(_userKey);
    if (userStr != null) {
      return jsonDecode(userStr) as Map<String, dynamic>;
    }
    return null;
  }

  // Device ID
  Future<void> saveDeviceId(String deviceId) async {
    await _prefs.setString(_deviceIdKey, deviceId);
  }

  String? getDeviceId() {
    return _prefs.getString(_deviceIdKey);
  }

  // Last Sync Time
  Future<void> saveLastSyncTime(DateTime time) async {
    await _prefs.setString(_lastSyncKey, time.toIso8601String());
  }

  DateTime? getLastSyncTime() {
    final timeStr = _prefs.getString(_lastSyncKey);
    if (timeStr != null) {
      return DateTime.parse(timeStr);
    }
    return null;
  }

  // Clear Storage
  Future<void> clearAll() async {
    await _prefs.clear();
  }

  Future<void> clearAuth() async {
    await _prefs.remove(_tokenKey);
    await _prefs.remove(_refreshTokenKey);
    await _prefs.remove(_userKey);
  }
}
