import 'dart:async';
import 'dart:io';
import 'package:dio/dio.dart';
import 'package:injectable/injectable.dart';
import 'package:smart_campus/core/errors/app_error.dart';
import 'package:smart_campus/core/services/auth_service.dart';
import 'package:smart_campus/core/services/logger_service.dart';
import 'package:smart_campus/core/services/connectivity_service.dart';

@lazySingleton
class ApiClient {
  final Dio _dio;
  final LoggerService _logger;
  final ConnectivityService _connectivityService;
  final AuthService _authService;
  static const _defaultTimeout = Duration(seconds: 30);
  static const _maxRetries = 3;

  ApiClient(
    this._logger,
    this._connectivityService,
    this._authService,
  ) : _dio = Dio() {
    _dio.options
      ..baseUrl = 'https://api.smartcampus.com/v1' // Replace with your API URL
      ..connectTimeout = _defaultTimeout
      ..receiveTimeout = _defaultTimeout
      ..sendTimeout = _defaultTimeout
      ..validateStatus = (status) => status != null && status < 500;

    // Add interceptors
    _dio.interceptors.addAll([
      _createAuthInterceptor(),
      _createRetryInterceptor(),
      _createLoggerInterceptor(),
    ]);
  }

  Interceptor _createAuthInterceptor() {
    return InterceptorsWrapper(
      onRequest: (options, handler) async {
        final token = await _authService.getToken();
        if (token != null) {
          options.headers['Authorization'] = 'Bearer $token';
        }
        return handler.next(options);
      },
      onError: (error, handler) async {
        if (error.response?.statusCode == 401) {
          try {
            final newToken = await _authService.refreshToken();
            if (newToken != null) {
              error.requestOptions.headers['Authorization'] =
                  'Bearer $newToken';
              return handler.resolve(await _dio.fetch(error.requestOptions));
            }
          } catch (e) {
            _logger.error('Token refresh failed', e);
          }
        }
        return handler.next(error);
      },
    );
  }

  Interceptor _createRetryInterceptor() {
    return InterceptorsWrapper(
      onError: (error, handler) async {
        if (!_shouldRetry(error)) {
          return handler.next(error);
        }

        var retryCount = 0;
        while (retryCount < _maxRetries) {
          try {
            if (!await _connectivityService.isConnected()) {
              await _connectivityService.waitForConnection();
            }

            final response = await _dio.fetch(error.requestOptions);
            return handler.resolve(response);
          } catch (e) {
            retryCount++;
            if (retryCount == _maxRetries) {
              return handler.next(error);
            }
            await Future.delayed(Duration(seconds: retryCount * 2));
          }
        }
      },
    );
  }

  Interceptor _createLoggerInterceptor() {
    return InterceptorsWrapper(
      onRequest: (options, handler) {
        _logger.info(
          'API Request: ${options.method} ${options.uri}',
          {'headers': options.headers, 'data': options.data},
        );
        return handler.next(options);
      },
      onResponse: (response, handler) {
        _logger.info(
          'API Response: ${response.statusCode} ${response.requestOptions.uri}',
          {'headers': response.headers.map, 'data': response.data},
        );
        return handler.next(response);
      },
      onError: (error, handler) {
        _logger.error(
          'API Error: ${error.response?.statusCode} ${error.requestOptions.uri}',
          error,
        );
        return handler.next(error);
      },
    );
  }

  bool _shouldRetry(DioException error) {
    return error.type == DioExceptionType.connectionTimeout ||
        error.type == DioExceptionType.sendTimeout ||
        error.type == DioExceptionType.receiveTimeout ||
        (error.error is SocketException) ||
        (error.response?.statusCode == 503);
  }

  Future<T> get<T>(
    String path, {
    Map<String, dynamic>? queryParameters,
    Options? options,
    CancelToken? cancelToken,
  }) async {
    try {
      final response = await _dio.get<T>(
        path,
        queryParameters: queryParameters,
        options: options,
        cancelToken: cancelToken,
      );
      return _handleResponse<T>(response);
    } on DioException catch (e) {
      throw _handleError(e);
    }
  }

  Future<T> post<T>(
    String path, {
    dynamic data,
    Map<String, dynamic>? queryParameters,
    Options? options,
    CancelToken? cancelToken,
  }) async {
    try {
      final response = await _dio.post<T>(
        path,
        data: data,
        queryParameters: queryParameters,
        options: options,
        cancelToken: cancelToken,
      );
      return _handleResponse<T>(response);
    } on DioException catch (e) {
      throw _handleError(e);
    }
  }

  Future<T> put<T>(
    String path, {
    dynamic data,
    Map<String, dynamic>? queryParameters,
    Options? options,
    CancelToken? cancelToken,
  }) async {
    try {
      final response = await _dio.put<T>(
        path,
        data: data,
        queryParameters: queryParameters,
        options: options,
        cancelToken: cancelToken,
      );
      return _handleResponse<T>(response);
    } on DioException catch (e) {
      throw _handleError(e);
    }
  }

  Future<T> delete<T>(
    String path, {
    dynamic data,
    Map<String, dynamic>? queryParameters,
    Options? options,
    CancelToken? cancelToken,
  }) async {
    try {
      final response = await _dio.delete<T>(
        path,
        data: data,
        queryParameters: queryParameters,
        options: options,
        cancelToken: cancelToken,
      );
      return _handleResponse<T>(response);
    } on DioException catch (e) {
      throw _handleError(e);
    }
  }

  T _handleResponse<T>(Response<T> response) {
    if (response.statusCode == 401) {
      throw AppError.unauthorized();
    }

    if (response.statusCode == 403) {
      throw AppError('Access forbidden');
    }

    if (response.statusCode == 404) {
      throw AppError('Resource not found');
    }

    if (response.statusCode! >= 400) {
      final data = response.data as Map<String, dynamic>?;
      throw AppError(
        data?['message'] as String? ?? 'Unknown error occurred',
        code: data?['code'] as String?,
      );
    }

    return response.data as T;
  }

  AppError _handleError(DioException error) {
    if (error.type == DioExceptionType.connectionTimeout ||
        error.type == DioExceptionType.sendTimeout ||
        error.type == DioExceptionType.receiveTimeout) {
      return AppError.network('Connection timeout');
    }

    if (error.type == DioExceptionType.connectionError) {
      return AppError.network('No internet connection');
    }

    if (error.response != null) {
      final statusCode = error.response!.statusCode;
      final data = error.response!.data as Map<String, dynamic>?;

      if (statusCode == 401) {
        return AppError.unauthorized();
      }

      if (statusCode == 403) {
        return AppError('Access forbidden');
      }

      if (statusCode == 404) {
        return AppError('Resource not found');
      }

      if (statusCode! >= 500) {
        return AppError.server('Server error occurred');
      }

      return AppError(
        data?['message'] as String? ?? 'Unknown error occurred',
        code: data?['code'] as String?,
      );
    }

    return AppError('Unknown error occurred');
  }
}
