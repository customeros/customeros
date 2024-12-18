import { match } from 'ts-pattern';
import { observer } from 'mobx-react-lite';

import { cn } from '@ui/utils/cn';
import { DateTimeUtils } from '@utils/date.ts';
import { Flag04 } from '@ui/media/icons/Flag04';
import { useStore } from '@shared/hooks/useStore';
import { Trophy01 } from '@ui/media/icons/Trophy01';
import { FeaturedIcon } from '@ui/media/Icon/FeaturedIcon';
import { useDisclosure } from '@ui/utils/hooks/useDisclosure';
import { OnboardingStatus as OnboardingStatusEnum } from '@graphql/types';

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
  id: string;
}

export const OnboardingStatus = observer(({ id }: OnboardingStatusProps) => {
  const store = useStore();
  const organization = store.organizations.value.get(id);
  const onboardingStatus = organization?.value?.onboardingStatus;
  const onboardingComments = organization?.value?.onboardingComments;
  const onboardingUpdatedAt = organization?.value?.onboardingStatusUpdatedAt;

  const { open, onClose, onOpen } = useDisclosure();

  const timeElapsed = match(onboardingStatus)
    .with(
      OnboardingStatusEnum.NotApplicable,
      OnboardingStatusEnum.Successful,
      () => '',
    )
    .otherwise(() => {
      if (!onboardingUpdatedAt) return '';

      return match(DateTimeUtils.getDifferenceFromNow(onboardingUpdatedAt))
        .with([null, 'today'], () => {
          const [value, unit] =
            DateTimeUtils.getDifferenceInMinutesOrHours(onboardingUpdatedAt);

          return `for ${Math.abs(value as number)} ${unit}`;
        })
        .otherwise(
          ([value, unit]) => `for ${Math.abs(value as number)} ${unit}`,
        );
    });

  const label =
    labelMap[
      organization?.value?.onboardingStatus ??
        OnboardingStatusEnum.NotApplicable
    ].toLowerCase();

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const colorScheme: any = match(onboardingStatus)
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

  return (
    <>
      <div
        onClick={onOpen}
        className={cn(
          onboardingComments ? 'justify-start' : 'justify-center',
          'flex mt-1 ml-[15px] gap-4 w-full items-center cursor-pointer overflow-visible justify-start opacity-100',
        )}
      >
        <FeaturedIcon colorScheme={colorScheme}>
          {onboardingStatus === OnboardingStatusEnum.Successful ? (
            <Trophy01 />
          ) : (
            <Flag04 />
          )}
        </FeaturedIcon>
        <div className='flex-col inline-grid'>
          <div className='flex'>
            <span className='ml-1 mr-1 font-semibold'>Onboarding</span>
            <span className='text-gray-500'>{`${label} ${
              store.organizations?.isLoading ? '' : timeElapsed
            }`}</span>
          </div>
          {onboardingComments && (
            <span className='line-clamp-2 text-gray-500 text-sm'>{`“${onboardingComments}”`}</span>
          )}
        </div>
      </div>
      {open && <OnboardingStatusModal isOpen={open} onClose={onClose} />}
    </>
  );
});
