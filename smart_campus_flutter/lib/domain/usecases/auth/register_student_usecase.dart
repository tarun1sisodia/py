import 'package:dartz/dartz.dart';
import 'package:equatable/equatable.dart';
import '../../entities/user.dart';
import '../../repositories/auth_repository.dart';
import '../../../core/error/failures.dart';
import '../../../core/usecases/usecase.dart';

class RegisterStudentUseCase implements UseCase<User, RegisterStudentParams> {
  final AuthRepository repository;

  RegisterStudentUseCase(this.repository);

  @override
  Future<Either<Failure, User>> call(RegisterStudentParams params) async {
    return await repository.registerStudent(
      fullName: params.fullName,
      rollNumber: params.rollNumber,
      course: params.course,
      academicYear: params.academicYear,
      phoneNumber: params.phoneNumber,
      password: params.password,
    );
  }
}

class RegisterStudentParams extends Equatable {
  final String fullName;
  final String rollNumber;
  final String course;
  final String academicYear;
  final String phoneNumber;
  final String password;

  const RegisterStudentParams({
    required this.fullName,
    required this.rollNumber,
    required this.course,
    required this.academicYear,
    required this.phoneNumber,
    required this.password,
  });

  @override
  List<Object> get props => [
    fullName,
    rollNumber,
    course,
    academicYear,
    phoneNumber,
    password,
  ];
}
