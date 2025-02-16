import 'package:injectable/injectable.dart';
import 'package:smart_campus/core/services/database_service.dart';
import 'package:smart_campus/core/services/logger_service.dart';
import 'package:smart_campus/data/models/course_model.dart';
import 'package:smart_campus/data/repositories/base_repository.dart';

@lazySingleton
class CourseRepository extends BaseRepository<CourseModel> {
  static const _tableName = 'courses';

  CourseRepository(
    DatabaseService databaseService,
    LoggerService logger,
  ) : super(databaseService, logger, _tableName);

  @override
  CourseModel fromDatabase(Map<String, dynamic> data) =>
      CourseModel.fromDatabase(data);

  Future<CourseModel?> getCourseByCode(String code) async {
    try {
      final results = await getAll(
        where: 'code = ?',
        whereArgs: [code],
      );

      return results.isEmpty ? null : results.first;
    } catch (e) {
      logger.error('Error getting course by code', e);
      rethrow;
    }
  }

  Future<List<CourseModel>> getCoursesByDepartment(String department) async {
    try {
      return await getAll(
        where: 'department = ?',
        whereArgs: [department],
        orderBy: 'year_of_study ASC, semester ASC, name ASC',
      );
    } catch (e) {
      logger.error('Error getting courses by department', e);
      rethrow;
    }
  }

  Future<List<CourseModel>> getCoursesByYear(int yearOfStudy) async {
    try {
      return await getAll(
        where: 'year_of_study = ?',
        whereArgs: [yearOfStudy],
        orderBy: 'semester ASC, name ASC',
      );
    } catch (e) {
      logger.error('Error getting courses by year', e);
      rethrow;
    }
  }

  Future<List<CourseModel>> getCoursesBySemester(
    int yearOfStudy,
    int semester,
  ) async {
    try {
      return await getAll(
        where: 'year_of_study = ? AND semester = ?',
        whereArgs: [yearOfStudy, semester],
        orderBy: 'name ASC',
      );
    } catch (e) {
      logger.error('Error getting courses by semester', e);
      rethrow;
    }
  }

  Future<bool> isCourseCodeTaken(String code) async {
    try {
      final course = await getCourseByCode(code);
      return course != null;
    } catch (e) {
      logger.error('Error checking course code', e);
      rethrow;
    }
  }
}
