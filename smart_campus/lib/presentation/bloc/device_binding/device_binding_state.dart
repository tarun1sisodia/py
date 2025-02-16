import 'package:equatable/equatable.dart';

enum DeviceBindingStatus {
  initial,
  binding,
  bound,
  unbound,
  failed,
  developerMode,
}

class DeviceBindingState extends Equatable {
  final DeviceBindingStatus status;
  final String? deviceId;
  final Map<String, dynamic>? deviceInfo;
  final bool isDeveloperMode;
  final String? errorMessage;
  final DateTime? bindingTime;

  const DeviceBindingState({
    this.status = DeviceBindingStatus.initial,
    this.deviceId,
    this.deviceInfo,
    this.isDeveloperMode = false,
    this.errorMessage,
    this.bindingTime,
  });

  bool get isBound => status == DeviceBindingStatus.bound;

  DeviceBindingState copyWith({
    DeviceBindingStatus? status,
    String? deviceId,
    Map<String, dynamic>? deviceInfo,
    bool? isDeveloperMode,
    String? errorMessage,
    DateTime? bindingTime,
  }) {
    return DeviceBindingState(
      status: status ?? this.status,
      deviceId: deviceId ?? this.deviceId,
      deviceInfo: deviceInfo ?? this.deviceInfo,
      isDeveloperMode: isDeveloperMode ?? this.isDeveloperMode,
      errorMessage: errorMessage ?? this.errorMessage,
      bindingTime: bindingTime ?? this.bindingTime,
    );
  }

  @override
  List<Object?> get props => [
        status,
        deviceId,
        deviceInfo,
        isDeveloperMode,
        errorMessage,
        bindingTime,
      ];
}
