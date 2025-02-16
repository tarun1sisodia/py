// dart format width=80
// GENERATED CODE - DO NOT MODIFY BY HAND

// **************************************************************************
// InjectableConfigGenerator
// **************************************************************************

// ignore_for_file: type=lint
// coverage:ignore-file

// ignore_for_file: no_leading_underscores_for_library_prefixes
import 'package:get_it/get_it.dart' as _i174;
import 'package:injectable/injectable.dart' as _i526;
import 'package:shared_preferences/shared_preferences.dart' as _i460;
import 'package:sqflite/sqflite.dart' as _i1;

import '../core/services/attendance_verification_service.dart' as _i151;
import '../core/services/database_service.dart' as _i162;
import '../core/services/device_binding_service.dart' as _i548;
import '../core/services/location_service.dart' as _i848;
import '../core/services/logger_service.dart' as _i470;
import '../core/services/network_service.dart' as _i759;
import '../core/services/storage_service.dart' as _i736;
import '../core/services/sync_service.dart' as _i827;
import '../core/services/wifi_service.dart' as _i123;
import '../data/datasources/local/database_helper.dart' as _i808;
import '../data/datasources/local/session_local_datasource.dart' as _i371;
import '../data/repositories/course_repository_impl.dart' as _i1039;
import '../data/repositories/session_repository_impl.dart' as _i165;
import '../domain/repositories/course_repository.dart' as _i720;
import '../domain/repositories/session_repository.dart' as _i1010;
import '../presentation/bloc/attendance_verification/attendance_verification_bloc.dart'
    as _i444;
import '../presentation/bloc/device_binding/device_binding_bloc.dart' as _i609;
import '../services/device_service.dart' as _i738;
import 'package:smart_campus/data/datasources/local/attendance_local_datasource.dart'
    as _i371;
import 'package:smart_campus/core/services/connectivity_service.dart' as _i895;

extension GetItInjectableX on _i174.GetIt {
// initializes the registration of main-scope dependencies inside of GetIt
  _i174.GetIt init({
    String? environment,
    _i526.EnvironmentFilter? environmentFilter,
  }) {
    final gh = _i526.GetItHelper(
      this,
      environment,
      environmentFilter,
    );
    gh.lazySingleton<_i808.DatabaseHelper>(
      () => _i808.DatabaseHelper(gh<_i1.Database>()),
    );
    gh.lazySingleton<_i470.LoggerService>(() => _i470.LoggerService());
    gh.lazySingleton<_i548.DeviceBindingService>(
        () => _i548.DeviceBindingService(
              gh<_i460.SharedPreferences>(),
              gh<_i470.LoggerService>(),
            ));
    gh.lazySingleton<_i736.StorageService>(
        () => _i736.StorageService(gh<_i460.SharedPreferences>()));
    gh.factory<_i609.DeviceBindingBloc>(() => _i609.DeviceBindingBloc(
          gh<_i548.DeviceBindingService>(),
          gh<_i470.LoggerService>(),
        ));
    gh.lazySingleton<_i738.DeviceService>(
        () => _i738.DeviceService(gh<_i808.DatabaseHelper>()));
    gh.lazySingleton<_i162.DatabaseService>(
        () => _i162.DatabaseService(gh<_i470.LoggerService>()));
    gh.lazySingleton<_i123.WifiService>(
        () => _i123.WifiService(gh<_i470.LoggerService>()));
    gh.lazySingleton<_i848.LocationService>(
        () => _i848.LocationService(gh<_i470.LoggerService>()));
    gh.lazySingleton<_i371.SessionLocalDataSource>(
        () => _i371.SessionLocalDataSourceImpl(gh<_i808.DatabaseHelper>()));
    gh.lazySingleton<_i759.NetworkService>(
        () => _i759.NetworkService(gh<_i736.StorageService>()));
    gh.lazySingleton<_i151.AttendanceVerificationService>(
        () => _i151.AttendanceVerificationService(
              gh<_i848.LocationService>(),
              gh<_i123.WifiService>(),
              gh<_i548.DeviceBindingService>(),
              gh<_i470.LoggerService>(),
            ));
    gh.lazySingleton<_i1010.SessionRepository>(
        () => _i165.SessionRepositoryImpl(
              gh<_i162.DatabaseService>(),
              gh<_i759.NetworkService>(),
              gh<_i848.LocationService>(),
              gh<_i123.WifiService>(),
              gh<_i470.LoggerService>(),
              gh<_i371.SessionLocalDataSource>(),
              gh<_i738.DeviceService>(),
              gh<_i808.DatabaseHelper>(),
            ));
    gh.lazySingleton<_i827.SyncService>(() => _i827.SyncService(
          localDatasource: gh<_i371.AttendanceLocalDatasource>(),
          sessionRepository: gh<_i1010.SessionRepository>(),
          connectivityService: gh<_i895.ConnectivityService>(),
          logger: gh<_i470.LoggerService>(),
        ));
    gh.lazySingleton<_i720.CourseRepository>(() => _i1039.CourseRepositoryImpl(
          gh<_i162.DatabaseService>(),
          gh<_i759.NetworkService>(),
          gh<_i470.LoggerService>(),
        ));
    gh.factory<_i444.AttendanceVerificationBloc>(
        () => _i444.AttendanceVerificationBloc(
              gh<_i151.AttendanceVerificationService>(),
              gh<_i470.LoggerService>(),
            ));
    return this;
  }
}
