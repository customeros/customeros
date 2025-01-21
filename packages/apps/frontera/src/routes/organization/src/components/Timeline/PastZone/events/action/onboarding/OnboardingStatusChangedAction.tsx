import { useMemo } from 'react';

import { Flag04 } from '@ui/media/icons/Flag04';
import { Action, OnboardingStatus } from '@graphql/types';
import { getMetadata } from '@organization/components/Timeline/PastZone/events/action/utils';
import { useTimelineEventPreviewMethodsContext } from '@organization/components/Timeline/shared/TimelineEventPreview/context/TimelineEventPreviewContext';

import { getIconColor } from './util';

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
  const iconClassName = getIconColor(status);

  return (
    <div
      onClick={() => openModal(data.id)}
      className='flex items-center cursor-pointer'
    >
      <Flag04 className={iconClassName} />

      <p className='max-w-[500px] ml-2 text-sm text-gray-700 line-clamp-2'>
        {content}
        <span className='font-semibold ml-1'>{statusLabel}</span>
      </p>
    </div>
  );
};
