import 'package:flutter/material.dart';

class SnackbarUtils {
  static void showSuccess(
    BuildContext context, {
    required String message,
    Duration? duration,
  }) {
    _show(
      context,
      message: message,
      backgroundColor: Colors.green,
      icon: const Icon(
        Icons.check_circle_outline,
        color: Colors.white,
      ),
      duration: duration,
    );
  }

  static void showError(
    BuildContext context, {
    required String message,
    Duration? duration,
  }) {
    _show(
      context,
      message: message,
      backgroundColor: Colors.red,
      icon: const Icon(
        Icons.error_outline,
        color: Colors.white,
      ),
      duration: duration,
    );
  }

  static void showInfo(
    BuildContext context, {
    required String message,
    Duration? duration,
  }) {
    _show(
      context,
      message: message,
      backgroundColor: Colors.blue,
      icon: const Icon(
        Icons.info_outline,
        color: Colors.white,
      ),
      duration: duration,
    );
  }

  static void showWarning(
    BuildContext context, {
    required String message,
    Duration? duration,
  }) {
    _show(
      context,
      message: message,
      backgroundColor: Colors.orange,
      icon: const Icon(
        Icons.warning_amber_outlined,
        color: Colors.white,
      ),
      duration: duration,
    );
  }

  static void _show(
    BuildContext context, {
    required String message,
    required Color backgroundColor,
    required Widget icon,
    Duration? duration,
  }) {
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(
        content: Row(
          children: [
            icon,
            const SizedBox(width: 8),
            Expanded(
              child: Text(
                message,
                style: const TextStyle(color: Colors.white),
              ),
            ),
          ],
        ),
        backgroundColor: backgroundColor,
        duration: duration ?? const Duration(seconds: 4),
        behavior: SnackBarBehavior.floating,
        dismissDirection: DismissDirection.horizontal,
        action: SnackBarAction(
          label: 'Dismiss',
          textColor: Colors.white,
          onPressed: () {
            ScaffoldMessenger.of(context).hideCurrentSnackBar();
          },
        ),
      ),
    );
  }
}
