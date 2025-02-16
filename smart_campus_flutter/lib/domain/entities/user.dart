import 'package:equatable/equatable.dart';

class User extends Equatable {
  final String id;
  final String fullName;
  final String email;
  final String phoneNumber;
  final String role; // 'teacher' or 'student'
  final String? academicYear; // Only for students
  final String? course; // Only for students
  final String? highestDegree; // Only for teachers
  final int? experience; // Only for teachers

  const User({
    required this.id,
    required this.fullName,
    required this.email,
    required this.phoneNumber,
    required this.role,
    this.academicYear,
    this.course,
    this.highestDegree,
    this.experience,
  });

  @override
  List<Object?> get props => [
    id,
    fullName,
    email,
    phoneNumber,
    role,
    academicYear,
    course,
    highestDegree,
    experience,
  ];

  bool get isTeacher => role == 'teacher';
  bool get isStudent => role == 'student';

  User copyWith({
    String? id,
    String? fullName,
    String? email,
    String? phoneNumber,
    String? role,
    String? academicYear,
    String? course,
    String? highestDegree,
    int? experience,
  }) {
    return User(
      id: id ?? this.id,
      fullName: fullName ?? this.fullName,
      email: email ?? this.email,
      phoneNumber: phoneNumber ?? this.phoneNumber,
      role: role ?? this.role,
      academicYear: academicYear ?? this.academicYear,
      course: course ?? this.course,
      highestDegree: highestDegree ?? this.highestDegree,
      experience: experience ?? this.experience,
    );
  }
}
