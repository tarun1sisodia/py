enum Environment { dev, staging, prod }

class AppConfig {
  static late final Environment environment;
  static late final String apiBaseUrl;
  static late final String appName;

  static void initialize(Environment env) {
    environment = env;

    switch (environment) {
      case Environment.dev:
        apiBaseUrl = 'http://localhost:8080/api';
        appName = 'Smart Campus (Dev)';
        break;
      case Environment.staging:
        apiBaseUrl = 'https://staging-api.smartcampus.com/api';
        appName = 'Smart Campus (Staging)';
        break;
      case Environment.prod:
        apiBaseUrl = 'https://api.smartcampus.com/api';
        appName = 'Smart Campus';
        break;
    }
  }

  static bool get isDevelopment => environment == Environment.dev;
  static bool get isStaging => environment == Environment.staging;
  static bool get isProduction => environment == Environment.prod;

  // API Endpoints
  static String get loginEndpoint => '$apiBaseUrl/auth/login';
  static String get registerEndpoint => '$apiBaseUrl/auth/register';
  static String get refreshTokenEndpoint => '$apiBaseUrl/auth/refresh';
  static String get logoutEndpoint => '$apiBaseUrl/auth/logout';

  // Session Endpoints
  static String get createSessionEndpoint => '$apiBaseUrl/sessions/create';
  static String get activeSessionsEndpoint => '$apiBaseUrl/sessions/active';
  static String get markAttendanceEndpoint => '$apiBaseUrl/attendance/mark';

  // Feature Flags
  static const bool enableOfflineMode = true;
  static const bool enableBiometricAuth = true;
  static const bool enableLocationVerification = true;
  static const bool enableWifiVerification = true;

  // Timeouts
  static const int connectionTimeout = 30000; // 30 seconds
  static const int receiveTimeout = 30000; // 30 seconds

  // Session Configuration
  static const int defaultSessionDuration = 300; // 5 minutes
  static const int locationUpdateInterval = 60; // 1 minute
  static const int maxRetryAttempts = 3;

  // Cache Configuration
  static const int cacheDuration = 7; // 7 days

  // Location Configuration
  static const int defaultLocationRadius = 100; // meters
  static const double defaultLocationAccuracy = 50.0; // meters

  // Debug Configuration
  static bool get enableLogging => !isProduction;
  static bool get enableCrashlytics => isProduction;
}
