import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:smart_campus/domain/repositories/session_repository.dart';
import 'package:smart_campus/domain/entities/attendance_statistics.dart';
import 'package:smart_campus/core/services/sync_service.dart';
import 'package:smart_campus/core/services/connectivity_service.dart';
import 'package:smart_campus/presentation/bloc/student/student_event.dart';
import 'package:smart_campus/presentation/bloc/student/student_state.dart';

class StudentBloc extends Bloc<StudentEvent, StudentState> {
  final SessionRepository _sessionRepository;
  final SyncService _syncService;
  final ConnectivityService _connectivityService;

  StudentBloc({
    required SessionRepository sessionRepository,
    required SyncService syncService,
    required ConnectivityService connectivityService,
  })  : _sessionRepository = sessionRepository,
        _syncService = syncService,
        _connectivityService = connectivityService,
        super(const StudentInitial()) {
    on<LoadStudentDashboard>(_onLoadStudentDashboard);
    on<LoadStudentSessions>(_onLoadStudentSessions);
    on<MarkAttendanceRequested>(_onMarkAttendanceRequested);
  }

  Future<void> _onLoadStudentDashboard(
    LoadStudentDashboard event,
    Emitter<StudentState> emit,
  ) async {
    try {
      emit(const StudentLoading());
      final activeSessions = await _sessionRepository.getActiveSessions();
      final courses = await _sessionRepository.getCourses();
      final attendanceStatsJson =
          await _sessionRepository.getAttendanceStatistics();
      final attendanceStats =
          AttendanceStatistics.fromJson(attendanceStatsJson);

      emit(StudentDashboardLoaded(
        activeSessions: activeSessions,
        courses: courses,
        attendanceStats: attendanceStats,
      ));
    } catch (e) {
      emit(StudentFailure(e.toString()));
    }
  }

  Future<void> _onLoadStudentSessions(
    LoadStudentSessions event,
    Emitter<StudentState> emit,
  ) async {
    try {
      emit(const StudentLoading());
      final sessions = await _sessionRepository.getSessions(
        startDate: event.startDate,
        endDate: event.endDate,
        course: event.course,
      );
      final courses = await _sessionRepository.getCourses();

      emit(StudentSessionsLoaded(
        sessions: sessions,
        courses: courses,
      ));
    } catch (e) {
      emit(StudentFailure(e.toString()));
    }
  }

  Future<void> _onMarkAttendanceRequested(
    MarkAttendanceRequested event,
    Emitter<StudentState> emit,
  ) async {
    try {
      emit(const StudentLoading());

      final isOnline = await _connectivityService.isOnline();

      if (isOnline) {
        // Online mode: Mark attendance directly
        final record = await _sessionRepository.markAttendance(
          sessionId: event.sessionId,
          studentId: event.studentId,
          locationLatitude: event.locationLatitude,
          locationLongitude: event.locationLongitude,
          wifiSSID: event.wifiSSID,
          wifiBSSID: event.wifiBSSID,
          deviceId: event.deviceId,
        );

        emit(AttendanceMarked(record: record));
      } else {
        // Offline mode: Store attendance record locally
        final record = await _syncService.storeOfflineAttendance(
          sessionId: event.sessionId,
          studentId: event.studentId,
          studentName: event.studentName,
          locationLatitude: event.locationLatitude,
          locationLongitude: event.locationLongitude,
          wifiSSID: event.wifiSSID,
          wifiBSSID: event.wifiBSSID,
          deviceId: event.deviceId,
        );

        emit(AttendanceMarked(record: record, isOffline: true));
      }

      // Trigger background sync if there are pending records
      _syncService.startBackgroundSync();
    } catch (e) {
      emit(StudentFailure(e.toString()));
    }
  }
}
