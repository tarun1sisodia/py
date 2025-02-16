import 'package:flutter/material.dart';
import 'package:smart_attendance_app/presentation/screens/auth/splash_screen.dart';
import 'package:smart_attendance_app/presentation/screens/auth/welcome_screen.dart';
import 'package:smart_attendance_app/presentation/screens/auth/teacher_registration_screen.dart';
import 'package:smart_attendance_app/presentation/screens/auth/teacher_login_screen.dart';
import 'package:smart_attendance_app/presentation/screens/auth/student_registration_screen.dart';
import 'package:smart_attendance_app/presentation/screens/auth/student_login_screen.dart';
import 'package:smart_attendance_app/presentation/screens/auth/reset_password_screen.dart';

class AppRoutes {
  static const String splash = '/';
  static const String welcome = '/welcome';
  static const String teacherLogin = '/teacher/login';
  static const String studentLogin = '/student/login';
  static const String teacherRegister = '/teacher/register';
  static const String studentRegister = '/student/register';
  static const String otpVerification = '/otp-verification';
  static const String teacherDashboard = '/teacher/dashboard';
  static const String studentDashboard = '/student/dashboard';
  static const String sessionManagement = '/teacher/session';
  static const String attendanceMarking = '/student/attendance';
  static const String profile = '/profile';
  static const String settings = '/settings';
  static const String resetPassword = '/reset-password';

  static Map<String, WidgetBuilder> get routes => {
    splash: (context) => const SplashScreen(),
    welcome: (context) => const WelcomeScreen(),
    teacherRegister: (context) => const TeacherRegistrationScreen(),
    teacherLogin: (context) => const TeacherLoginScreen(),
    studentRegister: (context) => const StudentRegistrationScreen(),
    studentLogin: (context) => const StudentLoginScreen(),
    resetPassword: (context) => const ResetPasswordScreen(),
    // otpVerification route is handled in onGenerateRoute since it might require parameters
  };

  static Route<dynamic> generateRoute(RouteSettings settings) {
    // Handle dynamic route parameters here
    switch (settings.name) {
      case otpVerification:
        // Example of handling route with parameters:
        // final args = settings.arguments as Map<String, dynamic>;
        // return MaterialPageRoute(
        //   builder: (_) => OtpVerificationScreen(
        //     phoneNumber: args['phoneNumber'],
        //     userType: args['userType'],
        //   ),
        // );
        return MaterialPageRoute(
          builder:
              (_) => const Scaffold(
                body: Center(
                  child: Text('OTP Verification Screen - To be implemented'),
                ),
              ),
        );
      default:
        return MaterialPageRoute(
          builder:
              (_) => Scaffold(
                body: Center(
                  child: Text('No route defined for ${settings.name}'),
                ),
              ),
        );
    }
  }
}
