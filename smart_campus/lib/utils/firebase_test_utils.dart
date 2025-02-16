import 'package:firebase_auth/firebase_auth.dart';

class FirebaseTestUtils {
  static void setupTestPhoneAuth() {
    // Enable test mode for phone auth
    FirebaseAuth.instance.setSettings(
      appVerificationDisabledForTesting: true,
    );
    // In production, do not set test phone number and code
    // Ensure app verification is enabled
    const testPhoneNumber = null;
    const testVerificationCode = null;

    FirebaseAuth.instance.setSettings(
      phoneNumber: testPhoneNumber,
      smsCode: testVerificationCode,
    );
  }
}
