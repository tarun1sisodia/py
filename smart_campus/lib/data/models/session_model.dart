import 'package:smart_campus/domain/entities/session.dart';

class SessionModel extends Session {
  const SessionModel({
    required super.id,
    required super.teacherId,
    required super.courseId,
    required super.sessionDate,
    required super.startTime,
    required super.endTime,
    super.wifiSSID,
    super.wifiBSSID,
    super.locationLatitude,
    super.locationLongitude,
    super.locationRadius,
    required super.status,
    required super.createdAt,
    required super.updatedAt,
  });

  factory SessionModel.fromJson(Map<String, dynamic> json) {
    return SessionModel(
      id: json['id'] as String,
      teacherId: json['teacher_id'] as String,
      courseId: json['course_id'] as String,
      sessionDate: DateTime.parse(json['session_date'] as String),
      startTime: DateTime.parse(json['start_time'] as String),
      endTime: DateTime.parse(json['end_time'] as String),
      wifiSSID: json['wifi_ssid'] as String?,
      wifiBSSID: json['wifi_bssid'] as String?,
      locationLatitude: json['location_latitude'] as double?,
      locationLongitude: json['location_longitude'] as double?,
      locationRadius: json['location_radius'] as int?,
      status: SessionStatus.values.firstWhere(
        (e) => e.toString().split('.').last == json['status'],
      ),
      createdAt: DateTime.parse(json['created_at'] as String),
      updatedAt: DateTime.parse(json['updated_at'] as String),
    );
  }

  @override
  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'teacher_id': teacherId,
      'course_id': courseId,
      'session_date': sessionDate.toIso8601String(),
      'start_time': startTime.toIso8601String(),
      'end_time': endTime.toIso8601String(),
      'wifi_ssid': wifiSSID,
      'wifi_bssid': wifiBSSID,
      'location_latitude': locationLatitude,
      'location_longitude': locationLongitude,
      'location_radius': locationRadius,
      'status': status.toString().split('.').last,
      'created_at': createdAt.toIso8601String(),
      'updated_at': updatedAt.toIso8601String(),
    };
  }

  factory SessionModel.fromEntity(Session session) {
    return SessionModel(
      id: session.id,
      teacherId: session.teacherId,
      courseId: session.courseId,
      sessionDate: session.sessionDate,
      startTime: session.startTime,
      endTime: session.endTime,
      wifiSSID: session.wifiSSID,
      wifiBSSID: session.wifiBSSID,
      locationLatitude: session.locationLatitude,
      locationLongitude: session.locationLongitude,
      locationRadius: session.locationRadius,
      status: session.status,
      createdAt: session.createdAt,
      updatedAt: session.updatedAt,
    );
  }

  @override
  SessionModel copyWith({
    String? id,
    String? teacherId,
    String? courseId,
    DateTime? sessionDate,
    DateTime? startTime,
    DateTime? endTime,
    String? wifiSSID,
    String? wifiBSSID,
    double? locationLatitude,
    double? locationLongitude,
    int? locationRadius,
    SessionStatus? status,
    DateTime? createdAt,
    DateTime? updatedAt,
  }) {
    return SessionModel(
      id: id ?? this.id,
      teacherId: teacherId ?? this.teacherId,
      courseId: courseId ?? this.courseId,
      sessionDate: sessionDate ?? this.sessionDate,
      startTime: startTime ?? this.startTime,
      endTime: endTime ?? this.endTime,
      wifiSSID: wifiSSID ?? this.wifiSSID,
      wifiBSSID: wifiBSSID ?? this.wifiBSSID,
      locationLatitude: locationLatitude ?? this.locationLatitude,
      locationLongitude: locationLongitude ?? this.locationLongitude,
      locationRadius: locationRadius ?? this.locationRadius,
      status: status ?? this.status,
      createdAt: createdAt ?? this.createdAt,
      updatedAt: updatedAt ?? this.updatedAt,
    );
  }
}
