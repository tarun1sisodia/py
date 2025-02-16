import 'package:flutter/material.dart';
import 'package:smart_attendance_app/services/authentication_service.dart';

class OTPVerificationScreen extends StatefulWidget {
  final String userId;
  const OTPVerificationScreen({Key? key, required this.userId})
    : super(key: key);

  @override
  _OTPVerificationScreenState createState() => _OTPVerificationScreenState();
}

class _OTPVerificationScreenState extends State<OTPVerificationScreen> {
  final _formKey = GlobalKey<FormState>();
  final TextEditingController _otpController = TextEditingController();

  @override
  void dispose() {
    _otpController.dispose();
    super.dispose();
  }

  void _verifyOTP() async {
    if (_formKey.currentState!.validate()) {
      try {
        final authService = AuthenticationService();
        final response = await authService.verifyOTP(
          widget.userId,
          _otpController.text,
        );
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('OTP verified successfully')),
        );
        // TODO: Navigate to teacher/student dashboard after successful OTP verification
        // Navigator.pushReplacementNamed(context, '/dashboard');
      } catch (e) {
        ScaffoldMessenger.of(
          context,
        ).showSnackBar(SnackBar(content: Text('OTP verification failed: $e')));
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('OTP Verification')),
      body: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Form(
          key: _formKey,
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              const Text('Enter the OTP sent to your phone:'),
              const SizedBox(height: 20),
              TextFormField(
                controller: _otpController,
                keyboardType: TextInputType.number,
                decoration: const InputDecoration(labelText: 'OTP'),
                validator: (value) {
                  if (value == null || value.trim().isEmpty) {
                    return 'Please enter OTP';
                  }
                  if (value.trim().length != 6) {
                    return 'OTP must be 6 digits';
                  }
                  return null;
                },
              ),
              const SizedBox(height: 20),
              ElevatedButton(
                onPressed: _verifyOTP,
                child: const Text('Verify'),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
