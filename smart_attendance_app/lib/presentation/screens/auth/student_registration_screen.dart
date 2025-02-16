import 'package:flutter/material.dart';
import 'package:smart_attendance_app/services/authentication_service.dart';
import 'package:smart_attendance_app/presentation/screens/auth/otp_verification_screen.dart';

class StudentRegistrationScreen extends StatefulWidget {
  const StudentRegistrationScreen({Key? key}) : super(key: key);

  @override
  _StudentRegistrationScreenState createState() =>
      _StudentRegistrationScreenState();
}

class _StudentRegistrationScreenState extends State<StudentRegistrationScreen> {
  final _formKey = GlobalKey<FormState>();

  // Controllers for input fields
  final TextEditingController _fullNameController = TextEditingController();
  final TextEditingController _rollNumberController = TextEditingController();
  final TextEditingController _courseController = TextEditingController();
  final TextEditingController _academicYearController = TextEditingController();
  final TextEditingController _phoneController = TextEditingController();
  final TextEditingController _passwordController = TextEditingController();
  final TextEditingController _confirmPasswordController =
      TextEditingController();

  @override
  void dispose() {
    _fullNameController.dispose();
    _rollNumberController.dispose();
    _courseController.dispose();
    _academicYearController.dispose();
    _phoneController.dispose();
    _passwordController.dispose();
    _confirmPasswordController.dispose();
    super.dispose();
  }

  void _register() async {
    if (_formKey.currentState!.validate()) {
      final data = {
        'full_name': _fullNameController.text,
        'roll_number': _rollNumberController.text,
        'course': _courseController.text,
        'academic_year': _academicYearController.text,
        'phone': _phoneController.text,
        'password': _passwordController.text,
      };

      try {
        final authService = AuthenticationService();
        final response = await authService.studentRegister(data);
        final userId = response.data['userId'] ?? '0';

        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('Registration successful')),
        );

        Navigator.push(
          context,
          MaterialPageRoute(
            builder:
                (context) => OTPVerificationScreen(userId: userId.toString()),
          ),
        );
      } catch (e) {
        ScaffoldMessenger.of(
          context,
        ).showSnackBar(SnackBar(content: Text('Registration failed: $e')));
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Student Registration')),
      body: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Form(
          key: _formKey,
          child: ListView(
            children: [
              TextFormField(
                controller: _fullNameController,
                decoration: const InputDecoration(labelText: 'Full Name'),
                validator:
                    (value) =>
                        (value == null || value.isEmpty)
                            ? 'Please enter your full name'
                            : null,
              ),
              const SizedBox(height: 10),
              TextFormField(
                controller: _rollNumberController,
                decoration: const InputDecoration(labelText: 'Roll Number'),
                validator:
                    (value) =>
                        (value == null || value.isEmpty)
                            ? 'Please enter your roll number'
                            : null,
              ),
              const SizedBox(height: 10),
              TextFormField(
                controller: _courseController,
                decoration: const InputDecoration(labelText: 'Course'),
                validator:
                    (value) =>
                        (value == null || value.isEmpty)
                            ? 'Please enter your course'
                            : null,
              ),
              const SizedBox(height: 10),
              TextFormField(
                controller: _academicYearController,
                decoration: const InputDecoration(labelText: 'Academic Year'),
                validator:
                    (value) =>
                        (value == null || value.isEmpty)
                            ? 'Please enter your academic year'
                            : null,
              ),
              const SizedBox(height: 10),
              TextFormField(
                controller: _phoneController,
                decoration: const InputDecoration(labelText: 'Phone'),
                validator:
                    (value) =>
                        (value == null || value.isEmpty)
                            ? 'Please enter your phone number'
                            : null,
              ),
              const SizedBox(height: 10),
              TextFormField(
                controller: _passwordController,
                decoration: const InputDecoration(labelText: 'Password'),
                obscureText: true,
                validator:
                    (value) =>
                        (value == null || value.length < 6)
                            ? 'Password must be at least 6 characters'
                            : null,
              ),
              const SizedBox(height: 10),
              TextFormField(
                controller: _confirmPasswordController,
                decoration: const InputDecoration(
                  labelText: 'Confirm Password',
                ),
                obscureText: true,
                validator:
                    (value) =>
                        value != _passwordController.text
                            ? 'Passwords do not match'
                            : null,
              ),
              const SizedBox(height: 20),
              ElevatedButton(
                onPressed: _register,
                child: const Text('Register'),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
