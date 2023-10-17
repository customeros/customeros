import React from 'react';
import { convert } from 'html-to-text';
import copy from 'copy-to-clipboard';

import { CardHeader, CardBody } from '@ui/presentation/Card';
import { Heading } from '@ui/typography/Heading';
import { Text } from '@ui/typography/Text';
import { Flex } from '@ui/layout/Flex';
import { VStack } from '@ui/layout/Stack';
import { Tooltip } from '@ui/presentation/Tooltip';
import { IconButton } from '@ui/form/IconButton';
import { DateTimeUtils } from '@spaces/utils/date';
import { Link03 } from '@ui/media/icons/Link03';
import { XClose } from '@ui/media/icons/XClose';
import { getName } from '@spaces/utils/getParticipantsName';
import {
  ContactParticipant,
  InteractionEvent,
  InteractionEventParticipant,
  JobRoleParticipant,
  UserParticipant,
} from '@graphql/types';
import { Divider } from '@ui/presentation/Divider';
import { useGetTimelineEventsQuery } from '@organization/src/graphql/getTimelineEvents.generated';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';

import { MessageCardSkeleton } from '../../shared';
import { IntercomMessageCard } from './IntercomMessageCard';
import {
  useTimelineEventPreviewMethodsContext,
  useTimelineEventPreviewStateContext,
} from '@organization/src/components/Timeline/preview/context/TimelineEventPreviewContext';

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
      <CardHeader
        py='4'
        px='6'
        pb='1'
        position='sticky'
        top={0}
        borderRadius='xl'
      >
        <Flex
          direction='row'
          justifyContent='space-between'
          alignItems='center'
          maxH='calc(100vh - 5rem)'
        >
          <Flex mb={2} alignItems='center'>
            <Heading
              size='sm'
              fontSize='lg'
              noOfLines={1}
              maxW={event?.interactionSession?.name ? 'unset' : '248px'}
            >
              {title}
            </Heading>
          </Flex>
          <Flex direction='row' justifyContent='flex-end' alignItems='center'>
            <Tooltip label='Copy link to this thread' placement='bottom'>
              <IconButton
                variant='ghost'
                aria-label='Copy link to this thread'
                color='gray.500'
                size='sm'
                mr={1}
                icon={<Link03 color='gray.500' boxSize='4' />}
                onClick={() => copy(window.location.href)}
              />
            </Tooltip>
            <Tooltip label='Close' aria-label='close' placement='bottom'>
              <IconButton
                variant='ghost'
                aria-label='Close preview'
                color='gray.500'
                size='sm'
                icon={<XClose color='gray.500' boxSize='5' />}
                onClick={closeModal}
              />
            </Tooltip>
          </Flex>
        </Flex>
      </CardHeader>
      <CardBody mt={0} maxHeight='calc(100vh - 9rem)' overflow='auto' pb={6}>
        <IntercomMessageCard
          w='full'
          name={getName(intercomSender)}
          profilePhotoUrl={intercomSender?.profilePhotoUrl}
          sourceUrl={event?.externalLinks?.[0]?.externalUrl}
          content={event?.content || ''}
          // @ts-expect-error typescript does not work well with aliases
          date={DateTimeUtils.timeAgo(event?.date, { addSuffix: true })}
        />

        {isLoading && (
          <>
            <Flex marginY={2} alignItems='center'>
              <Text color='gray.400' fontSize='sm' whiteSpace='nowrap' mr={2}>
                {/* subtracting 2 for intercom because system messages are hidden */}
                {timelineEventsIds.length - 2}{' '}
                {timelineEventsIds.length - 2 === 1 ? 'reply' : 'replies'}
              </Text>
              <Divider />
            </Flex>
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
            <Flex marginY={2} alignItems='center'>
              <Text color='gray.400' fontSize='sm' whiteSpace='nowrap' mr={2}>
                {intercomEventReplies.length}{' '}
                {intercomEventReplies.length === 1 ? 'reply' : 'replies'}
              </Text>
              <Divider />
            </Flex>

            <Flex direction='column' gap={2}>
              {intercomEventReplies.map((reply) => {
                const sentBy = event?.interactionSession?.events?.find(
                  (e) => e.id === reply.id,
                )?.sentBy;

                const replyParticipant = getParticipant(sentBy);
                return (
                  <IntercomMessageCard
                    key={`intercom-event-thread-reply-preview-modal-${reply.id}`}
                    w='full'
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
            </Flex>
          </>
        )}
      </CardBody>
    </>
  );
};
