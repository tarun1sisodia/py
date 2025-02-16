import 'package:equatable/equatable.dart';

enum VerificationStatus {
  initial,
  verifying,
  verified,
  failed,
  monitoring,
  invalid,
}

class AttendanceVerificationState extends Equatable {
  final VerificationStatus status;
  final bool isLocationValid;
  final bool isWifiValid;
  final String? errorMessage;
  final DateTime? lastVerifiedAt;

  const AttendanceVerificationState({
    this.status = VerificationStatus.initial,
    this.isLocationValid = false,
    this.isWifiValid = false,
    this.errorMessage,
    this.lastVerifiedAt,
  });

  bool get isValid => isLocationValid && isWifiValid;

  AttendanceVerificationState copyWith({
    VerificationStatus? status,
    bool? isLocationValid,
    bool? isWifiValid,
    String? errorMessage,
    DateTime? lastVerifiedAt,
  }) {
    return AttendanceVerificationState(
      status: status ?? this.status,
      isLocationValid: isLocationValid ?? this.isLocationValid,
      isWifiValid: isWifiValid ?? this.isWifiValid,
      errorMessage: errorMessage ?? this.errorMessage,
      lastVerifiedAt: lastVerifiedAt ?? this.lastVerifiedAt,
    );
  }

  @override
  List<Object?> get props => [
        status,
        isLocationValid,
        isWifiValid,
        errorMessage,
        lastVerifiedAt,
      ];
}
