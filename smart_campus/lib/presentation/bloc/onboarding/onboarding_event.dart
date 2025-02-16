import 'package:equatable/equatable.dart';

abstract class OnboardingEvent extends Equatable {
  const OnboardingEvent();

  @override
  List<Object?> get props => [];
}

class OnboardingPageChanged extends OnboardingEvent {
  final int page;
  final bool isLastPage;

  const OnboardingPageChanged({
    required this.page,
    required this.isLastPage,
  });

  @override
  List<Object?> get props => [page, isLastPage];
}

class OnboardingCompleted extends OnboardingEvent {}
