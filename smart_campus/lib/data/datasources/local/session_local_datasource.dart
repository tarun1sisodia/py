import 'package:injectable/injectable.dart';
import 'package:smart_campus/data/datasources/local/database_helper.dart';
import 'package:smart_campus/domain/entities/session.dart';
import 'package:smart_campus/domain/entities/course.dart';
import 'package:smart_campus/domain/entities/attendance_record.dart';
import 'package:smart_campus/domain/entities/user.dart';
import 'package:smart_campus/core/error/app_error.dart';
import 'package:uuid/uuid.dart';
import 'dart:math';

abstract class SessionLocalDataSource {
  Future<List<Session>> getSessions({
    DateTime? startDate,
    DateTime? endDate,
    Course? course,
  });

  Future<List<Session>> getActiveSessions();

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
  });

  Future<Session> endSession(String sessionId);

  Future<Session> getSessionById(String sessionId);

  Future<List<AttendanceRecord>> getSessionAttendance(String sessionId);

  Future<AttendanceRecord> markAttendance({
    required String sessionId,
    required String studentId,
    double? locationLatitude,
    double? locationLongitude,
    String? wifiSSID,
    String? wifiBSSID,
    required String deviceId,
  });

  Future<List<Course>> getCourses();

  Future<Course> getCourseById(String courseId);

  Future<User> getStudentById(String studentId);

  Future<AttendanceRecord> updateAttendanceRecord(AttendanceRecord record);

  Future<Session> updateSession(Session session);

  Future<Session> cancelSession(String sessionId);

  Future<AttendanceRecord> verifyAttendance(String attendanceId);

  Future<AttendanceRecord> rejectAttendance(String attendanceId,
      {required String reason});

  Future<Map<String, dynamic>> getAttendanceStatistics({
    String? courseId,
    String? studentId,
    DateTime? startDate,
    DateTime? endDate,
  });
}

@LazySingleton(as: SessionLocalDataSource)
class SessionLocalDataSourceImpl implements SessionLocalDataSource {
  final DatabaseHelper _databaseHelper;
  final _uuid = const Uuid();

  SessionLocalDataSourceImpl(this._databaseHelper);

  @override
  Future<List<Session>> getSessions({
    DateTime? startDate,
    DateTime? endDate,
    Course? course,
  }) async {
    try {
      final conditions = <String>[];
      final args = <dynamic>[];

      if (startDate != null) {
        conditions.add('session_date >= ?');
        args.add(startDate.toIso8601String());
      }
      if (endDate != null) {
        conditions.add('session_date <= ?');
        args.add(endDate.toIso8601String());
      }
      if (course != null) {
        conditions.add('course_id = ?');
        args.add(course.id);
      }

      final results = await _databaseHelper.query(
        'attendance_sessions',
        where: conditions.isEmpty ? null : conditions.join(' AND '),
        whereArgs: args.isEmpty ? null : args,
        orderBy: 'session_date DESC',
      );

      return results
          .map((map) => Session.fromJson(_mapToSessionJson(map)))
          .toList();
    } catch (e) {
      throw AppError('Failed to get sessions: $e');
    }
  }

  @override
  Future<List<Session>> getActiveSessions() async {
    try {
      final results = await _databaseHelper.query(
        'attendance_sessions',
        where: 'status = ?',
        whereArgs: [SessionStatus.active.toString().split('.').last],
      );

      return results
          .map((map) => Session.fromJson(_mapToSessionJson(map)))
          .toList();
    } catch (e) {
      throw AppError('Failed to get active sessions: $e');
    }
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
    try {
      await _validateCourseExists(course.id);
      _validateWifiCredentials(wifiSSID, wifiBSSID);
      _validateLocation(locationLatitude, locationLongitude, locationRadius);

      if (sessionDate
          .isBefore(DateTime.now().subtract(const Duration(days: 1)))) {
        throw AppError('Session date cannot be in the past');
      }

      final session = Session(
        id: _uuid.v4(),
        teacherId:
            'teacher_1', // This should come from auth service in production
        courseId: course.id,
        sessionDate: sessionDate,
        startTime: startTime,
        endTime: endTime,
        locationLatitude: locationLatitude,
        locationLongitude: locationLongitude,
        locationRadius: locationRadius,
        wifiSSID: wifiSSID,
        wifiBSSID: wifiBSSID,
        status: SessionStatus.active,
        createdAt: DateTime.now(),
        updatedAt: DateTime.now(),
      );

      _validateSession(session);

      final result = await _databaseHelper.insert(
        'attendance_sessions',
        _sessionToMap(session),
      );

      if (result <= 0) {
        throw AppError('Failed to create session in database');
      }

      return session;
    } catch (e) {
      if (e is AppError) {
        rethrow;
      }
      throw AppError('Failed to create session: $e');
    }
  }

  @override
  Future<Session> endSession(String sessionId) async {
    try {
      final session = await getSessionById(sessionId);

      if (session.status != SessionStatus.active) {
        throw AppError('Session is not active');
      }

      final updatedSession = session.copyWith(
        status: SessionStatus.completed,
        updatedAt: DateTime.now(),
      );

      await _databaseHelper.update(
        'attendance_sessions',
        _sessionToMap(updatedSession),
        where: 'id = ?',
        whereArgs: [sessionId],
      );

      return updatedSession;
    } catch (e) {
      throw AppError('Failed to end session: $e');
    }
  }

  @override
  Future<Session> getSessionById(String sessionId) async {
    try {
      final results = await _databaseHelper.query(
        'attendance_sessions',
        where: 'id = ?',
        whereArgs: [sessionId],
      );

      if (results.isEmpty) {
        throw AppError('Session not found');
      }

      return Session.fromJson(_mapToSessionJson(results.first));
    } catch (e) {
      throw AppError('Failed to get session: $e');
    }
  }

  @override
  Future<List<AttendanceRecord>> getSessionAttendance(String sessionId) async {
    try {
      final results = await _databaseHelper.query(
        'attendance_records',
        where: 'session_id = ?',
        whereArgs: [sessionId],
        orderBy: 'marked_at DESC',
      );

      return results
          .map((map) =>
              AttendanceRecord.fromJson(_mapToAttendanceRecordJson(map)))
          .toList();
    } catch (e) {
      throw AppError('Failed to get session attendance: $e');
    }
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
    try {
      await _validateSessionExists(sessionId);
      _validateWifiCredentials(wifiSSID, wifiBSSID);
      _validateDeviceId(deviceId);
      _validateLocation(locationLatitude, locationLongitude, null);

      if (!await _isSessionActive(sessionId)) {
        throw AppError('Session is not active');
      }

      if (!await _isWithinSessionTimeWindow(sessionId)) {
        throw AppError(
            'Attendance can only be marked during the session time window');
      }

      if (!await _isWithinSessionLocation(
          sessionId, locationLatitude, locationLongitude)) {
        throw AppError(
            'Attendance must be marked within the session location radius');
      }

      const studentId =
          'student_1'; // This should come from auth service in production
      await _validateStudentExists(studentId);

      if (await _hasExistingAttendance(sessionId, studentId)) {
        throw AppError('Attendance already marked for this session');
      }

      // Get student name from user table
      final student = await getStudentById(studentId);
      final studentName = student.fullName;

      final attendanceRecord = AttendanceRecord(
        id: _uuid.v4(),
        sessionId: sessionId,
        studentId: studentId,
        studentName: studentName,
        markedAt: DateTime.now(),
        locationLatitude: locationLatitude,
        locationLongitude: locationLongitude,
        wifiSSID: wifiSSID,
        wifiBSSID: wifiBSSID,
        deviceId: deviceId,
        verificationStatus: VerificationStatus.pending.toString(),
        createdAt: DateTime.now(),
        updatedAt: DateTime.now(),
      );

      _validateAttendanceRecord(attendanceRecord);

      final result = await _databaseHelper.insert(
        'attendance_records',
        _attendanceRecordToMap(attendanceRecord),
      );

      if (result <= 0) {
        throw AppError('Failed to save attendance record in database');
      }

      return attendanceRecord;
    } catch (e) {
      if (e is AppError) {
        rethrow;
      }
      throw AppError('Failed to mark attendance: $e');
    }
  }

  @override
  Future<List<Course>> getCourses() async {
    try {
      final results = await _databaseHelper.query('courses');
      return results
          .map((map) => Course.fromJson(_mapToCourseJson(map)))
          .toList();
    } catch (e) {
      throw AppError('Failed to get courses: $e');
    }
  }

  @override
  Future<Course> getCourseById(String courseId) async {
    try {
      final results = await _databaseHelper.query(
        'courses',
        where: 'id = ?',
        whereArgs: [courseId],
      );

      if (results.isEmpty) {
        throw AppError('Course not found');
      }

      return Course.fromJson(_mapToCourseJson(results.first));
    } catch (e) {
      throw AppError('Failed to get course: $e');
    }
  }

  @override
  Future<User> getStudentById(String studentId) async {
    try {
      final results = await _databaseHelper.query(
        'users',
        where: 'id = ? AND role = ?',
        whereArgs: [studentId, 'student'],
      );

      if (results.isEmpty) {
        throw AppError('Student not found');
      }

      return User.fromJson(_mapToUserJson(results.first));
    } catch (e) {
      throw AppError('Failed to get student: $e');
    }
  }

  @override
  Future<AttendanceRecord> updateAttendanceRecord(
      AttendanceRecord record) async {
    try {
      // First check if the record exists
      final results = await _databaseHelper.query(
        'attendance_records',
        where: 'id = ?',
        whereArgs: [record.id],
      );

      if (results.isEmpty) {
        throw AppError('Attendance record not found');
      }

      final values = {
        'session_id': record.sessionId,
        'student_id': record.studentId,
        'marked_at': record.markedAt.toIso8601String(),
        'wifi_ssid': record.wifiSSID,
        'wifi_bssid': record.wifiBSSID,
        'location_latitude': record.locationLatitude,
        'location_longitude': record.locationLongitude,
        'device_id': record.deviceId,
        'verification_status':
            record.verificationStatus.toString().split('.').last,
        'rejection_reason': record.rejectionReason,
        'updated_at': DateTime.now().toIso8601String(),
      };

      await _databaseHelper.update(
        'attendance_records',
        values,
        where: 'id = ?',
        whereArgs: [record.id],
      );

      return record.copyWith(updatedAt: DateTime.now());
    } catch (e) {
      throw AppError('Failed to update attendance record: $e');
    }
  }

  @override
  Future<Session> updateSession(Session session) async {
    try {
      // First check if the session exists
      final results = await _databaseHelper.query(
        'attendance_sessions',
        where: 'id = ?',
        whereArgs: [session.id],
      );

      if (results.isEmpty) {
        throw AppError('Session not found');
      }

      final values = {
        'teacher_id': session.teacherId,
        'course_id': session.courseId,
        'session_date': session.sessionDate.toIso8601String(),
        'start_time': session.startTime.toIso8601String(),
        'end_time': session.endTime.toIso8601String(),
        'wifi_ssid': session.wifiSSID,
        'wifi_bssid': session.wifiBSSID,
        'location_latitude': session.locationLatitude,
        'location_longitude': session.locationLongitude,
        'location_radius': session.locationRadius,
        'status': session.status.toString().split('.').last,
        'updated_at': DateTime.now().toIso8601String(),
      };

      await _databaseHelper.update(
        'attendance_sessions',
        values,
        where: 'id = ?',
        whereArgs: [session.id],
      );

      return session.copyWith(updatedAt: DateTime.now());
    } catch (e) {
      throw AppError('Failed to update session: $e');
    }
  }

  @override
  Future<Session> cancelSession(String sessionId) async {
    try {
      final session = await getSessionById(sessionId);

      if (session.status == SessionStatus.cancelled) {
        throw AppError('Session is already cancelled');
      }

      final updatedSession = session.copyWith(
        status: SessionStatus.cancelled,
        updatedAt: DateTime.now(),
      );

      await _databaseHelper.update(
        'attendance_sessions',
        _sessionToMap(updatedSession),
        where: 'id = ?',
        whereArgs: [sessionId],
      );

      return updatedSession;
    } catch (e) {
      throw AppError('Failed to cancel session: $e');
    }
  }

  @override
  Future<AttendanceRecord> verifyAttendance(String attendanceId) async {
    try {
      await _validateAttendanceRecordExists(attendanceId);

      final results = await _databaseHelper.query(
        'attendance_records',
        where: 'id = ?',
        whereArgs: [attendanceId],
      );

      final record =
          AttendanceRecord.fromJson(_mapToAttendanceRecordJson(results.first));

      if (record.verificationStatus != VerificationStatus.pending) {
        throw AppError('Only pending attendance records can be verified');
      }

      // Check if the session is still valid
      final session = await getSessionById(record.sessionId);
      if (session.status == SessionStatus.cancelled) {
        throw AppError('Cannot verify attendance for a cancelled session');
      }

      final updatedRecord = record.copyWith(
        verificationStatus: VerificationStatus.verified.toString(),
        updatedAt: DateTime.now(),
      );

      final result = await _databaseHelper.update(
        'attendance_records',
        _attendanceRecordToMap(updatedRecord),
        where: 'id = ?',
        whereArgs: [attendanceId],
      );

      if (result <= 0) {
        throw AppError('Failed to update attendance record in database');
      }

      return updatedRecord;
    } catch (e) {
      if (e is AppError) {
        rethrow;
      }
      throw AppError('Failed to verify attendance: $e');
    }
  }

  @override
  Future<AttendanceRecord> rejectAttendance(
    String attendanceId, {
    required String reason,
  }) async {
    try {
      if (reason.trim().isEmpty) {
        throw AppError('Rejection reason cannot be empty');
      }

      await _validateAttendanceRecordExists(attendanceId);

      final results = await _databaseHelper.query(
        'attendance_records',
        where: 'id = ?',
        whereArgs: [attendanceId],
      );

      final record =
          AttendanceRecord.fromJson(_mapToAttendanceRecordJson(results.first));

      if (record.verificationStatus != VerificationStatus.pending) {
        throw AppError('Only pending attendance records can be rejected');
      }

      // Check if the session is still valid
      final session = await getSessionById(record.sessionId);
      if (session.status == SessionStatus.cancelled) {
        throw AppError('Cannot reject attendance for a cancelled session');
      }

      final updatedRecord = record.copyWith(
        verificationStatus: VerificationStatus.rejected.toString(),
        rejectionReason: reason.trim(),
        updatedAt: DateTime.now(),
      );

      final result = await _databaseHelper.update(
        'attendance_records',
        _attendanceRecordToMap(updatedRecord),
        where: 'id = ?',
        whereArgs: [attendanceId],
      );

      if (result <= 0) {
        throw AppError('Failed to update attendance record in database');
      }

      return updatedRecord;
    } catch (e) {
      if (e is AppError) {
        rethrow;
      }
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
      final sessions = await getSessions(
        startDate: startDate,
        endDate: endDate,
      );

      int totalSessions = 0;
      int totalPresent = 0;
      int totalRejected = 0;
      int totalPending = 0;

      for (final session in sessions) {
        if (courseId != null && session.courseId != courseId) {
          continue;
        }

        totalSessions++;
        final records = await getSessionAttendance(session.id);

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
        'attendancePercentage': totalSessions > 0
            ? (totalPresent / totalSessions * 100).toStringAsFixed(2)
            : '0.00',
      };
    } catch (e) {
      throw AppError('Failed to get attendance statistics: $e');
    }
  }

  // Helper methods to convert between database and entity formats
  Map<String, dynamic> _sessionToMap(Session session) {
    return {
      'id': session.id,
      'teacher_id': session.teacherId,
      'course_id': session.courseId,
      'session_date': session.sessionDate.toIso8601String(),
      'start_time': session.startTime.toIso8601String(),
      'end_time': session.endTime.toIso8601String(),
      'wifi_ssid': session.wifiSSID,
      'wifi_bssid': session.wifiBSSID,
      'location_latitude': session.locationLatitude,
      'location_longitude': session.locationLongitude,
      'location_radius': session.locationRadius,
      'status': session.status.toString().split('.').last,
      'created_at': session.createdAt.toIso8601String(),
      'updated_at': session.updatedAt.toIso8601String(),
    };
  }

  Map<String, dynamic> _attendanceRecordToMap(AttendanceRecord record) {
    return {
      'id': record.id,
      'session_id': record.sessionId,
      'student_id': record.studentId,
      'marked_at': record.markedAt.toIso8601String(),
      'wifi_ssid': record.wifiSSID,
      'wifi_bssid': record.wifiBSSID,
      'location_latitude': record.locationLatitude,
      'location_longitude': record.locationLongitude,
      'device_id': record.deviceId,
      'verification_status':
          record.verificationStatus.toString().split('.').last,
      'rejection_reason': record.rejectionReason,
      'created_at': record.createdAt.toIso8601String(),
      'updated_at': record.updatedAt.toIso8601String(),
    };
  }

  Map<String, dynamic> _mapToSessionJson(Map<String, dynamic> map) {
    return {
      'id': map['id'],
      'teacher_id': map['teacher_id'],
      'course_id': map['course_id'],
      'session_date': map['session_date'],
      'start_time': map['start_time'],
      'end_time': map['end_time'],
      'wifi_ssid': map['wifi_ssid'],
      'wifi_bssid': map['wifi_bssid'],
      'location_latitude': map['location_latitude'],
      'location_longitude': map['location_longitude'],
      'location_radius': map['location_radius'],
      'status': map['status'],
      'created_at': map['created_at'],
      'updated_at': map['updated_at'],
    };
  }

  Map<String, dynamic> _mapToAttendanceRecordJson(Map<String, dynamic> map) {
    return {
      'id': map['id'],
      'session_id': map['session_id'],
      'student_id': map['student_id'],
      'marked_at': map['marked_at'],
      'wifi_ssid': map['wifi_ssid'],
      'wifi_bssid': map['wifi_bssid'],
      'location_latitude': map['location_latitude'],
      'location_longitude': map['location_longitude'],
      'device_id': map['device_id'],
      'verification_status': map['verification_status'],
      'rejection_reason': map['rejection_reason'],
      'created_at': map['created_at'],
      'updated_at': map['updated_at'],
    };
  }

  Map<String, dynamic> _mapToCourseJson(Map<String, dynamic> map) {
    return {
      'id': map['id'],
      'course_code': map['code'],
      'course_name': map['name'],
      'department': map['department'],
      'year_of_study': map['year_of_study'],
      'semester': map['semester'],
      'created_at': map['created_at'],
      'updated_at': map['updated_at'],
    };
  }

  Map<String, dynamic> _mapToUserJson(Map<String, dynamic> map) {
    return {
      'id': map['id'],
      'email': map['email'],
      'full_name': map['full_name'],
      'role': map['role'],
      'department': map['department'],
      'year_of_study': map['year_of_study'],
      'enrollment_number': map['enrollment_number'],
      'employee_id': map['employee_id'],
      'device_id': map['device_id'],
      'created_at': map['created_at'],
      'updated_at': map['updated_at'],
    };
  }

  // Validation methods
  void _validateSession(Session session) {
    if (session.startTime.isAfter(session.endTime)) {
      throw AppError('Session end time must be after start time');
    }

    if (session.locationRadius != null && session.locationRadius! <= 0) {
      throw AppError('Location radius must be positive');
    }

    if (session.locationLatitude != null &&
        (session.locationLatitude! < -90 || session.locationLatitude! > 90)) {
      throw AppError('Invalid latitude value');
    }

    if (session.locationLongitude != null &&
        (session.locationLongitude! < -180 ||
            session.locationLongitude! > 180)) {
      throw AppError('Invalid longitude value');
    }
  }

  void _validateAttendanceRecord(AttendanceRecord record) {
    if (record.locationLatitude != null &&
        (record.locationLatitude! < -90 || record.locationLatitude! > 90)) {
      throw AppError('Invalid latitude value');
    }

    if (record.locationLongitude != null &&
        (record.locationLongitude! < -180 || record.locationLongitude! > 180)) {
      throw AppError('Invalid longitude value');
    }

    if (record.deviceId.isEmpty) {
      throw AppError('Device ID is required');
    }
  }

  Future<void> _validateSessionExists(String sessionId) async {
    final results = await _databaseHelper.query(
      'attendance_sessions',
      where: 'id = ?',
      whereArgs: [sessionId],
    );

    if (results.isEmpty) {
      throw AppError('Session not found');
    }
  }

  Future<void> _validateAttendanceRecordExists(String recordId) async {
    final results = await _databaseHelper.query(
      'attendance_records',
      where: 'id = ?',
      whereArgs: [recordId],
    );

    if (results.isEmpty) {
      throw AppError('Attendance record not found');
    }
  }

  Future<void> _validateStudentExists(String studentId) async {
    final results = await _databaseHelper.query(
      'users',
      where: 'id = ? AND role = ?',
      whereArgs: [studentId, 'student'],
    );

    if (results.isEmpty) {
      throw AppError('Student not found');
    }
  }

  Future<void> _validateCourseExists(String courseId) async {
    final results = await _databaseHelper.query(
      'courses',
      where: 'id = ?',
      whereArgs: [courseId],
    );

    if (results.isEmpty) {
      throw AppError('Course not found');
    }
  }

  Future<bool> _hasExistingAttendance(
      String sessionId, String studentId) async {
    final results = await _databaseHelper.query(
      'attendance_records',
      where: 'session_id = ? AND student_id = ?',
      whereArgs: [sessionId, studentId],
    );

    return results.isNotEmpty;
  }

  Future<bool> _isSessionActive(String sessionId) async {
    final results = await _databaseHelper.query(
      'attendance_sessions',
      where: 'id = ? AND status = ?',
      whereArgs: [sessionId, SessionStatus.active.toString().split('.').last],
    );

    return results.isNotEmpty;
  }

  Future<bool> _isWithinSessionTimeWindow(String sessionId) async {
    final results = await _databaseHelper.query(
      'attendance_sessions',
      where: 'id = ?',
      whereArgs: [sessionId],
    );

    if (results.isEmpty) {
      return false;
    }

    final session = Session.fromJson(_mapToSessionJson(results.first));
    final now = DateTime.now();

    return now.isAfter(session.startTime) && now.isBefore(session.endTime);
  }

  void _validateWifiCredentials(String? ssid, String? bssid) {
    if ((ssid != null && ssid.isEmpty) || (bssid != null && bssid.isEmpty)) {
      throw AppError('WiFi credentials cannot be empty if provided');
    }

    if ((ssid != null && bssid == null) || (ssid == null && bssid != null)) {
      throw AppError('Both SSID and BSSID must be provided together');
    }

    if (bssid != null &&
        !RegExp(r'^([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})$').hasMatch(bssid)) {
      throw AppError('Invalid BSSID format');
    }
  }

  void _validateDeviceId(String deviceId) {
    if (deviceId.trim().isEmpty) {
      throw AppError('Device ID is required');
    }

    // Add additional device ID format validation if needed
    if (!RegExp(r'^[a-zA-Z0-9-_]+$').hasMatch(deviceId)) {
      throw AppError('Invalid device ID format');
    }
  }

  void _validateLocation(double? latitude, double? longitude, int? radius) {
    if ((latitude != null && longitude == null) ||
        (latitude == null && longitude != null)) {
      throw AppError('Both latitude and longitude must be provided together');
    }

    if (latitude != null && (latitude < -90 || latitude > 90)) {
      throw AppError('Invalid latitude value');
    }

    if (longitude != null && (longitude < -180 || longitude > 180)) {
      throw AppError('Invalid longitude value');
    }

    if (radius != null) {
      if (radius <= 0) {
        throw AppError('Location radius must be positive');
      }
      if (latitude == null || longitude == null) {
        throw AppError(
            'Location coordinates are required when radius is specified');
      }
    }
  }

  Future<bool> _isWithinSessionLocation(
    String sessionId,
    double? latitude,
    double? longitude,
  ) async {
    if (latitude == null || longitude == null) {
      return true; // Skip location check if coordinates are not provided
    }

    final session = await getSessionById(sessionId);
    if (session.locationLatitude == null ||
        session.locationLongitude == null ||
        session.locationRadius == null) {
      return true; // Skip location check if session doesn't have location constraints
    }

    final distance = _calculateDistance(
      latitude,
      longitude,
      session.locationLatitude!,
      session.locationLongitude!,
    );

    return distance <= session.locationRadius!;
  }

  double _calculateDistance(
    double lat1,
    double lon1,
    double lat2,
    double lon2,
  ) {
    const double earthRadius = 6371000; // Earth's radius in meters
    final double lat1Rad = _degreesToRadians(lat1);
    final double lat2Rad = _degreesToRadians(lat2);
    final double deltaLat = _degreesToRadians(lat2 - lat1);
    final double deltaLon = _degreesToRadians(lon2 - lon1);

    final double a = sin(deltaLat / 2) * sin(deltaLat / 2) +
        cos(lat1Rad) * cos(lat2Rad) * sin(deltaLon / 2) * sin(deltaLon / 2);
    final double c = 2 * atan2(sqrt(a), sqrt(1 - a));

    return earthRadius * c;
  }

  double _degreesToRadians(double degrees) {
    return degrees * pi / 180;
  }
}
