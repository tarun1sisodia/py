import 'package:path/path.dart';
import 'package:sqflite/sqflite.dart';

class DatabaseHelper {
  static const _databaseName = "smart_campus.db";
  static const _databaseVersion = 1;

  // Singleton pattern
  DatabaseHelper._();
  static final DatabaseHelper instance = DatabaseHelper._();

  static Database? _database;

  Future<Database> get database async {
    _database ??= await _initDatabase();
    return _database!;
  }

  Future<Database> _initDatabase() async {
    final String path = join(await getDatabasesPath(), _databaseName);
    return await openDatabase(
      path,
      version: _databaseVersion,
      onCreate: _onCreate,
      onUpgrade: _onUpgrade,
    );
  }

  Future<void> _onCreate(Database db, int version) async {
    // Users table
    await db.execute('''
      CREATE TABLE users (
        id TEXT PRIMARY KEY,
        role TEXT NOT NULL,
        email TEXT UNIQUE NOT NULL,
        full_name TEXT NOT NULL,
        enrollment_number TEXT UNIQUE,
        employee_id TEXT UNIQUE,
        department TEXT,
        year_of_study INTEGER,
        device_id TEXT,
        created_at TEXT NOT NULL,
        updated_at TEXT NOT NULL,
        last_login TEXT,
        is_active INTEGER NOT NULL DEFAULT 1
      )
    ''');

    // Courses table
    await db.execute('''
      CREATE TABLE courses (
        id TEXT PRIMARY KEY,
        course_code TEXT UNIQUE NOT NULL,
        course_name TEXT NOT NULL,
        department TEXT NOT NULL,
        year_of_study INTEGER NOT NULL,
        semester INTEGER NOT NULL,
        created_at TEXT NOT NULL,
        updated_at TEXT NOT NULL
      )
    ''');

    // Teacher Course Assignments table
    await db.execute('''
      CREATE TABLE teacher_course_assignments (
        id TEXT PRIMARY KEY,
        teacher_id TEXT NOT NULL,
        course_id TEXT NOT NULL,
        academic_year TEXT NOT NULL,
        is_active INTEGER NOT NULL DEFAULT 1,
        created_at TEXT NOT NULL,
        updated_at TEXT NOT NULL,
        FOREIGN KEY (teacher_id) REFERENCES users (id),
        FOREIGN KEY (course_id) REFERENCES courses (id),
        UNIQUE (teacher_id, course_id, academic_year)
      )
    ''');

    // Attendance Sessions table
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
        FOREIGN KEY (teacher_id) REFERENCES users (id),
        FOREIGN KEY (course_id) REFERENCES courses (id)
      )
    ''');

    // Attendance Records table
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
        device_id TEXT,
        verification_status TEXT NOT NULL,
        rejection_reason TEXT,
        created_at TEXT NOT NULL,
        updated_at TEXT NOT NULL,
        FOREIGN KEY (session_id) REFERENCES attendance_sessions (id),
        FOREIGN KEY (student_id) REFERENCES users (id),
        UNIQUE (session_id, student_id)
      )
    ''');

    // Device Bindings table
    await db.execute('''
      CREATE TABLE device_bindings (
        id TEXT PRIMARY KEY,
        user_id TEXT NOT NULL,
        device_id TEXT NOT NULL,
        device_name TEXT,
        device_model TEXT,
        os_version TEXT,
        is_active INTEGER NOT NULL DEFAULT 1,
        bound_at TEXT NOT NULL,
        last_used_at TEXT,
        created_at TEXT NOT NULL,
        updated_at TEXT NOT NULL,
        FOREIGN KEY (user_id) REFERENCES users (id),
        UNIQUE (user_id, device_id)
      )
    ''');

    // Create indexes
    await db.execute('CREATE INDEX idx_users_role ON users (role)');
    await db.execute('CREATE INDEX idx_users_email ON users (email)');
    await db.execute('CREATE INDEX idx_courses_code ON courses (course_code)');
    await db.execute(
        'CREATE INDEX idx_sessions_date ON attendance_sessions (session_date)');
    await db.execute(
        'CREATE INDEX idx_sessions_status ON attendance_sessions (status)');
    await db.execute(
        'CREATE INDEX idx_attendance_verification ON attendance_records (verification_status)');
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

    // Offline attendance records table
    await db.execute('''
      CREATE TABLE offline_attendance_records (
        id TEXT PRIMARY KEY,
        session_id TEXT NOT NULL,
        student_id TEXT NOT NULL,
        marked_at TEXT NOT NULL,
        location_latitude REAL,
        location_longitude REAL,
        wifi_ssid TEXT,
        wifi_bssid TEXT,
        device_id TEXT NOT NULL,
        sync_status TEXT NOT NULL,
        error_message TEXT,
        created_at TEXT NOT NULL,
        updated_at TEXT NOT NULL,
        FOREIGN KEY (session_id) REFERENCES attendance_sessions (id) ON DELETE CASCADE,
        FOREIGN KEY (student_id) REFERENCES users (id) ON DELETE CASCADE
      )
    ''');

    // Create indexes for offline attendance records
    await db.execute(
        'CREATE INDEX idx_offline_records_status ON offline_attendance_records(sync_status)');
    await db.execute(
        'CREATE INDEX idx_offline_records_session ON offline_attendance_records(session_id)');
    await db.execute(
        'CREATE INDEX idx_offline_records_student ON offline_attendance_records(student_id)');
  }

  Future<void> _onUpgrade(Database db, int oldVersion, int newVersion) async {
    // Handle database upgrades here
    if (oldVersion < 2) {
      // Add upgrade logic for version 2
    }
  }

  // Helper methods for CRUD operations
  Future<int> insert(String table, Map<String, dynamic> row) async {
    final Database db = await database;
    return await db.insert(table, row);
  }

  Future<List<Map<String, dynamic>>> queryAllRows(String table) async {
    final Database db = await database;
    return await db.query(table);
  }

  Future<List<Map<String, dynamic>>> query(
    String table, {
    bool? distinct,
    List<String>? columns,
    String? where,
    List<Object?>? whereArgs,
    String? groupBy,
    String? having,
    String? orderBy,
    int? limit,
    int? offset,
  }) async {
    final Database db = await database;
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
  }

  Future<int> update(
    String table,
    Map<String, dynamic> row,
    String where,
    List<Object?> whereArgs,
  ) async {
    final Database db = await database;
    return await db.update(table, row, where: where, whereArgs: whereArgs);
  }

  Future<int> delete(
    String table,
    String where,
    List<Object?> whereArgs,
  ) async {
    final Database db = await database;
    return await db.delete(table, where: where, whereArgs: whereArgs);
  }

  Future<List<Map<String, dynamic>>> rawQuery(
    String sql, [
    List<Object?>? arguments,
  ]) async {
    final Database db = await database;
    return await db.rawQuery(sql, arguments);
  }

  Future<void> batch(Future<void> Function(Batch batch) operations) async {
    final Database db = await database;
    final batch = db.batch();
    await operations(batch);
    await batch.commit();
  }

  Future<void> close() async {
    if (_database != null) {
      await _database!.close();
      _database = null;
    }
  }
}
