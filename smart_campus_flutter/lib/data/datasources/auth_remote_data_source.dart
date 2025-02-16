import 'dart:convert';
import 'package:http/http.dart' as http;
import '../models/user_model.dart';

abstract class AuthRemoteDataSource {
  Future<UserModel> registerTeacher({
    required String id,
    required String fullName,
    required String username,
    required String email,
    required String phoneNumber,
    required String highestDegree,
    required int experience,
  });

  Future<UserModel> registerStudent({
    required String id,
    required String fullName,
    required String rollNumber,
    required String course,
    required String academicYear,
    required String phoneNumber,
  });

  Future<UserModel> getTeacherProfile({required String id});
  Future<UserModel> getStudentProfile({required String id});
  Future<UserModel> getCurrentUser();
  Future<void> logout();
}

class AuthRemoteDataSourceImpl implements AuthRemoteDataSource {
  final http.Client client;
  final String baseUrl;

  AuthRemoteDataSourceImpl({required this.client, required this.baseUrl});

  @override
  Future<UserModel> registerTeacher({
    required String id,
    required String fullName,
    required String username,
    required String email,
    required String phoneNumber,
    required String highestDegree,
    required int experience,
  }) async {
    final response = await client.post(
      Uri.parse('$baseUrl/auth/register/teacher'),
      body: json.encode({
        'id': id,
        'full_name': fullName,
        'username': username,
        'email': email,
        'phone': phoneNumber,
        'highest_degree': highestDegree,
        'experience': experience,
        'role': 'teacher',
      }),
      headers: {'Content-Type': 'application/json'},
    );

    if (response.statusCode == 201) {
      return UserModel.fromJson(json.decode(response.body));
    } else {
      throw Exception('Failed to register teacher');
    }
  }

  @override
  Future<UserModel> registerStudent({
    required String id,
    required String fullName,
    required String rollNumber,
    required String course,
    required String academicYear,
    required String phoneNumber,
  }) async {
    final response = await client.post(
      Uri.parse('$baseUrl/auth/register/student'),
      body: json.encode({
        'id': id,
        'full_name': fullName,
        'roll_number': rollNumber,
        'course': course,
        'academic_year': academicYear,
        'phone': phoneNumber,
        'role': 'student',
      }),
      headers: {'Content-Type': 'application/json'},
    );

    if (response.statusCode == 201) {
      return UserModel.fromJson(json.decode(response.body));
    } else {
      throw Exception('Failed to register student');
    }
  }

  @override
  Future<UserModel> getTeacherProfile({required String id}) async {
    final response = await client.get(
      Uri.parse('$baseUrl/teachers/$id'),
      headers: {'Content-Type': 'application/json'},
    );

    if (response.statusCode == 200) {
      return UserModel.fromJson(json.decode(response.body));
    } else {
      throw Exception('Failed to get teacher profile');
    }
  }

  @override
  Future<UserModel> getStudentProfile({required String id}) async {
    final response = await client.get(
      Uri.parse('$baseUrl/students/$id'),
      headers: {'Content-Type': 'application/json'},
    );

    if (response.statusCode == 200) {
      return UserModel.fromJson(json.decode(response.body));
    } else {
      throw Exception('Failed to get student profile');
    }
  }

  @override
  Future<UserModel> getCurrentUser() async {
    final response = await client.get(
      Uri.parse('$baseUrl/auth/me'),
      headers: {'Content-Type': 'application/json'},
    );

    if (response.statusCode == 200) {
      return UserModel.fromJson(json.decode(response.body));
    } else {
      throw Exception('Failed to get current user');
    }
  }

  @override
  Future<void> logout() async {
    final response = await client.post(
      Uri.parse('$baseUrl/auth/logout'),
      headers: {'Content-Type': 'application/json'},
    );

    if (response.statusCode != 200) {
      throw Exception('Failed to logout');
    }
  }
}
