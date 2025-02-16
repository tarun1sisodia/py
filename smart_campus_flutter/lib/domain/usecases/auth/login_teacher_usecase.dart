import 'package:dartz/dartz.dart';
import 'package:equatable/equatable.dart';
import '../../entities/user.dart';
import '../../repositories/auth_repository.dart';
import '../../../core/error/failures.dart';
import '../../../core/usecases/usecase.dart';

class LoginTeacherUseCase implements UseCase<User, LoginTeacherParams> {
  final AuthRepository repository;

  LoginTeacherUseCase(this.repository);

  @override
  Future<Either<Failure, User>> call(LoginTeacherParams params) async {
    return await repository.loginTeacher(
      email: params.email,
      password: params.password,
    );
  }
}

class LoginTeacherParams extends Equatable {
  final String email;
  final String password;

  const LoginTeacherParams({required this.email, required this.password});

  @override
  List<Object> get props => [email, password];
}
