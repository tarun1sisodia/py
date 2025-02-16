import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:smart_campus/core/widgets/loading_indicator.dart';
import 'package:smart_campus/domain/entities/course.dart';
import 'package:smart_campus/presentation/bloc/teacher/teacher_bloc.dart';
import 'package:smart_campus/config/routes.dart';
import 'package:geolocator/geolocator.dart';
import 'package:network_info_plus/network_info_plus.dart';

class CreateSessionScreen extends StatefulWidget {
  const CreateSessionScreen({super.key});

  @override
  State<CreateSessionScreen> createState() => _CreateSessionScreenState();
}

class _CreateSessionScreenState extends State<CreateSessionScreen> {
  final _formKey = GlobalKey<FormState>();
  Course? _selectedCourse;
  DateTime? _sessionDate;
  TimeOfDay? _startTime;
  TimeOfDay? _endTime;
  Position? _currentPosition;
  String? _wifiSSID;
  String? _wifiBSSID;
  bool _isLoading = false;
  final _radiusController =
      TextEditingController(text: '50'); // Default 50m radius

  @override
  void initState() {
    super.initState();
    _sessionDate = DateTime.now();
    _getCurrentLocation();
    _getWifiInfo();
  }

  Future<void> _getCurrentLocation() async {
    setState(() => _isLoading = true);
    try {
      final permission = await Geolocator.checkPermission();
      if (permission == LocationPermission.denied) {
        await Geolocator.requestPermission();
      }
      _currentPosition = await Geolocator.getCurrentPosition();
    } catch (e) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('Error getting location: $e')),
      );
    }
    setState(() => _isLoading = false);
  }

  Future<void> _getWifiInfo() async {
    try {
      final info = NetworkInfo();
      _wifiSSID = await info.getWifiName();
      _wifiBSSID = await info.getWifiBSSID();
    } catch (e) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('Error getting WiFi info: $e')),
      );
    }
  }

  Future<void> _selectDate() async {
    final picked = await showDatePicker(
      context: context,
      initialDate: _sessionDate ?? DateTime.now(),
      firstDate: DateTime.now(),
      lastDate: DateTime.now().add(const Duration(days: 30)),
    );
    if (picked != null) {
      setState(() => _sessionDate = picked);
    }
  }

  Future<void> _selectStartTime() async {
    final picked = await showTimePicker(
      context: context,
      initialTime: TimeOfDay.now(),
    );
    if (picked != null) {
      setState(() => _startTime = picked);
    }
  }

  Future<void> _selectEndTime() async {
    final picked = await showTimePicker(
      context: context,
      initialTime:
          _startTime?.replacing(hour: _startTime!.hour + 1) ?? TimeOfDay.now(),
    );
    if (picked != null) {
      setState(() => _endTime = picked);
    }
  }

  void _createSession() {
    if (!_formKey.currentState!.validate()) return;
    if (_currentPosition == null) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Location not available')),
      );
      return;
    }
    if (_wifiSSID == null || _wifiBSSID == null) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('WiFi information not available')),
      );
      return;
    }

    final startDateTime = DateTime(
      _sessionDate!.year,
      _sessionDate!.month,
      _sessionDate!.day,
      _startTime!.hour,
      _startTime!.minute,
    );

    final endDateTime = DateTime(
      _sessionDate!.year,
      _sessionDate!.month,
      _sessionDate!.day,
      _endTime!.hour,
      _endTime!.minute,
    );

    context.read<TeacherBloc>().add(
          CreateSessionRequested(
            courseId: _selectedCourse!.id,
            sessionDate: _sessionDate!,
            startTime: startDateTime,
            endTime: endDateTime,
            locationLatitude: _currentPosition!.latitude,
            locationLongitude: _currentPosition!.longitude,
            locationRadius: int.parse(_radiusController.text),
            wifiSSID: _wifiSSID!,
            wifiBSSID: _wifiBSSID!,
          ),
        );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Create Session'),
      ),
      body: BlocConsumer<TeacherBloc, TeacherState>(
        listener: (context, state) {
          if (state is SessionCreationSuccess) {
            ScaffoldMessenger.of(context).showSnackBar(
              const SnackBar(content: Text('Session created successfully')),
            );
            Routes.pop(context);
          } else if (state is TeacherError) {
            ScaffoldMessenger.of(context).showSnackBar(
              SnackBar(content: Text(state.message)),
            );
          }
        },
        builder: (context, state) {
          if (state is TeacherLoading || _isLoading) {
            return const Center(child: LoadingIndicator());
          }

          return SingleChildScrollView(
            padding: const EdgeInsets.all(16.0),
            child: Form(
              key: _formKey,
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.stretch,
                children: [
                  if (state is TeacherDashboardLoaded) ...[
                    DropdownButtonFormField<Course>(
                      value: _selectedCourse,
                      decoration: const InputDecoration(
                        labelText: 'Select Course',
                      ),
                      items: state.assignedCourses
                          .map((course) => DropdownMenuItem(
                                value: course,
                                child: Text(
                                    '${course.courseCode} - ${course.courseName}'),
                              ))
                          .toList(),
                      onChanged: (course) {
                        setState(() => _selectedCourse = course);
                      },
                      validator: (value) =>
                          value == null ? 'Please select a course' : null,
                    ),
                    const SizedBox(height: 16),
                  ],
                  ListTile(
                    title: Text(
                      'Session Date: ${_sessionDate?.toString().split(' ')[0] ?? 'Not selected'}',
                    ),
                    trailing: const Icon(Icons.calendar_today),
                    onTap: _selectDate,
                  ),
                  const SizedBox(height: 16),
                  ListTile(
                    title: Text(
                      'Start Time: ${_startTime?.format(context) ?? 'Not selected'}',
                    ),
                    trailing: const Icon(Icons.access_time),
                    onTap: _selectStartTime,
                  ),
                  const SizedBox(height: 16),
                  ListTile(
                    title: Text(
                      'End Time: ${_endTime?.format(context) ?? 'Not selected'}',
                    ),
                    trailing: const Icon(Icons.access_time),
                    onTap: _selectEndTime,
                  ),
                  const SizedBox(height: 16),
                  TextFormField(
                    controller: _radiusController,
                    decoration: const InputDecoration(
                      labelText: 'Allowed Radius (meters)',
                      suffixText: 'm',
                    ),
                    keyboardType: TextInputType.number,
                    validator: (value) {
                      if (value == null || value.isEmpty) {
                        return 'Please enter allowed radius';
                      }
                      final radius = int.tryParse(value);
                      if (radius == null || radius < 10 || radius > 200) {
                        return 'Radius must be between 10m and 200m';
                      }
                      return null;
                    },
                  ),
                  const SizedBox(height: 24),
                  FilledButton(
                    onPressed: _createSession,
                    child: const Text('Create Session'),
                  ),
                ],
              ),
            ),
          );
        },
      ),
    );
  }

  @override
  void dispose() {
    _radiusController.dispose();
    super.dispose();
  }
}
