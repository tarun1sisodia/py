import 'dart:async';
import 'package:injectable/injectable.dart';
import 'package:path/path.dart';
import 'package:sqflite/sqflite.dart';
import 'package:smart_campus/core/services/logger_service.dart';
import 'package:smart_campus/core/errors/app_error.dart';

@lazySingleton
class DatabaseService {
  static const String _databaseName = 'smart_campus.db';
  static const int _databaseVersion = 1;

  final LoggerService _logger;
  Database? _database;

  DatabaseService(this._logger);

  Future<Database> get database async {
    _database ??= await _initDatabase();
    return _database!;
  }

  Future<Database> _initDatabase() async {
    try {
      final databasePath = await getDatabasesPath();
      final path = join(databasePath, _databaseName);

      return await openDatabase(
        path,
        version: _databaseVersion,
        onCreate: _onCreate,
        onUpgrade: _onUpgrade,
        onDowngrade: onDatabaseDowngradeDelete,
      );
    } catch (e) {
      _logger.error('Error initializing database', e);
      throw AppError('Failed to initialize database: $e');
    }
  }

  Future<void> _onCreate(Database db, int version) async {
    try {
      await db.transaction((txn) async {
        // Users table
        await txn.execute('''
          CREATE TABLE users (
            id TEXT PRIMARY KEY,
            email TEXT UNIQUE NOT NULL,
            role TEXT NOT NULL,
            full_name TEXT NOT NULL,
            enrollment_number TEXT UNIQUE,
            created_at INTEGER NOT NULL,
            updated_at INTEGER NOT NULL
          )
        ''');

        // Courses table
        await txn.execute('''
          CREATE TABLE courses (
            id TEXT PRIMARY KEY,
            code TEXT UNIQUE NOT NULL,
            name TEXT NOT NULL,
            department TEXT NOT NULL,
            year_of_study INTEGER NOT NULL,
            semester INTEGER NOT NULL,
            created_at INTEGER NOT NULL,
            updated_at INTEGER NOT NULL
          )
        ''');

        // Attendance sessions table
        await txn.execute('''
          CREATE TABLE attendance_sessions (
            id TEXT PRIMARY KEY,
            course_id TEXT NOT NULL,
            teacher_id TEXT NOT NULL,
            date TEXT NOT NULL,
            start_time INTEGER NOT NULL,
            end_time INTEGER NOT NULL,
            latitude REAL NOT NULL,
            longitude REAL NOT NULL,
            radius INTEGER NOT NULL,
            wifi_ssid TEXT NOT NULL,
            wifi_bssid TEXT NOT NULL,
            status TEXT NOT NULL,
            created_at INTEGER NOT NULL,
            updated_at INTEGER NOT NULL,
            FOREIGN KEY (course_id) REFERENCES courses (id),
            FOREIGN KEY (teacher_id) REFERENCES users (id)
          )
        ''');

        // Attendance records table
        await txn.execute('''
          CREATE TABLE attendance_records (
            id TEXT PRIMARY KEY,
            session_id TEXT NOT NULL,
            student_id TEXT NOT NULL,
            marked_at INTEGER NOT NULL,
            latitude REAL NOT NULL,
            longitude REAL NOT NULL,
            wifi_ssid TEXT NOT NULL,
            wifi_bssid TEXT NOT NULL,
            device_id TEXT NOT NULL,
            sync_status TEXT NOT NULL,
            created_at INTEGER NOT NULL,
            updated_at INTEGER NOT NULL,
            FOREIGN KEY (session_id) REFERENCES attendance_sessions (id),
            FOREIGN KEY (student_id) REFERENCES users (id)
          )
        ''');

        // Device bindings table
        await txn.execute('''
          CREATE TABLE device_bindings (
            id TEXT PRIMARY KEY,
            user_id TEXT NOT NULL,
            device_id TEXT NOT NULL,
            device_info TEXT NOT NULL,
            bound_at INTEGER NOT NULL,
            created_at INTEGER NOT NULL,
            updated_at INTEGER NOT NULL,
            FOREIGN KEY (user_id) REFERENCES users (id)
          )
        ''');

        _logger.info('Database tables created successfully');
      });
    } catch (e) {
      _logger.error('Error creating database tables', e);
      throw AppError('Failed to create database tables: $e');
    }
  }

  Future<void> _onUpgrade(Database db, int oldVersion, int newVersion) async {
    try {
      // Handle database upgrades here
      _logger.info('Database upgraded from $oldVersion to $newVersion');
    } catch (e) {
      _logger.error('Error upgrading database', e);
      throw AppError('Failed to upgrade database: $e');
    }
  }

  Future<void> close() async {
    try {
      if (_database != null) {
        await _database!.close();
        _database = null;
      }
    } catch (e) {
      _logger.error('Error closing database', e);
      throw AppError('Failed to close database: $e');
    }
  }

  // Helper methods for CRUD operations
  Future<int> insert(String table, Map<String, dynamic> data) async {
    try {
      final db = await database;
      return await db.insert(
        table,
        {
          ...data,
          'created_at': DateTime.now().millisecondsSinceEpoch,
          'updated_at': DateTime.now().millisecondsSinceEpoch,
        },
        conflictAlgorithm: ConflictAlgorithm.replace,
      );
    } catch (e) {
      _logger.error('Error inserting data into $table', e);
      throw AppError('Failed to insert data: $e');
    }
  }

  Future<List<Map<String, dynamic>>> query(
    String table, {
    bool? distinct,
    List<String>? columns,
    String? where,
    List<dynamic>? whereArgs,
    String? groupBy,
    String? having,
    String? orderBy,
    int? limit,
    int? offset,
  }) async {
    try {
      final db = await database;
      return await db.query(
        table,
        distinct: distinct,
        columns: columns,
        where: where,
        whereArgs: whereArgs,
        groupBy: groupBy,
        having: having,
        orderBy: orderBy,
        limit: limit,
        offset: offset,
      );
    } catch (e) {
      _logger.error('Error querying $table', e);
      throw AppError('Failed to query data: $e');
    }
  }

  Future<int> update(
    String table,
    Map<String, dynamic> data, {
    String? where,
    List<dynamic>? whereArgs,
  }) async {
    try {
      final db = await database;
      return await db.update(
        table,
        {
          ...data,
          'updated_at': DateTime.now().millisecondsSinceEpoch,
        },
        where: where,
        whereArgs: whereArgs,
      );
    } catch (e) {
      _logger.error('Error updating data in $table', e);
      throw AppError('Failed to update data: $e');
    }
  }

  Future<int> delete(
    String table, {
    String? where,
    List<dynamic>? whereArgs,
  }) async {
    try {
      final db = await database;
      return await db.delete(
        table,
        where: where,
        whereArgs: whereArgs,
      );
    } catch (e) {
      _logger.error('Error deleting data from $table', e);
      throw AppError('Failed to delete data: $e');
    }
  }

  Future<T> transaction<T>(Future<T> Function(Transaction txn) action) async {
    try {
      final db = await database;
      return await db.transaction(action);
    } catch (e) {
      _logger.error('Error executing transaction', e);
      throw AppError('Failed to execute transaction: $e');
    }
  }
}
