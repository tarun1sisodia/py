import 'package:equatable/equatable.dart';
import 'package:smart_campus/domain/entities/session.dart';
import 'package:smart_campus/domain/entities/course.dart';
import 'package:smart_campus/domain/entities/attendance_record.dart';

abstract class SessionState extends Equatable {
  const SessionState();

  @override
  List<Object?> get props => [];
}

class SessionInitial extends SessionState {
  const SessionInitial();
}

class SessionLoading extends SessionState {
  const SessionLoading();
}

class SessionLoadSuccess extends SessionState {
  final List<Session> sessions;
  final List<Course> courses;

  const SessionLoadSuccess({
    required this.sessions,
    required this.courses,
  });

  @override
  List<Object> get props => [sessions, courses];
}

class ActiveSessionsLoaded extends SessionState {
  final List<Session> activeSessions;
  final List<Course> courses;

  const ActiveSessionsLoaded({
    required this.activeSessions,
    required this.courses,
  });

  @override
  List<Object> get props => [activeSessions, courses];
}

class SessionCreated extends SessionState {
  final Session session;
  final Course course;

  const SessionCreated({
    required this.session,
    required this.course,
  });

  @override
  List<Object> get props => [session, course];
}

class SessionEnded extends SessionState {
  final Session session;

  const SessionEnded(this.session);

  @override
  List<Object> get props => [session];
}

class SessionDetailsLoaded extends SessionState {
  final Session session;
  final Course course;
  final List<AttendanceRecord> attendanceRecords;

  const SessionDetailsLoaded({
    required this.session,
    required this.course,
    required this.attendanceRecords,
  });

  @override
  List<Object> get props => [session, course, attendanceRecords];
}

class AttendanceMarked extends SessionState {
  final AttendanceRecord attendanceRecord;

  const AttendanceMarked(this.attendanceRecord);

  @override
  List<Object> get props => [attendanceRecord];
}

class AttendanceMarkedOffline extends SessionState {
  const AttendanceMarkedOffline();
}

class AttendanceVerified extends SessionState {
  final AttendanceRecord attendanceRecord;

  const AttendanceVerified(this.attendanceRecord);

  @override
  List<Object> get props => [attendanceRecord];
}

class AttendanceRejected extends SessionState {
  final AttendanceRecord attendanceRecord;

  const AttendanceRejected(this.attendanceRecord);

  @override
  List<Object> get props => [attendanceRecord];
}

class SessionFailure extends SessionState {
  final String message;

  const SessionFailure(this.message);

  @override
  List<Object> get props => [message];
}
