import 'package:flutter/material.dart';
import 'package:flutter_spinkit/flutter_spinkit.dart';

enum LoadingStyle {
  threeBounce,
  circle,
  wave,
  doubleBounce,
  fadingCircle,
}

class LoadingIndicator extends StatelessWidget {
  final LoadingStyle style;
  final Color? color;
  final double size;

  const LoadingIndicator({
    super.key,
    this.style = LoadingStyle.threeBounce,
    this.color,
    this.size = 40.0,
  });

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final effectiveColor = color ?? theme.colorScheme.primary;

    switch (style) {
      case LoadingStyle.threeBounce:
        return SpinKitThreeBounce(
          color: effectiveColor,
          size: size,
        );
      case LoadingStyle.circle:
        return SpinKitCircle(
          color: effectiveColor,
          size: size,
        );
      case LoadingStyle.wave:
        return SpinKitWave(
          color: effectiveColor,
          size: size,
        );
      case LoadingStyle.doubleBounce:
        return SpinKitDoubleBounce(
          color: effectiveColor,
          size: size,
        );
      case LoadingStyle.fadingCircle:
        return SpinKitFadingCircle(
          color: effectiveColor,
          size: size,
        );
    }
  }
}

class LoadingOverlay extends StatelessWidget {
  final bool isLoading;
  final Widget child;
  final Color? barrierColor;
  final LoadingStyle loadingStyle;

  const LoadingOverlay({
    super.key,
    required this.isLoading,
    required this.child,
    this.barrierColor,
    this.loadingStyle = LoadingStyle.threeBounce,
  });

  @override
  Widget build(BuildContext context) {
    return Stack(
      children: [
        child,
        if (isLoading)
          Container(
            color: barrierColor ?? Colors.black26,
            child: Center(
              child: LoadingIndicator(
                style: loadingStyle,
              ),
            ),
          ),
      ],
    );
  }
}
