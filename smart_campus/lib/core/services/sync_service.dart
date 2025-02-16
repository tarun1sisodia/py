import 'dart:async';
import 'package:injectable/injectable.dart';
import 'package:smart_campus/domain/entities/attendance_record.dart';
import 'package:smart_campus/data/datasources/local/attendance_local_datasource.dart';
import 'package:smart_campus/domain/repositories/session_repository.dart';
import 'package:smart_campus/core/services/connectivity_service.dart';
import 'package:smart_campus/core/services/logger_service.dart';

@injectable
class SyncService {
  final AttendanceLocalDatasource _localDatasource;
  final SessionRepository _sessionRepository;
  final ConnectivityService _connectivityService;
  final LoggerService _logger;
  Timer? _syncTimer;

  SyncService({
    required AttendanceLocalDatasource localDatasource,
    required SessionRepository sessionRepository,
    required ConnectivityService connectivityService,
    required LoggerService logger,
  })  : _localDatasource = localDatasource,
        _sessionRepository = sessionRepository,
        _connectivityService = connectivityService,
        _logger = logger;

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
      final record = await _localDatasource.storeOfflineAttendance(
        sessionId: sessionId,
        studentId: studentId,
        studentName: studentName,
        locationLatitude: locationLatitude,
        locationLongitude: locationLongitude,
        wifiSSID: wifiSSID,
        wifiBSSID: wifiBSSID,
        deviceId: deviceId,
      );
      _logger.info('Stored offline attendance for session $sessionId');
      return record;
    } catch (e) {
      _logger.error('Failed to store offline attendance', e);
      rethrow;
    }
  }

  Future<List<AttendanceRecord>> getPendingRecords() async {
    try {
      return await _localDatasource.getPendingRecords();
    } catch (e) {
      _logger.error('Failed to get pending records', e);
      return [];
    }
  }

  Future<List<AttendanceRecord>> getFailedRecords() async {
    try {
      return await _localDatasource.getFailedRecords();
    } catch (e) {
      _logger.error('Failed to get failed records', e);
      return [];
    }
  }

  Future<void> retryFailedRecord(String sessionId, String studentId) async {
    try {
      final record = await _localDatasource.getRecord(sessionId, studentId);
      if (record == null) {
        _logger.warning(
            'No record found for session $sessionId and student $studentId');
        return;
      }

      await _sessionRepository.markAttendance(
        sessionId: sessionId,
        studentId: studentId,
        locationLatitude: record.locationLatitude,
        locationLongitude: record.locationLongitude,
        wifiSSID: record.wifiSSID,
        wifiBSSID: record.wifiBSSID,
        deviceId: record.deviceId,
      );

      await _localDatasource.deleteRecord(sessionId, studentId);
      _logger.info('Successfully retried failed record for session $sessionId');
    } catch (e) {
      _logger.error('Failed to retry record', e);
      await _localDatasource.markRecordAsFailed(
        sessionId,
        studentId,
        e.toString(),
      );
      rethrow;
    }
  }

  void startBackgroundSync() {
    _syncTimer?.cancel();
    _syncTimer = Timer.periodic(const Duration(minutes: 15), (_) => _sync());
    _logger.info('Background sync started');
  }

  void stopBackgroundSync() {
    _syncTimer?.cancel();
    _syncTimer = null;
    _logger.info('Background sync stopped');
  }

  Future<void> _sync() async {
    if (!await _connectivityService.isOnline()) {
      _logger.info('Skipping sync - device is offline');
      return;
    }

    try {
      final pendingRecords = await getPendingRecords();
      _logger.info('Found ${pendingRecords.length} pending records to sync');

      for (final record in pendingRecords) {
        try {
          await _sessionRepository.markAttendance(
            sessionId: record.sessionId,
            studentId: record.studentId,
            locationLatitude: record.locationLatitude,
            locationLongitude: record.locationLongitude,
            wifiSSID: record.wifiSSID,
            wifiBSSID: record.wifiBSSID,
            deviceId: record.deviceId,
          );

          await _localDatasource.deleteRecord(
              record.sessionId, record.studentId);
          _logger.info(
              'Successfully synced record for session ${record.sessionId}');
        } catch (e) {
          _logger.error('Failed to sync record', e);
          await _localDatasource.markRecordAsFailed(
            record.sessionId,
            record.studentId,
            e.toString(),
          );
        }
      }
    } catch (e) {
      _logger.error('Sync operation failed', e);
    }
  }

  @disposeMethod
  void dispose() {
    stopBackgroundSync();
  }
}
