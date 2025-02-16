import 'package:smart_campus/domain/entities/session.dart';
import 'package:smart_campus/domain/entities/course.dart';
import 'package:smart_campus/domain/entities/attendance_record.dart';

abstract class SessionRepository {
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

  Future<List<Session>> getSessionHistory({
    String? teacherId,
    String? courseId,
    DateTime? startDate,
    DateTime? endDate,
  });

  Future<Session> updateSession(Session session);

  Future<void> cancelSession(String sessionId);

  Future<AttendanceRecord> verifyAttendance(String attendanceId);

  Future<AttendanceRecord> rejectAttendance(
    String attendanceId, {
    required String reason,
  });

  Future<Map<String, dynamic>> getAttendanceStatistics({
    String? courseId,
    String? studentId,
    DateTime? startDate,
    DateTime? endDate,
  });

  Future<bool> isStudentEligibleForAttendance({
    required String sessionId,
    required String studentId,
  });

  Future<bool> isLocationValid({
    required String sessionId,
    required double latitude,
    required double longitude,
  });

  Future<bool> isWifiValid({
    required String sessionId,
    required String ssid,
    required String bssid,
  });

  Future<bool> isDeviceValid({
    required String sessionId,
    required String studentId,
    required String deviceId,
  });
}
