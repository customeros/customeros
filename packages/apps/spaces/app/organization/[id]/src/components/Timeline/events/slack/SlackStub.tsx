'use client';
import React, { FC } from 'react';

import { Flex } from '@ui/layout/Flex';
import { Button } from '@ui/form/Button';
import { Avatar } from '@ui/media/Avatar';
import { User02 } from '@ui/media/icons/User02';
import { DateTimeUtils } from '@spaces/utils/date';
import { getName } from '@spaces/utils/getParticipantsName';
import { InteractionEventWithDate } from '@organization/src/components/Timeline/types';
import { SlackMessageCard } from '@organization/src/components/Timeline/events/slack/SlackMessageCard';
import {
  UserParticipant,
  EmailParticipant,
  ContactParticipant,
  JobRoleParticipant,
} from '@graphql/types';
import { useTimelineEventPreviewMethodsContext } from '@organization/src/components/Timeline/preview/context/TimelineEventPreviewContext';

export const SlackStub: FC<{ slackEvent: InteractionEventWithDate }> = ({
  slackEvent,
}) => {
  const { openModal } = useTimelineEventPreviewMethodsContext();

  const slackSender =
    (slackEvent?.sentBy?.[0] as ContactParticipant)?.contactParticipant ||
    (slackEvent?.sentBy?.[0] as JobRoleParticipant)?.jobRoleParticipant
      ?.contact ||
    (slackEvent?.sentBy?.[0] as UserParticipant)?.userParticipant ||
    (slackEvent?.sentBy?.[0] as EmailParticipant)?.emailParticipant;
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
        (e?.sentBy?.[0] as UserParticipant)?.userParticipant ||
        (e?.sentBy?.[0] as EmailParticipant)?.emailParticipant
      );
    })
    ?.filter((v, i, a) => a.findIndex((t) => !!t && t?.id === v?.id) === i);

  return (
    <SlackMessageCard
      name={slackSender ? getName(slackSender) : 'Unknown'}
      profilePhotoUrl={slackSender?.profilePhotoUrl}
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
