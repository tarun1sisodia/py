import 'package:equatable/equatable.dart';

class Course extends Equatable {
  final String id;
  final String courseCode;
  final String courseName;
  final String department;
  final int yearOfStudy;
  final int semester;
  final DateTime createdAt;
  final DateTime updatedAt;

  const Course({
    required this.id,
    required this.courseCode,
    required this.courseName,
    required this.department,
    required this.yearOfStudy,
    required this.semester,
    required this.createdAt,
    required this.updatedAt,
  });

  Course copyWith({
    String? id,
    String? courseCode,
    String? courseName,
    String? department,
    int? yearOfStudy,
    int? semester,
    DateTime? createdAt,
    DateTime? updatedAt,
  }) {
    return Course(
      id: id ?? this.id,
      courseCode: courseCode ?? this.courseCode,
      courseName: courseName ?? this.courseName,
      department: department ?? this.department,
      yearOfStudy: yearOfStudy ?? this.yearOfStudy,
      semester: semester ?? this.semester,
      createdAt: createdAt ?? this.createdAt,
      updatedAt: updatedAt ?? this.updatedAt,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'courseCode': courseCode,
      'courseName': courseName,
      'department': department,
      'yearOfStudy': yearOfStudy,
      'semester': semester,
      'createdAt': createdAt.toIso8601String(),
      'updatedAt': updatedAt.toIso8601String(),
    };
  }

  factory Course.fromJson(Map<String, dynamic> json) {
    return Course(
      id: json['id'] as String,
      courseCode: json['courseCode'] as String,
      courseName: json['courseName'] as String,
      department: json['department'] as String,
      yearOfStudy: json['yearOfStudy'] as int,
      semester: json['semester'] as int,
      createdAt: DateTime.parse(json['createdAt'] as String),
      updatedAt: DateTime.parse(json['updatedAt'] as String),
    );
  }

  @override
  List<Object?> get props => [
        id,
        courseCode,
        courseName,
        department,
        yearOfStudy,
        semester,
        createdAt,
        updatedAt,
      ];
}
