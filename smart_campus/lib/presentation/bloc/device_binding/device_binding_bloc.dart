import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:injectable/injectable.dart';
import 'package:smart_campus/core/services/device_binding_service.dart';
import 'package:smart_campus/core/services/logger_service.dart';
import 'package:smart_campus/presentation/bloc/device_binding/device_binding_event.dart';
import 'package:smart_campus/presentation/bloc/device_binding/device_binding_state.dart';

@injectable
class DeviceBindingBloc extends Bloc<DeviceBindingEvent, DeviceBindingState> {
  final DeviceBindingService _deviceBindingService;
  final LoggerService _logger;

  DeviceBindingBloc(
    this._deviceBindingService,
    this._logger,
  ) : super(const DeviceBindingState()) {
    on<InitializeDeviceBinding>(_onInitializeDeviceBinding);
    on<BindDevice>(_onBindDevice);
    on<UnbindDevice>(_onUnbindDevice);
    on<VerifyDeviceBinding>(_onVerifyDeviceBinding);
    on<CheckDeviceBindingStatus>(_onCheckDeviceBindingStatus);
    on<ResetDeviceBinding>(_onResetDeviceBinding);
  }

  Future<void> _onInitializeDeviceBinding(
    InitializeDeviceBinding event,
    Emitter<DeviceBindingState> emit,
  ) async {
    try {
      final isDeveloperMode =
          await _deviceBindingService.isDeviceDeveloperModeEnabled();
      if (isDeveloperMode) {
        emit(state.copyWith(
          status: DeviceBindingStatus.developerMode,
          isDeveloperMode: true,
          errorMessage: 'Developer mode is enabled',
        ));
        return;
      }

      final deviceId = await _deviceBindingService.getDeviceId();
      final deviceInfo = await _deviceBindingService.getDeviceInfo();
      final bindingStatus = await _deviceBindingService.getBindingStatus();

      emit(state.copyWith(
        status: bindingStatus['isBound']
            ? DeviceBindingStatus.bound
            : DeviceBindingStatus.unbound,
        deviceId: deviceId,
        deviceInfo: deviceInfo,
        isDeveloperMode: false,
        bindingTime: bindingStatus['bindingTime'] != null
            ? DateTime.parse(bindingStatus['bindingTime'])
            : null,
      ));
    } catch (e) {
      _logger.error('Error initializing device binding', e);
      emit(state.copyWith(
        status: DeviceBindingStatus.failed,
        errorMessage: e.toString(),
      ));
    }
  }

  Future<void> _onBindDevice(
    BindDevice event,
    Emitter<DeviceBindingState> emit,
  ) async {
    try {
      emit(state.copyWith(status: DeviceBindingStatus.binding));

      final isDeveloperMode =
          await _deviceBindingService.isDeviceDeveloperModeEnabled();
      if (isDeveloperMode) {
        emit(state.copyWith(
          status: DeviceBindingStatus.developerMode,
          isDeveloperMode: true,
          errorMessage: 'Cannot bind device with developer mode enabled',
        ));
        return;
      }

      await _deviceBindingService.bindDevice(event.userId);
      final deviceId = await _deviceBindingService.getDeviceId();
      final deviceInfo = await _deviceBindingService.getDeviceInfo();

      emit(state.copyWith(
        status: DeviceBindingStatus.bound,
        deviceId: deviceId,
        deviceInfo: deviceInfo,
        isDeveloperMode: false,
        bindingTime: DateTime.now(),
        errorMessage: null,
      ));
    } catch (e) {
      _logger.error('Error binding device', e);
      emit(state.copyWith(
        status: DeviceBindingStatus.failed,
        errorMessage: e.toString(),
      ));
    }
  }

  Future<void> _onUnbindDevice(
    UnbindDevice event,
    Emitter<DeviceBindingState> emit,
  ) async {
    try {
      await _deviceBindingService.unbindDevice();
      emit(const DeviceBindingState(status: DeviceBindingStatus.unbound));
    } catch (e) {
      _logger.error('Error unbinding device', e);
      emit(state.copyWith(
        status: DeviceBindingStatus.failed,
        errorMessage: e.toString(),
      ));
    }
  }

  Future<void> _onVerifyDeviceBinding(
    VerifyDeviceBinding event,
    Emitter<DeviceBindingState> emit,
  ) async {
    try {
      final isDeveloperMode =
          await _deviceBindingService.isDeviceDeveloperModeEnabled();
      if (isDeveloperMode) {
        emit(state.copyWith(
          status: DeviceBindingStatus.developerMode,
          isDeveloperMode: true,
          errorMessage: 'Developer mode is enabled',
        ));
        return;
      }

      final isValid =
          await _deviceBindingService.verifyDeviceBinding(event.userId);
      if (isValid) {
        final deviceId = await _deviceBindingService.getDeviceId();
        final deviceInfo = await _deviceBindingService.getDeviceInfo();
        emit(state.copyWith(
          status: DeviceBindingStatus.bound,
          deviceId: deviceId,
          deviceInfo: deviceInfo,
          isDeveloperMode: false,
          errorMessage: null,
        ));
      } else {
        emit(state.copyWith(
          status: DeviceBindingStatus.unbound,
          errorMessage: 'Device binding verification failed',
        ));
      }
    } catch (e) {
      _logger.error('Error verifying device binding', e);
      emit(state.copyWith(
        status: DeviceBindingStatus.failed,
        errorMessage: e.toString(),
      ));
    }
  }

  Future<void> _onCheckDeviceBindingStatus(
    CheckDeviceBindingStatus event,
    Emitter<DeviceBindingState> emit,
  ) async {
    try {
      final bindingStatus = await _deviceBindingService.getBindingStatus();

      emit(state.copyWith(
        status: bindingStatus['isBound']
            ? DeviceBindingStatus.bound
            : DeviceBindingStatus.unbound,
        deviceId: bindingStatus['deviceId'],
        isDeveloperMode: bindingStatus['isDeveloperMode'],
        bindingTime: bindingStatus['bindingTime'] != null
            ? DateTime.parse(bindingStatus['bindingTime'])
            : null,
      ));
    } catch (e) {
      _logger.error('Error checking device binding status', e);
      emit(state.copyWith(
        status: DeviceBindingStatus.failed,
        errorMessage: e.toString(),
      ));
    }
  }

  void _onResetDeviceBinding(
    ResetDeviceBinding event,
    Emitter<DeviceBindingState> emit,
  ) {
    emit(const DeviceBindingState());
  }
}
