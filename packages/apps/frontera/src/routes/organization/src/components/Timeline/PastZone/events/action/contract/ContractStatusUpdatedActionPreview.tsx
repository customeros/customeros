import { useMemo } from 'react';

import { cn } from '@ui/utils/cn';
import { Action } from '@graphql/types';
import { File02 } from '@ui/media/icons/File02';
import { FeaturedIcon } from '@ui/media/Icon/FeaturedIcon';
import { Card, CardFooter, CardContent } from '@ui/presentation/Card/Card';
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
          <div className='inline-block w-[30px]'>
            <FeaturedIcon
              size='md'
              // eslint-disable-next-line @typescript-eslint/no-explicit-any
              colorScheme={iconsByStatus[status].colorScheme as any}
            >
              {iconsByStatus[status].icon}
            </FeaturedIcon>
          </div>
          <p className='my-1 max-w[500px] line-clamp-2 ml-2 text-sm text-gray-700'>
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
        {metadata?.comment && (
          <CardFooter className='flex p-0 pt-3 mt-4 items-center border-t border-gray-200'>
            <File02 className='text-gray-400' />
            <p className='max-w-[500px] line-clamp-2 ml-2 text-sm text-gray-500'>
              {metadata.content}
            </p>
          </CardFooter>
        )}
      </Card>
    </>
  );
};
