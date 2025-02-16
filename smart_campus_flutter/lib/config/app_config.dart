enum Environment { dev, staging, prod }

class AppConfig {
  static late final Environment environment;
  static late final String apiBaseUrl;
  static late final bool enableLogging;

  static void initialize(Environment env) {
    environment = env;

    switch (env) {
      case Environment.dev:
        apiBaseUrl = 'http://localhost:8080';
        enableLogging = true;
        break;
      case Environment.staging:
        apiBaseUrl = 'https://staging-api.smartcampus.com';
        enableLogging = true;
        break;
      case Environment.prod:
        apiBaseUrl = 'https://api.smartcampus.com';
        enableLogging = false;
        break;
    }
  }

  static bool get isDevelopment => environment == Environment.dev;
  static bool get isStaging => environment == Environment.staging;
  static bool get isProduction => environment == Environment.prod;
}
