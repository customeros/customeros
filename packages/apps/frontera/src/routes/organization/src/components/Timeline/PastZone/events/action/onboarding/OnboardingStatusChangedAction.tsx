import { useMemo } from 'react';

import { Flag04 } from '@ui/media/icons/Flag04';
import { Action, OnboardingStatus } from '@graphql/types';
import { FeaturedIcon } from '@ui/media/Icon/FeaturedIcon';
import { getMetadata } from '@organization/components/Timeline/PastZone/events/action/utils';
import { useTimelineEventPreviewMethodsContext } from '@organization/components/Timeline/shared/TimelineEventPreview/context/TimelineEventPreviewContext';

import { getColorScheme } from './util';

const statusMap = {
  [OnboardingStatus.Late]: 'Late',
  [OnboardingStatus.OnTrack]: 'On track',
  [OnboardingStatus.Done]: 'Done',
  [OnboardingStatus.Stuck]: 'Stuck',
  [OnboardingStatus.NotStarted]: 'Not started',
  [OnboardingStatus.NotApplicable]: 'Not applicable',
  [OnboardingStatus.Successful]: 'Successful',
};

interface OnboardingStatusChangedActionProps {
  data: Action;
}

export const OnboardingStatusChangedAction = ({
  data,
}: OnboardingStatusChangedActionProps) => {
  const { openModal } = useTimelineEventPreviewMethodsContext();
  const status = useMemo(() => {
    return getMetadata(data?.metadata)?.status;
  }, [data?.metadata]) as OnboardingStatus;

  // handle this situation
  if (!data.content || !status) return null;

  const statusLabel = statusMap[status];
  const content = data?.content.replace(statusLabel, '').trimEnd();
  const colorScheme = getColorScheme(status);

  return (
    <div
      onClick={() => openModal(data.id)}
      className='flex items-center cursor-pointer min-h-[40px]'
    >
      <div className='inline w-[30px]'>
        <FeaturedIcon
          size='md'
          // eslint-disable-next-line @typescript-eslint/no-explicit-any
          colorScheme={colorScheme as any}
        >
          <Flag04 />
        </FeaturedIcon>
      </div>

      <p className='my-1 max-w-[500px] ml-2 text-sm text-gray-700 line-clamp-2'>
        {content}
        <span className='font-semibold ml-1'>{statusLabel}</span>
      </p>
    </div>
  );
};
