import 'package:json_annotation/json_annotation.dart';
import 'base_model.dart';

part 'attendance_session_model.g.dart';

enum AttendanceSessionStatus {
  @JsonValue('scheduled')
  scheduled,
  @JsonValue('active')
  active,
  @JsonValue('completed')
  completed,
  @JsonValue('cancelled')
  cancelled,
}

@JsonSerializable()
class AttendanceSessionModel extends BaseModel {
  @JsonKey(name: 'course_id')
  final String courseId;
  @JsonKey(name: 'teacher_id')
  final String teacherId;
  @JsonKey(name: 'session_date')
  final String sessionDate;
  @JsonKey(name: 'start_time')
  final String startTime;
  @JsonKey(name: 'end_time')
  final String endTime;
  @JsonKey(name: 'wifi_ssid')
  final String? wifiSSID;
  @JsonKey(name: 'wifi_bssid')
  final String? wifiBSSID;
  @JsonKey(name: 'location_latitude')
  final double? locationLatitude;
  @JsonKey(name: 'location_longitude')
  final double? locationLongitude;
  @JsonKey(name: 'location_radius')
  final int? locationRadius;
  @JsonKey(name: 'location_name')
  final String? locationName;
  final AttendanceSessionStatus status;
  @JsonKey(name: 'total_students')
  final int? totalStudents;
  @JsonKey(name: 'present_count')
  final int? presentCount;
  @JsonKey(name: 'late_count')
  final int? lateCount;
  @JsonKey(name: 'absent_count')
  final int? absentCount;
  @JsonKey(name: 'course_name')
  final String? courseName;
  @JsonKey(name: 'teacher_name')
  final String? teacherName;

  const AttendanceSessionModel({
    required super.id,
    required this.courseId,
    required this.teacherId,
    required this.sessionDate,
    required this.startTime,
    required this.endTime,
    this.wifiSSID,
    this.wifiBSSID,
    this.locationLatitude,
    this.locationLongitude,
    this.locationRadius,
    this.locationName,
    required this.status,
    this.totalStudents,
    this.presentCount,
    this.lateCount,
    this.absentCount,
    this.courseName,
    this.teacherName,
    required super.createdAt,
    required super.updatedAt,
  });

  factory AttendanceSessionModel.fromJson(Map<String, dynamic> json) =>
      _$AttendanceSessionModelFromJson(json);

  @override
  Map<String, dynamic> toJson() => _$AttendanceSessionModelToJson(this);

  @override
  Map<String, dynamic> toDatabase() => {
        'id': id,
        'course_id': courseId,
        'teacher_id': teacherId,
        'session_date': sessionDate,
        'start_time': startTime,
        'end_time': endTime,
        'wifi_ssid': wifiSSID,
        'wifi_bssid': wifiBSSID,
        'location_latitude': locationLatitude,
        'location_longitude': locationLongitude,
        'location_radius': locationRadius,
        'location_name': locationName,
        'status': status.toString().split('.').last,
        'total_students': totalStudents,
        'present_count': presentCount,
        'late_count': lateCount,
        'absent_count': absentCount,
        'course_name': courseName,
        'teacher_name': teacherName,
        'created_at': createdAt,
        'updated_at': updatedAt,
      };

  factory AttendanceSessionModel.fromDatabase(Map<String, dynamic> data) =>
      AttendanceSessionModel(
        id: data['id'] as String,
        courseId: data['course_id'] as String,
        teacherId: data['teacher_id'] as String,
        sessionDate: data['session_date'] as String,
        startTime: data['start_time'] as String,
        endTime: data['end_time'] as String,
        wifiSSID: data['wifi_ssid'] as String?,
        wifiBSSID: data['wifi_bssid'] as String?,
        locationLatitude: data['location_latitude'] as double?,
        locationLongitude: data['location_longitude'] as double?,
        locationRadius: data['location_radius'] as int?,
        locationName: data['location_name'] as String?,
        status: AttendanceSessionStatus.values.firstWhere(
          (status) => status.toString().split('.').last == data['status'],
          orElse: () => AttendanceSessionStatus.scheduled,
        ),
        totalStudents: data['total_students'] as int?,
        presentCount: data['present_count'] as int?,
        lateCount: data['late_count'] as int?,
        absentCount: data['absent_count'] as int?,
        courseName: data['course_name'] as String?,
        teacherName: data['teacher_name'] as String?,
        createdAt: data['created_at'] as String,
        updatedAt: data['updated_at'] as String,
      );

  AttendanceSessionModel copyWith({
    String? id,
    String? courseId,
    String? teacherId,
    String? sessionDate,
    String? startTime,
    String? endTime,
    String? wifiSSID,
    String? wifiBSSID,
    double? locationLatitude,
    double? locationLongitude,
    int? locationRadius,
    String? locationName,
    AttendanceSessionStatus? status,
    int? totalStudents,
    int? presentCount,
    int? lateCount,
    int? absentCount,
    String? courseName,
    String? teacherName,
    String? createdAt,
    String? updatedAt,
  }) {
    return AttendanceSessionModel(
      id: id ?? this.id,
      courseId: courseId ?? this.courseId,
      teacherId: teacherId ?? this.teacherId,
      sessionDate: sessionDate ?? this.sessionDate,
      startTime: startTime ?? this.startTime,
      endTime: endTime ?? this.endTime,
      wifiSSID: wifiSSID ?? this.wifiSSID,
      wifiBSSID: wifiBSSID ?? this.wifiBSSID,
      locationLatitude: locationLatitude ?? this.locationLatitude,
      locationLongitude: locationLongitude ?? this.locationLongitude,
      locationRadius: locationRadius ?? this.locationRadius,
      locationName: locationName ?? this.locationName,
      status: status ?? this.status,
      totalStudents: totalStudents ?? this.totalStudents,
      presentCount: presentCount ?? this.presentCount,
      lateCount: lateCount ?? this.lateCount,
      absentCount: absentCount ?? this.absentCount,
      courseName: courseName ?? this.courseName,
      teacherName: teacherName ?? this.teacherName,
      createdAt: createdAt ?? this.createdAt,
      updatedAt: updatedAt ?? this.updatedAt,
    );
  }

  bool get isScheduled => status == AttendanceSessionStatus.scheduled;
  bool get isActive => status == AttendanceSessionStatus.active;
  bool get isCompleted => status == AttendanceSessionStatus.completed;
  bool get isCancelled => status == AttendanceSessionStatus.cancelled;

  double get attendanceRate {
    if (totalStudents == null || totalStudents == 0) return 0.0;
    final present = (presentCount ?? 0) + (lateCount ?? 0);
    return (present / totalStudents!) * 100;
  }
}
