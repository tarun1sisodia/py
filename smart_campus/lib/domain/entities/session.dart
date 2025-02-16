import 'package:equatable/equatable.dart';

enum SessionStatus { active, completed, cancelled }

class Session extends Equatable {
  final String id;
  final String teacherId;
  final String courseId;
  final DateTime sessionDate;
  final DateTime startTime;
  final DateTime endTime;
  final String? wifiSSID;
  final String? wifiBSSID;
  final double? locationLatitude;
  final double? locationLongitude;
  final int? locationRadius;
  final SessionStatus status;
  final DateTime createdAt;
  final DateTime updatedAt;

  bool get isActive => status == SessionStatus.active;
  bool get isCompleted => status == SessionStatus.completed;
  bool get isPending => status == SessionStatus.active;
  bool get isRejected => status == SessionStatus.cancelled;
  bool get isEnded => status == SessionStatus.completed;

  const Session({
    required this.id,
    required this.teacherId,
    required this.courseId,
    required this.sessionDate,
    required this.startTime,
    required this.endTime,
    this.wifiSSID,
    this.wifiBSSID,
    this.locationLatitude,
    this.locationLongitude,
    this.locationRadius,
    required this.status,
    required this.createdAt,
    required this.updatedAt,
  });

  @override
  List<Object?> get props => [
        id,
        teacherId,
        courseId,
        sessionDate,
        startTime,
        endTime,
        wifiSSID,
        wifiBSSID,
        locationLatitude,
        locationLongitude,
        locationRadius,
        status,
        createdAt,
        updatedAt,
      ];

  Session copyWith({
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
    return Session(
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

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'teacherId': teacherId,
      'courseId': courseId,
      'sessionDate': sessionDate.toIso8601String(),
      'startTime': startTime.toIso8601String(),
      'endTime': endTime.toIso8601String(),
      'wifiSSID': wifiSSID,
      'wifiBSSID': wifiBSSID,
      'locationLatitude': locationLatitude,
      'locationLongitude': locationLongitude,
      'locationRadius': locationRadius,
      'status': status.toString().split('.').last,
      'createdAt': createdAt.toIso8601String(),
      'updatedAt': updatedAt.toIso8601String(),
      'isActive': isActive,
      'isCompleted': isCompleted,
    };
  }

  factory Session.fromJson(Map<String, dynamic> json) {
    return Session(
      id: json['id'] as String,
      teacherId: json['teacherId'] as String,
      courseId: json['courseId'] as String,
      sessionDate: DateTime.parse(json['sessionDate'] as String),
      startTime: DateTime.parse(json['startTime'] as String),
      endTime: DateTime.parse(json['endTime'] as String),
      wifiSSID: json['wifiSSID'] as String?,
      wifiBSSID: json['wifiBSSID'] as String?,
      locationLatitude: json['locationLatitude'] as double?,
      locationLongitude: json['locationLongitude'] as double?,
      locationRadius: json['locationRadius'] as int?,
      status: SessionStatus.values.firstWhere(
        (status) => status.toString().split('.').last == json['status'],
      ),
      createdAt: DateTime.parse(json['createdAt'] as String),
      updatedAt: DateTime.parse(json['updatedAt'] as String),
    );
  }
}
