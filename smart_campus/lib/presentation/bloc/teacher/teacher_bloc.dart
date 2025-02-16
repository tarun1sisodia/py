import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:injectable/injectable.dart';
import 'package:smart_campus/domain/repositories/session_repository.dart';
import 'package:smart_campus/domain/entities/session.dart';
import 'package:smart_campus/domain/entities/course.dart';
import 'package:smart_campus/domain/entities/attendance_record.dart';

// Events
abstract class TeacherEvent {
  const TeacherEvent();
}

class LoadTeacherDashboard extends TeacherEvent {
  const LoadTeacherDashboard();
}

class LoadSessionDetails extends TeacherEvent {
  final String sessionId;

  const LoadSessionDetails(this.sessionId);
}

class EndSessionRequested extends TeacherEvent {
  final String sessionId;

  const EndSessionRequested(this.sessionId);
}

class VerifyAttendanceRequested extends TeacherEvent {
  final String attendanceId;

  const VerifyAttendanceRequested(this.attendanceId);
}

class RejectAttendanceRequested extends TeacherEvent {
  final String attendanceId;
  final String reason;

  const RejectAttendanceRequested({
    required this.attendanceId,
    required this.reason,
  });
}

class CreateSessionRequested extends TeacherEvent {
  final String courseId;
  final DateTime sessionDate;
  final DateTime startTime;
  final DateTime endTime;
  final double locationLatitude;
  final double locationLongitude;
  final int locationRadius;
  final String wifiSSID;
  final String wifiBSSID;

  const CreateSessionRequested({
    required this.courseId,
    required this.sessionDate,
    required this.startTime,
    required this.endTime,
    required this.locationLatitude,
    required this.locationLongitude,
    required this.locationRadius,
    required this.wifiSSID,
    required this.wifiBSSID,
  });
}

// States
abstract class TeacherState {
  const TeacherState();
}

class TeacherInitial extends TeacherState {
  const TeacherInitial();
}

class TeacherLoading extends TeacherState {
  const TeacherLoading();
}

class TeacherDashboardLoaded extends TeacherState {
  final List<Course> assignedCourses;
  final List<Session> activeSessions;
  final int totalSessionsToday;
  final int totalStudentsPresent;

  const TeacherDashboardLoaded({
    required this.assignedCourses,
    required this.activeSessions,
    required this.totalSessionsToday,
    required this.totalStudentsPresent,
  });
}

class SessionDetailsLoaded extends TeacherState {
  final Session session;
  final Course course;
  final List<AttendanceRecord> attendanceRecords;

  const SessionDetailsLoaded({
    required this.session,
    required this.course,
    required this.attendanceRecords,
  });
}

class SessionCreationSuccess extends TeacherState {
  final Session session;
  final Course course;

  const SessionCreationSuccess({
    required this.session,
    required this.course,
  });
}

class TeacherError extends TeacherState {
  final String message;

  const TeacherError(this.message);
}

// Bloc
@injectable
class TeacherBloc extends Bloc<TeacherEvent, TeacherState> {
  final SessionRepository _sessionRepository;

  TeacherBloc({
    required SessionRepository sessionRepository,
  })  : _sessionRepository = sessionRepository,
        super(const TeacherInitial()) {
    on<LoadTeacherDashboard>(_onLoadTeacherDashboard);
    on<CreateSessionRequested>(_onCreateSessionRequested);
    on<LoadSessionDetails>(_onLoadSessionDetails);
    on<EndSessionRequested>(_onEndSessionRequested);
    on<VerifyAttendanceRequested>(_onVerifyAttendanceRequested);
    on<RejectAttendanceRequested>(_onRejectAttendanceRequested);
  }

  Future<void> _onLoadTeacherDashboard(
    LoadTeacherDashboard event,
    Emitter<TeacherState> emit,
  ) async {
    try {
      emit(const TeacherLoading());

      // Get assigned courses
      final courses = await _sessionRepository.getCourses();

      // Get active sessions
      final activeSessions = await _sessionRepository.getActiveSessions();

      // Get today's sessions count
      final today = DateTime.now();
      final todaySessions = await _sessionRepository.getSessions(
        startDate: DateTime(today.year, today.month, today.day),
        endDate: DateTime(today.year, today.month, today.day, 23, 59, 59),
      );

      // Calculate total students present in active sessions
      int totalPresent = 0;
      for (final session in activeSessions) {
        final attendance =
            await _sessionRepository.getSessionAttendance(session.id);
        totalPresent += attendance.length;
      }

      emit(TeacherDashboardLoaded(
        assignedCourses: courses,
        activeSessions: activeSessions,
        totalSessionsToday: todaySessions.length,
        totalStudentsPresent: totalPresent,
      ));
    } catch (e) {
      emit(TeacherError(e.toString()));
    }
  }

  Future<void> _onCreateSessionRequested(
    CreateSessionRequested event,
    Emitter<TeacherState> emit,
  ) async {
    try {
      emit(const TeacherLoading());

      // Get course object first
      final course = await _sessionRepository.getCourseById(event.courseId);

      final session = await _sessionRepository.createSession(
        course: course,
        sessionDate: event.sessionDate,
        startTime: event.startTime,
        endTime: event.endTime,
        locationLatitude: event.locationLatitude,
        locationLongitude: event.locationLongitude,
        locationRadius: event.locationRadius,
        wifiSSID: event.wifiSSID,
        wifiBSSID: event.wifiBSSID,
      );

      emit(SessionCreationSuccess(
        session: session,
        course: course,
      ));

      // Reload dashboard after successful session creation
      add(const LoadTeacherDashboard());
    } catch (e) {
      emit(TeacherError(e.toString()));
    }
  }

  Future<void> _onLoadSessionDetails(
    LoadSessionDetails event,
    Emitter<TeacherState> emit,
  ) async {
    try {
      emit(const TeacherLoading());

      final session = await _sessionRepository.getSessionById(event.sessionId);
      final course = await _sessionRepository.getCourseById(session.courseId);
      final attendanceRecords =
          await _sessionRepository.getSessionAttendance(event.sessionId);

      emit(SessionDetailsLoaded(
        session: session,
        course: course,
        attendanceRecords: attendanceRecords,
      ));
    } catch (e) {
      emit(TeacherError(e.toString()));
    }
  }

  Future<void> _onEndSessionRequested(
    EndSessionRequested event,
    Emitter<TeacherState> emit,
  ) async {
    try {
      emit(const TeacherLoading());

      final session = await _sessionRepository.endSession(event.sessionId);
      final course = await _sessionRepository.getCourseById(session.courseId);
      final attendanceRecords =
          await _sessionRepository.getSessionAttendance(event.sessionId);

      emit(SessionDetailsLoaded(
        session: session,
        course: course,
        attendanceRecords: attendanceRecords,
      ));
    } catch (e) {
      emit(TeacherError(e.toString()));
    }
  }

  Future<void> _onVerifyAttendanceRequested(
    VerifyAttendanceRequested event,
    Emitter<TeacherState> emit,
  ) async {
    try {
      emit(const TeacherLoading());

      final record =
          await _sessionRepository.verifyAttendance(event.attendanceId);
      final session = await _sessionRepository.getSessionById(record.sessionId);
      final course = await _sessionRepository.getCourseById(session.courseId);
      final attendanceRecords =
          await _sessionRepository.getSessionAttendance(session.id);

      emit(SessionDetailsLoaded(
        session: session,
        course: course,
        attendanceRecords: attendanceRecords,
      ));
    } catch (e) {
      emit(TeacherError(e.toString()));
    }
  }

  Future<void> _onRejectAttendanceRequested(
    RejectAttendanceRequested event,
    Emitter<TeacherState> emit,
  ) async {
    try {
      emit(const TeacherLoading());

      final record = await _sessionRepository.rejectAttendance(
        event.attendanceId,
        reason: event.reason,
      );
      final session = await _sessionRepository.getSessionById(record.sessionId);
      final course = await _sessionRepository.getCourseById(session.courseId);
      final attendanceRecords =
          await _sessionRepository.getSessionAttendance(session.id);

      emit(SessionDetailsLoaded(
        session: session,
        course: course,
        attendanceRecords: attendanceRecords,
      ));
    } catch (e) {
      emit(TeacherError(e.toString()));
    }
  }
}
