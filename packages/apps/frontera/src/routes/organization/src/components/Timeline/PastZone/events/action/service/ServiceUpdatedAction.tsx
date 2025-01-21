import React from 'react';

import { cn } from '@ui/utils/cn';
import { Action } from '@graphql/types';
import { XCircle } from '@ui/media/icons/XCircle';
import { DotSingle } from '@ui/media/icons/DotSingle';
import { getMetadata } from '@organization/components/Timeline/PastZone/events/action/utils';
import { useTimelineEventPreviewMethodsContext } from '@organization/components/Timeline/shared/TimelineEventPreview/context/TimelineEventPreviewContext';

import { formatString } from './utils.tsx';

interface ServiceUpdatedActionProps {
  data: Action;
  mode?: 'created' | 'updated' | 'removed';
}

export const ServiceUpdatedAction = ({
  data,
  mode = 'updated',
}: ServiceUpdatedActionProps) => {
  const { openModal } = useTimelineEventPreviewMethodsContext();
  const metadata = getMetadata(data?.metadata);

  if (!data.content) return null;
  const isTemporary = data.appSource === 'customeros-optimistic-update';
  const formattedContent = formatString(
    data.content,
    metadata?.billedType,
    metadata?.currency ?? 'USD',
  );

  return (
    <div
      onClick={() => !isTemporary && openModal(data.id)}
      className={cn(
        isTemporary
          ? 'opacity-50 cursor-progress'
          : 'opacity-100 cursor-pointer',
        'flex items-start',
      )}
    >
      {mode === 'removed' ? (
        <XCircle className='mt-0.5 text-gray-500' />
      ) : (
        <DotSingle
          className={cn('mt-0.5', {
            'text-primary-600': mode === 'created',
            'text-gray-500': mode !== 'created',
          })}
        />
      )}

      <p className='max-w-[500px] line-clamp-2 ml-2 text-sm text-gray-700'>
        {formattedContent}
      </p>
    </div>
  );
};
