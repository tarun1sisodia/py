import 'package:injectable/injectable.dart';
import 'package:sqflite/sqflite.dart';

@lazySingleton
class DatabaseHelper {
  static const String _databaseName = 'smart_campus.db';
  static const int _databaseVersion = 1;

  final Database _database;

  DatabaseHelper(this._database);

  Future<Database> get database async {
    return _database;
  }

  Future<void> _onCreate(Database db, int version) async {
    // Users table
    await db.execute('''
      CREATE TABLE users (
        id TEXT PRIMARY KEY,
        email TEXT NOT NULL UNIQUE,
        full_name TEXT NOT NULL,
        role TEXT NOT NULL,
        department TEXT,
        year_of_study INTEGER,
        enrollment_number TEXT,
        employee_id TEXT,
        device_id TEXT,
        created_at TEXT NOT NULL,
        updated_at TEXT NOT NULL
      )
    ''');

    // Courses table
    await db.execute('''
      CREATE TABLE courses (
        id TEXT PRIMARY KEY,
        code TEXT NOT NULL UNIQUE,
        name TEXT NOT NULL,
        department TEXT NOT NULL,
        year_of_study INTEGER NOT NULL,
        semester INTEGER NOT NULL,
        created_at TEXT NOT NULL,
        updated_at TEXT NOT NULL
      )
    ''');

    // Teacher-Course assignments
    await db.execute('''
      CREATE TABLE teacher_course_assignments (
        id TEXT PRIMARY KEY,
        teacher_id TEXT NOT NULL,
        course_id TEXT NOT NULL,
        created_at TEXT NOT NULL,
        updated_at TEXT NOT NULL,
        FOREIGN KEY (teacher_id) REFERENCES users (id) ON DELETE CASCADE,
        FOREIGN KEY (course_id) REFERENCES courses (id) ON DELETE CASCADE
      )
    ''');

    // Attendance sessions
    await db.execute('''
      CREATE TABLE attendance_sessions (
        id TEXT PRIMARY KEY,
        teacher_id TEXT NOT NULL,
        course_id TEXT NOT NULL,
        session_date TEXT NOT NULL,
        start_time TEXT NOT NULL,
        end_time TEXT NOT NULL,
        wifi_ssid TEXT,
        wifi_bssid TEXT,
        location_latitude REAL,
        location_longitude REAL,
        location_radius INTEGER,
        status TEXT NOT NULL,
        created_at TEXT NOT NULL,
        updated_at TEXT NOT NULL,
        FOREIGN KEY (teacher_id) REFERENCES users (id) ON DELETE CASCADE,
        FOREIGN KEY (course_id) REFERENCES courses (id) ON DELETE CASCADE
      )
    ''');

    // Attendance records
    await db.execute('''
      CREATE TABLE attendance_records (
        id TEXT PRIMARY KEY,
        session_id TEXT NOT NULL,
        student_id TEXT NOT NULL,
        marked_at TEXT NOT NULL,
        wifi_ssid TEXT,
        wifi_bssid TEXT,
        location_latitude REAL,
        location_longitude REAL,
        device_id TEXT NOT NULL,
        verification_status TEXT NOT NULL,
        rejection_reason TEXT,
        created_at TEXT NOT NULL,
        updated_at TEXT NOT NULL,
        FOREIGN KEY (session_id) REFERENCES attendance_sessions (id) ON DELETE CASCADE,
        FOREIGN KEY (student_id) REFERENCES users (id) ON DELETE CASCADE
      )
    ''');

    // Device bindings
    await db.execute('''
      CREATE TABLE device_bindings (
        id TEXT PRIMARY KEY,
        student_id TEXT NOT NULL,
        device_id TEXT NOT NULL,
        device_name TEXT NOT NULL,
        device_model TEXT NOT NULL,
        is_active INTEGER NOT NULL DEFAULT 1,
        is_blacklisted INTEGER NOT NULL DEFAULT 0,
        created_at TEXT NOT NULL,
        updated_at TEXT NOT NULL,
        FOREIGN KEY (student_id) REFERENCES users (id) ON DELETE CASCADE
      )
    ''');

    // Create indexes
    await db.execute(
        'CREATE INDEX idx_attendance_sessions_course ON attendance_sessions(course_id)');
    await db.execute(
        'CREATE INDEX idx_attendance_sessions_teacher ON attendance_sessions(teacher_id)');
    await db.execute(
        'CREATE INDEX idx_attendance_records_session ON attendance_records(session_id)');
    await db.execute(
        'CREATE INDEX idx_attendance_records_student ON attendance_records(student_id)');
    await db.execute(
        'CREATE INDEX idx_device_bindings_student ON device_bindings(student_id)');
    await db.execute(
        'CREATE INDEX idx_device_bindings_device ON device_bindings(device_id)');
  }

  Future<void> _onUpgrade(Database db, int oldVersion, int newVersion) async {
    // Handle database schema upgrades here
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
    return await _database.query(
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
  }

  Future<int> insert(
    String table,
    Map<String, dynamic> values, {
    ConflictAlgorithm? conflictAlgorithm,
  }) async {
    return await _database.insert(
      table,
      values,
      conflictAlgorithm: conflictAlgorithm,
    );
  }

  Future<int> update(
    String table,
    Map<String, dynamic> values, {
    String? where,
    List<dynamic>? whereArgs,
  }) async {
    return await _database.update(
      table,
      values,
      where: where,
      whereArgs: whereArgs,
    );
  }

  Future<int> delete(
    String table, {
    String? where,
    List<dynamic>? whereArgs,
  }) async {
    return await _database.delete(
      table,
      where: where,
      whereArgs: whereArgs,
    );
  }

  Future<T> transaction<T>(Future<T> Function(Transaction txn) action) async {
    return await _database.transaction(action);
  }

  Future<void> close() async {
    await _database.close();
  }
}
