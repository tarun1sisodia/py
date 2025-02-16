import 'package:injectable/injectable.dart';
import 'package:smart_campus/core/services/encryption_service.dart';
import 'package:smart_campus/core/services/secure_storage_service.dart';
import 'package:smart_campus/core/services/device_binding_service.dart';

@lazySingleton
class SecurityService {
  final EncryptionService _encryptionService;
  final SecureStorageService _secureStorage;
  final DeviceBindingService _deviceBinding;

  SecurityService(
    this._encryptionService,
    this._secureStorage,
    this._deviceBinding,
  );

  Future<bool> isDeviceSecure() async {
    final bindingStatus = await _deviceBinding.getBindingStatus();
    return bindingStatus['isBound'] == true &&
        !bindingStatus['isDeveloperMode'];
  }

  Future<void> storeSecureData(
      {required String key, required String value}) async {
    final encrypted = await _encryptionService.encrypt(value);
    await _secureStorage.write(key: key, value: encrypted);
  }

  Future<String?> getSecureData(String key) async {
    final encrypted = await _secureStorage.read(key);
    if (encrypted == null) return null;
    return await _encryptionService.decrypt(encrypted);
  }

  Future<void> clearSecureData() async {
    await _secureStorage.deleteAll();
  }

  Future<bool> verifyDeviceBinding(String userId) async {
    return await _deviceBinding.verifyDeviceBinding(userId);
  }

  Future<void> bindDevice(String userId) async {
    await _deviceBinding.bindDevice(userId);
  }

  Future<void> unbindDevice() async {
    await _deviceBinding.unbindDevice();
  }
}
