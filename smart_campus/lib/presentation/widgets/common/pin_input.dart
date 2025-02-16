import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:smart_campus/config/theme.dart';

class PinInput extends StatelessWidget {
  final TextEditingController controller;
  final int length;
  final void Function(String)? onCompleted;
  final double spacing;
  final double size;

  const PinInput({
    super.key,
    required this.controller,
    required this.length,
    this.onCompleted,
    this.spacing = 8.0,
    this.size = 50.0,
  });

  @override
  Widget build(BuildContext context) {
    return Form(
      child: Row(
        mainAxisAlignment: MainAxisAlignment.center,
        children: List.generate(
          length,
          (index) => Container(
            width: size,
            height: size,
            margin: EdgeInsets.symmetric(horizontal: spacing / 2),
            child: TextFormField(
              controller: TextEditingController(
                text: controller.text.length > index
                    ? controller.text[index]
                    : '',
              ),
              onChanged: (value) {
                if (value.isNotEmpty) {
                  if (index < length - 1) {
                    FocusScope.of(context).nextFocus();
                  } else {
                    FocusScope.of(context).unfocus();
                  }

                  // Update the main controller
                  final currentText = controller.text;
                  if (currentText.length <= index) {
                    controller.text = currentText + value;
                  } else {
                    controller.text = currentText.substring(0, index) +
                        value +
                        currentText.substring(index + 1);
                  }

                  // Check if completed
                  if (controller.text.length == length) {
                    onCompleted?.call(controller.text);
                  }
                } else if (value.isEmpty && index > 0) {
                  FocusScope.of(context).previousFocus();
                  // Update the main controller
                  final currentText = controller.text;
                  if (currentText.isNotEmpty) {
                    controller.text = currentText.substring(0, index);
                  }
                }
              },
              decoration: InputDecoration(
                counterText: '',
                border: OutlineInputBorder(
                  borderRadius: BorderRadius.circular(12),
                  borderSide: const BorderSide(color: AppTheme.primaryColor),
                ),
                focusedBorder: OutlineInputBorder(
                  borderRadius: BorderRadius.circular(12),
                  borderSide: const BorderSide(
                    color: AppTheme.primaryColor,
                    width: 2,
                  ),
                ),
              ),
              style: AppTheme.headlineMedium,
              keyboardType: TextInputType.number,
              textAlign: TextAlign.center,
              inputFormatters: [
                LengthLimitingTextInputFormatter(1),
                FilteringTextInputFormatter.digitsOnly,
              ],
            ),
          ),
        ),
      ),
    );
  }
}
