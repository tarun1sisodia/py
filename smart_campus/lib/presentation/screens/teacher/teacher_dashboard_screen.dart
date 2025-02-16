import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:smart_campus/config/routes.dart';
import 'package:smart_campus/domain/entities/session.dart';
import 'package:smart_campus/domain/entities/course.dart';
import 'package:smart_campus/presentation/bloc/teacher/teacher_bloc.dart';
import 'package:smart_campus/core/widgets/loading_indicator.dart';

class TeacherDashboardScreen extends StatelessWidget {
  const TeacherDashboardScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return BlocBuilder<TeacherBloc, TeacherState>(
      builder: (context, state) {
        if (state is TeacherLoading) {
          return const Center(
            child: LoadingIndicator(),
          );
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

        if (state is TeacherDashboardLoaded) {
          return SingleChildScrollView(
            padding: const EdgeInsets.all(16.0),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                _buildHeader(context, state),
                const SizedBox(height: 24),
                _buildStatCards(context, state),
                const SizedBox(height: 24),
                _buildActiveSessions(context, state),
                const SizedBox(height: 24),
                _buildAssignedCourses(context, state),
              ],
            ),
          );
        }

        return const SizedBox.shrink();
      },
    );
  }

  Widget _buildHeader(BuildContext context, TeacherDashboardLoaded state) {
    return Row(
      mainAxisAlignment: MainAxisAlignment.spaceBetween,
      children: [
        Text(
          'Dashboard',
          style: Theme.of(context).textTheme.headlineMedium,
        ),
        FilledButton.icon(
          onPressed: () => Routes.navigateTo(context, Routes.createSession),
          icon: const Icon(Icons.add),
          label: const Text('New Session'),
        ),
      ],
    );
  }

  Widget _buildStatCards(BuildContext context, TeacherDashboardLoaded state) {
    return Row(
      children: [
        Expanded(
          child: _StatCard(
            title: 'Today\'s Sessions',
            value: state.totalSessionsToday.toString(),
            icon: Icons.calendar_today,
          ),
        ),
        const SizedBox(width: 16),
        Expanded(
          child: _StatCard(
            title: 'Students Present',
            value: state.totalStudentsPresent.toString(),
            icon: Icons.people,
          ),
        ),
      ],
    );
  }

  Widget _buildActiveSessions(
    BuildContext context,
    TeacherDashboardLoaded state,
  ) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          'Active Sessions',
          style: Theme.of(context).textTheme.titleLarge,
        ),
        const SizedBox(height: 16),
        if (state.activeSessions.isEmpty)
          Center(
            child: Text(
              'No active sessions',
              style: Theme.of(context).textTheme.bodyLarge?.copyWith(
                    color: Theme.of(context).colorScheme.outline,
                  ),
            ),
          )
        else
          ListView.builder(
            shrinkWrap: true,
            physics: const NeverScrollableScrollPhysics(),
            itemCount: state.activeSessions.length,
            itemBuilder: (context, index) {
              return _SessionCard(session: state.activeSessions[index]);
            },
          ),
      ],
    );
  }

  Widget _buildAssignedCourses(
    BuildContext context,
    TeacherDashboardLoaded state,
  ) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          'Assigned Courses',
          style: Theme.of(context).textTheme.titleLarge,
        ),
        const SizedBox(height: 16),
        if (state.assignedCourses.isEmpty)
          Center(
            child: Text(
              'No courses assigned',
              style: Theme.of(context).textTheme.bodyLarge?.copyWith(
                    color: Theme.of(context).colorScheme.outline,
                  ),
            ),
          )
        else
          ListView.builder(
            shrinkWrap: true,
            physics: const NeverScrollableScrollPhysics(),
            itemCount: state.assignedCourses.length,
            itemBuilder: (context, index) {
              return _CourseCard(course: state.assignedCourses[index]);
            },
          ),
      ],
    );
  }
}

class _StatCard extends StatelessWidget {
  final String title;
  final String value;
  final IconData icon;

  const _StatCard({
    required this.title,
    required this.value,
    required this.icon,
  });

  @override
  Widget build(BuildContext context) {
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Icon(
              icon,
              color: Theme.of(context).colorScheme.primary,
            ),
            const SizedBox(height: 8),
            Text(
              value,
              style: Theme.of(context).textTheme.headlineMedium?.copyWith(
                    color: Theme.of(context).colorScheme.primary,
                  ),
            ),
            const SizedBox(height: 4),
            Text(
              title,
              style: Theme.of(context).textTheme.bodyMedium,
            ),
          ],
        ),
      ),
    );
  }
}

class _SessionCard extends StatelessWidget {
  final Session session;

  const _SessionCard({required this.session});

  @override
  Widget build(BuildContext context) {
    return Card(
      child: ListTile(
        title: Text(session.courseId), // TODO: Get course name from courseId
        subtitle: Text('Started at ${session.startTime}'),
        trailing: FilledButton(
          onPressed: () {
            // TODO: Navigate to session details
          },
          child: const Text('View'),
        ),
      ),
    );
  }
}

class _CourseCard extends StatelessWidget {
  final Course course;

  const _CourseCard({required this.course});

  @override
  Widget build(BuildContext context) {
    return Card(
      child: ListTile(
        title: Text(course.courseName),
        subtitle: Text('${course.courseCode} - Year ${course.yearOfStudy}'),
        trailing: IconButton(
          onPressed: () {
            // TODO: Navigate to course details
          },
          icon: const Icon(Icons.chevron_right),
        ),
      ),
    );
  }
}
