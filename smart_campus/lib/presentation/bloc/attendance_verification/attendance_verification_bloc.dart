import 'dart:async';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:injectable/injectable.dart';
import 'package:smart_campus/core/services/attendance_verification_service.dart';
import 'package:smart_campus/core/services/logger_service.dart';
import 'package:smart_campus/presentation/bloc/attendance_verification/attendance_verification_event.dart';
import 'package:smart_campus/presentation/bloc/attendance_verification/attendance_verification_state.dart';

@injectable
class AttendanceVerificationBloc
    extends Bloc<AttendanceVerificationEvent, AttendanceVerificationState> {
  final AttendanceVerificationService _verificationService;
  final LoggerService _logger;
  StreamSubscription<bool>? _monitoringSubscription;

  AttendanceVerificationBloc(
    this._verificationService,
    this._logger,
  ) : super(const AttendanceVerificationState()) {
    on<StartVerification>(_onStartVerification);
    on<StartMonitoring>(_onStartMonitoring);
    on<StopMonitoring>(_onStopMonitoring);
    on<ResetVerification>(_onResetVerification);
  }

  Future<void> _onStartVerification(
    StartVerification event,
    Emitter<AttendanceVerificationState> emit,
  ) async {
    try {
      emit(state.copyWith(status: VerificationStatus.verifying));

      final isValid = await _verificationService.verifyAttendance(
        userId: event.userId,
        sessionLatitude: event.sessionLatitude,
        sessionLongitude: event.sessionLongitude,
        allowedRadius: event.allowedRadius.round(),
        sessionSSID: event.sessionSSID,
        sessionBSSID: event.sessionBSSID,
      );

      if (isValid) {
        emit(state.copyWith(
          status: VerificationStatus.verified,
          isLocationValid: true,
          isWifiValid: true,
          lastVerifiedAt: DateTime.now(),
          errorMessage: null,
        ));
      } else {
        emit(state.copyWith(
          status: VerificationStatus.failed,
          errorMessage: 'Location or WiFi verification failed',
        ));
      }
    } catch (e) {
      _logger.error('Error during verification', e);
      emit(state.copyWith(
        status: VerificationStatus.failed,
        errorMessage: e.toString(),
      ));
    }
  }

  Future<void> _onStartMonitoring(
    StartMonitoring event,
    Emitter<AttendanceVerificationState> emit,
  ) async {
    await _monitoringSubscription?.cancel();

    emit(state.copyWith(status: VerificationStatus.monitoring));

    _monitoringSubscription = _verificationService
        .monitorAttendance(
      userId: event.userId,
      sessionLatitude: event.sessionLatitude,
      sessionLongitude: event.sessionLongitude,
      allowedRadius: event.allowedRadius.round(),
      sessionSSID: event.sessionSSID,
      sessionBSSID: event.sessionBSSID,
    )
        .listen(
      (isValid) => add(
        isValid
            ? StartVerification(
                userId: event.userId,
                sessionLatitude: event.sessionLatitude,
                sessionLongitude: event.sessionLongitude,
                allowedRadius: event.allowedRadius,
                sessionSSID: event.sessionSSID,
                sessionBSSID: event.sessionBSSID,
              )
            : ResetVerification(),
      ),
      onError: (error) {
        _logger.error('Error during monitoring', error);
        emit(state.copyWith(
          status: VerificationStatus.failed,
          errorMessage: error.toString(),
        ));
      },
    );
  }

  Future<void> _onStopMonitoring(
    StopMonitoring event,
    Emitter<AttendanceVerificationState> emit,
  ) async {
    await _monitoringSubscription?.cancel();
    _monitoringSubscription = null;
    emit(state.copyWith(status: VerificationStatus.initial));
  }

  void _onResetVerification(
    ResetVerification event,
    Emitter<AttendanceVerificationState> emit,
  ) {
    emit(const AttendanceVerificationState());
  }

  @override
  Future<void> close() async {
    await _monitoringSubscription?.cancel();
    return super.close();
  }
}
