import { SelectOption } from '@ui/utils/types';
import { OnboardingStatus } from '@graphql/types';

export const options: SelectOption<OnboardingStatus>[] = [
  { label: 'Not applicable', value: OnboardingStatus.NotApplicable },
  { label: 'Not started', value: OnboardingStatus.NotStarted },
  { label: 'Late', value: OnboardingStatus.Late },
  { label: 'Stuck', value: OnboardingStatus.Stuck },
  { label: 'On track', value: OnboardingStatus.OnTrack },
  { label: 'Done', value: OnboardingStatus.Done },
  { label: 'Success', value: OnboardingStatus.Successful },
];
