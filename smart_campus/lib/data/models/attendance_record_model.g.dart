// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'attendance_record_model.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

AttendanceRecordModel _$AttendanceRecordModelFromJson(
        Map<String, dynamic> json) =>
    $checkedCreate(
      'AttendanceRecordModel',
      json,
      ($checkedConvert) {
        $checkKeys(
          json,
          allowedKeys: const [
            'id',
            'session_id',
            'student_id',
            'marked_at',
            'latitude',
            'longitude',
            'wifi_ssid',
            'wifi_bssid',
            'device_id',
            'sync_status',
            'created_at',
            'updated_at'
          ],
        );
        final val = AttendanceRecordModel(
          id: $checkedConvert('id', (v) => v as String),
          sessionId: $checkedConvert('session_id', (v) => v as String),
          studentId: $checkedConvert('student_id', (v) => v as String),
          markedAt: $checkedConvert('marked_at', (v) => v as String),
          locationLatitude:
              $checkedConvert('location_latitude', (v) => (v as num).toDouble()),
          locationLongitude:
              $checkedConvert('location_longitude', (v) => (v as num).toDouble()),
          wifiSSID: $checkedConvert('wifi_ssid', (v) => v as String),
          wifiBSSID: $checkedConvert('wifi_bssid', (v) => v as String),
          deviceId: $checkedConvert('device_id', (v) => v as String),
          syncStatus: $checkedConvert(
              'sync_status', (v) => $enumDecode(_$SyncStatusEnumMap, v)),
          createdAt: $checkedConvert('created_at', (v) => v as String),
          updatedAt: $checkedConvert('updated_at', (v) => v as String),
          status: AttendanceStatus.present,
          verificationStatus: VerificationStatus.pending,
          verificationLog: null,
          studentName: null,
          courseName: null,
        );
        return val;
      },
      fieldKeyMap: const {
        'sessionId': 'session_id',
        'studentId': 'student_id',
        'markedAt': 'marked_at',
        'wifiSSID': 'wifi_ssid',
        'wifiBSSID': 'wifi_bssid',
        'deviceId': 'device_id',
        'syncStatus': 'sync_status',
        'createdAt': 'created_at',
        'updatedAt': 'updated_at'
      },
    );

Map<String, dynamic> _$AttendanceRecordModelToJson(
        AttendanceRecordModel instance) =>
    <String, dynamic>{
      'id': instance.id,
      'session_id': instance.sessionId,
      'student_id': instance.studentId,
      'marked_at': instance.markedAt,
      'location_latitude': instance.locationLatitude,
      'location_longitude': instance.locationLongitude,
      'wifi_ssid': instance.wifiSSID,
      'wifi_bssid': instance.wifiBSSID,
      'device_id': instance.deviceId,
      'sync_status': _$SyncStatusEnumMap[instance.syncStatus]!,
      'created_at': instance.createdAt,
      'updated_at': instance.updatedAt,
    };

const _$SyncStatusEnumMap = {
  SyncStatus.pending: 'pending',
  SyncStatus.synced: 'synced',
  SyncStatus.failed: 'failed',
};
