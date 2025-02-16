import 'dart:async';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:smart_campus/config/app_config.dart';
import 'package:smart_campus/di/injection.dart';
import 'package:sentry_flutter/sentry_flutter.dart';
import 'package:smart_campus/presentation/screens/auth/login_screen.dart';
import 'package:smart_campus/presentation/screens/auth/register_screen.dart';
import 'package:smart_campus/presentation/screens/student/student_dashboard_screen.dart';
import 'package:smart_campus/presentation/screens/teacher/teacher_dashboard_screen.dart';
import 'package:smart_campus/presentation/screens/common/about_screen.dart';
import 'package:smart_campus/presentation/screens/common/profile_screen.dart';
import 'package:smart_campus/presentation/screens/common/settings_screen.dart';
import 'package:firebase_core/firebase_core.dart';
import 'package:smart_campus/presentation/screens/auth/phone_auth_screen.dart';
import 'package:smart_campus/utils/firebase_test_utils.dart';

void main() async {
  await runZonedGuarded(() async {
    WidgetsFlutterBinding.ensureInitialized();
    await Firebase.initializeApp();

    // Enable test mode for phone auth in development
    if (const bool.fromEnvironment('dart.vm.product') == false) {
      FirebaseTestUtils.setupTestPhoneAuth();
    }

    // Force portrait orientation
    await SystemChrome.setPreferredOrientations([
      DeviceOrientation.portraitUp,
      DeviceOrientation.portraitDown,
    ]);

    // Initialize app configuration
    AppConfig.initialize(Environment.dev);

    // Initialize dependencies
    await configureDependencies();

    runApp(const MyApp());
  }, (error, stackTrace) async {
    // Log fatal errors to Sentry
    await Sentry.captureException(
      error,
      stackTrace: stackTrace,
    );
  });
}

class MyApp extends StatelessWidget {
  const MyApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Smart Campus',
      theme: ThemeData(
        primarySwatch: Colors.blue,
        useMaterial3: true,
      ),
      routes: {
        '/': (context) => const PhoneAuthScreen(),
        // Add other routes as needed
      },
    );
  }
}

class ScreenSelector extends StatelessWidget {
  const ScreenSelector({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Screen Selector'),
      ),
      body: ListView(
        padding: const EdgeInsets.all(16.0),
        children: [
          _buildSection('Authentication Screens', [
            _buildScreenButton(
              context,
              'Login Screen',
              const LoginScreen(),
            ),
            _buildScreenButton(
              context,
              'Register Screen',
              const RegisterScreen(),
            ),
          ]),
          _buildSection('Student Screens', [
            _buildScreenButton(
              context,
              'Student Dashboard',
              const StudentDashboardScreen(),
            ),
          ]),
          _buildSection('Teacher Screens', [
            _buildScreenButton(
              context,
              'Teacher Dashboard',
              const TeacherDashboardScreen(),
            ),
          ]),
          _buildSection('Common Screens', [
            _buildScreenButton(
              context,
              'Profile Screen',
              const ProfileScreen(),
            ),
            _buildScreenButton(
              context,
              'Settings Screen',
              const SettingsScreen(),
            ),
            _buildScreenButton(
              context,
              'About Screen',
              const AboutScreen(),
            ),
          ]),
        ],
      ),
    );
  }

  Widget _buildSection(String title, List<Widget> children) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Padding(
          padding: const EdgeInsets.symmetric(vertical: 8.0),
          child: Text(
            title,
            style: const TextStyle(
              fontSize: 20,
              fontWeight: FontWeight.bold,
            ),
          ),
        ),
        Card(
          child: Padding(
            padding: const EdgeInsets.all(8.0),
            child: Column(
              children: children,
            ),
          ),
        ),
        const SizedBox(height: 16),
      ],
    );
  }

  Widget _buildScreenButton(BuildContext context, String title, Widget screen) {
    return SizedBox(
      width: double.infinity,
      child: ElevatedButton(
        style: ElevatedButton.styleFrom(
          padding: const EdgeInsets.symmetric(vertical: 12),
        ),
        onPressed: () {
          Navigator.push(
            context,
            MaterialPageRoute(builder: (context) => screen),
          );
        },
        child: Text(title),
      ),
    );
  }
}
