import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:smart_campus/core/widgets/loading_indicator.dart';
import 'package:smart_campus/domain/entities/session.dart';
import 'package:smart_campus/domain/entities/course.dart';
import 'package:smart_campus/domain/entities/attendance_record.dart';
import 'package:smart_campus/presentation/bloc/teacher/teacher_bloc.dart';

class SessionDetailsScreen extends StatefulWidget {
  final String sessionId;

  const SessionDetailsScreen({
    super.key,
    required this.sessionId,
  });

  @override
  State<SessionDetailsScreen> createState() => _SessionDetailsScreenState();
}

class _SessionDetailsScreenState extends State<SessionDetailsScreen> {
  @override
  void initState() {
    super.initState();
    context.read<TeacherBloc>().add(LoadSessionDetails(widget.sessionId));
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Session Details'),
        actions: [
          IconButton(
            icon: const Icon(Icons.refresh),
            onPressed: () {
              context
                  .read<TeacherBloc>()
                  .add(LoadSessionDetails(widget.sessionId));
            },
          ),
        ],
      ),
      body: BlocBuilder<TeacherBloc, TeacherState>(
        builder: (context, state) {
          if (state is TeacherLoading) {
            return const Center(child: LoadingIndicator());
          }

          if (state is TeacherError) {
            return Center(
              child: Text(
                state.message,
                style: Theme.of(context).textTheme.bodyLarge?.copyWith(
                      color: Theme.of(context).colorScheme.error,
                    ),
              ),
            );
          }

          if (state is SessionDetailsLoaded) {
            return _buildSessionDetails(context, state);
          }

          return const SizedBox.shrink();
        },
      ),
      floatingActionButton: BlocBuilder<TeacherBloc, TeacherState>(
        builder: (context, state) {
          if (state is SessionDetailsLoaded && !state.session.isEnded) {
            return FloatingActionButton.extended(
              onPressed: () {
                context
                    .read<TeacherBloc>()
                    .add(EndSessionRequested(widget.sessionId));
              },
              icon: const Icon(Icons.stop),
              label: const Text('End Session'),
            );
          }
          return const SizedBox.shrink();
        },
      ),
    );
  }

  Widget _buildSessionDetails(
      BuildContext context, SessionDetailsLoaded state) {
    return SingleChildScrollView(
      padding: const EdgeInsets.all(16.0),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          _buildSessionInfo(context, state.session, state.course),
          const SizedBox(height: 24),
          _buildAttendanceStats(context, state),
          const SizedBox(height: 24),
          _buildAttendanceList(context, state.attendanceRecords),
        ],
      ),
    );
  }

  Widget _buildSessionInfo(
    BuildContext context,
    Session session,
    Course course,
  ) {
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              course.courseName,
              style: Theme.of(context).textTheme.titleLarge,
            ),
            const SizedBox(height: 8),
            Text(
              course.courseCode,
              style: Theme.of(context).textTheme.bodyLarge,
            ),
            const SizedBox(height: 16),
            _buildInfoRow(
              context,
              'Status',
              session.isEnded ? 'Ended' : 'Active',
              session.isEnded
                  ? Icons.check_circle_outline
                  : Icons.radio_button_checked,
              session.isEnded
                  ? Theme.of(context).colorScheme.outline
                  : Theme.of(context).colorScheme.primary,
            ),
            const SizedBox(height: 8),
            _buildInfoRow(
              context,
              'Date',
              session.sessionDate.toString().split(' ')[0],
              Icons.calendar_today,
              Theme.of(context).colorScheme.secondary,
            ),
            const SizedBox(height: 8),
            _buildInfoRow(
              context,
              'Time',
              '${session.startTime.hour}:${session.startTime.minute.toString().padLeft(2, '0')} - ${session.endTime.hour}:${session.endTime.minute.toString().padLeft(2, '0')}',
              Icons.access_time,
              Theme.of(context).colorScheme.secondary,
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildInfoRow(
    BuildContext context,
    String label,
    String value,
    IconData icon,
    Color color,
  ) {
    return Row(
      children: [
        Icon(icon, size: 20, color: color),
        const SizedBox(width: 8),
        Text(
          '$label: ',
          style: Theme.of(context).textTheme.bodyMedium,
        ),
        Text(
          value,
          style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                fontWeight: FontWeight.bold,
                color: color,
              ),
        ),
      ],
    );
  }

  Widget _buildAttendanceStats(
      BuildContext context, SessionDetailsLoaded state) {
    final totalStudents = state.attendanceRecords.length;
    final verifiedCount = state.attendanceRecords
        .where((record) => record.verificationStatus == 'verified')
        .length;
    final pendingCount = state.attendanceRecords
        .where((record) => record.verificationStatus == 'pending')
        .length;
    final rejectedCount = state.attendanceRecords
        .where((record) => record.verificationStatus == 'rejected')
        .length;

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          'Attendance Overview',
          style: Theme.of(context).textTheme.titleLarge,
        ),
        const SizedBox(height: 16),
        Row(
          children: [
            Expanded(
              child: _buildStatCard(
                context,
                'Total',
                totalStudents.toString(),
                Icons.people,
                Theme.of(context).colorScheme.primary,
              ),
            ),
            const SizedBox(width: 8),
            Expanded(
              child: _buildStatCard(
                context,
                'Verified',
                verifiedCount.toString(),
                Icons.check_circle,
                Colors.green,
              ),
            ),
          ],
        ),
        const SizedBox(height: 8),
        Row(
          children: [
            Expanded(
              child: _buildStatCard(
                context,
                'Pending',
                pendingCount.toString(),
                Icons.pending,
                Colors.orange,
              ),
            ),
            const SizedBox(width: 8),
            Expanded(
              child: _buildStatCard(
                context,
                'Rejected',
                rejectedCount.toString(),
                Icons.cancel,
                Theme.of(context).colorScheme.error,
              ),
            ),
          ],
        ),
      ],
    );
  }

  Widget _buildStatCard(
    BuildContext context,
    String label,
    String value,
    IconData icon,
    Color color,
  ) {
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Column(
          children: [
            Icon(icon, color: color),
            const SizedBox(height: 8),
            Text(
              value,
              style: Theme.of(context).textTheme.headlineMedium?.copyWith(
                    color: color,
                    fontWeight: FontWeight.bold,
                  ),
            ),
            Text(
              label,
              style: Theme.of(context).textTheme.bodyMedium,
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildAttendanceList(
    BuildContext context,
    List<AttendanceRecord> records,
  ) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          'Attendance Records',
          style: Theme.of(context).textTheme.titleLarge,
        ),
        const SizedBox(height: 16),
        ListView.builder(
          shrinkWrap: true,
          physics: const NeverScrollableScrollPhysics(),
          itemCount: records.length,
          itemBuilder: (context, index) {
            final record = records[index];
            return Card(
              child: ListTile(
                leading: _buildStatusIcon(context, record.verificationStatus),
                title: Text(record.studentName),
                subtitle: Text(
                  'Marked at ${record.markedAt.hour}:${record.markedAt.minute.toString().padLeft(2, '0')}',
                ),
                trailing: record.verificationStatus == 'pending'
                    ? Row(
                        mainAxisSize: MainAxisSize.min,
                        children: [
                          IconButton(
                            icon: const Icon(Icons.check),
                            color: Colors.green,
                            onPressed: () {
                              context.read<TeacherBloc>().add(
                                    VerifyAttendanceRequested(record.id),
                                  );
                            },
                          ),
                          IconButton(
                            icon: const Icon(Icons.close),
                            color: Theme.of(context).colorScheme.error,
                            onPressed: () {
                              _showRejectDialog(context, record.id);
                            },
                          ),
                        ],
                      )
                    : null,
              ),
            );
          },
        ),
      ],
    );
  }

  Widget _buildStatusIcon(BuildContext context, String status) {
    switch (status) {
      case 'verified':
        return const CircleAvatar(
          backgroundColor: Colors.green,
          child: Icon(Icons.check, color: Colors.white),
        );
      case 'rejected':
        return CircleAvatar(
          backgroundColor: Theme.of(context).colorScheme.error,
          child: const Icon(Icons.close, color: Colors.white),
        );
      default:
        return const CircleAvatar(
          backgroundColor: Colors.orange,
          child: Icon(Icons.pending, color: Colors.white),
        );
    }
  }

  Future<void> _showRejectDialog(
      BuildContext context, String attendanceId) async {
    final reasonController = TextEditingController();
    return showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Reject Attendance'),
        content: TextField(
          controller: reasonController,
          decoration: const InputDecoration(
            labelText: 'Reason for rejection',
            hintText: 'Enter reason',
          ),
          maxLines: 3,
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('Cancel'),
          ),
          FilledButton(
            onPressed: () {
              if (reasonController.text.isNotEmpty) {
                context.read<TeacherBloc>().add(
                      RejectAttendanceRequested(
                        attendanceId: attendanceId,
                        reason: reasonController.text,
                      ),
                    );
                Navigator.pop(context);
              }
            },
            child: const Text('Reject'),
          ),
        ],
      ),
    );
  }
}
