import 'package:injectable/injectable.dart';
import 'package:smart_campus/core/services/database_service.dart';
import 'package:smart_campus/core/services/logger_service.dart';
import 'package:smart_campus/data/models/attendance_record_model.dart';
import 'package:smart_campus/data/repositories/base_repository.dart';

@lazySingleton
class AttendanceRecordRepository extends BaseRepository<AttendanceRecordModel> {
  static const _tableName = 'attendance_records';

  AttendanceRecordRepository(
    DatabaseService databaseService,
    LoggerService logger,
  ) : super(databaseService, logger, _tableName);

  @override
  AttendanceRecordModel fromDatabase(Map<String, dynamic> data) =>
      AttendanceRecordModel.fromDatabase(data);

  Future<List<AttendanceRecordModel>> getRecordsBySession(
    String sessionId,
  ) async {
    try {
      return await getAll(
        where: 'session_id = ?',
        whereArgs: [sessionId],
        orderBy: 'marked_at DESC',
      );
    } catch (e) {
      logger.error('Error getting records by session', e);
      rethrow;
    }
  }

  Future<List<AttendanceRecordModel>> getRecordsByStudent(
    String studentId, {
    String? sessionId,
    SyncStatus? syncStatus,
  }) async {
    try {
      String where = 'student_id = ?';
      final whereArgs = [studentId];

      if (sessionId != null) {
        where += ' AND session_id = ?';
        whereArgs.add(sessionId);
      }

      if (syncStatus != null) {
        where += ' AND sync_status = ?';
        whereArgs.add(syncStatus.toString().split('.').last);
      }

      return await getAll(
        where: where,
        whereArgs: whereArgs,
        orderBy: 'marked_at DESC',
      );
    } catch (e) {
      logger.error('Error getting records by student', e);
      rethrow;
    }
  }

  Future<bool> hasMarkedAttendance(
    String sessionId,
    String studentId,
  ) async {
    try {
      final records = await getAll(
        where: 'session_id = ? AND student_id = ?',
        whereArgs: [sessionId, studentId],
      );
      return records.isNotEmpty;
    } catch (e) {
      logger.error('Error checking attendance record', e);
      rethrow;
    }
  }

  Future<List<AttendanceRecordModel>> getPendingSyncRecords() async {
    try {
      return await getAll(
        where: 'sync_status = ?',
        whereArgs: ['pending'],
        orderBy: 'marked_at ASC',
      );
    } catch (e) {
      logger.error('Error getting pending sync records', e);
      rethrow;
    }
  }

  Future<void> updateSyncStatus(
    String recordId,
    SyncStatus status,
  ) async {
    try {
      await update(
        AttendanceRecordModel.fromDatabase({
          'id': recordId,
          'sync_status': status.toString().split('.').last,
          'updated_at': DateTime.now().millisecondsSinceEpoch,
        }),
      );
    } catch (e) {
      logger.error('Error updating sync status', e);
      rethrow;
    }
  }

  Future<void> markAsSynced(String recordId) async {
    try {
      await updateSyncStatus(recordId, SyncStatus.synced);
    } catch (e) {
      logger.error('Error marking record as synced', e);
      rethrow;
    }
  }

  Future<void> markAsFailed(String recordId) async {
    try {
      await updateSyncStatus(recordId, SyncStatus.failed);
    } catch (e) {
      logger.error('Error marking record as failed', e);
      rethrow;
    }
  }
}
