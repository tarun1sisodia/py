import 'package:equatable/equatable.dart';

enum SyncStatus { pending, completed, failed }

class OfflineAttendanceRecord extends Equatable {
  final String id;
  final String sessionId;
  final String studentId;
  final DateTime markedAt;
  final double? locationLatitude;
  final double? locationLongitude;
  final String? wifiSSID;
  final String? wifiBSSID;
  final String deviceId;
  final SyncStatus syncStatus;
  final String? errorMessage;
  final DateTime createdAt;
  final DateTime updatedAt;

  const OfflineAttendanceRecord({
    required this.id,
    required this.sessionId,
    required this.studentId,
    required this.markedAt,
    this.locationLatitude,
    this.locationLongitude,
    this.wifiSSID,
    this.wifiBSSID,
    required this.deviceId,
    required this.syncStatus,
    this.errorMessage,
    required this.createdAt,
    required this.updatedAt,
  });

  factory OfflineAttendanceRecord.fromJson(Map<String, dynamic> json) {
    return OfflineAttendanceRecord(
      id: json['id'] as String,
      sessionId: json['session_id'] as String,
      studentId: json['student_id'] as String,
      markedAt: DateTime.parse(json['marked_at'] as String),
      locationLatitude: json['location_latitude'] as double?,
      locationLongitude: json['location_longitude'] as double?,
      wifiSSID: json['wifi_ssid'] as String?,
      wifiBSSID: json['wifi_bssid'] as String?,
      deviceId: json['device_id'] as String,
      syncStatus: SyncStatus.values.firstWhere(
        (status) => status.toString().split('.').last == json['sync_status'],
      ),
      errorMessage: json['error_message'] as String?,
      createdAt: DateTime.parse(json['created_at'] as String),
      updatedAt: DateTime.parse(json['updated_at'] as String),
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'session_id': sessionId,
      'student_id': studentId,
      'marked_at': markedAt.toIso8601String(),
      'location_latitude': locationLatitude,
      'location_longitude': locationLongitude,
      'wifi_ssid': wifiSSID,
      'wifi_bssid': wifiBSSID,
      'device_id': deviceId,
      'sync_status': syncStatus.toString().split('.').last,
      'error_message': errorMessage,
      'created_at': createdAt.toIso8601String(),
      'updated_at': updatedAt.toIso8601String(),
    };
  }

  OfflineAttendanceRecord copyWith({
    String? id,
    String? sessionId,
    String? studentId,
    DateTime? markedAt,
    double? locationLatitude,
    double? locationLongitude,
    String? wifiSSID,
    String? wifiBSSID,
    String? deviceId,
    SyncStatus? syncStatus,
    String? errorMessage,
    DateTime? createdAt,
    DateTime? updatedAt,
  }) {
    return OfflineAttendanceRecord(
      id: id ?? this.id,
      sessionId: sessionId ?? this.sessionId,
      studentId: studentId ?? this.studentId,
      markedAt: markedAt ?? this.markedAt,
      locationLatitude: locationLatitude ?? this.locationLatitude,
      locationLongitude: locationLongitude ?? this.locationLongitude,
      wifiSSID: wifiSSID ?? this.wifiSSID,
      wifiBSSID: wifiBSSID ?? this.wifiBSSID,
      deviceId: deviceId ?? this.deviceId,
      syncStatus: syncStatus ?? this.syncStatus,
      errorMessage: errorMessage ?? this.errorMessage,
      createdAt: createdAt ?? this.createdAt,
      updatedAt: updatedAt ?? this.updatedAt,
    );
  }

  @override
  List<Object?> get props => [
        id,
        sessionId,
        studentId,
        markedAt,
        locationLatitude,
        locationLongitude,
        wifiSSID,
        wifiBSSID,
        deviceId,
        syncStatus,
        errorMessage,
        createdAt,
        updatedAt,
      ];
}
