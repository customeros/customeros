import { File02 } from '@ui/media/icons/File02';
import { Flag04 } from '@ui/media/icons/Flag04';
import { Action, OnboardingStatus } from '@graphql/types';
import { Card, CardFooter, CardContent } from '@ui/presentation/Card/Card';
import { getMetadata } from '@organization/components/Timeline/PastZone/events/action/utils';
import { TimelineEventPreviewHeader } from '@organization/components/Timeline/shared/TimelineEventPreview/header/TimelineEventPreviewHeader';
import {
  useTimelineEventPreviewStateContext,
  useTimelineEventPreviewMethodsContext,
} from '@organization/components/Timeline/shared/TimelineEventPreview/context/TimelineEventPreviewContext';

import { getIconColor } from './util';

export const OnboardingStatusChangedActionPreview = () => {
  const { modalContent } = useTimelineEventPreviewStateContext();
  const { closeModal } = useTimelineEventPreviewMethodsContext();

  const event = modalContent as Action;
  const metadata = getMetadata(event?.metadata);
  const reasonForChange = metadata?.comments;
  const status = metadata?.status as OnboardingStatus;
  const iconClassName = getIconColor(status);

  return (
    <>
      <TimelineEventPreviewHeader
        onClose={closeModal}
        date={event?.createdAt}
        name='Onboarding status changed'
        copyLabel='Copy link to this event'
      />
      <Card className='m-6 mt-3 p-4 shadow-xs'>
        <CardContent className='flex p-0 items-center'>
          <Flag04 className={iconClassName} />
          <p className='max-w-[500px] line-clamp-2 ml-2 text-sm text-gray-700'>
            {event?.content}
          </p>
        </CardContent>

        {metadata?.comments && (
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
