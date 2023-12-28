import { match } from 'ts-pattern';

import { Flex } from '@ui/layout/Flex';
import { Text } from '@ui/typography/Text';
import { OnboardingStatus } from '@graphql/types';
import {
  getDifferenceFromNow,
  getDifferenceInMinutesOrHours,
} from '@shared/util/date';

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

      return match(getDifferenceFromNow(updatedAt))
        .with(['0', 'days'], () => {
          const [value, unit] = getDifferenceInMinutesOrHours(updatedAt);

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
      () => 'success.500',
    )
    .with(OnboardingStatus.Late, OnboardingStatus.Stuck, () => 'warning.500')
    .otherwise(() => 'gray.500');

  const label = labelMap[status];

  return (
    <Flex flexDir='column'>
      <Text color={color}>{label}</Text>
      <Text color='gray.500'>{timeElapsed}</Text>
    </Flex>
  );
};
