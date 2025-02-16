import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:intl/intl.dart';
import 'package:smart_campus/config/theme.dart';
import 'package:smart_campus/domain/entities/course.dart';
import 'package:smart_campus/domain/entities/session.dart';
import 'package:smart_campus/presentation/bloc/student/student_bloc.dart';
import 'package:smart_campus/presentation/bloc/student/student_event.dart';
import 'package:smart_campus/presentation/bloc/student/student_state.dart';

class AttendanceHistoryScreen extends StatefulWidget {
  const AttendanceHistoryScreen({super.key});

  @override
  State<AttendanceHistoryScreen> createState() =>
      _AttendanceHistoryScreenState();
}

class _AttendanceHistoryScreenState extends State<AttendanceHistoryScreen> {
  List<Session> _sessions = [];
  List<Course> _courses = [];
  bool _isLoading = true;
  Course? _selectedCourse;
  DateTime? _startDate;
  DateTime? _endDate;

  @override
  void initState() {
    super.initState();
    _loadAttendanceHistory();
  }

  void _loadAttendanceHistory() {
    context.read<StudentBloc>().add(LoadStudentSessions(
          startDate: _startDate,
          endDate: _endDate,
          course: _selectedCourse,
        ));
  }

  @override
  Widget build(BuildContext context) {
    return BlocConsumer<StudentBloc, StudentState>(
      listener: (context, state) {
        if (state is StudentSessionsLoaded) {
          setState(() {
            _sessions = state.sessions;
            _courses = state.courses;
            _isLoading = false;
          });
        }
      },
      builder: (context, state) {
        return Scaffold(
          appBar: AppBar(
            title: const Text('Attendance History'),
          ),
          body: _isLoading
              ? const Center(child: CircularProgressIndicator())
              : Column(
                  children: [
                    _buildFilters(),
                    Expanded(
                      child: _sessions.isEmpty
                          ? const Center(
                              child: Text('No attendance records found'),
                            )
                          : ListView.builder(
                              padding: const EdgeInsets.all(16.0),
                              itemCount: _sessions.length,
                              itemBuilder: (context, index) {
                                final session = _sessions[index];
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
                    ),
                  ],
                ),
        );
      },
    );
  }

  Widget _buildFilters() {
    return Card(
      margin: const EdgeInsets.all(16.0),
      child: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              'Filters',
              style: AppTheme.titleLarge,
            ),
            const SizedBox(height: 16),
            Row(
              children: [
                Expanded(
                  child: DropdownButtonFormField<Course>(
                    value: _selectedCourse,
                    decoration: const InputDecoration(
                      labelText: 'Course',
                      border: OutlineInputBorder(),
                    ),
                    items: [
                      const DropdownMenuItem<Course>(
                        value: null,
                        child: Text('All Courses'),
                      ),
                      ..._courses.map((course) {
                        return DropdownMenuItem<Course>(
                          value: course,
                          child: Text(course.courseCode),
                        );
                      }).toList(),
                    ],
                    onChanged: (Course? value) {
                      setState(() {
                        _selectedCourse = value;
                      });
                      _loadAttendanceHistory();
                    },
                  ),
                ),
                const SizedBox(width: 16),
                IconButton(
                  onPressed: () async {
                    final DateTimeRange? dateRange = await showDateRangePicker(
                      context: context,
                      firstDate: DateTime(2020),
                      lastDate: DateTime.now(),
                      initialDateRange: _startDate != null && _endDate != null
                          ? DateTimeRange(
                              start: _startDate!,
                              end: _endDate!,
                            )
                          : null,
                    );

                    if (dateRange != null) {
                      setState(() {
                        _startDate = dateRange.start;
                        _endDate = dateRange.end;
                      });
                      _loadAttendanceHistory();
                    }
                  },
                  icon: const Icon(Icons.calendar_month),
                ),
              ],
            ),
            if (_startDate != null && _endDate != null) ...[
              const SizedBox(height: 8),
              Text(
                'Date Range: ${DateFormat('dd/MM/yyyy').format(_startDate!)} - ${DateFormat('dd/MM/yyyy').format(_endDate!)}',
                style: AppTheme.bodyMedium,
              ),
            ],
          ],
        ),
      ),
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
        subtitle: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              DateFormat('dd/MM/yyyy').format(session.sessionDate),
              style: AppTheme.bodyMedium,
            ),
            Text(
              '${session.startTime.hour}:${session.startTime.minute.toString().padLeft(2, '0')} - ${session.endTime.hour}:${session.endTime.minute.toString().padLeft(2, '0')}',
              style: AppTheme.bodyMedium,
            ),
          ],
        ),
        trailing: _buildAttendanceStatus(session),
      ),
    );
  }

  Widget _buildAttendanceStatus(Session session) {
    Color color;
    String text;

    if (session.isCompleted) {
      color = Colors.green;
      text = 'Present';
    } else if (session.isPending) {
      color = Colors.orange;
      text = 'Pending';
    } else if (session.isRejected) {
      color = Colors.red;
      text = 'Rejected';
    } else {
      color = Colors.grey;
      text = 'Absent';
    }

    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
      decoration: BoxDecoration(
        color: color.withOpacity(0.1),
        borderRadius: BorderRadius.circular(12),
        border: Border.all(color: color),
      ),
      child: Text(
        text,
        style: AppTheme.bodyMedium.copyWith(color: color),
      ),
    );
  }
}
