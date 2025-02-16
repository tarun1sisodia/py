import 'package:dio/dio.dart';

class AuthenticationService {
  final Dio _dio;
  // Updated backend URL
  static const String baseUrl = 'http://localhost:8080/api/v1/auth';

  AuthenticationService()
    : _dio = Dio(
        BaseOptions(
          baseUrl: baseUrl,
          connectTimeout: Duration(milliseconds: 5000),
          receiveTimeout: Duration(milliseconds: 3000),
        ),
      );

  Future<Response> teacherRegister(Map<String, dynamic> data) async {
    try {
      // Updated endpoint
      Response response = await _dio.post('/register/teacher', data: data);
      return response;
    } catch (e) {
      throw Exception('Teacher registration failed: $e');
    }
  }

  Future<Response> teacherLogin(String email, String password) async {
    try {
      // Updated endpoint
      Response response = await _dio.post(
        '/login/teacher',
        data: {'email': email, 'password': password},
      );
      return response;
    } catch (e) {
      throw Exception('Teacher login failed: $e');
    }
  }

  Future<Response> studentRegister(Map<String, dynamic> data) async {
    try {
      // Updated endpoint
      Response response = await _dio.post('/register/student', data: data);
      return response;
    } catch (e) {
      throw Exception('Student registration failed: $e');
    }
  }

  Future<Response> studentLogin(String rollNumber, String password) async {
    try {
      // Updated endpoint
      Response response = await _dio.post(
        '/login/student',
        data: {'roll_number': rollNumber, 'password': password},
      );
      return response;
    } catch (e) {
      throw Exception('Student login failed: $e');
    }
  }

  Future<Response> resetPassword(Map<String, dynamic> data) async {
    try {
      Response response = await _dio.post('/reset-password', data: data);
      return response;
    } catch (e) {
      throw Exception('Password reset failed: $e');
    }
  }

  Future<Response> verifyOTP(String userId, String otp) async {
    try {
      Response response = await _dio.post(
        '/verify-otp',
        data: {'user_id': userId, 'otp': otp},
      );
      return response;
    } catch (e) {
      throw Exception('OTP verification failed: $e');
    }
  }
}
