import 'package:equatable/equatable.dart';

abstract class AuthEvent extends Equatable {
  const AuthEvent();

  @override
  List<Object?> get props => [];
}

class AuthRegisterTeacher extends AuthEvent {
  final String fullName;
  final String username;
  final String email;
  final String phone;
  final String password;
  final String highestDegree;
  final String experience;

  const AuthRegisterTeacher({
    required this.fullName,
    required this.username,
    required this.email,
    required this.phone,
    required this.password,
    required this.highestDegree,
    required this.experience,
  });

  @override
  List<Object?> get props => [
    fullName,
    username,
    email,
    phone,
    password,
    highestDegree,
    experience,
  ];
}

class AuthRegisterStudent extends AuthEvent {
  final String fullName;
  final String rollNumber;
  final String course;
  final String academicYear;
  final String phone;
  final String password;

  const AuthRegisterStudent({
    required this.fullName,
    required this.rollNumber,
    required this.course,
    required this.academicYear,
    required this.phone,
    required this.password,
  });

  @override
  List<Object?> get props => [
    fullName,
    rollNumber,
    course,
    academicYear,
    phone,
    password,
  ];
}

class AuthLoginTeacher extends AuthEvent {
  final String email;
  final String password;

  const AuthLoginTeacher({required this.email, required this.password});

  @override
  List<Object?> get props => [email, password];
}

class AuthLoginStudent extends AuthEvent {
  final String rollNumber;
  final String password;

  const AuthLoginStudent({required this.rollNumber, required this.password});

  @override
  List<Object?> get props => [rollNumber, password];
}

class AuthVerifyPhone extends AuthEvent {
  final String verificationId;
  final String otp;

  const AuthVerifyPhone({required this.verificationId, required this.otp});

  @override
  List<Object?> get props => [verificationId, otp];
}

class AuthSendOTP extends AuthEvent {
  final String phoneNumber;

  const AuthSendOTP({required this.phoneNumber});

  @override
  List<Object?> get props => [phoneNumber];
}

class AuthResetPassword extends AuthEvent {
  final String email;

  const AuthResetPassword({required this.email});

  @override
  List<Object?> get props => [email];
}

class AuthLogout extends AuthEvent {}

class AuthCheckStatus extends AuthEvent {}
