import 'package:injectable/injectable.dart';
import 'package:uuid/uuid.dart';
import 'package:smart_campus/core/errors/app_error.dart';
import 'package:smart_campus/core/services/database_service.dart';
import 'package:smart_campus/core/services/logger_service.dart';
import 'package:smart_campus/core/services/network_service.dart';
import 'package:smart_campus/domain/entities/course.dart';
import 'package:smart_campus/domain/repositories/course_repository.dart';

@LazySingleton(as: CourseRepository)
class CourseRepositoryImpl implements CourseRepository {
  final DatabaseService _databaseService;
  final NetworkService _networkService;
  final LoggerService _logger;
  final _uuid = const Uuid();

  CourseRepositoryImpl(
    this._databaseService,
    this._networkService,
    this._logger,
  );

  @override
  Future<List<Course>> getCourses({
    String? department,
    int? yearOfStudy,
    int? semester,
  }) async {
    try {
      // Build query conditions
      final conditions = <String>[];
      final args = <dynamic>[];

      if (department != null) {
        conditions.add('department = ?');
        args.add(department);
      }
      if (yearOfStudy != null) {
        conditions.add('year_of_study = ?');
        args.add(yearOfStudy);
      }
      if (semester != null) {
        conditions.add('semester = ?');
        args.add(semester);
      }

      final where = conditions.isEmpty ? null : conditions.join(' AND ');

      final results = await _databaseService.query(
        'courses',
        where: where,
        whereArgs: args.isEmpty ? null : args,
      );

      return results
          .map((map) => Course.fromJson(_mapToCourseJson(map)))
          .toList();
    } catch (e) {
      _logger.error('Error getting courses', e);
      throw AppError('Failed to get courses: $e');
    }
  }

  @override
  Future<Course> getCourseById(String id) async {
    try {
      final results = await _databaseService.query(
        'courses',
        where: 'id = ?',
        whereArgs: [id],
      );

      if (results.isEmpty) {
        throw AppError('Course not found');
      }

      return Course.fromJson(_mapToCourseJson(results.first));
    } catch (e) {
      _logger.error('Error getting course by ID', e);
      throw AppError('Failed to get course: $e');
    }
  }

  @override
  Future<Course> getCourseByCode(String courseCode) async {
    try {
      final results = await _databaseService.query(
        'courses',
        where: 'course_code = ?',
        whereArgs: [courseCode],
      );

      if (results.isEmpty) {
        throw AppError('Course not found');
      }

      return Course.fromJson(_mapToCourseJson(results.first));
    } catch (e) {
      _logger.error('Error getting course by code', e);
      throw AppError('Failed to get course: $e');
    }
  }

  @override
  Future<List<Course>> getTeacherCourses(String teacherId) async {
    try {
      final results = await _databaseService.query(
        'courses',
        where: '''
          id IN (
            SELECT course_id 
            FROM teacher_course_assignments 
            WHERE teacher_id = ? AND is_active = 1
          )
        ''',
        whereArgs: [teacherId],
      );

      return results
          .map((map) => Course.fromJson(_mapToCourseJson(map)))
          .toList();
    } catch (e) {
      _logger.error('Error getting teacher courses', e);
      throw AppError('Failed to get teacher courses: $e');
    }
  }

  @override
  Future<List<Course>> getStudentCourses(String studentId) async {
    try {
      // Get student details first
      final studentResults = await _databaseService.query(
        'users',
        where: 'id = ?',
        whereArgs: [studentId],
      );

      if (studentResults.isEmpty) {
        throw AppError('Student not found');
      }

      final yearOfStudy = studentResults.first['year_of_study'] as int;
      final department = studentResults.first['department'] as String;

      // Get courses for student's year and department
      final results = await _databaseService.query(
        'courses',
        where: 'year_of_study = ? AND department = ?',
        whereArgs: [yearOfStudy, department],
      );

      return results
          .map((map) => Course.fromJson(_mapToCourseJson(map)))
          .toList();
    } catch (e) {
      _logger.error('Error getting student courses', e);
      throw AppError('Failed to get student courses: $e');
    }
  }

  @override
  Future<Course> createCourse(Course course) async {
    try {
      final newCourse = course.copyWith(
        id: _uuid.v4(),
        createdAt: DateTime.now(),
        updatedAt: DateTime.now(),
      );

      await _databaseService.insert(
        'courses',
        _courseToMap(newCourse),
      );

      return newCourse;
    } catch (e) {
      _logger.error('Error creating course', e);
      throw AppError('Failed to create course: $e');
    }
  }

  @override
  Future<Course> updateCourse(Course course) async {
    try {
      final updatedCourse = course.copyWith(
        updatedAt: DateTime.now(),
      );

      final count = await _databaseService.update(
        'courses',
        _courseToMap(updatedCourse),
        where: 'id = ?',
        whereArgs: [course.id],
      );

      if (count == 0) {
        throw AppError('Course not found');
      }

      return updatedCourse;
    } catch (e) {
      _logger.error('Error updating course', e);
      throw AppError('Failed to update course: $e');
    }
  }

  @override
  Future<void> deleteCourse(String id) async {
    try {
      final count = await _databaseService.delete(
        'courses',
        where: 'id = ?',
        whereArgs: [id],
      );

      if (count == 0) {
        throw AppError('Course not found');
      }
    } catch (e) {
      _logger.error('Error deleting course', e);
      throw AppError('Failed to delete course: $e');
    }
  }

  @override
  Future<void> assignTeacherToCourse({
    required String teacherId,
    required String courseId,
    required String academicYear,
  }) async {
    try {
      await _databaseService.insert(
        'teacher_course_assignments',
        {
          'id': _uuid.v4(),
          'teacher_id': teacherId,
          'course_id': courseId,
          'academic_year': academicYear,
          'is_active': 1,
          'created_at': DateTime.now().toIso8601String(),
          'updated_at': DateTime.now().toIso8601String(),
        },
      );
    } catch (e) {
      _logger.error('Error assigning teacher to course', e);
      throw AppError('Failed to assign teacher to course: $e');
    }
  }

  @override
  Future<void> removeTeacherFromCourse({
    required String teacherId,
    required String courseId,
    required String academicYear,
  }) async {
    try {
      final count = await _databaseService.update(
        'teacher_course_assignments',
        {'is_active': 0},
        where: 'teacher_id = ? AND course_id = ? AND academic_year = ?',
        whereArgs: [teacherId, courseId, academicYear],
      );

      if (count == 0) {
        throw AppError('Assignment not found');
      }
    } catch (e) {
      _logger.error('Error removing teacher from course', e);
      throw AppError('Failed to remove teacher from course: $e');
    }
  }

  @override
  Future<List<Course>> searchCourses(String query) async {
    try {
      final results = await _databaseService.query(
        'courses',
        where: 'course_code LIKE ? OR course_name LIKE ?',
        whereArgs: ['%$query%', '%$query%'],
      );

      return results
          .map((map) => Course.fromJson(_mapToCourseJson(map)))
          .toList();
    } catch (e) {
      _logger.error('Error searching courses', e);
      throw AppError('Failed to search courses: $e');
    }
  }

  @override
  Future<bool> isCourseCodeUnique(String courseCode) async {
    try {
      final results = await _databaseService.query(
        'courses',
        where: 'course_code = ?',
        whereArgs: [courseCode],
      );

      return results.isEmpty;
    } catch (e) {
      _logger.error('Error checking course code uniqueness', e);
      throw AppError('Failed to check course code uniqueness: $e');
    }
  }

  // Helper methods to convert between database and entity formats
  Map<String, dynamic> _courseToMap(Course course) {
    return {
      'id': course.id,
      'course_code': course.courseCode,
      'course_name': course.courseName,
      'department': course.department,
      'year_of_study': course.yearOfStudy,
      'semester': course.semester,
      'created_at': course.createdAt.toIso8601String(),
      'updated_at': course.updatedAt.toIso8601String(),
    };
  }

  Map<String, dynamic> _mapToCourseJson(Map<String, dynamic> map) {
    return {
      'id': map['id'],
      'courseCode': map['course_code'],
      'courseName': map['course_name'],
      'department': map['department'],
      'yearOfStudy': map['year_of_study'],
      'semester': map['semester'],
      'createdAt': map['created_at'],
      'updatedAt': map['updated_at'],
    };
  }
}
