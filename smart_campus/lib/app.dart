import 'package:flutter/material.dart';
import 'package:smart_campus/config/app_config.dart';
import 'package:smart_campus/config/routes.dart';
import 'package:smart_campus/config/theme.dart';

class SmartCampusApp extends StatelessWidget {
  const SmartCampusApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: AppConfig.appName,
      theme: AppTheme.lightTheme,
      onGenerateRoute: Routes.generateRoute,
      debugShowCheckedModeBanner: false,
      builder: (context, child) {
        return MediaQuery(
          // Prevent system text scaling
          data: MediaQuery.of(context).copyWith(
            textScaler: const TextScaler.linear(1.0),
          ),
          child: child!,
        );
      },
    );
  }
}
