import 'package:injectable/injectable.dart';
import 'package:path/path.dart';
import 'package:sqflite/sqflite.dart';
import 'package:smart_campus/data/datasources/local/attendance_local_datasource.dart';

@lazySingleton
class DatabaseService {
  static const String _databaseName = 'smart_campus.db';
  static const int _databaseVersion = 1;

  static Future<Database> init() async {
    final databasePath = await getDatabasesPath();
    final path = join(databasePath, _databaseName);

    return openDatabase(
      path,
      version: _databaseVersion,
      onCreate: _onCreate,
      onUpgrade: _onUpgrade,
    );
  }

  static Future<void> _onCreate(Database db, int version) async {
    // Create tables
    await AttendanceLocalDatasourceImpl.createTable(db);
  }

  static Future<void> _onUpgrade(
    Database db,
    int oldVersion,
    int newVersion,
  ) async {
    // Handle database upgrades
    if (oldVersion < newVersion) {
      // Add upgrade logic here when needed
    }
  }

  static Future<void> deleteDatabase() async {
    final databasePath = await getDatabasesPath();
    final path = join(databasePath, _databaseName);
    await databaseFactory.deleteDatabase(path);
  }
}
