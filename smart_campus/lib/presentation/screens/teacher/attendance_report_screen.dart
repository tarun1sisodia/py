import 'package:flutter/material.dart';
import 'package:smart_campus/config/theme.dart';
import 'package:smart_campus/domain/entities/course.dart';
import 'package:smart_campus/domain/entities/session.dart';

class AttendanceReportScreen extends StatefulWidget {
  const AttendanceReportScreen({super.key});

  @override
  State<AttendanceReportScreen> createState() => _AttendanceReportScreenState();
}

class _AttendanceReportScreenState extends State<AttendanceReportScreen> {
  bool _isLoading = false;
  Course? _selectedCourse;
  DateTime? _startDate;
  DateTime? _endDate;

  @override
  void initState() {
    super.initState();
    _loadInitialData();
  }

  Future<void> _loadInitialData() async {
    setState(() {
      _isLoading = true;
    });

    try {
      // TODO: Load courses and initial report data
    } catch (e) {
      // TODO: Handle error
    } finally {
      setState(() {
        _isLoading = false;
      });
    }
  }

  Future<void> _selectDateRange() async {
    final DateTimeRange? picked = await showDateRangePicker(
      context: context,
      firstDate: DateTime(2020),
      lastDate: DateTime.now(),
      initialDateRange: _startDate != null && _endDate != null
          ? DateTimeRange(start: _startDate!, end: _endDate!)
          : null,
    );

    if (picked != null) {
      setState(() {
        _startDate = picked.start;
        _endDate = picked.end;
      });
      // TODO: Reload report data with new date range
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Attendance Report'),
      ),
      body: _isLoading
          ? const Center(child: CircularProgressIndicator())
          : SingleChildScrollView(
              padding: const EdgeInsets.all(16.0),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  _buildFilters(),
                  const SizedBox(height: 24),
                  _buildOverallStats(),
                  const SizedBox(height: 24),
                  _buildSessionList(),
                ],
              ),
            ),
    );
  }

  Widget _buildFilters() {
    return Card(
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
            DropdownButtonFormField<Course>(
              value: _selectedCourse,
              decoration: const InputDecoration(
                labelText: 'Course',
                prefixIcon: Icon(Icons.book),
              ),
              items: const [], // TODO: Add course items
              onChanged: (Course? value) {
                setState(() {
                  _selectedCourse = value;
                });
                // TODO: Reload report data with new course
              },
            ),
            const SizedBox(height: 16),
            TextFormField(
              readOnly: true,
              decoration: const InputDecoration(
                labelText: 'Date Range',
                prefixIcon: Icon(Icons.calendar_today),
              ),
              controller: TextEditingController(
                text: _startDate != null && _endDate != null
                    ? '${_startDate!.day}/${_startDate!.month}/${_startDate!.year} - ${_endDate!.day}/${_endDate!.month}/${_endDate!.year}'
                    : 'Select date range',
              ),
              onTap: _selectDateRange,
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildOverallStats() {
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              'Overall Statistics',
              style: AppTheme.titleLarge,
            ),
            const SizedBox(height: 16),
            Row(
              children: [
                Expanded(
                  child: _buildStatCard(
                    'Total Sessions',
                    '0',
                    Icons.calendar_month,
                  ),
                ),
                const SizedBox(width: 16),
                Expanded(
                  child: _buildStatCard(
                    'Average Attendance',
                    '0%',
                    Icons.people,
                  ),
                ),
              ],
            ),
            const SizedBox(height: 16),
            Row(
              children: [
                Expanded(
                  child: _buildStatCard(
                    'Total Students',
                    '0',
                    Icons.school,
                  ),
                ),
                const SizedBox(width: 16),
                Expanded(
                  child: _buildStatCard(
                    'Active Sessions',
                    '0',
                    Icons.timer,
                  ),
                ),
              ],
            ),
          ],
        ),
      ),
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

  Widget _buildSessionList() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          'Recent Sessions',
          style: AppTheme.titleLarge,
        ),
        const SizedBox(height: 16),
        // TODO: Replace with actual session list
        const Center(
          child: Text('No sessions found'),
        ),
      ],
    );
  }

  Widget _buildSessionCard(Session session) {
    return Card(
      child: ListTile(
        leading: CircleAvatar(
          backgroundColor: Theme.of(context).colorScheme.primary,
          child: Text(
            '${session.startTime.day}',
            style: TextStyle(
              color: Theme.of(context).colorScheme.onPrimary,
            ),
          ),
        ),
        title: Text(
          '${session.startTime.hour}:${session.startTime.minute} - ${session.endTime.hour}:${session.endTime.minute}',
          style: AppTheme.bodyLarge,
        ),
        subtitle: Text(
          session.isActive ? 'Active' : 'Completed',
          style: AppTheme.bodyMedium.copyWith(
            color: session.isActive ? Colors.green : null,
          ),
        ),
        trailing: Text(
          '0/0', // TODO: Replace with actual attendance count
          style: AppTheme.bodyLarge.copyWith(fontWeight: FontWeight.w600),
        ),
        onTap: () {
          // TODO: Navigate to session details
        },
      ),
    );
  }
}
