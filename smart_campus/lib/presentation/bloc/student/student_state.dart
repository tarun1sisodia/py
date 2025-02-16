import 'package:equatable/equatable.dart';
import 'package:smart_campus/domain/entities/session.dart';
import 'package:smart_campus/domain/entities/course.dart';
import 'package:smart_campus/domain/entities/attendance_statistics.dart';
import 'package:smart_campus/domain/entities/attendance_record.dart';

abstract class StudentState extends Equatable {
  const StudentState();

  @override
  List<Object?> get props => [];
}

class StudentInitial extends StudentState {
  const StudentInitial();
}

class StudentLoading extends StudentState {
  const StudentLoading();
}

class StudentDashboardLoaded extends StudentState {
  final List<Session> activeSessions;
  final List<Course> courses;
  final AttendanceStatistics attendanceStats;

  const StudentDashboardLoaded({
    required this.activeSessions,
    required this.courses,
    required this.attendanceStats,
  });

  @override
  List<Object> get props => [activeSessions, courses, attendanceStats];
}

class StudentSessionsLoaded extends StudentState {
  final List<Session> sessions;
  final List<Course> courses;

  const StudentSessionsLoaded({
    required this.sessions,
    required this.courses,
  });

  @override
  List<Object> get props => [sessions, courses];
}

class AttendanceMarked extends StudentState {
  final AttendanceRecord record;
  final bool isOffline;

  const AttendanceMarked({
    required this.record,
    this.isOffline = false,
  });

  @override
  List<Object> get props => [record, isOffline];
}

class StudentFailure extends StudentState {
  final String message;

  const StudentFailure(this.message);

  @override
  List<Object> get props => [message];
}
