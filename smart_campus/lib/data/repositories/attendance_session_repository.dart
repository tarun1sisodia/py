import 'package:injectable/injectable.dart';
import 'package:smart_campus/core/services/database_service.dart';
import 'package:smart_campus/core/services/logger_service.dart';
import 'package:smart_campus/data/models/attendance_session_model.dart';
import 'package:smart_campus/data/repositories/base_repository.dart';

@lazySingleton
class AttendanceSessionRepository
    extends BaseRepository<AttendanceSessionModel> {
  static const _tableName = 'attendance_sessions';

  AttendanceSessionRepository(
    DatabaseService databaseService,
    LoggerService logger,
  ) : super(databaseService, logger, _tableName);

  @override
  AttendanceSessionModel fromDatabase(Map<String, dynamic> data) =>
      AttendanceSessionModel.fromDatabase(data);

  Future<List<AttendanceSessionModel>> getActiveSessions() async {
    try {
      return await getAll(
        where: 'status = ?',
        whereArgs: ['active'],
        orderBy: 'start_time DESC',
      );
    } catch (e) {
      logger.error('Error getting active sessions', e);
      rethrow;
    }
  }

  Future<List<AttendanceSessionModel>> getSessionsByTeacher(
    String teacherId, {
    AttendanceSessionStatus? status,
    String? date,
  }) async {
    try {
      String where = 'teacher_id = ?';
      final whereArgs = [teacherId];

      if (status != null) {
        where += ' AND status = ?';
        whereArgs.add(status.toString().split('.').last);
      }

      if (date != null) {
        where += ' AND date = ?';
        whereArgs.add(date);
      }

      return await getAll(
        where: where,
        whereArgs: whereArgs,
        orderBy: 'start_time DESC',
      );
    } catch (e) {
      logger.error('Error getting sessions by teacher', e);
      rethrow;
    }
  }

  Future<List<AttendanceSessionModel>> getSessionsByCourse(
    String courseId, {
    AttendanceSessionStatus? status,
    String? date,
  }) async {
    try {
      String where = 'course_id = ?';
      final whereArgs = [courseId];

      if (status != null) {
        where += ' AND status = ?';
        whereArgs.add(status.toString().split('.').last);
      }

      if (date != null) {
        where += ' AND date = ?';
        whereArgs.add(date);
      }

      return await getAll(
        where: where,
        whereArgs: whereArgs,
        orderBy: 'start_time DESC',
      );
    } catch (e) {
      logger.error('Error getting sessions by course', e);
      rethrow;
    }
  }

  Future<List<AttendanceSessionModel>> getSessionsByDate(
    String date, {
    AttendanceSessionStatus? status,
  }) async {
    try {
      String where = 'date = ?';
      final whereArgs = [date];

      if (status != null) {
        where += ' AND status = ?';
        whereArgs.add(status.toString().split('.').last);
      }

      return await getAll(
        where: where,
        whereArgs: whereArgs,
        orderBy: 'start_time DESC',
      );
    } catch (e) {
      logger.error('Error getting sessions by date', e);
      rethrow;
    }
  }

  Future<bool> hasActiveSession(String courseId) async {
    try {
      final sessions = await getAll(
        where: 'course_id = ? AND status = ?',
        whereArgs: [courseId, 'active'],
      );
      return sessions.isNotEmpty;
    } catch (e) {
      logger.error('Error checking active session', e);
      rethrow;
    }
  }

  Future<void> endSession(String sessionId) async {
    try {
      await update(
        AttendanceSessionModel.fromDatabase({
          'id': sessionId,
          'status': 'completed',
          'updated_at': DateTime.now().millisecondsSinceEpoch,
        }),
      );
    } catch (e) {
      logger.error('Error ending session', e);
      rethrow;
    }
  }

  Future<void> cancelSession(String sessionId) async {
    try {
      await update(
        AttendanceSessionModel.fromDatabase({
          'id': sessionId,
          'status': 'cancelled',
          'updated_at': DateTime.now().millisecondsSinceEpoch,
        }),
      );
    } catch (e) {
      logger.error('Error cancelling session', e);
      rethrow;
    }
  }
}
