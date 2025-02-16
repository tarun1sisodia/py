import 'package:equatable/equatable.dart';

enum UserRole { teacher, student }

class User extends Equatable {
  final String id;
  final String email;
  final String fullName;
  final UserRole role;
  final String? enrollmentNumber;
  final String? employeeId;
  final String? department;
  final int? yearOfStudy;
  final String? deviceId;
  final DateTime createdAt;
  final DateTime updatedAt;
  final DateTime? lastLogin;
  final bool isActive;

  const User({
    required this.id,
    required this.email,
    required this.fullName,
    required this.role,
    this.enrollmentNumber,
    this.employeeId,
    this.department,
    this.yearOfStudy,
    this.deviceId,
    required this.createdAt,
    required this.updatedAt,
    this.lastLogin,
    this.isActive = true,
  });

  bool get isTeacher => role == UserRole.teacher;
  bool get isStudent => role == UserRole.student;

  User copyWith({
    String? id,
    String? email,
    String? fullName,
    UserRole? role,
    String? enrollmentNumber,
    String? employeeId,
    String? department,
    int? yearOfStudy,
    String? deviceId,
    DateTime? createdAt,
    DateTime? updatedAt,
    DateTime? lastLogin,
    bool? isActive,
  }) {
    return User(
      id: id ?? this.id,
      email: email ?? this.email,
      fullName: fullName ?? this.fullName,
      role: role ?? this.role,
      enrollmentNumber: enrollmentNumber ?? this.enrollmentNumber,
      employeeId: employeeId ?? this.employeeId,
      department: department ?? this.department,
      yearOfStudy: yearOfStudy ?? this.yearOfStudy,
      deviceId: deviceId ?? this.deviceId,
      createdAt: createdAt ?? this.createdAt,
      updatedAt: updatedAt ?? this.updatedAt,
      lastLogin: lastLogin ?? this.lastLogin,
      isActive: isActive ?? this.isActive,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'email': email,
      'fullName': fullName,
      'role': role.toString().split('.').last,
      'enrollmentNumber': enrollmentNumber,
      'employeeId': employeeId,
      'department': department,
      'yearOfStudy': yearOfStudy,
      'deviceId': deviceId,
      'createdAt': createdAt.toIso8601String(),
      'updatedAt': updatedAt.toIso8601String(),
      'lastLogin': lastLogin?.toIso8601String(),
      'isActive': isActive,
    };
  }

  factory User.fromJson(Map<String, dynamic> json) {
    return User(
      id: json['id'] as String,
      email: json['email'] as String,
      fullName: json['fullName'] as String,
      role: UserRole.values.firstWhere(
        (role) => role.toString().split('.').last == json['role'],
      ),
      enrollmentNumber: json['enrollmentNumber'] as String?,
      employeeId: json['employeeId'] as String?,
      department: json['department'] as String?,
      yearOfStudy: json['yearOfStudy'] as int?,
      deviceId: json['deviceId'] as String?,
      createdAt: DateTime.parse(json['createdAt'] as String),
      updatedAt: DateTime.parse(json['updatedAt'] as String),
      lastLogin: json['lastLogin'] != null
          ? DateTime.parse(json['lastLogin'] as String)
          : null,
      isActive: json['isActive'] as bool? ?? true,
    );
  }

  @override
  List<Object?> get props => [
        id,
        email,
        fullName,
        role,
        enrollmentNumber,
        employeeId,
        department,
        yearOfStudy,
        deviceId,
        createdAt,
        updatedAt,
        lastLogin,
        isActive,
      ];
}
