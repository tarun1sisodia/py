import 'package:json_annotation/json_annotation.dart';

part 'teacher_course_assignment_model.g.dart';

@JsonSerializable()
class TeacherCourseAssignmentModel {
	final String id;
	@JsonKey(name: 'teacher_id')
	final String teacherId;
	@JsonKey(name: 'course_id')
	final String courseId;
	@JsonKey(name: 'academic_year')
	final String academicYear;
	@JsonKey(name: 'is_active')
	final bool isActive;
	@JsonKey(name: 'created_at')
	final String createdAt;
	@JsonKey(name: 'updated_at')
	final String updatedAt;

	const TeacherCourseAssignmentModel({
		required this.id,
		required this.teacherId,
		required this.courseId,
		required this.academicYear,
		this.isActive = true,
		required this.createdAt,
		required this.updatedAt,
	});

	factory TeacherCourseAssignmentModel.fromJson(Map<String, dynamic> json) =>
			_$TeacherCourseAssignmentModelFromJson(json);

	Map<String, dynamic> toJson() => _$TeacherCourseAssignmentModelToJson(this);

	Map<String, dynamic> toDatabase() => {
				'id': id,
				'teacher_id': teacherId,
				'course_id': courseId,
				'academic_year': academicYear,
				'is_active': isActive ? 1 : 0,
				'created_at': createdAt,
				'updated_at': updatedAt,
			};

	factory TeacherCourseAssignmentModel.fromDatabase(Map<String, dynamic> data) =>
			TeacherCourseAssignmentModel(
				id: data['id'] as String,
				teacherId: data['teacher_id'] as String,
				courseId: data['course_id'] as String,
				academicYear: data['academic_year'] as String,
				isActive: (data['is_active'] as int) == 1,
				createdAt: data['created_at'] as String,
				updatedAt: data['updated_at'] as String,
			);

	TeacherCourseAssignmentModel copyWith({
		String? id,
		String? teacherId,
		String? courseId,
		String? academicYear,
		bool? isActive,
		String? createdAt,
		String? updatedAt,
	}) {
		return TeacherCourseAssignmentModel(
			id: id ?? this.id,
			teacherId: teacherId ?? this.teacherId,
			courseId: courseId ?? this.courseId,
			academicYear: academicYear ?? this.academicYear,
			isActive: isActive ?? this.isActive,
			createdAt: createdAt ?? this.createdAt,
			updatedAt: updatedAt ?? this.updatedAt,
		);
	}
}