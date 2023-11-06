import React from 'react';

import copy from 'copy-to-clipboard';

import { Flex } from '@ui/layout/Flex';
import { VStack } from '@ui/layout/Stack';
import { Text } from '@ui/typography/Text';
import { Link03 } from '@ui/media/icons/Link03';
import { XClose } from '@ui/media/icons/XClose';
import { Heading } from '@ui/typography/Heading';
import { IconButton } from '@ui/form/IconButton';
import { Tooltip } from '@ui/presentation/Tooltip';
import { DateTimeUtils } from '@spaces/utils/date';
import { Divider } from '@ui/presentation/Divider';
import { getName } from '@spaces/utils/getParticipantsName';
import { CardBody, CardHeader } from '@ui/presentation/Card';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useGetTimelineEventsQuery } from '@organization/src/graphql/getTimelineEvents.generated';
import { SlackMessageCard } from '@organization/src/components/Timeline/events/slack/SlackMessageCard';
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
} from '@organization/src/components/Timeline/preview/context/TimelineEventPreviewContext';

import { MessageCardSkeleton } from '../../shared';

const getParticipant = (sentBy?: InteractionEventParticipant[]) => {
  const sender =
    (sentBy?.[0] as ContactParticipant)?.contactParticipant ||
    (sentBy?.[0] as JobRoleParticipant)?.jobRoleParticipant?.contact ||
    (sentBy?.[0] as UserParticipant)?.userParticipant;

  return sender;
};
export const SlackThreadPreviewModal: React.FC = () => {
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
        >
          <Flex mb={2} alignItems='center'>
            <Heading size='sm' fontSize='lg'>
              {event?.interactionSession?.name || 'Thread'}
            </Heading>
            {/* todo uncomment when channel data is available  */}
            {/*{channel && (*/}
            {/*  <Text color='gray.500' ml={2} fontSize='sm'>*/}
            {/*    {channel}*/}
            {/*  </Text>*/}
            {/*)}*/}
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
      <CardBody mt={0} maxHeight='50%' overflow='auto' pb={6} pt={0}>
        <SlackMessageCard
          w='full'
          name={getName(slackSender)}
          profilePhotoUrl={slackSender?.profilePhotoUrl}
          sourceUrl={event?.externalLinks?.[0]?.externalUrl}
          content={event?.content || ''}
          // @ts-expect-error typescript does not work well with aliases
          date={DateTimeUtils.timeAgo(event?.date, { addSuffix: true })}
        />

        {isLoading && (
          <>
            <Flex marginY={2} alignItems='center'>
              <Text color='gray.400' fontSize='sm' whiteSpace='nowrap' mr={2}>
                {timelineEventsIds.length - 1}{' '}
                {timelineEventsIds.length - 1 === 1 ? 'reply' : 'replies'}
              </Text>
              <Divider />
            </Flex>
            <VStack w='full'>
              {Array.from({ length: timelineEventsIds.length - 1 }).map(
                (_, idx) => (
                  <MessageCardSkeleton key={idx} />
                ),
              )}
            </VStack>
          </>
        )}
        {!!slackEventReplies.length && (
          <>
            <Flex marginY={2} alignItems='center'>
              <Text color='gray.400' fontSize='sm' whiteSpace='nowrap' mr={2}>
                {slackEventReplies.length}{' '}
                {slackEventReplies.length === 1 ? 'reply' : 'replies'}
              </Text>
              <Divider />
            </Flex>

            <Flex direction='column' gap={2}>
              {slackEventReplies.map((reply) => {
                const sentBy = event?.interactionSession?.events?.find(
                  (e) => e.id === reply.id,
                )?.sentBy;

                const replyParticipant = getParticipant(sentBy);

                return (
                  <SlackMessageCard
                    key={`slack-event-thread-reply-preview-modal-${reply.id}`}
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
