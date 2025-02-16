import 'package:equatable/equatable.dart';

enum VerificationStatus { pending, verified, rejected }

class AttendanceRecord extends Equatable {
  final String id;
  final String sessionId;
  final String studentId;
  final String studentName;
  final DateTime markedAt;
  final double? locationLatitude;
  final double? locationLongitude;
  final String? wifiSSID;
  final String? wifiBSSID;
  final String deviceId;
  final String verificationStatus;
  final String? rejectionReason;
  final DateTime createdAt;
  final DateTime updatedAt;

  VerificationStatus get status => VerificationStatus.values.firstWhere(
        (e) => e.toString().split('.').last == verificationStatus,
        orElse: () => VerificationStatus.pending,
      );

  const AttendanceRecord({
    required this.id,
    required this.sessionId,
    required this.studentId,
    required this.studentName,
    required this.markedAt,
    this.locationLatitude,
    this.locationLongitude,
    this.wifiSSID,
    this.wifiBSSID,
    required this.deviceId,
    required this.verificationStatus,
    this.rejectionReason,
    required this.createdAt,
    required this.updatedAt,
  });

  @override
  List<Object?> get props => [
        id,
        sessionId,
        studentId,
        studentName,
        markedAt,
        locationLatitude,
        locationLongitude,
        wifiSSID,
        wifiBSSID,
        deviceId,
        verificationStatus,
        rejectionReason,
        createdAt,
        updatedAt,
      ];

  AttendanceRecord copyWith({
    String? id,
    String? sessionId,
    String? studentId,
    String? studentName,
    DateTime? markedAt,
    double? locationLatitude,
    double? locationLongitude,
    String? wifiSSID,
    String? wifiBSSID,
    String? deviceId,
    String? verificationStatus,
    String? rejectionReason,
    DateTime? createdAt,
    DateTime? updatedAt,
  }) {
    return AttendanceRecord(
      id: id ?? this.id,
      sessionId: sessionId ?? this.sessionId,
      studentId: studentId ?? this.studentId,
      studentName: studentName ?? this.studentName,
      markedAt: markedAt ?? this.markedAt,
      locationLatitude: locationLatitude ?? this.locationLatitude,
      locationLongitude: locationLongitude ?? this.locationLongitude,
      wifiSSID: wifiSSID ?? this.wifiSSID,
      wifiBSSID: wifiBSSID ?? this.wifiBSSID,
      deviceId: deviceId ?? this.deviceId,
      verificationStatus: verificationStatus ?? this.verificationStatus,
      rejectionReason: rejectionReason ?? this.rejectionReason,
      createdAt: createdAt ?? this.createdAt,
      updatedAt: updatedAt ?? this.updatedAt,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'sessionId': sessionId,
      'studentId': studentId,
      'studentName': studentName,
      'markedAt': markedAt.toIso8601String(),
      'locationLatitude': locationLatitude,
      'locationLongitude': locationLongitude,
      'wifiSSID': wifiSSID,
      'wifiBSSID': wifiBSSID,
      'deviceId': deviceId,
      'verificationStatus': verificationStatus,
      'rejectionReason': rejectionReason,
      'createdAt': createdAt.toIso8601String(),
      'updatedAt': updatedAt.toIso8601String(),
    };
  }

  factory AttendanceRecord.fromJson(Map<String, dynamic> json) {
    return AttendanceRecord(
      id: json['id'] as String,
      sessionId: json['sessionId'] as String,
      studentId: json['studentId'] as String,
      studentName: json['studentName'] as String,
      markedAt: DateTime.parse(json['markedAt'] as String),
      locationLatitude: json['locationLatitude'] as double?,
      locationLongitude: json['locationLongitude'] as double?,
      wifiSSID: json['wifiSSID'] as String?,
      wifiBSSID: json['wifiBSSID'] as String?,
      deviceId: json['deviceId'] as String,
      verificationStatus: json['verificationStatus'] as String,
      rejectionReason: json['rejectionReason'] as String?,
      createdAt: DateTime.parse(json['createdAt'] as String),
      updatedAt: DateTime.parse(json['updatedAt'] as String),
    );
  }
}
