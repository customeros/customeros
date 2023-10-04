'use client';
import React, { FC } from 'react';
import {
  ContactParticipant,
  JobRoleParticipant,
  UserParticipant,
} from '@graphql/types';
import { useTimelineEventPreviewContext } from '@organization/src/components/Timeline/preview/context/TimelineEventPreviewContext';
import { Avatar } from '@ui/media/Avatar';
import { Flex } from '@ui/layout/Flex';
import { getName } from '@spaces/utils/getParticipantsName';
import { Button } from '@ui/form/Button';
import { SlackMessageCard } from '@organization/src/components/Timeline/events/slack/SlackMessageCard';
import { DateTimeUtils } from '@spaces/utils/date';
import { InteractionEventWithDate } from '@organization/src/components/Timeline/types';
import { User02 } from '@ui/media/icons/User02';

export const SlackStub: FC<{ slackEvent: InteractionEventWithDate }> = ({
  slackEvent,
}) => {
  const { openModal } = useTimelineEventPreviewContext();

  const slackSender =
    (slackEvent?.sentBy?.[0] as ContactParticipant)?.contactParticipant ||
    (slackEvent?.sentBy?.[0] as JobRoleParticipant)?.jobRoleParticipant
      ?.contact ||
    (slackEvent?.sentBy?.[0] as UserParticipant)?.userParticipant;
  const isSentByTenantUser =
    slackEvent?.sentBy?.[0]?.__typename === 'UserParticipant';
  const slackEventReplies = slackEvent.interactionSession?.events?.filter(
    (e) => e?.id !== slackEvent?.id,
  );
  const uniqThreadParticipants = slackEventReplies
    ?.map((e) => {
      return (
        (e?.sentBy?.[0] as ContactParticipant)?.contactParticipant ||
        (e?.sentBy?.[0] as JobRoleParticipant)?.jobRoleParticipant?.contact ||
        (e?.sentBy?.[0] as UserParticipant)?.userParticipant
      );
    })
    ?.filter((v, i, a) => a.findIndex((t) => !!t && t?.id === v?.id) === i);

  return (
    <SlackMessageCard
      name={slackSender ? getName(slackSender) : 'Unknown'}
      profilePhotoUrl={slackSender?.profilePhotoUrl}
      sourceUrl={slackEvent?.externalLinks?.[0]?.externalUrl}
      content={slackEvent?.content || ''}
      onClick={() => openModal(slackEvent)}
      date={DateTimeUtils.formatTime(slackEvent?.date)}
      showDateOnHover
      ml={isSentByTenantUser ? 6 : 0}
    >
      {!!slackEventReplies?.length && (
        <Flex mt={1}>
          <Flex columnGap={1} mr={1}>
            {uniqThreadParticipants?.map(
              ({ id, name, firstName, lastName, profilePhotoUrl }) => {
                const displayName =
                  name ?? [firstName, lastName].filter(Boolean).join(' ');

                return (
                  <Avatar
                    size='xs'
                    name={displayName}
                    variant='roundedSquareSmall'
                    icon={<User02 color='primary.700' />}
                    src={profilePhotoUrl ? profilePhotoUrl : undefined}
                    key={`uniq-slack-thread-participant-${slackEvent.id}-${id}`}
                  />
                );
              },
            )}
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
