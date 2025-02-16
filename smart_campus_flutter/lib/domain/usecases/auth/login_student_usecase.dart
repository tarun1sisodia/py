import 'package:dartz/dartz.dart';
import 'package:equatable/equatable.dart';
import '../../entities/user.dart';
import '../../repositories/auth_repository.dart';
import '../../../core/error/failures.dart';
import '../../../core/usecases/usecase.dart';

class LoginStudentUseCase implements UseCase<User, LoginStudentParams> {
  final AuthRepository repository;

  LoginStudentUseCase(this.repository);

  @override
  Future<Either<Failure, User>> call(LoginStudentParams params) async {
    return await repository.loginStudent(
      rollNumber: params.rollNumber,
      password: params.password,
    );
  }
}

class LoginStudentParams extends Equatable {
  final String rollNumber;
  final String password;

  const LoginStudentParams({required this.rollNumber, required this.password});

  @override
  List<Object> get props => [rollNumber, password];
}
