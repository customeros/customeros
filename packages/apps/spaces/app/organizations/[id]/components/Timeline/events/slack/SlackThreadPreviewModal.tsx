import React from 'react';
import { CardHeader, CardBody } from '@ui/presentation/Card';
import { Heading } from '@ui/typography/Heading';
import { Text } from '@ui/typography/Text';
import { Flex } from '@ui/layout/Flex';
import { Tooltip } from '@ui/presentation/Tooltip';
import { IconButton } from '@ui/form/IconButton';
import { useTimelineEventPreviewContext } from '@organization/components/Timeline/preview/TimelineEventsPreviewContext/TimelineEventPreviewContext';
import { DateTimeUtils } from '@spaces/utils/date';
import CopyLink from '@spaces/atoms/icons/CopyLink';
import Times from '@spaces/atoms/icons/Times';
import copy from 'copy-to-clipboard';
import { SlackMessageCard } from '@organization/components/Timeline/events/slack/SlackMessageCard';
import { getName } from '@spaces/utils/getParticipantsName';
import {
  ContactParticipant,
  InteractionEventParticipant,
  JobRoleParticipant,
  UserParticipant,
} from '@graphql/types';
import { Divider } from '@ui/presentation/Divider';

const getParticipant = (sentBy?: InteractionEventParticipant[]) => {
  const sender =
    (sentBy?.[0] as ContactParticipant)?.contactParticipant ||
    (sentBy?.[0] as JobRoleParticipant)?.jobRoleParticipant?.contact ||
    (sentBy?.[0] as UserParticipant)?.userParticipant;
  return sender;
};
export const SlackThreadPreviewModal: React.FC = () => {
  const { closeModal, modalContent } = useTimelineEventPreviewContext();
  const slackSender = getParticipant(modalContent?.sentBy);

  const slackEventReplies =
    modalContent?.interactionSession?.events?.filter(
      (e) => e?.id !== modalContent?.id,
    ) || [];

  return (
    <>
      <CardHeader pb={1} position='sticky' top={0} borderRadius='xl'>
        <Flex
          direction='row'
          justifyContent='space-between'
          alignItems='center'
        >
          <div>
            <Heading size='sm' mb={2}>
              Thread
            </Heading>
          </div>
          <Flex direction='row' justifyContent='flex-end' alignItems='center'>
            <Tooltip label='Copy link to this thread' placement='bottom'>
              <IconButton
                variant='ghost'
                aria-label='Copy link to this thread'
                color='gray.500'
                size='sm'
                mr={1}
                icon={<CopyLink color='gray.500' height='18px' />}
                onClick={() => copy(window.location.href)}
              />
            </Tooltip>
            <Tooltip label='Close' aria-label='close' placement='bottom'>
              <IconButton
                variant='ghost'
                aria-label='Close preview'
                color='gray.500'
                size='sm'
                icon={<Times color='gray.500' height='24px' />}
                onClick={closeModal}
              />
            </Tooltip>
          </Flex>
        </Flex>
      </CardHeader>
      <CardBody mt={0} maxHeight='50%' overflow='auto' pb={6}>
        <SlackMessageCard
          w='full'
          name={getName(slackSender)}
          profilePhotoUrl={slackSender?.profilePhotoUrl}
          content={modalContent?.content || ''}
          // @ts-expect-error typescript does not work well with aliases
          date={DateTimeUtils.timeAgo(modalContent?.date, { addSuffix: true })}
        />

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
                const replyParticipant = getParticipant(reply?.sentBy);
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
