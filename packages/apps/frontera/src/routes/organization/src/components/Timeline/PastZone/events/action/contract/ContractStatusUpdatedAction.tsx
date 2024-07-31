import React, { useMemo } from 'react';

import { cn } from '@ui/utils/cn';
import { Action } from '@graphql/types';
import { FeaturedIcon } from '@ui/media/Icon/FeaturedIcon';
import { getMetadata } from '@organization/components/Timeline/PastZone/events/action/utils';
import { iconsByStatus } from '@organization/components/Timeline/PastZone/events/action/contract/utils';
import { useTimelineEventPreviewMethodsContext } from '@organization/components/Timeline/shared/TimelineEventPreview/context/TimelineEventPreviewContext';

interface ContractStatusUpdatedActionProps {
  data: Action;
}

export const ContractStatusUpdatedAction: React.FC<
  ContractStatusUpdatedActionProps
> = ({ data }) => {
  const { openModal } = useTimelineEventPreviewMethodsContext();
  const status = useMemo(() => {
    return getMetadata(data?.metadata)?.status?.toLowerCase();
  }, [data?.metadata]);

  // handle this situation
  if (!data.content || !status) return null;
  const content = data.content.replace(status.split('_').join(' '), '');

  return (
    <div
      onClick={() => openModal(data.id)}
      className='flex items-center cursor-pointer min-h-[40px]'
    >
      <div className='inline w-[30px]'>
        <FeaturedIcon
          size='md'
          // eslint-disable-next-line @typescript-eslint/no-explicit-any
          colorScheme={iconsByStatus[status]?.colorScheme as any}
        >
          {iconsByStatus[status]?.icon}
        </FeaturedIcon>
      </div>

      <p className='my-1 max-w-[500px] ml-2 text-sm text-gray-700 line-clamp-2'>
        {content}
        <span
          className={cn(
            status === 'renewed' ? 'font-base' : 'font-semibold',
            'ml-1',
          )}
        >
          {status.split('_').join(' ')}
        </span>
      </p>
    </div>
  );
};
