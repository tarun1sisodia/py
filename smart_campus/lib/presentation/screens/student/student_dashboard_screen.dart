import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:smart_campus/config/routes.dart';
import 'package:smart_campus/config/theme.dart';
import 'package:smart_campus/domain/entities/course.dart';
import 'package:smart_campus/domain/entities/session.dart';
import 'package:smart_campus/presentation/bloc/student/student_bloc.dart';
import 'package:smart_campus/presentation/bloc/student/student_event.dart';
import 'package:smart_campus/presentation/bloc/student/student_state.dart';

class StudentDashboardScreen extends StatefulWidget {
  const StudentDashboardScreen({super.key});

  @override
  State<StudentDashboardScreen> createState() => _StudentDashboardScreenState();
}

class _StudentDashboardScreenState extends State<StudentDashboardScreen> {
  List<Session> _activeSessions = [];
  List<Course> _courses = [];
  bool _isLoading = true;

  @override
  void initState() {
    super.initState();
    _loadDashboardData();
  }

  void _loadDashboardData() {
    context.read<StudentBloc>().add(const LoadStudentDashboard());
  }

  @override
  Widget build(BuildContext context) {
    return BlocConsumer<StudentBloc, StudentState>(
      listener: (context, state) {
        if (state is StudentDashboardLoaded) {
          setState(() {
            _activeSessions = state.activeSessions;
            _courses = state.courses;
            _isLoading = false;
          });
        }
      },
      builder: (context, state) {
        return Scaffold(
          appBar: AppBar(
            title: const Text('Student Dashboard'),
            actions: [
              IconButton(
                icon: const Icon(Icons.person),
                onPressed: () =>
                    Routes.navigateTo(context, Routes.profile),
              ),
              IconButton(
                icon: const Icon(Icons.settings),
                onPressed: () =>
                    Routes.navigateTo(context, Routes.settings),
              ),
            ],
          ),
          body: RefreshIndicator(
            onRefresh: () async {
              _loadDashboardData();
            },
            child: _isLoading
                ? const Center(child: CircularProgressIndicator())
                : SingleChildScrollView(
                    padding: const EdgeInsets.all(16.0),
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        _buildWelcomeCard(),
                        const SizedBox(height: 24),
                        _buildAttendanceStats(),
                        const SizedBox(height: 24),
                        _buildActiveSessions(),
                        const SizedBox(height: 24),
                        _buildRecentHistory(),
                      ],
                    ),
                  ),
          ),
        );
      },
    );
  }

  Widget _buildWelcomeCard() {
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              'Welcome back, Student!',
              style: AppTheme.headlineMedium,
            ),
            const SizedBox(height: 8),
            Text(
              'You have ${_activeSessions.length} active sessions',
              style: AppTheme.bodyLarge,
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildAttendanceStats() {
    return BlocBuilder<StudentBloc, StudentState>(
      builder: (context, state) {
        if (state is! StudentDashboardLoaded) {
          return const SizedBox.shrink();
        }

        final stats = (state).attendanceStats;

        return Card(
          child: Padding(
            padding: const EdgeInsets.all(16.0),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  'Attendance Overview',
                  style: AppTheme.titleLarge,
                ),
                const SizedBox(height: 16),
                Row(
                  children: [
                    Expanded(
                      child: _buildStatCard(
                        'Overall',
                        '${stats.overallAttendance.toStringAsFixed(1)}%',
                        Icons.analytics,
                      ),
                    ),
                    const SizedBox(width: 16),
                    Expanded(
                      child: _buildStatCard(
                        'This Month',
                        '${stats.monthlyAttendance.toStringAsFixed(1)}%',
                        Icons.calendar_month,
                      ),
                    ),
                  ],
                ),
                const SizedBox(height: 16),
                Row(
                  children: [
                    Expanded(
                      child: _buildStatCard(
                        'Total Sessions',
                        stats.totalSessions.toString(),
                        Icons.history,
                      ),
                    ),
                    const SizedBox(width: 16),
                    Expanded(
                      child: _buildStatCard(
                        'Present',
                        stats.presentSessions.toString(),
                        Icons.check_circle,
                      ),
                    ),
                  ],
                ),
              ],
            ),
          ),
        );
      },
    );
  }

  Widget _buildStatCard(String label, String value, IconData icon) {
    return Container(
      padding: const EdgeInsets.all(16.0),
      decoration: BoxDecoration(
        color: Theme.of(context).colorScheme.surface,
        borderRadius: BorderRadius.circular(12),
        border: Border.all(
          color: Theme.of(context).colorScheme.outline.withOpacity(0.2),
        ),
      ),
      child: Column(
        children: [
          Icon(
            icon,
            size: 32,
            color: Theme.of(context).colorScheme.primary,
          ),
          const SizedBox(height: 8),
          Text(
            value,
            style: AppTheme.headlineMedium,
          ),
          const SizedBox(height: 4),
          Text(
            label,
            style: AppTheme.bodyMedium,
            textAlign: TextAlign.center,
          ),
        ],
      ),
    );
  }

  Widget _buildActiveSessions() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Row(
          mainAxisAlignment: MainAxisAlignment.spaceBetween,
          children: [
            Text(
              'Active Sessions',
              style: AppTheme.titleLarge,
            ),
            TextButton(
              onPressed: () =>
                  Routes.navigateTo(context, Routes.markAttendance),
              child: const Text('Mark Attendance'),
            ),
          ],
        ),
        const SizedBox(height: 16),
        if (_activeSessions.isEmpty)
          const Center(
            child: Text('No active sessions'),
          )
        else
          ListView.builder(
            shrinkWrap: true,
            physics: const NeverScrollableScrollPhysics(),
            itemCount: _activeSessions.length,
            itemBuilder: (context, index) {
              final session = _activeSessions[index];
              final course = _courses.firstWhere(
                (c) => c.id == session.courseId,
                orElse: () => Course(
                  id: 'unknown',
                  courseName: 'Unknown Course',
                  courseCode: 'N/A',
                  department: 'N/A',
                  yearOfStudy: 0,
                  semester: 0,
                  createdAt: DateTime.now(),
                  updatedAt: DateTime.now(),
                ),
              );
              return _buildSessionCard(session, course);
            },
          ),
      ],
    );
  }

  Widget _buildRecentHistory() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Row(
          mainAxisAlignment: MainAxisAlignment.spaceBetween,
          children: [
            Text(
              'Recent History',
              style: AppTheme.titleLarge,
            ),
            TextButton(
              onPressed: () =>
                  Routes.navigateTo(context, Routes.attendanceHistory),
              child: const Text('View All'),
            ),
          ],
        ),
        const SizedBox(height: 16),
        if (_activeSessions.isEmpty)
          const Center(
            child: Text('No recent history'),
          )
        else
          ListView.builder(
            shrinkWrap: true,
            physics: const NeverScrollableScrollPhysics(),
            itemCount: _activeSessions.length,
            itemBuilder: (context, index) {
              final session = _activeSessions[index];
              final course = _courses.firstWhere(
                (c) => c.id == session.courseId,
                orElse: () => Course(
                  id: 'unknown',
                  courseName: 'Unknown Course',
                  courseCode: 'N/A',
                  department: 'N/A',
                  yearOfStudy: 0,
                  semester: 0,
                  createdAt: DateTime.now(),
                  updatedAt: DateTime.now(),
                ),
              );
              return _buildSessionCard(session, course);
            },
          ),
      ],
    );
  }

  Widget _buildSessionCard(Session session, Course course) {
    return Card(
      child: ListTile(
        leading: CircleAvatar(
          backgroundColor: Theme.of(context).colorScheme.primary,
          child: Text(
            course.courseCode.substring(0, 2),
            style: TextStyle(
              color: Theme.of(context).colorScheme.onPrimary,
            ),
          ),
        ),
        title: Text(
          course.courseName,
          style: AppTheme.bodyLarge,
        ),
        subtitle: Text(
          '${session.startTime.hour}:${session.startTime.minute.toString().padLeft(2, '0')} - ${session.endTime.hour}:${session.endTime.minute.toString().padLeft(2, '0')}',
          style: AppTheme.bodyMedium,
        ),
        trailing: session.isActive
            ? ElevatedButton(
                onPressed: () =>
                    Routes.navigateTo(context, Routes.markAttendance),
                child: const Text('Mark'),
              )
            : Text(
                session.isCompleted ? 'Present' : 'Absent',
                style: AppTheme.bodyLarge.copyWith(
                  color: session.isCompleted ? Colors.green : Colors.red,
                ),
              ),
      ),
    );
  }
}
