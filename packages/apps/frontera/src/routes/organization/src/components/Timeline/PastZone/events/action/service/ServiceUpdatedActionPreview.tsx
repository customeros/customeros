import { FC } from 'react';

import { cn } from '@ui/utils/cn.ts';
import { Action } from '@graphql/types';
import { File02 } from '@ui/media/icons/File02';
import { DotSingle } from '@ui/media/icons/DotSingle';
import { Card, CardFooter, CardContent } from '@ui/presentation/Card/Card';
import { getMetadata } from '@organization/components/Timeline/PastZone/events/action/utils';
import { TimelineEventPreviewHeader } from '@organization/components/Timeline/shared/TimelineEventPreview/header/TimelineEventPreviewHeader';
import {
  useTimelineEventPreviewStateContext,
  useTimelineEventPreviewMethodsContext,
} from '@organization/components/Timeline/shared/TimelineEventPreview/context/TimelineEventPreviewContext';

import { formatString } from './utils.tsx';

export const ServiceUpdatedActionPreview: FC<{
  mode?: 'created' | 'updated';
}> = ({ mode = 'updated' }) => {
  const { modalContent } = useTimelineEventPreviewStateContext();
  const { closeModal } = useTimelineEventPreviewMethodsContext();

  const event = modalContent as Action;
  const metadata = getMetadata(event?.metadata);
  const reasonForChange = metadata?.reasonForChange;

  const formattedContent = formatString(
    event?.content ?? '',
    metadata?.billedType,
    metadata?.currency ?? 'USD',
  );

  return (
    <>
      <TimelineEventPreviewHeader
        onClose={closeModal}
        name='Service updated'
        date={event?.createdAt}
        copyLabel='Copy link to this event'
      />
      <Card className='m-6 mt-3 p-4 shadow-xs'>
        <CardContent className='flex p-0 items-start'>
          <DotSingle
            className={cn('size-4 min-w-4 mt-0.5', {
              'text-primary-600': mode === 'created',
              'text-gray-500': mode !== 'created',
            })}
          />
          <p className='max-w-[500px] line-clamp-2 ml-2 text-sm text-gray-700'>
            {formattedContent}
          </p>
        </CardContent>

        {metadata?.comment && (
          <CardFooter className='flex p-0 pt-3 mt-4 items-center border-t border-gray-200'>
            <File02 color='gray.400' />
            <p className='max-w-[500px] line-clamp-2 ml-2 text-sm text-gray-500'>
              {reasonForChange}
            </p>
          </CardFooter>
        )}
      </Card>
    </>
  );
};
