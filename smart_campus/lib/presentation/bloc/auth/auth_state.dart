import 'package:equatable/equatable.dart';
import 'package:smart_campus/domain/entities/user.dart';

abstract class AuthState extends Equatable {
  const AuthState();

  @override
  List<Object?> get props => [];
}

class AuthInitial extends AuthState {
  const AuthInitial();
}

class AuthLoading extends AuthState {
  const AuthLoading();
}

class AuthAuthenticated extends AuthState {
  final User user;

  const AuthAuthenticated(this.user);

  @override
  List<Object> get props => [user];
}

class AuthUnauthenticated extends AuthState {
  const AuthUnauthenticated();
}

class AuthFailure extends AuthState {
  final String message;

  const AuthFailure(this.message);

  @override
  List<Object> get props => [message];
}

class AuthOTPSent extends AuthState {
  final String email;

  const AuthOTPSent(this.email);

  @override
  List<Object> get props => [email];
}

class AuthOTPVerified extends AuthState {
  final String email;
  final String otp;

  const AuthOTPVerified({
    required this.email,
    required this.otp,
  });

  @override
  List<Object> get props => [email, otp];
}

class AuthPasswordChanged extends AuthState {
  const AuthPasswordChanged();
}

class AuthDeviceBound extends AuthState {
  final String deviceId;

  const AuthDeviceBound(this.deviceId);

  @override
  List<Object> get props => [deviceId];
}

class AuthDeviceUnbound extends AuthState {
  final String deviceId;

  const AuthDeviceUnbound(this.deviceId);

  @override
  List<Object> get props => [deviceId];
}
