import { match } from 'ts-pattern';

import { cn } from '@ui/utils/cn';
import { DateTimeUtils } from '@utils/date.ts';
import { OnboardingStatus } from '@graphql/types';

interface OnboardingCellProps {
  updatedAt?: string;
  status?: OnboardingStatus;
}

const labelMap: Record<OnboardingStatus, string> = {
  [OnboardingStatus.NotApplicable]: 'Not applicable',
  [OnboardingStatus.NotStarted]: 'Not started',
  [OnboardingStatus.Successful]: 'Successful',
  [OnboardingStatus.OnTrack]: 'On track',
  [OnboardingStatus.Late]: 'Late',
  [OnboardingStatus.Stuck]: 'Stuck',
  [OnboardingStatus.Done]: 'Done',
};

export const OnboardingCell = ({
  updatedAt,
  status = OnboardingStatus.NotApplicable,
}: OnboardingCellProps) => {
  const timeElapsed = match(status)
    .with(OnboardingStatus.NotApplicable, OnboardingStatus.Successful, () => '')
    .otherwise(() => {
      if (!updatedAt) return '';

      return match(DateTimeUtils.getDifferenceFromNow(updatedAt))
        .with([null, 'today'], () => {
          const [value, unit] =
            DateTimeUtils.getDifferenceInMinutesOrHours(updatedAt);

          return `for ${Math.abs(value as number)} ${unit}`;
        })
        .otherwise(
          ([value, unit]) => `for ${Math.abs(value as number)} ${unit}`,
        );
    });

  const color = match(status)
    .returnType<string>()
    .with(
      OnboardingStatus.Successful,
      OnboardingStatus.OnTrack,
      OnboardingStatus.Done,
      () => 'text-success-500',
    )
    .with(OnboardingStatus.NotApplicable, () => 'text-gray-400')
    .with(
      OnboardingStatus.Late,
      OnboardingStatus.Stuck,
      () => 'text-warning-500',
    )
    .otherwise(() => 'gray.500');

  const label = labelMap[status];

  return (
    <div className='flex items-center gap-1'>
      <p
        className={cn(color, 'leading-none')}
        data-test='organization-onboarding-in-all-orgs-table'
      >
        {label}
      </p>
      {timeElapsed && (
        <span className='text-gray-500 text-xs leading-none'>•</span>
      )}
      <p className='text-gray-500 text-xs leading-none'>{timeElapsed}</p>
    </div>
  );
};
