import 'package:equatable/equatable.dart';

class AttendanceStatistics extends Equatable {
  final int totalSessions;
  final int presentSessions;
  final int pendingSessions;
  final int rejectedSessions;
  final double overallAttendance;
  final double monthlyAttendance;

  const AttendanceStatistics({
    required this.totalSessions,
    required this.presentSessions,
    required this.pendingSessions,
    required this.rejectedSessions,
    required this.overallAttendance,
    required this.monthlyAttendance,
  });

  factory AttendanceStatistics.fromJson(Map<String, dynamic> json) {
    return AttendanceStatistics(
      totalSessions: json['totalSessions'] as int,
      presentSessions: json['presentSessions'] as int,
      pendingSessions: json['pendingSessions'] as int,
      rejectedSessions: json['rejectedSessions'] as int,
      overallAttendance: (json['overallAttendance'] as num).toDouble(),
      monthlyAttendance: (json['monthlyAttendance'] as num).toDouble(),
    );
  }

  Map<String, dynamic> toJson() => {
        'totalSessions': totalSessions,
        'presentSessions': presentSessions,
        'pendingSessions': pendingSessions,
        'rejectedSessions': rejectedSessions,
        'overallAttendance': overallAttendance,
        'monthlyAttendance': monthlyAttendance,
      };

  @override
  List<Object> get props => [
        totalSessions,
        presentSessions,
        pendingSessions,
        rejectedSessions,
        overallAttendance,
        monthlyAttendance,
      ];
}
