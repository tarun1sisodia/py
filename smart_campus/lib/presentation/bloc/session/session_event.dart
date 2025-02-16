import 'package:equatable/equatable.dart';
import 'package:smart_campus/domain/entities/course.dart';

abstract class SessionEvent extends Equatable {
  const SessionEvent();

  @override
  List<Object?> get props => [];
}

class LoadSessions extends SessionEvent {
  final DateTime? startDate;
  final DateTime? endDate;
  final Course? course;

  const LoadSessions({
    this.startDate,
    this.endDate,
    this.course,
  });

  @override
  List<Object?> get props => [startDate, endDate, course];
}

class LoadActiveSessions extends SessionEvent {
  const LoadActiveSessions();
}

class CreateSession extends SessionEvent {
  final Course course;
  final DateTime sessionDate;
  final DateTime startTime;
  final DateTime endTime;
  final double? locationLatitude;
  final double? locationLongitude;
  final int? locationRadius;
  final String? wifiSSID;
  final String? wifiBSSID;

  const CreateSession({
    required this.course,
    required this.sessionDate,
    required this.startTime,
    required this.endTime,
    this.locationLatitude,
    this.locationLongitude,
    this.locationRadius,
    this.wifiSSID,
    this.wifiBSSID,
  });

  @override
  List<Object?> get props => [
        course,
        sessionDate,
        startTime,
        endTime,
        locationLatitude,
        locationLongitude,
        locationRadius,
        wifiSSID,
        wifiBSSID,
      ];
}

class EndSession extends SessionEvent {
  final String sessionId;

  const EndSession(this.sessionId);

  @override
  List<Object> get props => [sessionId];
}

class LoadSessionDetails extends SessionEvent {
  final String sessionId;

  const LoadSessionDetails(this.sessionId);

  @override
  List<Object> get props => [sessionId];
}

class MarkAttendance extends SessionEvent {
  final String sessionId;
  final String studentId;
  final String studentName;
  final double? locationLatitude;
  final double? locationLongitude;
  final String? wifiSSID;
  final String? wifiBSSID;
  final String deviceId;

  const MarkAttendance({
    required this.sessionId,
    required this.studentId,
    required this.studentName,
    this.locationLatitude,
    this.locationLongitude,
    this.wifiSSID,
    this.wifiBSSID,
    required this.deviceId,
  });

  @override
  List<Object?> get props => [
        sessionId,
        studentId,
        studentName,
        locationLatitude,
        locationLongitude,
        wifiSSID,
        wifiBSSID,
        deviceId,
      ];
}

class VerifyAttendance extends SessionEvent {
  final String attendanceId;

  const VerifyAttendance(this.attendanceId);

  @override
  List<Object> get props => [attendanceId];
}

class RejectAttendance extends SessionEvent {
  final String attendanceId;
  final String reason;

  const RejectAttendance({
    required this.attendanceId,
    required this.reason,
  });

  @override
  List<Object> get props => [attendanceId, reason];
}
