import 'package:equatable/equatable.dart';

abstract class AttendanceVerificationEvent extends Equatable {
  const AttendanceVerificationEvent();

  @override
  List<Object?> get props => [];
}

class StartVerification extends AttendanceVerificationEvent {
  final String userId;
  final double sessionLatitude;
  final double sessionLongitude;
  final double allowedRadius;
  final String sessionSSID;
  final String sessionBSSID;

  const StartVerification({
    required this.userId,
    required this.sessionLatitude,
    required this.sessionLongitude,
    required this.allowedRadius,
    required this.sessionSSID,
    required this.sessionBSSID,
  });

  @override
  List<Object> get props => [
        userId,
        sessionLatitude,
        sessionLongitude,
        allowedRadius,
        sessionSSID,
        sessionBSSID,
      ];
}

class StartMonitoring extends AttendanceVerificationEvent {
  final String userId;
  final double sessionLatitude;
  final double sessionLongitude;
  final double allowedRadius;
  final String sessionSSID;
  final String sessionBSSID;

  const StartMonitoring({
    required this.userId,
    required this.sessionLatitude,
    required this.sessionLongitude,
    required this.allowedRadius,
    required this.sessionSSID,
    required this.sessionBSSID,
  });

  @override
  List<Object> get props => [
        userId,
        sessionLatitude,
        sessionLongitude,
        allowedRadius,
        sessionSSID,
        sessionBSSID,
      ];
}

class StopMonitoring extends AttendanceVerificationEvent {}

class ResetVerification extends AttendanceVerificationEvent {}
