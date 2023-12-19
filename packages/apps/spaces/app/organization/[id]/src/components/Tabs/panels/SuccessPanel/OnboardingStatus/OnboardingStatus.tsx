'use client';

import { match } from 'ts-pattern';

import { Flex } from '@ui/layout/Flex';
import { useDisclosure } from '@ui/utils';
import { Text } from '@ui/typography/Text';
import { Flag04 } from '@ui/media/icons/Flag04';
import { getDifferenceFromNow } from '@shared/util/date';
import { FeaturedIcon } from '@ui/media/Icon/FeaturedIcon';
import {
  OnboardingDetails,
  OnboardingStatus as OnboardingStatusEnum,
} from '@graphql/types';

import { OnboardingStatusModal } from './OnboardingStatusModal';

const labelMap: Record<OnboardingStatusEnum, string> = {
  [OnboardingStatusEnum.NotApplicable]: 'Not applicable',
  [OnboardingStatusEnum.NotStarted]: 'Not started',
  [OnboardingStatusEnum.Successful]: 'Successful',
  [OnboardingStatusEnum.OnTrack]: 'On track',
  [OnboardingStatusEnum.Late]: 'Late',
  [OnboardingStatusEnum.Stuck]: 'Stuck',
  [OnboardingStatusEnum.Done]: 'Done',
};

interface OnboardingStatusProps {
  data?: OnboardingDetails | null;
}

export const OnboardingStatus = ({ data }: OnboardingStatusProps) => {
  const { isOpen, onOpen, onClose } = useDisclosure();

  const timeElapsed = match(data?.status)
    .returnType<string>()
    .with(
      OnboardingStatusEnum.NotApplicable,
      OnboardingStatusEnum.Successful,
      () => '',
    )
    .otherwise(() => {
      const [value, unit] = getDifferenceFromNow(data?.updatedAt);
      if (value === '0' && unit === 'days') return 'for today';

      return `for ${value} ${unit}`;
    });

  const label =
    labelMap[data?.status ?? OnboardingStatusEnum.NotApplicable].toLowerCase();

  const colorScheme = match(data?.status)
    .returnType<string>()
    .with(
      OnboardingStatusEnum.Successful,
      OnboardingStatusEnum.OnTrack,
      OnboardingStatusEnum.Done,
      () => 'success',
    )
    .with(
      OnboardingStatusEnum.Late,
      OnboardingStatusEnum.Stuck,
      () => 'warning',
    )
    .otherwise(() => 'gray');

  const reason = data?.comments;

  return (
    <>
      <Flex
        mt='1'
        gap='4'
        w='full'
        align='center'
        onClick={onOpen}
        cursor='pointer'
        overflow='visible'
        justify='flex-start'
      >
        <FeaturedIcon colorScheme={colorScheme}>
          <Flag04 />
        </FeaturedIcon>

        <Flex flexDir='column'>
          <Flex>
            <Text mr='1' fontWeight='semibold'>
              Onboarding
            </Text>
            <Text color='gray.500'>{`${label} ${timeElapsed}`}</Text>
          </Flex>
          {reason && (
            <Text color='gray.500' fontSize='sm'>{`“${reason}”`}</Text>
          )}
        </Flex>
      </Flex>

      <OnboardingStatusModal data={data} isOpen={isOpen} onClose={onClose} />
    </>
  );
};
