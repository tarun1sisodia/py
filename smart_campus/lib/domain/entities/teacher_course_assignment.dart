class TeacherCourseAssignment {
	final String id;
	final String teacherId;
	final String courseId;
	final String academicYear;
	final bool isActive;
	final String createdAt;
	final String updatedAt;

	const TeacherCourseAssignment({
		required this.id,
		required this.teacherId,
		required this.courseId,
		required this.academicYear,
		this.isActive = true,
		required this.createdAt,
		required this.updatedAt,
	});

	TeacherCourseAssignment copyWith({
		String? id,
		String? teacherId,
		String? courseId,
		String? academicYear,
		bool? isActive,
		String? createdAt,
		String? updatedAt,
	}) {
		return TeacherCourseAssignment(
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