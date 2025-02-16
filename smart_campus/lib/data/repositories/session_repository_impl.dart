import 'package:injectable/injectable.dart';
import 'package:uuid/uuid.dart';
import 'package:smart_campus/core/errors/app_error.dart';
import 'package:smart_campus/core/services/database_service.dart';
import 'package:smart_campus/core/services/location_service.dart';
import 'package:smart_campus/core/services/logger_service.dart';
import 'package:smart_campus/core/services/network_service.dart';
import 'package:smart_campus/core/services/wifi_service.dart';
import 'package:smart_campus/domain/entities/session.dart';
import 'package:smart_campus/domain/entities/attendance_record.dart';
import 'package:smart_campus/domain/repositories/session_repository.dart';
import 'package:smart_campus/data/datasources/local/session_local_datasource.dart';
import 'package:smart_campus/domain/entities/course.dart';
import 'package:smart_campus/services/device_service.dart';
import 'package:smart_campus/data/datasources/local/database_helper.dart';

@LazySingleton(as: SessionRepository)
class SessionRepositoryImpl implements SessionRepository {
  final DatabaseService _databaseService;
  final NetworkService _networkService;
  final LocationService _locationService;
  final WifiService _wifiService;
  final LoggerService _logger;
  final _uuid = const Uuid();
  final SessionLocalDataSource _localDataSource;
  final DeviceService _deviceService;
  final DatabaseHelper _databaseHelper;

  SessionRepositoryImpl(
    this._databaseService,
    this._networkService,
    this._locationService,
    this._wifiService,
    this._logger,
    this._localDataSource,
    this._deviceService,
    this._databaseHelper,
  );

  @override
  Future<List<Session>> getSessions({
    DateTime? startDate,
    DateTime? endDate,
    Course? course,
  }) async {
    return await _localDataSource.getSessions(
      startDate: startDate,
      endDate: endDate,
      course: course,
    );
  }

  @override
  Future<List<Session>> getActiveSessions() async {
    return await _localDataSource.getActiveSessions();
  }

  @override
  Future<Session> createSession({
    required Course course,
    required DateTime sessionDate,
    required DateTime startTime,
    required DateTime endTime,
    double? locationLatitude,
    double? locationLongitude,
    int? locationRadius,
    String? wifiSSID,
    String? wifiBSSID,
  }) async {
    return await _localDataSource.createSession(
      course: course,
      sessionDate: sessionDate,
      startTime: startTime,
      endTime: endTime,
      locationLatitude: locationLatitude,
      locationLongitude: locationLongitude,
      locationRadius: locationRadius,
      wifiSSID: wifiSSID,
      wifiBSSID: wifiBSSID,
    );
  }

  @override
  Future<Session> endSession(String sessionId) async {
    return await _localDataSource.endSession(sessionId);
  }

  @override
  Future<Session> getSessionById(String sessionId) async {
    return await _localDataSource.getSessionById(sessionId);
  }

  @override
  Future<List<AttendanceRecord>> getSessionAttendance(String sessionId) async {
    return await _localDataSource.getSessionAttendance(sessionId);
  }

  @override
  Future<AttendanceRecord> markAttendance({
    required String sessionId,
    required String studentId,
    double? locationLatitude,
    double? locationLongitude,
    String? wifiSSID,
    String? wifiBSSID,
    required String deviceId,
  }) async {
    return await _localDataSource.markAttendance(
      sessionId: sessionId,
      studentId: studentId,
      locationLatitude: locationLatitude,
      locationLongitude: locationLongitude,
      wifiSSID: wifiSSID,
      wifiBSSID: wifiBSSID,
      deviceId: deviceId,
    );
  }

  @override
  Future<List<Course>> getCourses() async {
    return await _localDataSource.getCourses();
  }

  @override
  Future<Course> getCourseById(String courseId) async {
    return await _localDataSource.getCourseById(courseId);
  }

  @override
  Future<List<Session>> getSessionHistory({
    String? teacherId,
    String? courseId,
    DateTime? startDate,
    DateTime? endDate,
  }) async {
    try {
      // Use getSessions from local data source with filters
      return await _localDataSource.getSessions(
        startDate: startDate,
        endDate: endDate,
        // Additional filters can be applied in the local data source
      );
    } catch (e) {
      throw AppError('Failed to get session history: $e');
    }
  }

  @override
  Future<Session> updateSession(Session session) async {
    try {
      // Get the current session to verify it exists
      await _localDataSource.getSessionById(session.id);

      // Update the session using the local data source
      // This will need to be implemented in the local data source
      final updatedSession = await _localDataSource.updateSession(session);

      return updatedSession;
    } catch (e) {
      throw AppError('Failed to update session: $e');
    }
  }

  @override
  Future<void> cancelSession(String sessionId) async {
    try {
      // Get the current session
      final session = await getSessionById(sessionId);

      // Check if session can be cancelled
      if (session.status != SessionStatus.active) {
        throw AppError('Only active sessions can be cancelled');
      }

      // Update session status to cancelled
      final updatedSession = session.copyWith(
        status: SessionStatus.cancelled,
        updatedAt: DateTime.now(),
      );

      await updateSession(updatedSession);
    } catch (e) {
      throw AppError('Failed to cancel session: $e');
    }
  }

  @override
  Future<AttendanceRecord> verifyAttendance(String attendanceId) async {
    try {
      // Get the attendance record
      final attendanceRecords =
          await _localDataSource.getSessionAttendance(attendanceId);
      final record = attendanceRecords.firstWhere(
        (record) => record.id == attendanceId,
        orElse: () => throw AppError('Attendance record not found'),
      );

      // Check if the record is pending verification
      if (record.verificationStatus !=
          VerificationStatus.pending.toString().split('.').last) {
        throw AppError('Attendance record is not pending verification');
      }

      // Get the session
      final session = await _localDataSource.getSessionById(record.sessionId);

      // Verify location if required
      if (session.locationLatitude != null &&
          session.locationLongitude != null &&
          session.locationRadius != null &&
          record.locationLatitude != null &&
          record.locationLongitude != null) {
        final locationVerified = await isLocationValid(
          sessionId: session.id,
          latitude: record.locationLatitude!,
          longitude: record.locationLongitude!,
        );

        if (!locationVerified) {
          throw AppError('Location verification failed');
        }
      }

      // Verify WiFi if required
      if (session.wifiSSID != null &&
          session.wifiBSSID != null &&
          record.wifiSSID != null &&
          record.wifiBSSID != null) {
        final wifiVerified = await isWifiValid(
          sessionId: session.id,
          ssid: record.wifiSSID!,
          bssid: record.wifiBSSID!,
        );

        if (!wifiVerified) {
          throw AppError('WiFi verification failed');
        }
      }

      // Update the attendance record
      final updatedRecord = record.copyWith(
        verificationStatus:
            VerificationStatus.verified.toString().split('.').last,
        updatedAt: DateTime.now(),
      );

      // Save the updated record
      await _localDataSource.updateAttendanceRecord(updatedRecord);

      return updatedRecord;
    } catch (e) {
      throw AppError('Failed to verify attendance: $e');
    }
  }

  @override
  Future<AttendanceRecord> rejectAttendance(
    String attendanceId, {
    required String reason,
  }) async {
    try {
      // Get the attendance record
      final attendanceRecords =
          await _localDataSource.getSessionAttendance(attendanceId);
      final record = attendanceRecords.firstWhere(
        (record) => record.id == attendanceId,
        orElse: () => throw AppError('Attendance record not found'),
      );

      // Check if the record can be rejected
      if (record.verificationStatus !=
          VerificationStatus.pending.toString().split('.').last) {
        throw AppError('Only pending attendance records can be rejected');
      }

      // Update the attendance record
      final updatedRecord = record.copyWith(
        verificationStatus:
            VerificationStatus.rejected.toString().split('.').last,
        rejectionReason: reason,
        updatedAt: DateTime.now(),
      );

      // Save the updated record
      await _localDataSource.updateAttendanceRecord(updatedRecord);

      return updatedRecord;
    } catch (e) {
      throw AppError('Failed to reject attendance: $e');
    }
  }

  @override
  Future<Map<String, dynamic>> getAttendanceStatistics({
    String? courseId,
    String? studentId,
    DateTime? startDate,
    DateTime? endDate,
  }) async {
    try {
      // Get all relevant sessions
      final sessions = await _localDataSource.getSessions(
        startDate: startDate,
        endDate: endDate,
      );

      // Get attendance records for these sessions
      int totalSessions = 0;
      int totalPresent = 0;
      int totalRejected = 0;
      int totalPending = 0;

      for (final session in sessions) {
        if (courseId != null && session.courseId != courseId) {
          continue;
        }

        totalSessions++;
        final records = await _localDataSource.getSessionAttendance(session.id);

        for (final record in records) {
          if (studentId != null && record.studentId != studentId) {
            continue;
          }

          switch (record.verificationStatus) {
            case VerificationStatus.verified:
              totalPresent++;
              break;
            case VerificationStatus.rejected:
              totalRejected++;
              break;
            case VerificationStatus.pending:
              totalPending++;
              break;
          }
        }
      }

      return {
        'totalSessions': totalSessions,
        'totalPresent': totalPresent,
        'totalRejected': totalRejected,
        'totalPending': totalPending,
        'attendancePercentage':
            totalSessions > 0 ? (totalPresent / totalSessions) * 100 : 0.0,
      };
    } catch (e) {
      throw AppError('Failed to get attendance statistics: $e');
    }
  }

  @override
  Future<bool> isStudentEligibleForAttendance({
    required String sessionId,
    required String studentId,
  }) async {
    try {
      // Get the session
      final session = await _localDataSource.getSessionById(sessionId);

      // Check if session is active
      if (session.status != SessionStatus.active) {
        return false;
      }

      // Check if student has already marked attendance
      final attendanceRecords =
          await _localDataSource.getSessionAttendance(sessionId);
      final hasMarkedAttendance =
          attendanceRecords.any((record) => record.studentId == studentId);

      if (hasMarkedAttendance) {
        return false;
      }

      // Get the course for the session
      final course = await _localDataSource.getCourseById(session.courseId);

      // Get the student
      final student = await _localDataSource.getStudentById(studentId);

      // Check if student is in the same department and year of study as the course
      return student.department == course.department &&
          student.yearOfStudy == course.yearOfStudy;
    } catch (e) {
      throw AppError('Failed to check student eligibility: $e');
    }
  }

  @override
  Future<bool> isLocationValid({
    required String sessionId,
    required double latitude,
    required double longitude,
  }) async {
    try {
      final session = await _localDataSource.getSessionById(sessionId);

      // If session doesn't require location validation
      if (session.locationLatitude == null ||
          session.locationLongitude == null ||
          session.locationRadius == null) {
        return true;
      }

      return _locationService.isWithinRadius(
        targetLatitude: session.locationLatitude!,
        targetLongitude: session.locationLongitude!,
        currentLatitude: latitude,
        currentLongitude: longitude,
        radiusInMeters: session.locationRadius!,
      );
    } catch (e) {
      throw AppError('Failed to validate location: $e');
    }
  }

  @override
  Future<bool> isWifiValid({
    required String sessionId,
    required String ssid,
    required String bssid,
  }) async {
    try {
      final session = await _localDataSource.getSessionById(sessionId);

      // If session doesn't require WiFi validation
      if (session.wifiSSID == null || session.wifiBSSID == null) {
        return true;
      }

      // Verify both SSID and BSSID match
      return session.wifiSSID == ssid && session.wifiBSSID == bssid;
    } catch (e) {
      throw AppError('Failed to validate WiFi connection: $e');
    }
  }

  @override
  Future<bool> isDeviceValid({
    required String sessionId,
    required String studentId,
    required String deviceId,
  }) async {
    try {
      // Check if device is registered to the student
      final deviceBindings = await _databaseHelper.query(
        'device_bindings',
        where: 'user_id = ? AND device_id = ? AND is_active = 1',
        whereArgs: [studentId, deviceId],
      );

      if (deviceBindings.isEmpty) {
        return false;
      }

      // Check if this device has been used for attendance in this session
      final existingAttendance = await _databaseHelper.query(
        'attendance_records',
        where: 'session_id = ? AND device_id = ?',
        whereArgs: [sessionId, deviceId],
      );

      // Return true if device is registered and hasn't been used in this session
      return existingAttendance.isEmpty;
    } catch (e) {
      throw AppError('Failed to verify device: $e');
    }
  }
}
