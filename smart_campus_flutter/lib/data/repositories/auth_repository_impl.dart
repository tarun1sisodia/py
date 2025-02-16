import 'package:dartz/dartz.dart';
import 'package:firebase_auth/firebase_auth.dart' as firebase_auth;
import 'package:smart_campus_flutter/core/error/failures.dart';
import 'package:smart_campus_flutter/data/datasources/auth_remote_data_source.dart';
import 'package:smart_campus_flutter/domain/entities/user.dart';
import 'package:smart_campus_flutter/domain/repositories/auth_repository.dart';

class AuthRepositoryImpl implements AuthRepository {
  final AuthRemoteDataSource remoteDataSource;
  final firebase_auth.FirebaseAuth firebaseAuth;

  AuthRepositoryImpl({
    required this.remoteDataSource,
    required this.firebaseAuth,
  });

  @override
  Future<Either<Failure, User>> registerTeacher({
    required String fullName,
    required String username,
    required String email,
    required String phoneNumber,
    required String password,
    required String highestDegree,
    required int experience,
  }) async {
    try {
      // Create user with Firebase Auth
      final firebaseUser = await firebaseAuth.createUserWithEmailAndPassword(
        email: email,
        password: password,
      );

      if (firebaseUser.user == null) {
        return const Left(AuthFailure('Failed to create user'));
      }

      // Register user with backend
      final user = await remoteDataSource.registerTeacher(
        id: firebaseUser.user!.uid,
        fullName: fullName,
        username: username,
        email: email,
        phoneNumber: phoneNumber,
        highestDegree: highestDegree,
        experience: experience,
      );

      return Right(user);
    } on firebase_auth.FirebaseAuthException catch (e) {
      return Left(AuthFailure(e.message ?? 'Authentication failed'));
    } catch (e) {
      return Left(UnexpectedFailure(e.toString()));
    }
  }

  @override
  Future<Either<Failure, User>> registerStudent({
    required String fullName,
    required String rollNumber,
    required String course,
    required String academicYear,
    required String phoneNumber,
    required String password,
  }) async {
    try {
      // Create user with Firebase Auth using roll number as email
      final firebaseUser = await firebaseAuth.createUserWithEmailAndPassword(
        email:
            '$rollNumber@student.smartcampus.edu', // Using roll number as email
        password: password,
      );

      if (firebaseUser.user == null) {
        return const Left(AuthFailure('Failed to create user'));
      }

      // Register user with backend
      final user = await remoteDataSource.registerStudent(
        id: firebaseUser.user!.uid,
        fullName: fullName,
        rollNumber: rollNumber,
        course: course,
        academicYear: academicYear,
        phoneNumber: phoneNumber,
      );

      return Right(user);
    } on firebase_auth.FirebaseAuthException catch (e) {
      return Left(AuthFailure(e.message ?? 'Authentication failed'));
    } catch (e) {
      return Left(UnexpectedFailure(e.toString()));
    }
  }

  @override
  Future<Either<Failure, User>> loginTeacher({
    required String email,
    required String password,
  }) async {
    try {
      // Sign in with Firebase
      final firebaseUser = await firebaseAuth.signInWithEmailAndPassword(
        email: email,
        password: password,
      );

      if (firebaseUser.user == null) {
        return const Left(AuthFailure('Failed to login'));
      }

      // Get user details from backend
      final user = await remoteDataSource.getTeacherProfile(
        id: firebaseUser.user!.uid,
      );

      return Right(user);
    } on firebase_auth.FirebaseAuthException catch (e) {
      return Left(AuthFailure(e.message ?? 'Authentication failed'));
    } catch (e) {
      return Left(UnexpectedFailure(e.toString()));
    }
  }

  @override
  Future<Either<Failure, User>> loginStudent({
    required String rollNumber,
    required String password,
  }) async {
    try {
      // Sign in with Firebase using roll number as email
      final firebaseUser = await firebaseAuth.signInWithEmailAndPassword(
        email: '$rollNumber@student.smartcampus.edu',
        password: password,
      );

      if (firebaseUser.user == null) {
        return const Left(AuthFailure('Failed to login'));
      }

      // Get user details from backend
      final user = await remoteDataSource.getStudentProfile(
        id: firebaseUser.user!.uid,
      );

      return Right(user);
    } on firebase_auth.FirebaseAuthException catch (e) {
      return Left(AuthFailure(e.message ?? 'Authentication failed'));
    } catch (e) {
      return Left(UnexpectedFailure(e.toString()));
    }
  }

  @override
  Future<Either<Failure, String>> verifyPhone({
    required String phoneNumber,
  }) async {
    try {
      String verificationId = '';

      await firebaseAuth.verifyPhoneNumber(
        phoneNumber: phoneNumber,
        verificationCompleted: (
          firebase_auth.PhoneAuthCredential credential,
        ) async {
          // Auto-verification handled by Firebase
        },
        verificationFailed: (firebase_auth.FirebaseAuthException e) {
          throw e;
        },
        codeSent: (String vId, int? resendToken) {
          verificationId = vId;
        },
        codeAutoRetrievalTimeout: (String vId) {
          verificationId = vId;
        },
      );

      return Right(verificationId);
    } on firebase_auth.FirebaseAuthException catch (e) {
      return Left(AuthFailure(e.message ?? 'Phone verification failed'));
    } catch (e) {
      return Left(UnexpectedFailure(e.toString()));
    }
  }

  @override
  Future<Either<Failure, User>> verifyOTP({
    required String verificationId,
    required String otp,
  }) async {
    try {
      // Create credential
      final credential = firebase_auth.PhoneAuthProvider.credential(
        verificationId: verificationId,
        smsCode: otp,
      );

      // Sign in with credential
      final firebaseUser = await firebaseAuth.signInWithCredential(credential);

      if (firebaseUser.user == null) {
        return const Left(AuthFailure('OTP verification failed'));
      }

      // Get current user from backend
      final user = await remoteDataSource.getCurrentUser();
      return Right(user);
    } on firebase_auth.FirebaseAuthException catch (e) {
      return Left(AuthFailure(e.message ?? 'OTP verification failed'));
    } catch (e) {
      return Left(UnexpectedFailure(e.toString()));
    }
  }

  @override
  Future<Either<Failure, void>> resetPassword({required String email}) async {
    try {
      await firebaseAuth.sendPasswordResetEmail(email: email);
      return const Right(null);
    } on firebase_auth.FirebaseAuthException catch (e) {
      return Left(AuthFailure(e.message ?? 'Password reset failed'));
    } catch (e) {
      return Left(UnexpectedFailure(e.toString()));
    }
  }

  @override
  Future<Either<Failure, void>> logout() async {
    try {
      await Future.wait([firebaseAuth.signOut(), remoteDataSource.logout()]);
      return const Right(null);
    } catch (e) {
      return Left(UnexpectedFailure(e.toString()));
    }
  }

  @override
  Future<Either<Failure, User?>> getCurrentUser() async {
    try {
      final firebaseUser = firebaseAuth.currentUser;
      if (firebaseUser == null) {
        return const Right(null);
      }

      final user = await remoteDataSource.getCurrentUser();
      return Right(user);
    } catch (e) {
      return Left(UnexpectedFailure(e.toString()));
    }
  }
}
