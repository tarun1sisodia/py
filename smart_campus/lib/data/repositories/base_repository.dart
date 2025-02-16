import 'package:smart_campus/core/errors/app_error.dart';
import 'package:smart_campus/core/services/database_service.dart';
import 'package:smart_campus/core/services/logger_service.dart';

abstract class BaseRepository<T> {
  final DatabaseService _databaseService;
  final LoggerService logger;
  final String tableName;

  BaseRepository(
    this._databaseService,
    this.logger,
    this.tableName,
  );

  Future<T> create(T model) async {
    try {
      final data = (model as dynamic).toDatabase();
      await _databaseService.insert(tableName, data);
      return model;
    } catch (e) {
      logger.error('Error creating $tableName', e);
      throw AppError('Failed to create $tableName: $e');
    }
  }

  Future<T?> getById(String id) async {
    try {
      final results = await _databaseService.query(
        tableName,
        where: 'id = ?',
        whereArgs: [id],
      );

      if (results.isEmpty) {
        return null;
      }

      return fromDatabase(results.first);
    } catch (e) {
      logger.error('Error getting $tableName by id', e);
      throw AppError('Failed to get $tableName: $e');
    }
  }

  Future<List<T>> getAll({
    String? where,
    List<dynamic>? whereArgs,
    String? orderBy,
    int? limit,
    int? offset,
  }) async {
    try {
      final results = await _databaseService.query(
        tableName,
        where: where,
        whereArgs: whereArgs,
        orderBy: orderBy,
        limit: limit,
        offset: offset,
      );

      return results.map((data) => fromDatabase(data)).toList();
    } catch (e) {
      logger.error('Error getting all $tableName', e);
      throw AppError('Failed to get all $tableName: $e');
    }
  }

  Future<bool> update(T model) async {
    try {
      final data = (model as dynamic).toDatabase();
      final count = await _databaseService.update(
        tableName,
        data,
        where: 'id = ?',
        whereArgs: [data['id']],
      );
      return count > 0;
    } catch (e) {
      logger.error('Error updating $tableName', e);
      throw AppError('Failed to update $tableName: $e');
    }
  }

  Future<bool> delete(String id) async {
    try {
      final count = await _databaseService.delete(
        tableName,
        where: 'id = ?',
        whereArgs: [id],
      );
      return count > 0;
    } catch (e) {
      logger.error('Error deleting $tableName', e);
      throw AppError('Failed to delete $tableName: $e');
    }
  }

  Future<void> deleteAll() async {
    try {
      await _databaseService.delete(tableName);
    } catch (e) {
      logger.error('Error deleting all $tableName', e);
      throw AppError('Failed to delete all $tableName: $e');
    }
  }

  T fromDatabase(Map<String, dynamic> data);
}
