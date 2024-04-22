import React from 'react';

import copy from 'copy-to-clipboard';
import { convert } from 'html-to-text';

import { cn } from '@ui/utils/cn';
import { VStack } from '@ui/layout/Stack';
import { Link03 } from '@ui/media/icons/Link03';
import { XClose } from '@ui/media/icons/XClose';
import { DateTimeUtils } from '@spaces/utils/date';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';
import { Divider } from '@ui/presentation/Divider/Divider';
import { IconButton } from '@ui/form/IconButton/IconButton';
import { getName } from '@spaces/utils/getParticipantsName';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { CardHeader, CardContent } from '@ui/presentation/Card/Card';
import { useGetTimelineEventsQuery } from '@organization/src/graphql/getTimelineEvents.generated';
import {
  UserParticipant,
  InteractionEvent,
  ContactParticipant,
  JobRoleParticipant,
  InteractionEventParticipant,
} from '@graphql/types';
import {
  useTimelineEventPreviewStateContext,
  useTimelineEventPreviewMethodsContext,
} from '@organization/src/components/Timeline/shared/TimelineEventPreview/context/TimelineEventPreviewContext';

import { MessageCardSkeleton } from '../../../shared';
import { IntercomMessageCard } from './IntercomMessageCard';

const getParticipant = (sentBy?: InteractionEventParticipant[]) => {
  const sender =
    (sentBy?.[0] as ContactParticipant)?.contactParticipant ||
    (sentBy?.[0] as JobRoleParticipant)?.jobRoleParticipant?.contact ||
    (sentBy?.[0] as UserParticipant)?.userParticipant;

  return sender;
};
export const IntercomThreadPreviewModal: React.FC = () => {
  const client = getGraphQLClient();
  const { modalContent } = useTimelineEventPreviewStateContext();
  const { closeModal } = useTimelineEventPreviewMethodsContext();

  const event = modalContent as InteractionEvent;

  const timelineEventsIds =
    event?.interactionSession?.events?.map((e) => e.id) || [];
  const { data, isLoading } = useGetTimelineEventsQuery(client, {
    ids: timelineEventsIds,
  });

  const intercomSender = getParticipant(event?.sentBy);
  const intercomEventReplies =
    (data?.timelineEvents as InteractionEvent[] | undefined)
      ?.filter((e) => e?.id !== event?.id)
      // TODO: remove this filter when we have a better way to handle this
      .filter(
        (e) =>
          !e.content?.includes('You’ll get replies here and in your email:'),
      ) || [];
  const title = (() => {
    const titleString = event?.interactionSession?.name || event?.content || '';

    return convert(`<p>${titleString}</p>`, {
      preserveNewlines: true,
      selectors: [
        {
          selector: 'a',
          options: { hideLinkHrefIfSameAsText: true, ignoreHref: true },
        },
      ],
    });
  })();

  return (
    <>
      <CardHeader className='py-4 px-6 pb-1 sticky top-0 rounded-xl'>
        <div className='flex justify-between items-center max-w-[calc(100vh-5rem)] flex-row'>
          <div className='flex mb-2 items-center'>
            <h2
              className={cn(
                'text-lg line-clamp-2',
                event?.interactionSession?.name ? '' : 'max-w-[248px]',
              )}
            >
              {title}
            </h2>
          </div>
          <div className='flex justify-end items-center flex-row'>
            <Tooltip
              label='Copy link to this thread'
              side='bottom'
              asChild={false}
            >
              <div>
                <IconButton
                  className='mr-1'
                  variant='ghost'
                  aria-label='Copy link to this thread'
                  color='gray.500'
                  size='sm'
                  icon={<Link03 className='text-gray-500' />}
                  onClick={() => copy(window.location.href)}
                />
              </div>
            </Tooltip>
            <Tooltip
              label='Close'
              aria-label='close'
              side='bottom'
              asChild={false}
            >
              <div>
                <IconButton
                  variant='ghost'
                  aria-label='Close preview'
                  color='gray.500'
                  size='sm'
                  icon={<XClose className='text-gray-500 size-5' />}
                  onClick={closeModal}
                />
              </div>
            </Tooltip>
          </div>
        </div>
      </CardHeader>
      <CardContent className='mt-0 max-h-[calc(100vh-9rem)] overflow-auto pb-6'>
        <IntercomMessageCard
          className='w-full'
          name={getName(intercomSender)}
          profilePhotoUrl={intercomSender?.profilePhotoUrl}
          sourceUrl={event?.externalLinks?.[0]?.externalUrl}
          content={event?.content || ''}
          // @ts-expect-error typescript does not work well with aliases
          date={DateTimeUtils.timeAgo(event?.date, { addSuffix: true })}
        />

        {isLoading && (
          <>
            <div className='flex my-2 items-center'>
              <p className='text-gray-400 text-sm whitespace-nowrap mr-2'>
                {/* subtracting 2 for intercom because system messages are hidden */}
                {timelineEventsIds.length - 2}{' '}
                {timelineEventsIds.length - 2 === 1 ? 'reply' : 'replies'}
              </p>
              <Divider />
            </div>
            <VStack w='full'>
              {Array.from({ length: timelineEventsIds.length - 2 }).map(
                (_, idx) => (
                  <MessageCardSkeleton key={idx} />
                ),
              )}
            </VStack>
          </>
        )}
        {!!intercomEventReplies.length && (
          <>
            <div className='flex my-2 items-center'>
              <p className='text-gray-400 text-sm whitespace-nowrap mr-2'>
                {intercomEventReplies.length}{' '}
                {intercomEventReplies.length === 1 ? 'reply' : 'replies'}
              </p>
              <Divider />
            </div>

            <div className='flex flex-col gap-2'>
              {intercomEventReplies.map((reply) => {
                const sentBy = event?.interactionSession?.events?.find(
                  (e) => e.id === reply.id,
                )?.sentBy;

                const replyParticipant = getParticipant(sentBy);

                return (
                  <IntercomMessageCard
                    key={`intercom-event-thread-reply-preview-modal-${reply.id}`}
                    className='w-full'
                    name={getName(replyParticipant)}
                    profilePhotoUrl={replyParticipant?.profilePhotoUrl}
                    content={reply?.content || ''}
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
