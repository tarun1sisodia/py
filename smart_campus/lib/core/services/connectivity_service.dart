import 'dart:async';
import 'package:connectivity_plus/connectivity_plus.dart';
import 'package:injectable/injectable.dart';
import 'package:smart_campus/core/services/logger_service.dart';

@lazySingleton
class ConnectivityService {
  final Connectivity _connectivity;
  final LoggerService _logger;
  final _connectivityController = StreamController<bool>.broadcast();

  ConnectivityService(this._logger) : _connectivity = Connectivity() {
    _initConnectivityStream();
  }

  // Expose raw Connectivity instance for specific use cases
  Connectivity get connectivity => _connectivity;

  void _initConnectivityStream() {
    _connectivity.onConnectivityChanged.listen((result) {
      final isConnected = result != ConnectivityResult.none;
      _logger.info('Connectivity changed: $isConnected');
      _connectivityController.add(isConnected);
    });
  }

  Future<bool> isConnected() async {
    try {
      final result = await _connectivity.checkConnectivity();
      return result != ConnectivityResult.none;
    } catch (e) {
      _logger.error('Error checking connectivity', e);
      return false;
    }
  }

  Future<void> waitForConnection() async {
    if (await isConnected()) return;

    await _connectivityController.stream
        .firstWhere((connected) => connected)
        .timeout(const Duration(minutes: 5));
  }

  void dispose() {
    _connectivityController.close();
  }

  Future<bool> isOnline() async {
    final result = await _connectivity.checkConnectivity();
    return result != ConnectivityResult.none;
  }
}
