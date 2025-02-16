import 'package:flutter/material.dart';
import 'package:smart_campus/config/theme.dart';

class AboutScreen extends StatelessWidget {
  const AboutScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('About'),
      ),
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(16.0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            _buildHeader(context),
            const SizedBox(height: 24),
            _buildDescription(),
            const SizedBox(height: 24),
            _buildFeatures(),
            const SizedBox(height: 24),
            _buildTechnologies(),
            const SizedBox(height: 24),
            _buildDevelopers(),
          ],
        ),
      ),
    );
  }

  Widget _buildHeader(BuildContext context) {
    return Center(
      child: Column(
        children: [
          Icon(
            Icons.school,
            size: 80,
            color: Theme.of(context).colorScheme.primary,
          ),
          const SizedBox(height: 16),
          Text(
            'Smart Campus',
            style: AppTheme.headlineLarge,
          ),
          const SizedBox(height: 8),
          Text(
            'Version 1.0.0',
            style: AppTheme.bodyLarge,
          ),
        ],
      ),
    );
  }

  Widget _buildDescription() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          'About the App',
          style: AppTheme.titleLarge,
        ),
        const SizedBox(height: 16),
        Text(
          'Smart Campus is a modern attendance management system designed to '
          'streamline the process of marking and tracking attendance in '
          'educational institutions. The app uses advanced technologies like '
          'location verification and WiFi authentication to ensure accurate '
          'attendance records.',
          style: AppTheme.bodyLarge,
        ),
      ],
    );
  }

  Widget _buildFeatures() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          'Key Features',
          style: AppTheme.titleLarge,
        ),
        const SizedBox(height: 16),
        _buildFeatureItem(
          Icons.location_on,
          'Location Verification',
          'Ensures students are physically present in the classroom',
        ),
        _buildFeatureItem(
          Icons.wifi,
          'WiFi Authentication',
          'Verifies connection to campus WiFi network',
        ),
        _buildFeatureItem(
          Icons.analytics,
          'Real-time Analytics',
          'Track attendance patterns and generate reports',
        ),
        _buildFeatureItem(
          Icons.offline_bolt,
          'Offline Support',
          'Works even without internet connectivity',
        ),
      ],
    );
  }

  Widget _buildFeatureItem(IconData icon, String title, String description) {
    return Padding(
      padding: const EdgeInsets.only(bottom: 16.0),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Icon(icon, size: 24),
          const SizedBox(width: 16),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  title,
                  style:
                      AppTheme.bodyLarge.copyWith(fontWeight: FontWeight.w600),
                ),
                const SizedBox(height: 4),
                Text(
                  description,
                  style: AppTheme.bodyMedium,
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildTechnologies() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          'Technologies Used',
          style: AppTheme.titleLarge,
        ),
        const SizedBox(height: 16),
        Wrap(
          spacing: 8,
          runSpacing: 8,
          children: [
            _buildTechChip('Flutter'),
            _buildTechChip('Dart'),
            _buildTechChip('Firebase'),
            _buildTechChip('MySQL'),
            _buildTechChip('Node.js'),
            _buildTechChip('Express'),
          ],
        ),
      ],
    );
  }

  Widget _buildTechChip(String label) {
    return Chip(
      label: Text(label),
      padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 0),
    );
  }

  Widget _buildDevelopers() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          'Development Team',
          style: AppTheme.titleLarge,
        ),
        const SizedBox(height: 16),
        _buildDeveloperItem(
          'John Doe',
          'Lead Developer',
          'Responsible for app architecture and core features',
        ),
        _buildDeveloperItem(
          'Jane Smith',
          'UI/UX Designer',
          'Created the beautiful and intuitive user interface',
        ),
        _buildDeveloperItem(
          'Mike Johnson',
          'Backend Developer',
          'Developed the robust server infrastructure',
        ),
      ],
    );
  }

  Widget _buildDeveloperItem(String name, String role, String description) {
    return Padding(
      padding: const EdgeInsets.only(bottom: 16.0),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            name,
            style: AppTheme.bodyLarge.copyWith(fontWeight: FontWeight.w600),
          ),
          const SizedBox(height: 4),
          Text(
            role,
            style: AppTheme.bodyMedium.copyWith(fontWeight: FontWeight.w500),
          ),
          const SizedBox(height: 4),
          Text(
            description,
            style: AppTheme.bodyMedium,
          ),
        ],
      ),
    );
  }
}
