import 'package:flutter/material.dart';
import 'package:smart_attendance_app/services/authentication_service.dart';
import 'package:smart_attendance_app/presentation/screens/auth/otp_verification_screen.dart';

class TeacherRegistrationScreen extends StatefulWidget {
  const TeacherRegistrationScreen({Key? key}) : super(key: key);

  @override
  _TeacherRegistrationScreenState createState() =>
      _TeacherRegistrationScreenState();
}

class _TeacherRegistrationScreenState extends State<TeacherRegistrationScreen> {
  final _formKey = GlobalKey<FormState>();

  // Controllers for input fields
  final TextEditingController _fullNameController = TextEditingController();
  final TextEditingController _usernameController = TextEditingController();
  final TextEditingController _emailController = TextEditingController();
  final TextEditingController _phoneController = TextEditingController();
  final TextEditingController _highestDegreeController =
      TextEditingController();
  final TextEditingController _experienceController = TextEditingController();
  final TextEditingController _passwordController = TextEditingController();
  final TextEditingController _confirmPasswordController =
      TextEditingController();

  @override
  void dispose() {
    _fullNameController.dispose();
    _usernameController.dispose();
    _emailController.dispose();
    _phoneController.dispose();
    _highestDegreeController.dispose();
    _experienceController.dispose();
    _passwordController.dispose();
    _confirmPasswordController.dispose();
    super.dispose();
  }

  void _register() async {
    if (_formKey.currentState!.validate()) {
      final data = {
        'full_name': _fullNameController.text,
        'username': _usernameController.text,
        'email': _emailController.text,
        'phone': _phoneController.text,
        'highest_degree': _highestDegreeController.text,
        'experience': _experienceController.text,
        'password': _passwordController.text,
      };

      try {
        final authService = AuthenticationService();
        final response = await authService.teacherRegister(data);
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
      appBar: AppBar(title: const Text('Teacher Registration')),
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
                        value!.isEmpty ? 'Please enter your full name' : null,
              ),
              const SizedBox(height: 10),
              TextFormField(
                controller: _usernameController,
                decoration: const InputDecoration(labelText: 'Username'),
                validator:
                    (value) =>
                        value!.isEmpty ? 'Please enter a username' : null,
              ),
              const SizedBox(height: 10),
              TextFormField(
                controller: _emailController,
                decoration: const InputDecoration(labelText: 'Email'),
                validator:
                    (value) =>
                        (value == null || value.isEmpty)
                            ? 'Please enter your email'
                            : null,
              ),
              const SizedBox(height: 10),
              TextFormField(
                controller: _phoneController,
                decoration: const InputDecoration(labelText: 'Phone'),
                validator:
                    (value) =>
                        value!.isEmpty
                            ? 'Please enter your phone number'
                            : null,
              ),
              const SizedBox(height: 10),
              TextFormField(
                controller: _highestDegreeController,
                decoration: const InputDecoration(labelText: 'Highest Degree'),
                validator:
                    (value) =>
                        value!.isEmpty
                            ? 'Please enter your highest degree'
                            : null,
              ),
              const SizedBox(height: 10),
              TextFormField(
                controller: _experienceController,
                decoration: const InputDecoration(labelText: 'Experience'),
                validator:
                    (value) =>
                        value!.isEmpty
                            ? 'Please enter your teaching experience'
                            : null,
              ),
              const SizedBox(height: 10),
              TextFormField(
                controller: _passwordController,
                decoration: const InputDecoration(labelText: 'Password'),
                obscureText: true,
                validator:
                    (value) =>
                        (value != null && value.length < 6)
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
