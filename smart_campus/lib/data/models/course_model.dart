import 'package:json_annotation/json_annotation.dart';
import 'base_model.dart';

part 'course_model.g.dart';

enum CourseStatus {
  @JsonValue('active')
  active,
  @JsonValue('inactive')
  inactive,
  @JsonValue('archived')
  archived,
}

@JsonSerializable()
class CourseModel extends BaseModel {
  @JsonKey(name: 'course_code')
  final String courseCode;
  @JsonKey(name: 'course_name')
  final String courseName;
  final String department;
  @JsonKey(name: 'year_of_study')
  final int yearOfStudy;
  final int semester;
  final CourseStatus status;
  @JsonKey(name: 'total_students')
  final int? totalStudents;
  @JsonKey(name: 'total_sessions')
  final int? totalSessions;
  @JsonKey(name: 'active_sessions')
  final int? activeSessions;
  @JsonKey(name: 'completed_sessions')
  final int? completedSessions;
  @JsonKey(name: 'average_attendance')
  final double? averageAttendance;
  @JsonKey(name: 'assigned_teachers')
  final List<String>? assignedTeachers;

  const CourseModel({
    required super.id,
    required this.courseCode,
    required this.courseName,
    required this.department,
    required this.yearOfStudy,
    required this.semester,
    required this.status,
    this.totalStudents,
    this.totalSessions,
    this.activeSessions,
    this.completedSessions,
    this.averageAttendance,
    this.assignedTeachers,
    required super.createdAt,
    required super.updatedAt,
  });

  factory CourseModel.fromJson(Map<String, dynamic> json) =>
      _$CourseModelFromJson(json);

  @override
  Map<String, dynamic> toJson() => _$CourseModelToJson(this);

  @override
  Map<String, dynamic> toDatabase() => {
        'id': id,
        'course_code': courseCode,
        'course_name': courseName,
        'department': department,
        'year_of_study': yearOfStudy,
        'semester': semester,
        'status': status.toString().split('.').last,
        'total_students': totalStudents,
        'total_sessions': totalSessions,
        'active_sessions': activeSessions,
        'completed_sessions': completedSessions,
        'average_attendance': averageAttendance,
        'assigned_teachers': assignedTeachers?.join(','),
        'created_at': createdAt,
        'updated_at': updatedAt,
      };

  factory CourseModel.fromDatabase(Map<String, dynamic> data) => CourseModel(
        id: data['id'] as String,
        courseCode: data['course_code'] as String,
        courseName: data['course_name'] as String,
        department: data['department'] as String,
        yearOfStudy: data['year_of_study'] as int,
        semester: data['semester'] as int,
        status: CourseStatus.values.firstWhere(
          (status) => status.toString().split('.').last == data['status'],
          orElse: () => CourseStatus.inactive,
        ),
        totalStudents: data['total_students'] as int?,
        totalSessions: data['total_sessions'] as int?,
        activeSessions: data['active_sessions'] as int?,
        completedSessions: data['completed_sessions'] as int?,
        averageAttendance: data['average_attendance'] as double?,
        assignedTeachers: data['assigned_teachers'] != null
            ? (data['assigned_teachers'] as String).split(',')
            : null,
        createdAt: data['created_at'] as String,
        updatedAt: data['updated_at'] as String,
      );

  CourseModel copyWith({
    String? id,
    String? courseCode,
    String? courseName,
    String? department,
    int? yearOfStudy,
    int? semester,
    CourseStatus? status,
    int? totalStudents,
    int? totalSessions,
    int? activeSessions,
    int? completedSessions,
    double? averageAttendance,
    List<String>? assignedTeachers,
    String? createdAt,
    String? updatedAt,
  }) {
    return CourseModel(
      id: id ?? this.id,
      courseCode: courseCode ?? this.courseCode,
      courseName: courseName ?? this.courseName,
      department: department ?? this.department,
      yearOfStudy: yearOfStudy ?? this.yearOfStudy,
      semester: semester ?? this.semester,
      status: status ?? this.status,
      totalStudents: totalStudents ?? this.totalStudents,
      totalSessions: totalSessions ?? this.totalSessions,
      activeSessions: activeSessions ?? this.activeSessions,
      completedSessions: completedSessions ?? this.completedSessions,
      averageAttendance: averageAttendance ?? this.averageAttendance,
      assignedTeachers: assignedTeachers ?? this.assignedTeachers,
      createdAt: createdAt ?? this.createdAt,
      updatedAt: updatedAt ?? this.updatedAt,
    );
  }

  bool get isActive => status == CourseStatus.active;
  bool get isInactive => status == CourseStatus.inactive;
  bool get isArchived => status == CourseStatus.archived;

  bool get hasActiveSessions => (activeSessions ?? 0) > 0;
  bool get hasAssignedTeachers => (assignedTeachers?.isNotEmpty ?? false);
}
