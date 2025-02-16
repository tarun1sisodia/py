import 'package:injectable/injectable.dart';
import 'package:sqflite/sqflite.dart';
import 'package:uuid/uuid.dart';
import 'package:smart_campus/domain/entities/attendance_record.dart';
import 'package:smart_campus/data/datasources/local/database_helper.dart';

abstract class AttendanceLocalDatasource {
  Future<AttendanceRecord> storeOfflineAttendance({
    required String sessionId,
    required String studentId,
    required String studentName,
    double? locationLatitude,
    double? locationLongitude,
    String? wifiSSID,
    String? wifiBSSID,
    required String deviceId,
  });

  Future<List<AttendanceRecord>> getPendingRecords();
  Future<List<AttendanceRecord>> getFailedRecords();
  Future<AttendanceRecord?> getRecord(String sessionId, String studentId);
  Future<void> deleteRecord(String sessionId, String studentId);
  Future<void> markRecordAsFailed(
      String sessionId, String studentId, String error);
}

@LazySingleton(as: AttendanceLocalDatasource)
class AttendanceLocalDatasourceImpl implements AttendanceLocalDatasource {
  final DatabaseHelper _databaseHelper;
  static const String tableName = 'offline_attendance_records';

  AttendanceLocalDatasourceImpl(this._databaseHelper);

  static Future<void> createTable(Database db) async {
    await db.execute('''
      CREATE TABLE IF NOT EXISTS $tableName (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        sessionId TEXT NOT NULL,
        studentId TEXT NOT NULL,
        studentName TEXT NOT NULL,
        locationLatitude REAL,
        locationLongitude REAL,
        wifiSSID TEXT,
        wifiBSSID TEXT,
        deviceId TEXT NOT NULL,
        syncStatus TEXT NOT NULL,
        errorMessage TEXT,
        createdAt TEXT NOT NULL,
        updatedAt TEXT NOT NULL,
        UNIQUE(sessionId, studentId)
      )
    ''');
  }

  @override
  Future<AttendanceRecord> storeOfflineAttendance({
    required String sessionId,
    required String studentId,
    required String studentName,
    double? locationLatitude,
    double? locationLongitude,
    String? wifiSSID,
    String? wifiBSSID,
    required String deviceId,
  }) async {
    try {
      final now = DateTime.now();
      final record = {
        'id': const Uuid().v4(),
        'sessionId': sessionId,
        'studentId': studentId,
        'studentName': studentName,
        'markedAt': now.toIso8601String(),
        'locationLatitude': locationLatitude,
        'locationLongitude': locationLongitude,
        'wifiSSID': wifiSSID,
        'wifiBSSID': wifiBSSID,
        'deviceId': deviceId,
        'verificationStatus': 'pending',
        'createdAt': now.toIso8601String(),
        'updatedAt': now.toIso8601String(),
      };

      final id = await _databaseHelper.insert(
        tableName,
        record,
        conflictAlgorithm: ConflictAlgorithm.replace,
      );

      return AttendanceRecord.fromJson(record);
    } catch (e) {
      rethrow;
    }
  }

  @override
  Future<List<AttendanceRecord>> getPendingRecords() async {
    try {
      final records = await _databaseHelper.query(
        tableName,
        where: 'verificationStatus = ?',
        whereArgs: ['pending'],
      );
      return records.map((r) => AttendanceRecord.fromJson(r)).toList();
    } catch (e) {
      return [];
    }
  }

  @override
  Future<List<AttendanceRecord>> getFailedRecords() async {
    try {
      final records = await _databaseHelper.query(
        tableName,
        where: 'verificationStatus = ?',
        whereArgs: ['failed'],
      );
      return records.map((r) => AttendanceRecord.fromJson(r)).toList();
    } catch (e) {
      return [];
    }
  }

  @override
  Future<AttendanceRecord?> getRecord(
      String sessionId, String studentId) async {
    try {
      final records = await _databaseHelper.query(
        tableName,
        where: 'sessionId = ? AND studentId = ?',
        whereArgs: [sessionId, studentId],
        limit: 1,
      );
      return records.isNotEmpty
          ? AttendanceRecord.fromJson(records.first)
          : null;
    } catch (e) {
      return null;
    }
  }

  @override
  Future<void> deleteRecord(String sessionId, String studentId) async {
    try {
      await _databaseHelper.delete(
        tableName,
        where: 'sessionId = ? AND studentId = ?',
        whereArgs: [sessionId, studentId],
      );
    } catch (e) {
      rethrow;
    }
  }

  @override
  Future<void> markRecordAsFailed(
    String sessionId,
    String studentId,
    String error,
  ) async {
    try {
      await _databaseHelper.update(
        tableName,
        {
          'verificationStatus': 'failed',
          'rejectionReason': error,
          'updatedAt': DateTime.now().toIso8601String(),
        },
        where: 'sessionId = ? AND studentId = ?',
        whereArgs: [sessionId, studentId],
      );
    } catch (e) {
      rethrow;
    }
  }
}
