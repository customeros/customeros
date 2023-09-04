'use client';
import React, { FC } from 'react';
import { Card, CardBody } from '@ui/presentation/Card';
import {
  ContactParticipant,
  JobRoleParticipant,
  UserParticipant,
} from '@graphql/types';
import { useTimelineEventPreviewContext } from '@organization/components/Timeline/preview/TimelineEventsPreviewContext/TimelineEventPreviewContext';
import { Avatar } from '@ui/media/Avatar';
import { Flex } from '@ui/layout/Flex';
import { Text } from '@ui/typography/Text';
import { getName } from '@spaces/utils/getParticipantsName';
import Slack from '@spaces/atoms/icons/Slack';
import { Button } from '@ui/form/Button';
import { SlackMessageCard } from '@organization/components/Timeline/events/slack/SlackMessageCard';
import { DateTimeUtils } from '@spaces/utils/date';
import { InteractionEventWithDate } from '@organization/components/Timeline/types';

export const SlackStub: FC<{ slackEvent: InteractionEventWithDate }> = ({
  slackEvent,
}) => {
  const { openModal } = useTimelineEventPreviewContext();

  const slackSender =
    (slackEvent?.sentBy?.[0] as ContactParticipant)?.contactParticipant ||
    (slackEvent?.sentBy?.[0] as JobRoleParticipant)?.jobRoleParticipant
      ?.contact ||
    (slackEvent?.sentBy?.[0] as UserParticipant)?.userParticipant;

  if (!slackSender) {
    return (
      <Card
        variant='outline'
        size='md'
        fontSize='14px'
        background='white'
        flexDirection='row'
        maxWidth={549}
        position='unset'
        cursor='pointer'
        boxShadow='xs'
        borderColor='gray.100'
      >
        <CardBody p={3} overflow={'hidden'}>
          <Flex gap={3} flex={1} alignItems='center'>
            <Avatar
              variant='roundedSquare'
              size='sm'
              colorScheme='gray'
              bg='white'
              border='1px solid var(--chakra-colors-gray-200)'
              icon={<Slack height='1.8rem' />}
            />
            <Text color='gray.700' as='span' fontWeight={600}>
              {slackEvent?.content}
            </Text>
          </Flex>
        </CardBody>
      </Card>
    );
  }

  const slackEventReplies = slackEvent.interactionSession?.events?.filter(
    (e) => e?.id !== slackEvent?.id,
  );
  const uniqThreadParticipants = slackEventReplies
    ?.map((e) => {
      const threadSender =
        (e?.sentBy?.[0] as ContactParticipant)?.contactParticipant ||
        (e?.sentBy?.[0] as JobRoleParticipant)?.jobRoleParticipant?.contact ||
        (e?.sentBy?.[0] as UserParticipant)?.userParticipant;

      return threadSender;
    })
    ?.filter((v, i, a) => a.findIndex((t) => !!t && t?.id === v?.id) === i);

  return (
    <>
      <SlackMessageCard
        name={getName(slackSender)}
        profilePhotoUrl={slackSender?.profilePhotoUrl}
        sourceUrl={slackEvent?.externalLinks?.[0]?.externalUrl}
        content={slackEvent?.content || ''}
        onClick={() => openModal(slackEvent)}
        date={DateTimeUtils.formatTime(slackEvent?.date)}
        showDateOnHover
      >
        {!!slackEventReplies?.length && (
          <Flex mt={1}>
            <Flex columnGap={1} mr={1}>
              {uniqThreadParticipants?.map(
                ({ id, name, firstName, lastName, ...rest }) => (
                  <Avatar
                    name={name || `${firstName} ${lastName}`}
                    key={`uniq-slack-thread-participant-${slackEvent.id}-${id}`}
                    variant='roundedSquareSmall'
                    size='xs'
                    src={rest?.profilePhotoUrl || undefined}
                  />
                ),
              )}
            </Flex>
            <Button variant='link' fontSize='sm' size='sm'>
              {slackEventReplies.length}{' '}
              {slackEventReplies.length === 1 ? 'reply' : 'replies'}
            </Button>
          </Flex>
        )}
      </SlackMessageCard>
    </>
  );
};
