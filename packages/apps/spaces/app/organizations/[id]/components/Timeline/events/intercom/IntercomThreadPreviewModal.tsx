import React, { useCallback } from 'react';
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
import { getName } from '@spaces/utils/getParticipantsName';
import {
  ContactParticipant,
  InteractionEvent,
  InteractionEventParticipant,
  JobRoleParticipant,
  UserParticipant,
} from '@graphql/types';
import { Divider } from '@ui/presentation/Divider';
import { IntercomMessageCard } from './IntercomMessageCard';
import { convert } from 'html-to-text';

const getParticipant = (sentBy?: InteractionEventParticipant[]) => {
  const sender =
    (sentBy?.[0] as ContactParticipant)?.contactParticipant ||
    (sentBy?.[0] as JobRoleParticipant)?.jobRoleParticipant?.contact ||
    (sentBy?.[0] as UserParticipant)?.userParticipant;
  return sender;
};
export const IntercomThreadPreviewModal: React.FC = () => {
  const { closeModal, modalContent } = useTimelineEventPreviewContext();
  const event = modalContent as InteractionEvent;
  const intercomSender = getParticipant(event?.sentBy);
  const slackEventReplies =
    event?.interactionSession?.events?.filter((e) => e?.id !== event?.id) || [];
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
      <CardHeader pb={1} position='sticky' top={0} borderRadius='xl'>
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
