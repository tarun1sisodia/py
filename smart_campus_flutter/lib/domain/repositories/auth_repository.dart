import 'package:dartz/dartz.dart';
import '../entities/user.dart';
import '../../core/error/failures.dart';

abstract class AuthRepository {
  Future<Either<Failure, User>> registerTeacher({
    required String fullName,
    required String username,
    required String email,
    required String phoneNumber,
    required String password,
    required String highestDegree,
    required int experience,
  });

  Future<Either<Failure, User>> registerStudent({
    required String fullName,
    required String rollNumber,
    required String course,
    required String academicYear,
    required String phoneNumber,
    required String password,
  });

  Future<Either<Failure, User>> loginTeacher({
    required String email,
    required String password,
  });

  Future<Either<Failure, User>> loginStudent({
    required String rollNumber,
    required String password,
  });

  Future<Either<Failure, String>> verifyPhone({required String phoneNumber});

  Future<Either<Failure, User>> verifyOTP({
    required String verificationId,
    required String otp,
  });

  Future<Either<Failure, void>> resetPassword({required String email});

  Future<Either<Failure, void>> logout();

  Future<Either<Failure, User?>> getCurrentUser();
}
