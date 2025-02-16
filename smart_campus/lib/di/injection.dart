import 'package:get_it/get_it.dart';
import 'package:injectable/injectable.dart';
import 'package:sqflite/sqflite.dart';
import 'package:path/path.dart';
import 'package:smart_campus/di/injection.config.dart';
import 'package:smart_campus/presentation/bloc/auth/auth_bloc.dart';
import 'package:smart_campus/presentation/bloc/onboarding/onboarding_bloc.dart';
import 'package:smart_campus/presentation/bloc/session/session_bloc.dart';
import 'package:smart_campus/presentation/bloc/student/student_bloc.dart';
import 'package:smart_campus/presentation/bloc/teacher/teacher_bloc.dart';
import 'package:smart_campus/domain/repositories/session_repository.dart';
import 'package:smart_campus/domain/repositories/auth_repository.dart';
import 'package:smart_campus/core/services/sync_service.dart';
import 'package:smart_campus/core/services/connectivity_service.dart';
import 'package:smart_campus/core/services/logger_service.dart';
import 'package:smart_campus/core/services/error_reporting_service.dart';
import 'package:connectivity_plus/connectivity_plus.dart';
import 'package:smart_campus/data/repositories/session_repository_impl.dart';

final getIt = GetIt.instance;

@InjectableInit()
Future<void> configureDependencies() async {
  // Initialize database first
  final databasePath = await getDatabasesPath();
  final path = join(databasePath, 'smart_campus.db');
  final database = await openDatabase(
    path,
    version: 1,
    onCreate: (db, version) async {
      await db.execute('''
        CREATE TABLE IF NOT EXISTS attendance_records (
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
          updated_at TEXT NOT NULL
        )
      ''');
    },
  );

  // Register database instance
  getIt.registerSingleton<Database>(database);

  // Initialize other dependencies
  getIt.init();

  // Register core dependencies
  getIt.registerLazySingleton(() => Connectivity());
  getIt.registerLazySingleton(() => LoggerService());
  getIt.registerLazySingleton(() => ErrorReportingService(getIt()));

  // Register services
  getIt.registerLazySingleton(() => ConnectivityService(getIt()));
  getIt.registerLazySingleton(() => SyncService(
        localDatasource: getIt(),
        sessionRepository: getIt(),
        connectivityService: getIt(),
        logger: getIt(),
      ));

  // Register repositories
  getIt.registerLazySingleton<SessionRepository>(() => SessionRepositoryImpl(
        getIt(), // DatabaseService
        getIt(), // NetworkService
        getIt(), // LocationService
        getIt(), // WifiService
        getIt(), // LoggerService
        getIt(), // SessionLocalDataSource
        getIt(), // DeviceService
        getIt(), // DatabaseHelper
      ));

  // Register blocs
  getIt
      .registerFactory(() => AuthBloc(authRepository: getIt<AuthRepository>()));
  getIt.registerFactory(() => OnboardingBloc(getIt(), getIt()));
  getIt.registerFactory(() => SessionBloc(
        sessionRepository: getIt<SessionRepository>(),
        syncService: getIt<SyncService>(),
        connectivity: getIt<Connectivity>(),
      ));
  getIt.registerFactory(() => StudentBloc(
        sessionRepository: getIt<SessionRepository>(),
        syncService: getIt<SyncService>(),
        connectivityService: getIt<ConnectivityService>(),
      ));
  getIt.registerFactory(() => TeacherBloc(
        sessionRepository: getIt<SessionRepository>(),
      ));
}
