import 'dart:async';
import 'package:geolocator/geolocator.dart';
import 'package:injectable/injectable.dart';
import 'package:smart_campus/core/errors/app_error.dart';
import 'package:smart_campus/core/services/logger_service.dart';

@lazySingleton
class LocationService {
  final LoggerService _logger;

  LocationService(this._logger);

  Future<bool> isLocationServiceEnabled() async {
    return await Geolocator.isLocationServiceEnabled();
  }

  Future<LocationPermission> checkPermission() async {
    return await Geolocator.checkPermission();
  }

  Future<LocationPermission> requestPermission() async {
    return await Geolocator.requestPermission();
  }

  Future<Position> getCurrentPosition({
    LocationAccuracy desiredAccuracy = LocationAccuracy.high,
    bool forceAndroidLocationManager = false,
    Duration? timeLimit,
  }) async {
    try {
      // Check if location services are enabled
      final serviceEnabled = await isLocationServiceEnabled();
      if (!serviceEnabled) {
        throw AppError.locationService('Location services are disabled');
      }

      // Check location permission
      var permission = await checkPermission();
      if (permission == LocationPermission.denied) {
        permission = await requestPermission();
        if (permission == LocationPermission.denied) {
          throw AppError.locationService('Location permission denied');
        }
      }

      if (permission == LocationPermission.deniedForever) {
        throw AppError.locationService(
          'Location permissions are permanently denied',
        );
      }

      // Get current position
      final position = await Geolocator.getCurrentPosition(
        desiredAccuracy: desiredAccuracy,
        forceAndroidLocationManager: forceAndroidLocationManager,
        timeLimit: timeLimit,
      );

      _logger.info('Current position: $position');
      return position;
    } catch (e) {
      _logger.error('Error getting current position', e);
      if (e is AppError) {
        rethrow;
      }
      throw AppError.locationService(e.toString());
    }
  }

  Future<bool> isWithinRadius({
    required double targetLatitude,
    required double targetLongitude,
    required double currentLatitude,
    required double currentLongitude,
    required int radiusInMeters,
  }) async {
    try {
      final distance = Geolocator.distanceBetween(
        targetLatitude,
        targetLongitude,
        currentLatitude,
        currentLongitude,
      );

      final isWithin = distance <= radiusInMeters;
      _logger.info(
        'Distance: ${distance}m, Radius: ${radiusInMeters}m, IsWithin: $isWithin',
      );
      return isWithin;
    } catch (e) {
      _logger.error('Error calculating distance', e);
      throw AppError.locationService('Error calculating distance: $e');
    }
  }

  Stream<Position> getPositionStream({
    LocationAccuracy desiredAccuracy = LocationAccuracy.high,
    int distanceFilter = 0,
    bool forceAndroidLocationManager = false,
  }) {
    try {
      return Geolocator.getPositionStream(
        locationSettings: LocationSettings(
          accuracy: desiredAccuracy,
          distanceFilter: distanceFilter,
        ),
      );
    } catch (e) {
      _logger.error('Error getting position stream', e);
      throw AppError.locationService('Error getting position stream: $e');
    }
  }

  Future<double> calculateDistance({
    required double startLatitude,
    required double startLongitude,
    required double endLatitude,
    required double endLongitude,
  }) async {
    try {
      return Geolocator.distanceBetween(
        startLatitude,
        startLongitude,
        endLatitude,
        endLongitude,
      );
    } catch (e) {
      _logger.error('Error calculating distance', e);
      throw AppError.locationService('Error calculating distance: $e');
    }
  }

  /// Verifies if the user's location is valid for attendance marking
  Future<bool> verifyAttendanceLocation({
    required double sessionLatitude,
    required double sessionLongitude,
    required int allowedRadius,
    Duration? timeLimit,
  }) async {
    try {
      final position = await getCurrentPosition(
        desiredAccuracy: LocationAccuracy.high,
        timeLimit: timeLimit ?? const Duration(seconds: 10),
      );

      final isWithinRange = await isWithinRadius(
        targetLatitude: sessionLatitude,
        targetLongitude: sessionLongitude,
        currentLatitude: position.latitude,
        currentLongitude: position.longitude,
        radiusInMeters: allowedRadius,
      );

      if (!isWithinRange) {
        _logger.warning(
          'User location (${position.latitude}, ${position.longitude}) '
          'is outside the allowed radius of $allowedRadius meters '
          'from session location ($sessionLatitude, $sessionLongitude)',
        );
      }

      return isWithinRange;
    } catch (e) {
      _logger.error('Error verifying attendance location', e);
      throw AppError.locationService('Error verifying attendance location: $e');
    }
  }

  /// Starts monitoring user's location for attendance session
  Stream<bool> monitorAttendanceLocation({
    required double sessionLatitude,
    required double sessionLongitude,
    required int allowedRadius,
  }) {
    return getPositionStream(
      desiredAccuracy: LocationAccuracy.high,
      distanceFilter: 5, // Update every 5 meters
    ).map((position) {
      final distance = Geolocator.distanceBetween(
        sessionLatitude,
        sessionLongitude,
        position.latitude,
        position.longitude,
      );

      final isWithinRange = distance <= allowedRadius;

      if (!isWithinRange) {
        _logger.warning(
          'User moved outside the allowed radius. '
          'Distance: ${distance}m, Allowed: ${allowedRadius}m',
        );
      }

      return isWithinRange;
    });
  }
}
