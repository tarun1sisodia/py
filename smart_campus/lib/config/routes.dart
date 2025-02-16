import 'package:flutter/material.dart';
import 'package:smart_campus/presentation/screens/auth/login_screen.dart';
import 'package:smart_campus/presentation/screens/onboarding/onboarding_screen.dart';
import 'package:smart_campus/presentation/screens/splash_screen.dart';
import 'package:smart_campus/presentation/screens/teacher/teacher_dashboard_screen.dart';
import 'package:smart_campus/presentation/screens/teacher/create_session_screen.dart';
import 'package:smart_campus/presentation/screens/debug/sentry_test_screen.dart';

class Routes {
  // Core Routes
  static const String initial = '/';
  static const String onboarding = '/onboarding';

  // Auth Routes
  static const String login = '/login';
  static const String register = '/register';
  static const String verifyOtp = '/verify-otp';
  static const String forgotPassword = '/forgot-password';
  static const String resetPassword = '/reset-password';
  static const String deviceBinding = '/device-binding';

  // Teacher Routes
  static const String teacherDashboard = '/teacher/dashboard';
  static const String createSession = '/teacher/create-session';
  static const String sessionDetails = '/teacher/session-details';
  static const String attendanceReport = '/teacher/attendance-report';

  // Student Routes
  static const String studentDashboard = '/student/dashboard';
  static const String markAttendance = '/student/mark-attendance';
  static const String attendanceHistory = '/student/attendance-history';

  // Common Routes
  static const String profile = '/profile';
  static const String settings = '/settings';
  static const String syncStatus = '/sync-status';
  static const String about = '/about';

  // Debug Routes
  static const String sentryTest = '/debug/sentry-test';

  static Route<dynamic> generateRoute(RouteSettings settings) {
    switch (settings.name) {
      case initial:
        return MaterialPageRoute(builder: (_) => const SplashScreen());
      case onboarding:
        return MaterialPageRoute(builder: (_) => const OnboardingScreen());
      case login:
        return MaterialPageRoute(builder: (_) => const LoginScreen());
      case teacherDashboard:
        return MaterialPageRoute(
            builder: (_) => const TeacherDashboardScreen());
      case createSession:
        return MaterialPageRoute(builder: (_) => const CreateSessionScreen());
      case sentryTest:
        return MaterialPageRoute(builder: (_) => const SentryTestScreen());
      // Add other route handlers as screens are implemented
      default:
        return MaterialPageRoute(
          builder: (_) => Scaffold(
            body: Center(
              child: Text('No route defined for ${settings.name}'),
            ),
          ),
        );
    }
  }

  // Navigation helper methods
  static void navigateTo(BuildContext context, String routeName,
      {Object? arguments}) {
    Navigator.pushNamed(context, routeName, arguments: arguments);
  }

  static void navigateToAndRemove(BuildContext context, String routeName,
      {Object? arguments}) {
    Navigator.pushNamedAndRemoveUntil(
      context,
      routeName,
      (route) => false,
      arguments: arguments,
    );
  }

  static void navigateToAndReplace(BuildContext context, String routeName,
      {Object? arguments}) {
    Navigator.pushReplacementNamed(
      context,
      routeName,
      arguments: arguments,
    );
  }

  static void pop(BuildContext context, [dynamic result]) {
    Navigator.pop(context, result);
  }
}
