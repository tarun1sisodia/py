import 'package:injectable/injectable.dart';
import 'package:smart_campus/core/services/database_service.dart';
import 'package:smart_campus/core/services/logger_service.dart';
import 'package:smart_campus/data/models/user_model.dart';
import 'package:smart_campus/data/repositories/base_repository.dart';

@lazySingleton
class UserRepository extends BaseRepository<UserModel> {
  static const _tableName = 'users';

  UserRepository(
    DatabaseService databaseService,
    LoggerService logger,
  ) : super(databaseService, logger, _tableName);

  @override
  UserModel fromDatabase(Map<String, dynamic> data) =>
      UserModel.fromDatabase(data);

  Future<UserModel?> getUserByEmail(String email) async {
    try {
      final results = await getAll(
        where: 'email = ?',
        whereArgs: [email],
      );

      return results.isEmpty ? null : results.first;
    } catch (e) {
      logger.error('Error getting user by email', e);
      rethrow;
    }
  }

  Future<List<UserModel>> getStudents({
    String? department,
    int? yearOfStudy,
  }) async {
    try {
      String? where;
      List<dynamic>? whereArgs;

      if (department != null || yearOfStudy != null) {
        final conditions = <String>[];
        whereArgs = [];

        conditions.add('role = ?');
        whereArgs.add('student');

        if (department != null) {
          conditions.add('department = ?');
          whereArgs.add(department);
        }

        if (yearOfStudy != null) {
          conditions.add('year_of_study = ?');
          whereArgs.add(yearOfStudy);
        }

        where = conditions.join(' AND ');
      } else {
        where = 'role = ?';
        whereArgs = ['student'];
      }

      return await getAll(
        where: where,
        whereArgs: whereArgs,
        orderBy: 'full_name ASC',
      );
    } catch (e) {
      logger.error('Error getting students', e);
      rethrow;
    }
  }

  Future<List<UserModel>> getTeachers({String? department}) async {
    try {
      String where = 'role = ?';
      final whereArgs = ['teacher'];

      if (department != null) {
        where += ' AND department = ?';
        whereArgs.add(department);
      }

      return await getAll(
        where: where,
        whereArgs: whereArgs,
        orderBy: 'full_name ASC',
      );
    } catch (e) {
      logger.error('Error getting teachers', e);
      rethrow;
    }
  }

  Future<bool> isEmailTaken(String email) async {
    try {
      final user = await getUserByEmail(email);
      return user != null;
    } catch (e) {
      logger.error('Error checking email', e);
      rethrow;
    }
  }
}
