import 'package:flutter/material.dart';
import 'package:smart_attendance_app/config/routes.dart';

class WelcomeScreen extends StatelessWidget {
  const WelcomeScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: SafeArea(
        child: Padding(
          padding: const EdgeInsets.symmetric(horizontal: 24.0),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.stretch,
            children: [
              const Spacer(),
              const Icon(Icons.school, size: 100, color: Colors.blue),
              const SizedBox(height: 32),
              Text(
                'Welcome to\nSmart Attendance',
                style: Theme.of(context).textTheme.headlineMedium?.copyWith(
                  fontWeight: FontWeight.bold,
                ),
                textAlign: TextAlign.center,
              ),
              const SizedBox(height: 16),
              Text(
                'Choose your role to continue',
                style: Theme.of(
                  context,
                ).textTheme.bodyLarge?.copyWith(color: Colors.grey[600]),
                textAlign: TextAlign.center,
              ),
              const Spacer(),
              _RoleCard(
                title: 'Teacher',
                icon: Icons.person,
                onLogin:
                    () => Navigator.pushNamed(context, AppRoutes.teacherLogin),
                onRegister:
                    () =>
                        Navigator.pushNamed(context, AppRoutes.teacherRegister),
              ),
              const SizedBox(height: 16),
              _RoleCard(
                title: 'Student',
                icon: Icons.school,
                onLogin:
                    () => Navigator.pushNamed(context, AppRoutes.studentLogin),
                onRegister:
                    () =>
                        Navigator.pushNamed(context, AppRoutes.studentRegister),
              ),
              const Spacer(),
            ],
          ),
        ),
      ),
    );
  }
}

class _RoleCard extends StatelessWidget {
  final String title;
  final IconData icon;
  final VoidCallback onLogin;
  final VoidCallback onRegister;

  const _RoleCard({
    required this.title,
    required this.icon,
    required this.onLogin,
    required this.onRegister,
  });

  @override
  Widget build(BuildContext context) {
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Column(
          children: [
            Icon(icon, size: 48, color: Theme.of(context).colorScheme.primary),
            const SizedBox(height: 16),
            Text(
              title,
              style: Theme.of(
                context,
              ).textTheme.titleLarge?.copyWith(fontWeight: FontWeight.bold),
            ),
            const SizedBox(height: 16),
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceEvenly,
              children: [
                Expanded(
                  child: ElevatedButton(
                    onPressed: onLogin,
                    child: const Text('Login'),
                  ),
                ),
                const SizedBox(width: 16),
                Expanded(
                  child: OutlinedButton(
                    onPressed: onRegister,
                    child: const Text('Register'),
                  ),
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }
}
