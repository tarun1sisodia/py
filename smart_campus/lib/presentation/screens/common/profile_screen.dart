import 'package:flutter/material.dart';
import 'package:smart_campus/config/theme.dart';

class ProfileScreen extends StatefulWidget {
  const ProfileScreen({super.key});

  @override
  State<ProfileScreen> createState() => _ProfileScreenState();
}

class _ProfileScreenState extends State<ProfileScreen> {
  bool _isLoading = false;
  bool _isEditing = false;
  final _formKey = GlobalKey<FormState>();
  final _fullNameController = TextEditingController();
  final _emailController = TextEditingController();
  final _departmentController = TextEditingController();
  final _yearOfStudyController = TextEditingController();
  final _enrollmentNumberController = TextEditingController();
  final _employeeIdController = TextEditingController();

  @override
  void initState() {
    super.initState();
    _loadProfile();
  }

  @override
  void dispose() {
    _fullNameController.dispose();
    _emailController.dispose();
    _departmentController.dispose();
    _yearOfStudyController.dispose();
    _enrollmentNumberController.dispose();
    _employeeIdController.dispose();
    super.dispose();
  }

  Future<void> _loadProfile() async {
    setState(() {
      _isLoading = true;
    });

    try {
      // TODO: Load user profile
      // For now, using dummy data
      _fullNameController.text = 'John Doe';
      _emailController.text = 'john.doe@example.com';
      _departmentController.text = 'Computer Science';
      _yearOfStudyController.text = '3';
      _enrollmentNumberController.text = 'CS2020001';
      _employeeIdController.text = 'EMP001';
    } catch (e) {
      // TODO: Handle error
    } finally {
      setState(() {
        _isLoading = false;
      });
    }
  }

  Future<void> _updateProfile() async {
    if (!_formKey.currentState!.validate()) {
      return;
    }

    setState(() {
      _isLoading = true;
    });

    try {
      // TODO: Update user profile
      setState(() {
        _isEditing = false;
      });
    } catch (e) {
      // TODO: Handle error
    } finally {
      setState(() {
        _isLoading = false;
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Profile'),
        actions: [
          IconButton(
            icon: Icon(_isEditing ? Icons.close : Icons.edit),
            onPressed: () {
              setState(() {
                _isEditing = !_isEditing;
              });
            },
          ),
        ],
      ),
      body: _isLoading
          ? const Center(child: CircularProgressIndicator())
          : SingleChildScrollView(
              padding: const EdgeInsets.all(16.0),
              child: Form(
                key: _formKey,
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    _buildProfileHeader(),
                    const SizedBox(height: 24),
                    _buildProfileForm(),
                    if (_isEditing) ...[
                      const SizedBox(height: 24),
                      _buildUpdateButton(),
                    ],
                  ],
                ),
              ),
            ),
    );
  }

  Widget _buildProfileHeader() {
    return Center(
      child: Column(
        children: [
          CircleAvatar(
            radius: 50,
            backgroundColor: Theme.of(context).colorScheme.primary,
            child: Text(
              _fullNameController.text.isNotEmpty
                  ? _fullNameController.text.substring(0, 1).toUpperCase()
                  : '?',
              style: AppTheme.headlineLarge.copyWith(
                color: Theme.of(context).colorScheme.onPrimary,
              ),
            ),
          ),
          const SizedBox(height: 16),
          Text(
            _fullNameController.text,
            style: AppTheme.headlineMedium,
          ),
          const SizedBox(height: 8),
          Text(
            _emailController.text,
            style: AppTheme.bodyLarge,
          ),
        ],
      ),
    );
  }

  Widget _buildProfileForm() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          'Personal Information',
          style: AppTheme.titleLarge,
        ),
        const SizedBox(height: 16),
        TextFormField(
          controller: _fullNameController,
          enabled: _isEditing,
          decoration: const InputDecoration(
            labelText: 'Full Name',
            prefixIcon: Icon(Icons.person),
          ),
          validator: (value) {
            if (value == null || value.isEmpty) {
              return 'Please enter your full name';
            }
            return null;
          },
        ),
        const SizedBox(height: 16),
        TextFormField(
          controller: _emailController,
          enabled: false,
          decoration: const InputDecoration(
            labelText: 'Email',
            prefixIcon: Icon(Icons.email),
          ),
        ),
        const SizedBox(height: 16),
        TextFormField(
          controller: _departmentController,
          enabled: _isEditing,
          decoration: const InputDecoration(
            labelText: 'Department',
            prefixIcon: Icon(Icons.business),
          ),
          validator: (value) {
            if (value == null || value.isEmpty) {
              return 'Please enter your department';
            }
            return null;
          },
        ),
        const SizedBox(height: 16),
        TextFormField(
          controller: _yearOfStudyController,
          enabled: _isEditing,
          keyboardType: TextInputType.number,
          decoration: const InputDecoration(
            labelText: 'Year of Study',
            prefixIcon: Icon(Icons.school),
          ),
          validator: (value) {
            if (value == null || value.isEmpty) {
              return 'Please enter your year of study';
            }
            final year = int.tryParse(value);
            if (year == null || year < 1 || year > 5) {
              return 'Please enter a valid year (1-5)';
            }
            return null;
          },
        ),
        const SizedBox(height: 16),
        TextFormField(
          controller: _enrollmentNumberController,
          enabled: false,
          decoration: const InputDecoration(
            labelText: 'Enrollment Number',
            prefixIcon: Icon(Icons.numbers),
          ),
        ),
        const SizedBox(height: 16),
        TextFormField(
          controller: _employeeIdController,
          enabled: false,
          decoration: const InputDecoration(
            labelText: 'Employee ID',
            prefixIcon: Icon(Icons.badge),
          ),
        ),
      ],
    );
  }

  Widget _buildUpdateButton() {
    return SizedBox(
      width: double.infinity,
      child: ElevatedButton(
        onPressed: _updateProfile,
        child: const Text('Update Profile'),
      ),
    );
  }
}
