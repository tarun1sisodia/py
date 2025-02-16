import 'package:flutter_bloc/flutter_bloc.dart';
import '../../../core/usecases/usecase.dart';
import '../../../domain/usecases/auth/get_current_user_usecase.dart';
import '../../../domain/usecases/auth/login_student_usecase.dart';
import '../../../domain/usecases/auth/login_teacher_usecase.dart';
import '../../../domain/usecases/auth/logout_usecase.dart';
import '../../../domain/usecases/auth/register_student_usecase.dart';
import '../../../domain/usecases/auth/register_teacher_usecase.dart';
import '../../../domain/usecases/auth/reset_password_usecase.dart';
import '../../../domain/usecases/auth/verify_otp_usecase.dart';
import '../../../domain/usecases/auth/verify_phone_usecase.dart';
import 'auth_event.dart';
import 'auth_state.dart';

class AuthBloc extends Bloc<AuthEvent, AuthState> {
  final RegisterTeacherUseCase registerTeacher;
  final RegisterStudentUseCase registerStudent;
  final LoginTeacherUseCase loginTeacher;
  final LoginStudentUseCase loginStudent;
  final VerifyPhoneUseCase verifyPhone;
  final VerifyOTPUseCase verifyOTP;
  final ResetPasswordUseCase resetPassword;
  final LogoutUseCase logout;
  final GetCurrentUserUseCase getCurrentUser;

  AuthBloc({
    required this.registerTeacher,
    required this.registerStudent,
    required this.loginTeacher,
    required this.loginStudent,
    required this.verifyPhone,
    required this.verifyOTP,
    required this.resetPassword,
    required this.logout,
    required this.getCurrentUser,
  }) : super(AuthInitial()) {
    on<AuthRegisterTeacher>(_onRegisterTeacher);
    on<AuthRegisterStudent>(_onRegisterStudent);
    on<AuthLoginTeacher>(_onLoginTeacher);
    on<AuthLoginStudent>(_onLoginStudent);
    on<AuthVerifyPhone>(_onVerifyPhone);
    on<AuthSendOTP>(_onSendOTP);
    on<AuthResetPassword>(_onResetPassword);
    on<AuthLogout>(_onLogout);
    on<AuthCheckStatus>(_onCheckStatus);
  }

  Future<void> _onRegisterTeacher(
    AuthRegisterTeacher event,
    Emitter<AuthState> emit,
  ) async {
    emit(AuthLoading());

    final result = await registerTeacher(
      RegisterTeacherParams(
        fullName: event.fullName,
        username: event.username,
        email: event.email,
        phoneNumber: event.phone,
        password: event.password,
        highestDegree: event.highestDegree,
        experience: int.parse(event.experience),
      ),
    );

    result.fold(
      (failure) => emit(AuthError(failure.message)),
      (user) => emit(AuthAuthenticated(user)),
    );
  }

  Future<void> _onRegisterStudent(
    AuthRegisterStudent event,
    Emitter<AuthState> emit,
  ) async {
    emit(AuthLoading());

    final result = await registerStudent(
      RegisterStudentParams(
        fullName: event.fullName,
        rollNumber: event.rollNumber,
        course: event.course,
        academicYear: event.academicYear,
        phoneNumber: event.phone,
        password: event.password,
      ),
    );

    result.fold(
      (failure) => emit(AuthError(failure.message)),
      (user) => emit(AuthAuthenticated(user)),
    );
  }

  Future<void> _onLoginTeacher(
    AuthLoginTeacher event,
    Emitter<AuthState> emit,
  ) async {
    emit(AuthLoading());

    final result = await loginTeacher(
      LoginTeacherParams(email: event.email, password: event.password),
    );

    result.fold(
      (failure) => emit(AuthError(failure.message)),
      (user) => emit(AuthAuthenticated(user)),
    );
  }

  Future<void> _onLoginStudent(
    AuthLoginStudent event,
    Emitter<AuthState> emit,
  ) async {
    emit(AuthLoading());

    final result = await loginStudent(
      LoginStudentParams(
        rollNumber: event.rollNumber,
        password: event.password,
      ),
    );

    result.fold(
      (failure) => emit(AuthError(failure.message)),
      (user) => emit(AuthAuthenticated(user)),
    );
  }

  Future<void> _onSendOTP(AuthSendOTP event, Emitter<AuthState> emit) async {
    emit(AuthLoading());

    final result = await verifyPhone(
      VerifyPhoneParams(phoneNumber: event.phoneNumber),
    );

    result.fold(
      (failure) => emit(AuthError(failure.message)),
      (verificationId) => emit(AuthOTPSent(verificationId)),
    );
  }

  Future<void> _onVerifyPhone(
    AuthVerifyPhone event,
    Emitter<AuthState> emit,
  ) async {
    emit(AuthLoading());

    final result = await verifyOTP(
      VerifyOTPParams(verificationId: event.verificationId, otp: event.otp),
    );

    result.fold(
      (failure) => emit(AuthError(failure.message)),
      (user) => emit(AuthAuthenticated(user)),
    );
  }

  Future<void> _onResetPassword(
    AuthResetPassword event,
    Emitter<AuthState> emit,
  ) async {
    emit(AuthLoading());

    final result = await resetPassword(ResetPasswordParams(email: event.email));

    result.fold(
      (failure) => emit(AuthError(failure.message)),
      (_) => emit(const AuthPasswordResetSent()),
    );
  }

  Future<void> _onLogout(AuthLogout event, Emitter<AuthState> emit) async {
    emit(AuthLoading());

    final result = await logout(const NoParams());

    result.fold(
      (failure) => emit(AuthError(failure.message)),
      (_) => emit(AuthUnauthenticated()),
    );
  }

  Future<void> _onCheckStatus(
    AuthCheckStatus event,
    Emitter<AuthState> emit,
  ) async {
    emit(AuthLoading());

    final result = await getCurrentUser(const NoParams());

    result.fold(
      (failure) => emit(AuthError(failure.message)),
      (user) =>
          emit(user != null ? AuthAuthenticated(user) : AuthUnauthenticated()),
    );
  }
}
