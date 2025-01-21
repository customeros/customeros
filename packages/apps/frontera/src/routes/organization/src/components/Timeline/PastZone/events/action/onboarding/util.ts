import { match } from 'ts-pattern';

import { OnboardingStatus } from '@graphql/types';

export const getIconColor = (status: OnboardingStatus) =>
  match(status)
    .returnType<string>()
    .with(
      OnboardingStatus.Successful,
      OnboardingStatus.OnTrack,
      OnboardingStatus.Done,
      () => 'text-success-500',
    )
    .with(
      OnboardingStatus.Late,
      OnboardingStatus.Stuck,
      () => 'text-warning-500',
    )
    .otherwise(() => 'text-gray-500');
