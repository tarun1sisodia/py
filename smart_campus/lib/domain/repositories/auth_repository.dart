import 'package:smart_campus/domain/entities/user.dart';

abstract class AuthRepository {
  Future<User> login({
    required String email,
    required String password,
    required String deviceId,
  });

  Future<User> register({
    required String email,
    required String password,
    required String fullName,
    required UserRole role,
    String? enrollmentNumber,
    String? employeeId,
    String? department,
    int? yearOfStudy,
    required String deviceId,
  });

  Future<void> logout();

  Future<void> forgotPassword(String email);

  Future<void> verifyOTP({
    required String email,
    required String otp,
  });

  Future<void> resetPassword({
    required String email,
    required String otp,
    required String newPassword,
  });

  Future<void> changePassword({
    required String currentPassword,
    required String newPassword,
  });

  Future<User?> getCurrentUser();

  Future<String> refreshToken();

  Future<bool> isDeviceBound(String deviceId);

  Future<void> bindDevice(String deviceId);

  Future<void> unbindDevice(String deviceId);
}
