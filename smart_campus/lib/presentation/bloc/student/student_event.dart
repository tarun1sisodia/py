import 'package:equatable/equatable.dart';
import 'package:smart_campus/domain/entities/course.dart';

abstract class StudentEvent extends Equatable {
  const StudentEvent();

  @override
  List<Object?> get props => [];
}

class LoadStudentDashboard extends StudentEvent {
  const LoadStudentDashboard();
}

class LoadStudentSessions extends StudentEvent {
  final DateTime? startDate;
  final DateTime? endDate;
  final Course? course;

  const LoadStudentSessions({
    this.startDate,
    this.endDate,
    this.course,
  });

  @override
  List<Object?> get props => [startDate, endDate, course];
}

class MarkAttendanceRequested extends StudentEvent {
  final String sessionId;
  final String studentId;
  final String studentName;
  final double? locationLatitude;
  final double? locationLongitude;
  final String? wifiSSID;
  final String? wifiBSSID;
  final String deviceId;

  const MarkAttendanceRequested({
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
