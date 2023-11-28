'use client';
import React, { FC } from 'react';

import { Flex } from '@ui/layout/Flex';
import { Button } from '@ui/form/Button';
import { Avatar } from '@ui/media/Avatar';
import { User02 } from '@ui/media/icons/User02';
import { DateTimeUtils } from '@spaces/utils/date';
import { InteractionEventWithDate } from '@organization/src/components/Timeline/types';
import { SlackMessageCard } from '@organization/src/components/Timeline/events/slack/SlackMessageCard';
import { useTimelineEventPreviewMethodsContext } from '@organization/src/components/Timeline/preview/context/TimelineEventPreviewContext';

import { getDisplayNameAndAvatar } from './util';

export const SlackStub: FC<{ slackEvent: InteractionEventWithDate }> = ({
  slackEvent,
}) => {
  const { openModal } = useTimelineEventPreviewMethodsContext();

  const slackSender = getDisplayNameAndAvatar(slackEvent?.sentBy?.[0]);

  const isSentByTenantUser =
    slackEvent?.sentBy?.[0]?.__typename === 'UserParticipant';
  const slackEventReplies = slackEvent.interactionSession?.events?.filter(
    (e) => e?.id !== slackEvent?.id,
  );
  const uniqThreadParticipants = slackEventReplies
    ?.map((e) => {
      return getDisplayNameAndAvatar(e.sentBy?.[0]);
    })
    ?.filter((v, i, a) => a.findIndex((t) => !!t && t?.id === v?.id) === i);

  return (
    <SlackMessageCard
      name={slackSender.displayName || 'Unknown'}
      profilePhotoUrl={slackSender.photoUrl || undefined}
      sourceUrl={slackEvent?.externalLinks?.[0]?.externalUrl}
      content={slackEvent?.content || ''}
      onClick={() => openModal(slackEvent.id)}
      date={DateTimeUtils.formatTime(slackEvent?.date)}
      showDateOnHover
      ml={isSentByTenantUser ? 6 : 0}
    >
      {!!slackEventReplies?.length && (
        <Flex mt={1}>
          <Flex columnGap={1} mr={1}>
            {uniqThreadParticipants?.map(({ id, displayName, photoUrl }) => {
              return (
                <Avatar
                  size='xs'
                  name={displayName}
                  src={photoUrl ?? undefined}
                  variant='roundedSquareSmall'
                  icon={<User02 color='primary.700' />}
                  key={`uniq-slack-thread-participant-${slackEvent.id}-${id}`}
                />
              );
            })}
          </Flex>
          <Button variant='link' fontSize='sm' size='sm'>
            {slackEventReplies.length}{' '}
            {slackEventReplies.length === 1 ? 'reply' : 'replies'}
          </Button>
        </Flex>
      )}
    </SlackMessageCard>
  );
};
