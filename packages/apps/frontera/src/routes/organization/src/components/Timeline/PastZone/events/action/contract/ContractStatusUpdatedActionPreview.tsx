import { useMemo } from 'react';

import { cn } from '@ui/utils/cn';
import { Action } from '@graphql/types';
import { Card, CardContent } from '@ui/presentation/Card/Card';
import { getMetadata } from '@organization/components/Timeline/PastZone/events/action/utils';
import { iconsByStatus } from '@organization/components/Timeline/PastZone/events/action/contract/utils';
import { TimelineEventPreviewHeader } from '@organization/components/Timeline/shared/TimelineEventPreview/header/TimelineEventPreviewHeader';
import {
  useTimelineEventPreviewStateContext,
  useTimelineEventPreviewMethodsContext,
} from '@organization/components/Timeline/shared/TimelineEventPreview/context/TimelineEventPreviewContext';

export const ContractStatusUpdatedActionPreview = () => {
  const { modalContent } = useTimelineEventPreviewStateContext();
  const { closeModal } = useTimelineEventPreviewMethodsContext();
  const event = modalContent as Action;
  const metadata = useMemo(() => {
    return getMetadata(event?.metadata);
  }, [event?.metadata]);
  const status = metadata?.status?.toLowerCase();

  // todo remove when contract name is passed from BE in metadata
  const getContractName = () => {
    const content = event.content ?? '';
    const endIndex =
      content.lastIndexOf(' is now live') > -1
        ? content.lastIndexOf(' is now live')
        : content.lastIndexOf(' renewed') > -1
        ? content.lastIndexOf(' renewed')
        : content.lastIndexOf(' has ended') > -1
        ? content.lastIndexOf(' has ended')
        : content.length;

    return content.substring(0, endIndex).trim();
  };
  const content = (event.content ?? '').substring(
    0,
    event?.content?.lastIndexOf(' '),
  );

  return (
    <>
      <TimelineEventPreviewHeader
        onClose={closeModal}
        date={event?.createdAt}
        copyLabel='Copy link to this event'
        name={`${getContractName()} ${iconsByStatus[status].text} ${status}`}
      />
      <Card className='m-6 mt-3 p-4 shadow-xs'>
        <CardContent className='flex p-0 items-center'>
          {iconsByStatus[status].icon}
          <p className='max-w-[500px] line-clamp-2 ml-2 text-sm text-gray-700'>
            {content}
            <span
              className={cn(
                status === 'renewed' ? 'font-normal' : 'font-semibold',
                'ml-1',
              )}
            >
              {status}
            </span>
          </p>
        </CardContent>
      </Card>
    </>
  );
};
