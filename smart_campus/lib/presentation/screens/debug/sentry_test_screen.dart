import 'package:flutter/material.dart';
import 'package:sentry_flutter/sentry_flutter.dart';
import 'package:smart_campus/core/services/error_reporting_service.dart';
import 'package:smart_campus/di/injection.dart';

class SentryTestScreen extends StatefulWidget {
  const SentryTestScreen({super.key});

  @override
  State<SentryTestScreen> createState() => _SentryTestScreenState();
}

class _SentryTestScreenState extends State<SentryTestScreen> {
  final _errorReporting = getIt<ErrorReportingService>();
  bool _isPerformanceMonitoring = false;
  ISentrySpan? _currentTransaction;

  @override
  void initState() {
    super.initState();
    _setupTestUser();
  }

  Future<void> _setupTestUser() async {
    await _errorReporting.setUser(
      id: 'test-user-123',
      email: 'test@example.com',
      username: 'TestUser',
      data: {
        'role': 'tester',
        'testMode': true,
      },
    );
  }

  Future<void> _testError() async {
    try {
      await _errorReporting.addBreadcrumb(
        message: 'User initiated test error',
        category: 'test',
        data: {'buttonPressed': 'testError'},
      );

      // Simulate an error
      throw Exception('This is a test error from SentryTestScreen');
    } catch (error, stackTrace) {
      await _errorReporting.reportError(
        error,
        stackTrace,
        hint: 'Test error triggered manually',
        extras: {
          'testTimestamp': DateTime.now().toIso8601String(),
          'testType': 'manual',
        },
      );

      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(
            content: Text('Test error sent to Sentry'),
            backgroundColor: Colors.green,
          ),
        );
      }
    }
  }

  Future<void> _testUnhandledError() async {
    await _errorReporting.addBreadcrumb(
      message: 'User initiated unhandled error test',
      category: 'test',
      data: {'buttonPressed': 'testUnhandledError'},
    );

    // This will crash the app and be caught by the zone guard
    throw Exception('This is an unhandled test error');
  }

  Future<void> _togglePerformanceMonitoring() async {
    setState(() => _isPerformanceMonitoring = !_isPerformanceMonitoring);

    if (_isPerformanceMonitoring) {
      _currentTransaction = Sentry.startTransaction(
        'test-transaction',
        'test',
        bindToScope: true,
      );

      await _errorReporting.addBreadcrumb(
        message: 'Started performance monitoring',
        category: 'performance',
      );

      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(
            content: Text('Performance monitoring started'),
            backgroundColor: Colors.blue,
          ),
        );
      }
    } else if (_currentTransaction != null) {
      await _errorReporting.addBreadcrumb(
        message: 'Finished performance monitoring',
        category: 'performance',
        data: {
          'duration': _currentTransaction!.endTimestamp
              ?.difference(_currentTransaction!.startTimestamp)
              .inMilliseconds,
        },
      );

      _currentTransaction!.finish();
      _currentTransaction = null;

      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(
            content: Text('Performance data sent to Sentry'),
            backgroundColor: Colors.green,
          ),
        );
      }
    }
  }

  Future<void> _simulateUserAction() async {
    if (!_isPerformanceMonitoring || _currentTransaction == null) return;

    final span = _currentTransaction!.startChild(
      'user-action',
      description: 'Simulated user action',
    );

    // Simulate some work
    await Future.delayed(const Duration(seconds: 2));

    await _errorReporting.addBreadcrumb(
      message: 'User action completed',
      category: 'performance',
      data: {'actionType': 'simulation'},
    );

    span.finish();

    if (mounted) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(
          content: Text('User action recorded'),
          backgroundColor: Colors.blue,
        ),
      );
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Sentry Test'),
      ),
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(16.0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            const Card(
              child: Padding(
                padding: EdgeInsets.all(16.0),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      'Sentry Test Panel',
                      style: TextStyle(
                        fontSize: 20,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                    SizedBox(height: 8),
                    Text(
                      'Use these controls to test Sentry integration',
                      style: TextStyle(
                        color: Colors.grey,
                      ),
                    ),
                  ],
                ),
              ),
            ),
            const SizedBox(height: 16),
            ElevatedButton.icon(
              onPressed: _testError,
              icon: const Icon(Icons.bug_report),
              label: const Text('Test Handled Error'),
            ),
            const SizedBox(height: 8),
            ElevatedButton.icon(
              onPressed: _testUnhandledError,
              icon: const Icon(Icons.warning),
              label: const Text('Test Unhandled Error'),
              style: ElevatedButton.styleFrom(
                backgroundColor: Colors.red,
                foregroundColor: Colors.white,
              ),
            ),
            const SizedBox(height: 16),
            ElevatedButton.icon(
              onPressed: _togglePerformanceMonitoring,
              icon: Icon(
                  _isPerformanceMonitoring ? Icons.stop : Icons.play_arrow),
              label: Text(
                _isPerformanceMonitoring
                    ? 'Stop Performance Monitoring'
                    : 'Start Performance Monitoring',
              ),
              style: ElevatedButton.styleFrom(
                backgroundColor:
                    _isPerformanceMonitoring ? Colors.orange : Colors.blue,
                foregroundColor: Colors.white,
              ),
            ),
            if (_isPerformanceMonitoring) ...[
              const SizedBox(height: 8),
              ElevatedButton.icon(
                onPressed: _simulateUserAction,
                icon: const Icon(Icons.psychology),
                label: const Text('Simulate User Action'),
                style: ElevatedButton.styleFrom(
                  backgroundColor: Colors.green,
                  foregroundColor: Colors.white,
                ),
              ),
            ],
          ],
        ),
      ),
    );
  }

  @override
  void dispose() {
    _currentTransaction?.finish();
    super.dispose();
  }
}
