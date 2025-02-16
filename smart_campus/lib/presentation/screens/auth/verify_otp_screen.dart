import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:smart_campus/config/routes.dart';
import 'package:smart_campus/config/theme.dart';
import 'package:smart_campus/core/constants/app_constants.dart';
import 'package:smart_campus/core/utils/app_utils.dart';
import 'package:smart_campus/presentation/bloc/auth/auth_bloc.dart';
import 'package:smart_campus/presentation/bloc/auth/auth_event.dart';
import 'package:smart_campus/presentation/bloc/auth/auth_state.dart';
import 'package:smart_campus/presentation/widgets/common/loading_overlay.dart';
import 'package:smart_campus/presentation/widgets/common/pin_input.dart';

class VerifyOTPScreen extends StatefulWidget {
  final String email;

  const VerifyOTPScreen({
    super.key,
    required this.email,
  });

  @override
  State<VerifyOTPScreen> createState() => _VerifyOTPScreenState();
}

class _VerifyOTPScreenState extends State<VerifyOTPScreen> {
  final _formKey = GlobalKey<FormState>();
  final _otpController = TextEditingController();

  @override
  void dispose() {
    _otpController.dispose();
    super.dispose();
  }

  void _handleVerify() {
    if (_formKey.currentState?.validate() ?? false) {
      context.read<AuthBloc>().add(
            AuthVerifyOTPRequested(
              email: widget.email,
              otp: _otpController.text,
            ),
          );
    }
  }

  void _handleResendOTP() {
    context.read<AuthBloc>().add(
          AuthForgotPasswordRequested(widget.email),
        );
  }

  @override
  Widget build(BuildContext context) {
    return BlocListener<AuthBloc, AuthState>(
      listener: (context, state) {
        if (state is AuthFailure) {
          AppUtils.showSnackBar(context, state.message, isError: true);
        } else if (state is AuthOTPVerified) {
          Routes.navigateTo(
            context,
            Routes.resetPassword,
            arguments: {
              'email': state.email,
              'otp': state.otp,
            },
          );
        } else if (state is AuthOTPSent) {
          AppUtils.showSnackBar(
            context,
            'New OTP has been sent to your email',
            isError: false,
          );
        }
      },
      child: Scaffold(
        appBar: AppBar(
          title: const Text('Verify OTP'),
        ),
        body: LoadingOverlay(
          isLoading: context.watch<AuthBloc>().state is AuthLoading,
          child: SafeArea(
            child: SingleChildScrollView(
              padding: const EdgeInsets.all(24.0),
              child: Form(
                key: _formKey,
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.stretch,
                  children: [
                    Text(
                      'Enter Verification Code',
                      style: AppTheme.headlineLarge,
                      textAlign: TextAlign.center,
                    ),
                    const SizedBox(height: 8),
                    Text(
                      'We have sent a verification code to\n${widget.email}',
                      style: AppTheme.bodyMedium,
                      textAlign: TextAlign.center,
                    ),
                    const SizedBox(height: 48),
                    PinInput(
                      controller: _otpController,
                      length: AppConstants.otpLength,
                      onCompleted: (pin) {
                        _handleVerify();
                      },
                    ),
                    const SizedBox(height: 32),
                    ElevatedButton(
                      onPressed: _handleVerify,
                      child: const Text('Verify'),
                    ),
                    const SizedBox(height: 16),
                    Row(
                      mainAxisAlignment: MainAxisAlignment.center,
                      children: [
                        const Text("Didn't receive the code?"),
                        TextButton(
                          onPressed: _handleResendOTP,
                          child: const Text('Resend'),
                        ),
                      ],
                    ),
                  ],
                ),
              ),
            ),
          ),
        ),
      ),
    );
  }
}
