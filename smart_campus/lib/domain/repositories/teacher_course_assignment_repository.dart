import '../entities/teacher_course_assignment.dart';

abstract class TeacherCourseAssignmentRepository {
	Future<List<TeacherCourseAssignment>> getAssignments({
		String? teacherId,
		String? courseId,
		String? academicYear,
		bool? isActive,
	});

	Future<TeacherCourseAssignment?> getAssignment(String id);

	Future<TeacherCourseAssignment> createAssignment(TeacherCourseAssignment assignment);

	Future<bool> updateAssignment(TeacherCourseAssignment assignment);

	Future<bool> deleteAssignment(String id);

	Future<List<TeacherCourseAssignment>> getTeacherAssignments(String teacherId);

	Future<List<TeacherCourseAssignment>> getCourseAssignments(String courseId);

	Future<bool> isTeacherAssignedToCourse({
		required String teacherId,
		required String courseId,
		required String academicYear,
	});
}