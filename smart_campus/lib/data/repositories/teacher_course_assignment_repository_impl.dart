import 'package:injectable/injectable.dart';
import 'package:smart_campus/core/services/database_service.dart';
import 'package:smart_campus/core/services/logger_service.dart';
import 'package:smart_campus/data/models/teacher_course_assignment_model.dart';
import 'package:smart_campus/domain/entities/teacher_course_assignment.dart';
import 'package:smart_campus/domain/repositories/teacher_course_assignment_repository.dart';
import 'package:smart_campus/data/repositories/base_repository.dart';

@LazySingleton(as: TeacherCourseAssignmentRepository)
class TeacherCourseAssignmentRepositoryImpl extends BaseRepository<TeacherCourseAssignmentModel>
	implements TeacherCourseAssignmentRepository {
  TeacherCourseAssignmentRepositoryImpl(
	DatabaseService databaseService,
	LoggerService logger,
  ) : super(databaseService, logger, 'teacher_course_assignments');

  @override
  TeacherCourseAssignmentModel fromDatabase(Map<String, dynamic> data) =>
	  TeacherCourseAssignmentModel.fromDatabase(data);

  @override
  Future<List<TeacherCourseAssignment>> getAssignments({
	String? teacherId,
	String? courseId,
	String? academicYear,
	bool? isActive,
  }) async {
	final conditions = <String>[];
	final args = <dynamic>[];

	if (teacherId != null) {
	  conditions.add('teacher_id = ?');
	  args.add(teacherId);
	}
	if (courseId != null) {
	  conditions.add('course_id = ?');
	  args.add(courseId);
	}
	if (academicYear != null) {
	  conditions.add('academic_year = ?');
	  args.add(academicYear);
	}
	if (isActive != null) {
	  conditions.add('is_active = ?');
	  args.add(isActive ? 1 : 0);
	}

	final assignments = await getAll(
	  where: conditions.isEmpty ? null : conditions.join(' AND '),
	  whereArgs: args.isEmpty ? null : args,
	  orderBy: 'created_at DESC',
	);

	return assignments;
  }

  @override
  Future<TeacherCourseAssignment?> getAssignment(String id) => getById(id);

  @override
  Future<TeacherCourseAssignment> createAssignment(TeacherCourseAssignment assignment) async {
	final model = TeacherCourseAssignmentModel(
	  id: assignment.id,
	  teacherId: assignment.teacherId,
	  courseId: assignment.courseId,
	  academicYear: assignment.academicYear,
	  isActive: assignment.isActive,
	  createdAt: assignment.createdAt,
	  updatedAt: assignment.updatedAt,
	);
	return await create(model);
  }

  @override
  Future<bool> updateAssignment(TeacherCourseAssignment assignment) async {
	final model = TeacherCourseAssignmentModel(
	  id: assignment.id,
	  teacherId: assignment.teacherId,
	  courseId: assignment.courseId,
	  academicYear: assignment.academicYear,
	  isActive: assignment.isActive,
	  createdAt: assignment.createdAt,
	  updatedAt: assignment.updatedAt,
	);
	return await update(model);
  }

  @override
  Future<bool> deleteAssignment(String id) => delete(id);

  @override
  Future<List<TeacherCourseAssignment>> getTeacherAssignments(String teacherId) =>
	  getAssignments(teacherId: teacherId, isActive: true);

  @override
  Future<List<TeacherCourseAssignment>> getCourseAssignments(String courseId) =>
	  getAssignments(courseId: courseId, isActive: true);

  @override
  Future<bool> isTeacherAssignedToCourse({
	required String teacherId,
	required String courseId,
	required String academicYear,
  }) async {
	final assignments = await getAssignments(
	  teacherId: teacherId,
	  courseId: courseId,
	  academicYear: academicYear,
	  isActive: true,
	);
	return assignments.isNotEmpty;
  }
}