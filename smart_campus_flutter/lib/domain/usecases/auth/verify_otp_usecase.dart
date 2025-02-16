import 'package:dartz/dartz.dart';
import 'package:equatable/equatable.dart';
import '../../entities/user.dart';
import '../../repositories/auth_repository.dart';
import '../../../core/error/failures.dart';
import '../../../core/usecases/usecase.dart';

class VerifyOTPUseCase implements UseCase<User, VerifyOTPParams> {
  final AuthRepository repository;

  VerifyOTPUseCase(this.repository);

  @override
  Future<Either<Failure, User>> call(VerifyOTPParams params) async {
    return await repository.verifyOTP(
      verificationId: params.verificationId,
      otp: params.otp,
    );
  }
}

class VerifyOTPParams extends Equatable {
  final String verificationId;
  final String otp;

  const VerifyOTPParams({required this.verificationId, required this.otp});

  @override
  List<Object> get props => [verificationId, otp];
}
