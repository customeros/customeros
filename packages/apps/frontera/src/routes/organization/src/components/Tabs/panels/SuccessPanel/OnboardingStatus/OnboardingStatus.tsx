import { useState } from 'react';

import { match } from 'ts-pattern';

import { cn } from '@ui/utils/cn';
import { DateTimeUtils } from '@utils/date.ts';
import { Flag04 } from '@ui/media/icons/Flag04';
import { Trophy01 } from '@ui/media/icons/Trophy01';
import { FeaturedIcon } from '@ui/media/Icon/FeaturedIcon';
import { useDisclosure } from '@ui/utils/hooks/useDisclosure';
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
  isLoading?: boolean;
  data?: OnboardingDetails | null;
}

export const OnboardingStatus = ({
  data,
  isLoading,
}: OnboardingStatusProps) => {
  const { open, onClose, onOpen } = useDisclosure();
  const [isFetching, setIsFetching] = useState(false);

  const handleIsFetching = (status: boolean) => setIsFetching(status);

  const timeElapsed = match(data?.status)
    .with(
      OnboardingStatusEnum.NotApplicable,
      OnboardingStatusEnum.Successful,
      () => '',
    )
    .otherwise(() => {
      if (!data?.updatedAt) return '';

      return match(DateTimeUtils.getDifferenceFromNow(data?.updatedAt))
        .with([null, 'today'], () => {
          const [value, unit] = DateTimeUtils.getDifferenceInMinutesOrHours(
            data?.updatedAt,
          );

          return `for ${Math.abs(value as number)} ${unit}`;
        })
        .otherwise(
          ([value, unit]) => `for ${Math.abs(value as number)} ${unit}`,
        );
    });

  const label =
    labelMap[data?.status ?? OnboardingStatusEnum.NotApplicable].toLowerCase();

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const colorScheme: any = match(data?.status)
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
      <div
        className={cn(
          isFetching ? 'opacity-50' : 'opacity-100',
          reason ? 'justify-start' : 'justify-center',
          isFetching ? 'animate-pulseOpacity' : 'unset',
          'flex mt-1 ml-[15px] gap-4 w-full items-center cursor-pointer overflow-visible justify-start',
        )}
        onClick={onOpen}
      >
        <FeaturedIcon colorScheme={colorScheme}>
          {data?.status === OnboardingStatusEnum.Successful ? (
            <Trophy01 />
          ) : (
            <Flag04 />
          )}
        </FeaturedIcon>
        <div className='flex-col inline-grid'>
          <div className='flex'>
            <span className='ml-1 mr-1 font-semibold'>Onboarding</span>
            <span className='text-gray-500'>{`${label} ${
              isLoading ? '' : timeElapsed
            }`}</span>
          </div>
          {reason && (
            <span className='line-clamp-2 text-gray-500 text-sm'>{`“${reason}”`}</span>
          )}
        </div>
      </div>
      {open && (
        <OnboardingStatusModal
          isOpen={open}
          onClose={onClose}
          data={data}
          onFetching={handleIsFetching}
        />
      )}
    </>
  );
};
