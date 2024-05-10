import { match } from 'ts-pattern';

import { OnboardingStatus } from '@graphql/types';

export const getColorScheme = (status: OnboardingStatus) =>
  match(status)
    .returnType<string>()
    .with(
      OnboardingStatus.Successful,
      OnboardingStatus.OnTrack,
      OnboardingStatus.Done,
      () => 'success',
    )
    .with(OnboardingStatus.Late, OnboardingStatus.Stuck, () => 'warning')
    .otherwise(() => 'gray');
