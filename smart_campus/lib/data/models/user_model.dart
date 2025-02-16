import 'package:json_annotation/json_annotation.dart';
import 'base_model.dart';

part 'user_model.g.dart';

enum UserRole {
  @JsonValue('admin')
  admin,
  @JsonValue('student')
  student,
  @JsonValue('teacher')
  teacher,
  @JsonValue('staff')
  staff,
}

@JsonSerializable()
class UserModel extends BaseModel {
  final String email;
  final UserRole role;
  @JsonKey(name: 'full_name')
  final String fullName;
  @JsonKey(name: 'enrollment_number')
  final String? enrollmentNumber;
  @JsonKey(name: 'employee_id')
  final String? employeeId;
  final String? department;
  @JsonKey(name: 'year_of_study')
  final int? yearOfStudy;
  @JsonKey(name: 'device_id')
  final String? deviceId;
  @JsonKey(name: 'last_login')
  final String? lastLogin;
  @JsonKey(name: 'is_active')
  final bool isActive;

  const UserModel({
    required super.id,
    required this.email,
    required this.role,
    required this.fullName,
    this.enrollmentNumber,
    this.employeeId,
    this.department,
    this.yearOfStudy,
    this.deviceId,
    required super.createdAt,
    required super.updatedAt,
    this.lastLogin,
    this.isActive = true,
  });

  factory UserModel.fromJson(Map<String, dynamic> json) =>
      _$UserModelFromJson(json);

  @override
  Map<String, dynamic> toJson() => _$UserModelToJson(this);

  @override
  Map<String, dynamic> toDatabase() => {
        'id': id,
        'email': email,
        'role': role.toString().split('.').last,
        'full_name': fullName,
        'enrollment_number': enrollmentNumber,
        'employee_id': employeeId,
        'department': department,
        'year_of_study': yearOfStudy,
        'device_id': deviceId,
        'created_at': createdAt,
        'updated_at': updatedAt,
        'last_login': lastLogin,
        'is_active': isActive ? 1 : 0,
      };

  factory UserModel.fromDatabase(Map<String, dynamic> data) => UserModel(
        id: data['id'] as String,
        email: data['email'] as String,
        role: UserRole.values.firstWhere(
          (role) => role.toString().split('.').last == data['role'],
          orElse: () => UserRole.student,
        ),
        fullName: data['full_name'] as String,
        enrollmentNumber: data['enrollment_number'] as String?,
        employeeId: data['employee_id'] as String?,
        department: data['department'] as String?,
        yearOfStudy: data['year_of_study'] as int?,
        deviceId: data['device_id'] as String?,
        createdAt: data['created_at'] as String,
        updatedAt: data['updated_at'] as String,
        lastLogin: data['last_login'] as String?,
        isActive: (data['is_active'] as int?) == 1,
      );

  UserModel copyWith({
    String? id,
    String? email,
    UserRole? role,
    String? fullName,
    String? enrollmentNumber,
    String? employeeId,
    String? department,
    int? yearOfStudy,
    String? deviceId,
    String? createdAt,
    String? updatedAt,
    String? lastLogin,
    bool? isActive,
  }) {
    return UserModel(
      id: id ?? this.id,
      email: email ?? this.email,
      role: role ?? this.role,
      fullName: fullName ?? this.fullName,
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

  bool get isStudent => role == UserRole.student;
  bool get isTeacher => role == UserRole.teacher;
  bool get isAdmin => role == UserRole.admin;
  bool get isStaff => role == UserRole.staff;
}
