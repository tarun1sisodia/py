// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'attendance_session_model.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

AttendanceSessionModel _$AttendanceSessionModelFromJson(
        Map<String, dynamic> json) =>
    $checkedCreate(
      'AttendanceSessionModel',
      json,
      ($checkedConvert) {
        $checkKeys(
          json,
          allowedKeys: const [
            'id',
            'course_id',
            'teacher_id',
            'date',
            'start_time',
            'end_time',
            'latitude',
            'longitude',
            'radius',
            'wifi_ssid',
            'wifi_bssid',
            'status',
            'created_at',
            'updated_at'
          ],
        );
        final val = AttendanceSessionModel(
          id: $checkedConvert('id', (v) => v as String),
          courseId: $checkedConvert('course_id', (v) => v as String),
          teacherId: $checkedConvert('teacher_id', (v) => v as String),
          date: $checkedConvert('date', (v) => v as String),
          startTime: $checkedConvert('start_time', (v) => (v as num).toInt()),
          endTime: $checkedConvert('end_time', (v) => (v as num).toInt()),
          latitude: $checkedConvert('latitude', (v) => (v as num).toDouble()),
          longitude: $checkedConvert('longitude', (v) => (v as num).toDouble()),
          radius: $checkedConvert('radius', (v) => (v as num).toInt()),
          wifiSSID: $checkedConvert('wifi_ssid', (v) => v as String),
          wifiBSSID: $checkedConvert('wifi_bssid', (v) => v as String),
          status: $checkedConvert('status',
              (v) => $enumDecode(_$AttendanceSessionStatusEnumMap, v)),
          createdAt: $checkedConvert('created_at', (v) => (v as num).toInt()),
          updatedAt: $checkedConvert('updated_at', (v) => (v as num).toInt()),
        );
        return val;
      },
      fieldKeyMap: const {
        'courseId': 'course_id',
        'teacherId': 'teacher_id',
        'startTime': 'start_time',
        'endTime': 'end_time',
        'wifiSSID': 'wifi_ssid',
        'wifiBSSID': 'wifi_bssid',
        'createdAt': 'created_at',
        'updatedAt': 'updated_at'
      },
    );

Map<String, dynamic> _$AttendanceSessionModelToJson(
        AttendanceSessionModel instance) =>
    <String, dynamic>{
      'id': instance.id,
      'course_id': instance.courseId,
      'teacher_id': instance.teacherId,
      'date': instance.date,
      'start_time': instance.startTime,
      'end_time': instance.endTime,
      'latitude': instance.latitude,
      'longitude': instance.longitude,
      'radius': instance.radius,
      'wifi_ssid': instance.wifiSSID,
      'wifi_bssid': instance.wifiBSSID,
      'status': _$AttendanceSessionStatusEnumMap[instance.status]!,
      'created_at': instance.createdAt,
      'updated_at': instance.updatedAt,
    };

const _$AttendanceSessionStatusEnumMap = {
  AttendanceSessionStatus.active: 'active',
  AttendanceSessionStatus.completed: 'completed',
  AttendanceSessionStatus.cancelled: 'cancelled',
};
