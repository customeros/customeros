import copy from 'copy-to-clipboard';

import { DateTimeUtils } from '@utils/date';
import { Link03 } from '@ui/media/icons/Link03';
import { XClose } from '@ui/media/icons/XClose';
import { IconButton } from '@ui/form/IconButton';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';
import { Divider } from '@ui/presentation/Divider/Divider';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { CardHeader, CardContent } from '@ui/presentation/Card/Card';
import { InteractionEvent, InteractionEventParticipant } from '@graphql/types';
import { useGetTimelineEventsQuery } from '@organization/graphql/getTimelineEvents.generated';
import { SlackMessageCard } from '@organization/components/Timeline/PastZone/events/slack/SlackMessageCard';
import {
  useTimelineEventPreviewStateContext,
  useTimelineEventPreviewMethodsContext,
} from '@organization/components/Timeline/shared/TimelineEventPreview/context/TimelineEventPreviewContext';

import { getDisplayNameAndAvatar } from './util';
import { MessageCardSkeleton } from '../../../shared';

const getParticipant = (sentBy?: InteractionEventParticipant[]) => {
  return getDisplayNameAndAvatar(sentBy?.[0]);
};

export const SlackThreadPreviewModal = () => {
  const client = getGraphQLClient();
  const { modalContent } = useTimelineEventPreviewStateContext();
  const { closeModal } = useTimelineEventPreviewMethodsContext();
  const event = modalContent as InteractionEvent;

  const timelineEventsIds =
    event?.interactionSession?.events?.map((e) => e.id) || [];
  const { data, isLoading } = useGetTimelineEventsQuery(client, {
    ids: timelineEventsIds,
  });

  const slackSender = getParticipant(event?.sentBy);
  const slackEventReplies =
    (data?.timelineEvents as InteractionEvent[] | undefined)?.filter(
      (e) => e?.id !== event?.id,
    ) || [];

  return (
    <>
      <CardHeader className='py-4 px-6 pb-1 sticky top-0 rounded-xl bg-white z-[1]'>
        <div className='flex justify-between items-center'>
          <div className='flex mb-2 items-center'>
            <h2 className='text-lg font-semibold'>
              {event?.interactionSession?.name || 'Thread'}
            </h2>
          </div>
          <div className='flex justify-end items-center'>
            <Tooltip side='bottom' label='Copy link to this thread'>
              <div>
                <IconButton
                  size='xs'
                  variant='ghost'
                  color='gray.500'
                  className='mr-1'
                  aria-label='Copy link to this thread'
                  onClick={() => copy(window.location.href)}
                  icon={<Link03 className='text-gray-500 size-4' />}
                />
              </div>
            </Tooltip>
            <Tooltip label='Close' side='bottom' aria-label='close'>
              <div>
                <IconButton
                  size='xs'
                  variant='ghost'
                  color='gray.500'
                  onClick={closeModal}
                  aria-label='Close preview'
                  icon={<XClose className='text-gray-500 size-5' />}
                />
              </div>
            </Tooltip>
          </div>
        </div>
      </CardHeader>
      <CardContent className='mt-0 max-h-[calc(100vh-60px-56px)] pb-6 pt-0'>
        <SlackMessageCard
          className='w-full'
          name={slackSender.displayName}
          content={event?.content || ''}
          profilePhotoUrl={slackSender?.photoUrl}
          sourceUrl={event?.externalLinks?.[0]?.externalUrl}
          // @ts-expect-error typescript does not work well with aliases
          date={DateTimeUtils.timeAgo(event?.date, { addSuffix: true })}
        />

        {isLoading && (
          <>
            <div className='flex my-2 items-center'>
              <p className='text-gray-400 text-sm whitespace-nowrap mr-2'>
                {timelineEventsIds.length - 1}{' '}
                {timelineEventsIds.length - 1 === 1 ? 'reply' : 'replies'}
              </p>
              <Divider />
            </div>
            <div className='flex flex-col w-full space-y-2'>
              {Array.from({ length: timelineEventsIds.length - 1 }).map(
                (_, idx) => (
                  <MessageCardSkeleton key={idx} />
                ),
              )}
            </div>
          </>
        )}
        {!!slackEventReplies.length && (
          <>
            <div className='flex my-2 items-center'>
              <p className='text-gray-400 text-sm whitespace-nowrap mr-2'>
                {slackEventReplies.length}{' '}
                {slackEventReplies.length === 1 ? 'reply' : 'replies'}
              </p>
              <Divider />
            </div>

            <div className='flex flex-col gap-2'>
              {slackEventReplies.map((reply) => {
                const sentBy = event?.interactionSession?.events?.find(
                  (e) => e.id === reply.id,
                )?.sentBy;

                const replyParticipant = getParticipant(sentBy);

                return (
                  <SlackMessageCard
                    className='w-full'
                    content={reply?.content || ''}
                    name={replyParticipant?.displayName}
                    profilePhotoUrl={replyParticipant?.photoUrl}
                    key={`slack-event-thread-reply-preview-modal-${reply.id}`}
                    // @ts-expect-error typescript does not work well with aliases
                    date={DateTimeUtils.timeAgo(reply?.date, {
                      addSuffix: true,
                    })}
                  />
                );
              })}
            </div>
          </>
        )}
      </CardContent>
    </>
  );
};
