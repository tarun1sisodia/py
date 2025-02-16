import 'package:json_annotation/json_annotation.dart';
import 'base_model.dart';

part 'attendance_record_model.g.dart';

enum AttendanceStatus {
  @JsonValue('present')
  present,
  @JsonValue('absent')
  absent,
  @JsonValue('late')
  late,
}

enum VerificationStatus {
  @JsonValue('pending')
  pending,
  @JsonValue('verified')
  verified,
  @JsonValue('rejected')
  rejected,
}

enum SyncStatus {
  @JsonValue('pending')
  pending,
  @JsonValue('synced')
  synced,
  @JsonValue('failed')
  failed,
}

@JsonSerializable()
class AttendanceRecordModel extends BaseModel {
  @JsonKey(name: 'session_id')
  final String sessionId;
  @JsonKey(name: 'student_id')
  final String studentId;
  @JsonKey(name: 'marked_at')
  final String markedAt;
  @JsonKey(name: 'location_latitude')
  final double? locationLatitude;
  @JsonKey(name: 'location_longitude')
  final double? locationLongitude;
  @JsonKey(name: 'wifi_ssid')
  final String? wifiSSID;
  @JsonKey(name: 'wifi_bssid')
  final String? wifiBSSID;
  @JsonKey(name: 'device_id')
  final String deviceId;
  final AttendanceStatus status;
  @JsonKey(name: 'verification_status')
  final VerificationStatus verificationStatus;
  @JsonKey(name: 'sync_status')
  final SyncStatus syncStatus;
  @JsonKey(name: 'verification_log')
  final String? verificationLog;
  @JsonKey(name: 'student_name')
  final String? studentName;
  @JsonKey(name: 'course_name')
  final String? courseName;

  const AttendanceRecordModel({
    required super.id,
    required this.sessionId,
    required this.studentId,
    required this.markedAt,
    this.locationLatitude,
    this.locationLongitude,
    this.wifiSSID,
    this.wifiBSSID,
    required this.deviceId,
    required this.status,
    required this.verificationStatus,
    required this.syncStatus,
    this.verificationLog,
    this.studentName,
    this.courseName,
    required super.createdAt,
    required super.updatedAt,
  });

  factory AttendanceRecordModel.fromJson(Map<String, dynamic> json) =>
      _$AttendanceRecordModelFromJson(json);

  @override
  Map<String, dynamic> toJson() => _$AttendanceRecordModelToJson(this);

  @override
  Map<String, dynamic> toDatabase() => {
        'id': id,
        'session_id': sessionId,
        'student_id': studentId,
        'marked_at': markedAt,
        'location_latitude': locationLatitude,
        'location_longitude': locationLongitude,
        'wifi_ssid': wifiSSID,
        'wifi_bssid': wifiBSSID,
        'device_id': deviceId,
        'status': status.toString().split('.').last,
        'verification_status': verificationStatus.toString().split('.').last,
        'sync_status': syncStatus.toString().split('.').last,
        'verification_log': verificationLog,
        'student_name': studentName,
        'course_name': courseName,
        'created_at': createdAt,
        'updated_at': updatedAt,
      };

  factory AttendanceRecordModel.fromDatabase(Map<String, dynamic> data) =>
      AttendanceRecordModel(
        id: data['id'] as String,
        sessionId: data['session_id'] as String,
        studentId: data['student_id'] as String,
        markedAt: data['marked_at'] as String,
        locationLatitude: data['location_latitude'] as double?,
        locationLongitude: data['location_longitude'] as double?,
        wifiSSID: data['wifi_ssid'] as String?,
        wifiBSSID: data['wifi_bssid'] as String?,
        deviceId: data['device_id'] as String,
        status: AttendanceStatus.values.firstWhere(
          (status) => status.toString().split('.').last == data['status'],
          orElse: () => AttendanceStatus.absent,
        ),
        verificationStatus: VerificationStatus.values.firstWhere(
          (status) =>
              status.toString().split('.').last == data['verification_status'],
          orElse: () => VerificationStatus.pending,
        ),
        syncStatus: SyncStatus.values.firstWhere(
          (status) => status.toString().split('.').last == data['sync_status'],
          orElse: () => SyncStatus.pending,
        ),
        verificationLog: data['verification_log'] as String?,
        studentName: data['student_name'] as String?,
        courseName: data['course_name'] as String?,
        createdAt: data['created_at'] as String,
        updatedAt: data['updated_at'] as String,
      );

  AttendanceRecordModel copyWith({
    String? id,
    String? sessionId,
    String? studentId,
    String? markedAt,
    double? locationLatitude,
    double? locationLongitude,
    String? wifiSSID,
    String? wifiBSSID,
    String? deviceId,
    AttendanceStatus? status,
    VerificationStatus? verificationStatus,
    SyncStatus? syncStatus,
    String? verificationLog,
    String? studentName,
    String? courseName,
    String? createdAt,
    String? updatedAt,
  }) {
    return AttendanceRecordModel(
      id: id ?? this.id,
      sessionId: sessionId ?? this.sessionId,
      studentId: studentId ?? this.studentId,
      markedAt: markedAt ?? this.markedAt,
      locationLatitude: locationLatitude ?? this.locationLatitude,
      locationLongitude: locationLongitude ?? this.locationLongitude,
      wifiSSID: wifiSSID ?? this.wifiSSID,
      wifiBSSID: wifiBSSID ?? this.wifiBSSID,
      deviceId: deviceId ?? this.deviceId,
      status: status ?? this.status,
      verificationStatus: verificationStatus ?? this.verificationStatus,
      syncStatus: syncStatus ?? this.syncStatus,
      verificationLog: verificationLog ?? this.verificationLog,
      studentName: studentName ?? this.studentName,
      courseName: courseName ?? this.courseName,
      createdAt: createdAt ?? this.createdAt,
      updatedAt: updatedAt ?? this.updatedAt,
    );
  }

  bool get isPresent => status == AttendanceStatus.present;
  bool get isAbsent => status == AttendanceStatus.absent;
  bool get isLate => status == AttendanceStatus.late;

  bool get isPending => verificationStatus == VerificationStatus.pending;
  bool get isVerified => verificationStatus == VerificationStatus.verified;
  bool get isRejected => verificationStatus == VerificationStatus.rejected;

  bool get isSynced => syncStatus == SyncStatus.synced;
  bool get isSyncPending => syncStatus == SyncStatus.pending;
  bool get isSyncFailed => syncStatus == SyncStatus.failed;
}
