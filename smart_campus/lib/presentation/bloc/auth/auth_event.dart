import 'package:equatable/equatable.dart';
import 'package:smart_campus/domain/entities/user.dart';

abstract class AuthEvent extends Equatable {
  const AuthEvent();

  @override
  List<Object?> get props => [];
}

class AuthCheckRequested extends AuthEvent {
  const AuthCheckRequested();
}

class AuthLoginRequested extends AuthEvent {
  final String email;
  final String password;
  final String deviceId;

  const AuthLoginRequested({
    required this.email,
    required this.password,
    required this.deviceId,
  });

  @override
  List<Object> get props => [email, password, deviceId];
}

class AuthRegisterRequested extends AuthEvent {
  final String email;
  final String password;
  final String fullName;
  final UserRole role;
  final String? enrollmentNumber;
  final String? employeeId;
  final String? department;
  final int? yearOfStudy;
  final String deviceId;

  const AuthRegisterRequested({
    required this.email,
    required this.password,
    required this.fullName,
    required this.role,
    this.enrollmentNumber,
    this.employeeId,
    this.department,
    this.yearOfStudy,
    required this.deviceId,
  });

  @override
  List<Object?> get props => [
        email,
        password,
        fullName,
        role,
        enrollmentNumber,
        employeeId,
        department,
        yearOfStudy,
        deviceId,
      ];
}

class AuthLogoutRequested extends AuthEvent {
  const AuthLogoutRequested();
}

class AuthForgotPasswordRequested extends AuthEvent {
  final String email;

  const AuthForgotPasswordRequested(this.email);

  @override
  List<Object> get props => [email];
}

class AuthVerifyOTPRequested extends AuthEvent {
  final String email;
  final String otp;

  const AuthVerifyOTPRequested({
    required this.email,
    required this.otp,
  });

  @override
  List<Object> get props => [email, otp];
}

class AuthResetPasswordRequested extends AuthEvent {
  final String email;
  final String otp;
  final String newPassword;

  const AuthResetPasswordRequested({
    required this.email,
    required this.otp,
    required this.newPassword,
  });

  @override
  List<Object> get props => [email, otp, newPassword];
}

class AuthChangePasswordRequested extends AuthEvent {
  final String currentPassword;
  final String newPassword;

  const AuthChangePasswordRequested({
    required this.currentPassword,
    required this.newPassword,
  });

  @override
  List<Object> get props => [currentPassword, newPassword];
}

class AuthDeviceBindingRequested extends AuthEvent {
  final String deviceId;

  const AuthDeviceBindingRequested(this.deviceId);

  @override
  List<Object> get props => [deviceId];
}

class AuthDeviceUnbindingRequested extends AuthEvent {
  final String deviceId;

  const AuthDeviceUnbindingRequested(this.deviceId);

  @override
  List<Object> get props => [deviceId];
}
