import 'package:dartz/dartz.dart';
import 'package:equatable/equatable.dart';
import '../../repositories/auth_repository.dart';
import '../../../core/error/failures.dart';
import '../../../core/usecases/usecase.dart';

class VerifyPhoneUseCase implements UseCase<String, VerifyPhoneParams> {
  final AuthRepository repository;

  VerifyPhoneUseCase(this.repository);

  @override
  Future<Either<Failure, String>> call(VerifyPhoneParams params) async {
    return await repository.verifyPhone(phoneNumber: params.phoneNumber);
  }
}

class VerifyPhoneParams extends Equatable {
  final String phoneNumber;

  const VerifyPhoneParams({required this.phoneNumber});

  @override
  List<Object> get props => [phoneNumber];
}
