import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:injectable/injectable.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:smart_campus/core/services/logger_service.dart';
import 'package:smart_campus/presentation/bloc/onboarding/onboarding_event.dart';
import 'package:smart_campus/presentation/bloc/onboarding/onboarding_state.dart';

@injectable
class OnboardingBloc extends Bloc<OnboardingEvent, OnboardingState> {
  final SharedPreferences _prefs;
  final LoggerService _logger;
  static const _onboardingCompletedKey = 'onboarding_completed';

  OnboardingBloc(
    this._prefs,
    this._logger,
  ) : super(const OnboardingState()) {
    on<OnboardingPageChanged>(_onPageChanged);
    on<OnboardingCompleted>(_onCompleted);
  }

  void _onPageChanged(
    OnboardingPageChanged event,
    Emitter<OnboardingState> emit,
  ) {
    emit(state.copyWith(
      currentPage: event.page,
      isLastPage: event.isLastPage,
    ));
  }

  Future<void> _onCompleted(
    OnboardingCompleted event,
    Emitter<OnboardingState> emit,
  ) async {
    try {
      await _prefs.setBool(_onboardingCompletedKey, true);
      _logger.info('Onboarding completed');
    } catch (e) {
      _logger.error('Error completing onboarding', e);
    }
  }

  bool isOnboardingCompleted() {
    return _prefs.getBool(_onboardingCompletedKey) ?? false;
  }
}
