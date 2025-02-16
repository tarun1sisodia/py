import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:smart_campus/core/services/sync_service.dart';
import 'package:smart_campus/domain/repositories/session_repository.dart';
import 'package:smart_campus/presentation/bloc/session/session_event.dart';
import 'package:smart_campus/presentation/bloc/session/session_state.dart';
import 'package:connectivity_plus/connectivity_plus.dart';

class SessionBloc extends Bloc<SessionEvent, SessionState> {
  final SessionRepository _sessionRepository;
  final SyncService _syncService;
  final Connectivity _connectivity;

  SessionBloc({
    required SessionRepository sessionRepository,
    required SyncService syncService,
    required Connectivity connectivity,
  })  : _sessionRepository = sessionRepository,
        _syncService = syncService,
        _connectivity = connectivity,
        super(const SessionInitial()) {
    on<LoadSessions>(_onLoadSessions);
    on<LoadActiveSessions>(_onLoadActiveSessions);
    on<CreateSession>(_onCreateSession);
    on<EndSession>(_onEndSession);
    on<LoadSessionDetails>(_onLoadSessionDetails);
    on<MarkAttendance>(_onMarkAttendance);
    on<VerifyAttendance>(_onVerifyAttendance);
    on<RejectAttendance>(_onRejectAttendance);
  }

  Future<void> _onLoadSessions(
    LoadSessions event,
    Emitter<SessionState> emit,
  ) async {
    try {
      emit(const SessionLoading());
      final sessions = await _sessionRepository.getSessions(
        startDate: event.startDate,
        endDate: event.endDate,
        course: event.course,
      );
      final courses = await _sessionRepository.getCourses();
      emit(SessionLoadSuccess(sessions: sessions, courses: courses));
    } catch (e) {
      emit(SessionFailure(e.toString()));
    }
  }

  Future<void> _onLoadActiveSessions(
    LoadActiveSessions event,
    Emitter<SessionState> emit,
  ) async {
    try {
      emit(const SessionLoading());
      final activeSessions = await _sessionRepository.getActiveSessions();
      final courses = await _sessionRepository.getCourses();
      emit(ActiveSessionsLoaded(
        activeSessions: activeSessions,
        courses: courses,
      ));
    } catch (e) {
      emit(SessionFailure(e.toString()));
    }
  }

  Future<void> _onCreateSession(
    CreateSession event,
    Emitter<SessionState> emit,
  ) async {
    try {
      emit(const SessionLoading());
      final session = await _sessionRepository.createSession(
        course: event.course,
        sessionDate: event.sessionDate,
        startTime: event.startTime,
        endTime: event.endTime,
        locationLatitude: event.locationLatitude,
        locationLongitude: event.locationLongitude,
        locationRadius: event.locationRadius,
        wifiSSID: event.wifiSSID,
        wifiBSSID: event.wifiBSSID,
      );
      emit(SessionCreated(session: session, course: event.course));
    } catch (e) {
      emit(SessionFailure(e.toString()));
    }
  }

  Future<void> _onEndSession(
    EndSession event,
    Emitter<SessionState> emit,
  ) async {
    try {
      emit(const SessionLoading());
      final session = await _sessionRepository.endSession(event.sessionId);
      emit(SessionEnded(session));
    } catch (e) {
      emit(SessionFailure(e.toString()));
    }
  }

  Future<void> _onLoadSessionDetails(
    LoadSessionDetails event,
    Emitter<SessionState> emit,
  ) async {
    try {
      emit(const SessionLoading());
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
      emit(SessionFailure(e.toString()));
    }
  }

  Future<void> _onMarkAttendance(
    MarkAttendance event,
    Emitter<SessionState> emit,
  ) async {
    try {
      emit(const SessionLoading());

      // Verify device ID first
      final isDeviceValid = await _sessionRepository.isDeviceValid(
        sessionId: event.sessionId,
        studentId: event.studentId,
        deviceId: event.deviceId,
      );

      if (!isDeviceValid) {
        emit(const SessionFailure(
            'Invalid device ID. Please use your registered device.'));
        return;
      }

      final connectivityResult = await _connectivity.checkConnectivity();
      if (connectivityResult == ConnectivityResult.none) {
        // Save attendance record offline
        await _syncService.storeOfflineAttendance(
          sessionId: event.sessionId,
          studentId: event.studentId,
          studentName: event.studentName,
          locationLatitude: event.locationLatitude,
          locationLongitude: event.locationLongitude,
          wifiSSID: event.wifiSSID,
          wifiBSSID: event.wifiBSSID,
          deviceId: event.deviceId,
        );
        emit(const AttendanceMarkedOffline());
      } else {
        // Mark attendance online
        final attendanceRecord = await _sessionRepository.markAttendance(
          sessionId: event.sessionId,
          studentId: event.studentId,
          locationLatitude: event.locationLatitude,
          locationLongitude: event.locationLongitude,
          wifiSSID: event.wifiSSID,
          wifiBSSID: event.wifiBSSID,
          deviceId: event.deviceId,
        );
        emit(AttendanceMarked(attendanceRecord));
      }
    } catch (e) {
      emit(SessionFailure(e.toString()));
    }
  }

  Future<void> _onVerifyAttendance(
    VerifyAttendance event,
    Emitter<SessionState> emit,
  ) async {
    try {
      emit(const SessionLoading());
      final attendanceRecord =
          await _sessionRepository.verifyAttendance(event.attendanceId);
      emit(AttendanceVerified(attendanceRecord));
    } catch (e) {
      emit(SessionFailure(e.toString()));
    }
  }

  Future<void> _onRejectAttendance(
    RejectAttendance event,
    Emitter<SessionState> emit,
  ) async {
    try {
      emit(const SessionLoading());
      final attendanceRecord = await _sessionRepository.rejectAttendance(
        event.attendanceId,
        reason: event.reason,
      );
      emit(AttendanceRejected(attendanceRecord));
    } catch (e) {
      emit(SessionFailure(e.toString()));
    }
  }
}
