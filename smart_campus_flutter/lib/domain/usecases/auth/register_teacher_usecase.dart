import 'package:dartz/dartz.dart';
import 'package:equatable/equatable.dart';
import '../../entities/user.dart';
import '../../repositories/auth_repository.dart';
import '../../../core/error/failures.dart';
import '../../../core/usecases/usecase.dart';

class RegisterTeacherUseCase implements UseCase<User, RegisterTeacherParams> {
  final AuthRepository repository;

  RegisterTeacherUseCase(this.repository);

  @override
  Future<Either<Failure, User>> call(RegisterTeacherParams params) async {
    return await repository.registerTeacher(
      fullName: params.fullName,
      username: params.username,
      email: params.email,
      phoneNumber: params.phoneNumber,
      password: params.password,
      highestDegree: params.highestDegree,
      experience: params.experience,
    );
  }
}

class RegisterTeacherParams extends Equatable {
  final String fullName;
  final String username;
  final String email;
  final String phoneNumber;
  final String password;
  final String highestDegree;
  final int experience;

  const RegisterTeacherParams({
    required this.fullName,
    required this.username,
    required this.email,
    required this.phoneNumber,
    required this.password,
    required this.highestDegree,
    required this.experience,
  });

  @override
  List<Object> get props => [
    fullName,
    username,
    email,
    phoneNumber,
    password,
    highestDegree,
    experience,
  ];
}
