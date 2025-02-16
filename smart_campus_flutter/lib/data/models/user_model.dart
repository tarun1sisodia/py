import 'package:smart_campus_flutter/domain/entities/user.dart';

class UserModel extends User {
  const UserModel({
    required super.id,
    required super.fullName,
    required super.email,
    required super.phoneNumber,
    required super.role,
    super.academicYear,
    super.course,
    super.highestDegree,
    super.experience,
  });

  factory UserModel.fromJson(Map<String, dynamic> json) {
    return UserModel(
      id: json['id'],
      fullName: json['full_name'],
      email: json['email'] ?? '',
      phoneNumber: json['phone'],
      role: json['role'],
      academicYear: json['academic_year'],
      course: json['course'],
      highestDegree: json['highest_degree'],
      experience:
          json['experience'] != null
              ? int.parse(json['experience'].toString())
              : null,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'full_name': fullName,
      'email': email,
      'phone': phoneNumber,
      'role': role,
      if (academicYear != null) 'academic_year': academicYear,
      if (course != null) 'course': course,
      if (highestDegree != null) 'highest_degree': highestDegree,
      if (experience != null) 'experience': experience,
    };
  }
}
