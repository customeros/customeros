import React, { FC } from 'react';

import { File02 } from '@ui/media/icons/File02';
import { Action, BilledType } from '@graphql/types';
import { DotSingle } from '@ui/media/icons/DotSingle';
import { FeaturedIcon } from '@ui/media/Icon/FeaturedIcon2';
import { formatCurrency } from '@spaces/utils/getFormattedCurrencyNumber';
import { Card, CardFooter, CardContent } from '@ui/presentation/Card/Card';
import { getMetadata } from '@organization/src/components/Timeline/PastZone/events/action/utils';
import { TimelineEventPreviewHeader } from '@organization/src/components/Timeline/shared/TimelineEventPreview/header/TimelineEventPreviewHeader';
import {
  useTimelineEventPreviewStateContext,
  useTimelineEventPreviewMethodsContext,
} from '@organization/src/components/Timeline/shared/TimelineEventPreview/context/TimelineEventPreviewContext';

export const ServiceUpdatedActionPreview: FC<{
  mode?: 'created' | 'updated';
}> = ({ mode = 'updated' }) => {
  const { modalContent } = useTimelineEventPreviewStateContext();
  const { closeModal } = useTimelineEventPreviewMethodsContext();

  const event = modalContent as Action;
  const metadata = getMetadata(event?.metadata);
  const reasonForChange = metadata?.reasonForChange;

  const formattedContent = (event?.content ?? '')
    .replace(
      metadata?.price,
      formatCurrency(
        Number(metadata?.price),
        metadata?.billedType === BilledType.Usage ? 4 : 2,
      ),
    )
    .replace(
      metadata?.previousPrice,
      formatCurrency(
        Number(metadata?.previousPrice),
        metadata?.billedType === BilledType.Usage ? 4 : 2,
      ),
    );

  return (
    <>
      <TimelineEventPreviewHeader
        date={event?.createdAt}
        name='Service updated'
        onClose={closeModal}
        copyLabel='Copy link to this event'
      />
      <Card className='m-6 mt-3 p-4 shadow-xs'>
        <CardContent className='flex p-0 items-center'>
          <FeaturedIcon
            className='min-w-10'
            size='md'
            colorScheme={mode === 'created' ? 'primary' : 'gray'}
          >
            <DotSingle />
          </FeaturedIcon>
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
