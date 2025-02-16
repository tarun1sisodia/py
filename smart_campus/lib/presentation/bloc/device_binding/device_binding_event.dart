import 'package:equatable/equatable.dart';

abstract class DeviceBindingEvent extends Equatable {
  const DeviceBindingEvent();

  @override
  List<Object?> get props => [];
}

class InitializeDeviceBinding extends DeviceBindingEvent {}

class BindDevice extends DeviceBindingEvent {
  final String userId;

  const BindDevice({required this.userId});

  @override
  List<Object> get props => [userId];
}

class UnbindDevice extends DeviceBindingEvent {}

class VerifyDeviceBinding extends DeviceBindingEvent {
  final String userId;

  const VerifyDeviceBinding({required this.userId});

  @override
  List<Object> get props => [userId];
}

class CheckDeviceBindingStatus extends DeviceBindingEvent {}

class ResetDeviceBinding extends DeviceBindingEvent {}
