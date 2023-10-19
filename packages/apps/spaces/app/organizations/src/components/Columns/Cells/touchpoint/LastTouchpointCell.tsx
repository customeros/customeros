import React from 'react';
import { match } from 'ts-pattern';

import {
  Maybe,
  ActionType,
  TimelineEvent,
  ExternalSystemType,
} from '@graphql/types';
import { Flex } from '@ui/layout/Flex';
import { Text } from '@ui/typography/Text';
import { IconProps } from '@ui/media/Icon';
import { File02 } from '@ui/media/icons/File02';
import { Mail01 } from '@ui/media/icons/Mail01';
import { DateTimeUtils } from '@spaces/utils/date';
import { Calendar } from '@ui/media/icons/Calendar';
import { Ticket02 } from '@ui/media/icons/Ticket02';
import { Building07 } from '@ui/media/icons/Building07';
import { PhoneOutgoing02 } from '@ui/media/icons/PhoneOutgoing02';
import { MessageChatSquare } from '@ui/media/icons/MessageChatSquare';

export const LastTouchpointCell = ({
  lastTouchPointAt,
  lastTouchPointTimelineEvent,
}: {
  lastTouchPointAt: any;
  lastTouchPointTimelineEvent: Maybe<TimelineEvent> | undefined;
}) => {
  const [label, Icon] = match(lastTouchPointTimelineEvent)
    .returnType<[string, (props: IconProps) => JSX.Element]>()
    .with({ __typename: 'Action', actionType: ActionType.Created }, () => [
      'Created',
      Building07,
    ])
    .with({ __typename: 'PageView' }, () => ['Page View', () => <></>])
    .with({ __typename: 'Issue' }, () => ['Issue', Ticket02])
    .with({ __typename: 'Note' }, (event) => [
      `Note by ${event.createdBy?.firstName} ${event.createdBy?.lastName}`,
      File02,
    ])
    .with({ __typename: 'LogEntry' }, (event) => [
      `Log entry${
        !!event.createdBy?.firstName || !!event.createdBy?.lastName
          ? ` by ${[event.createdBy?.firstName, event.createdBy?.lastName]
              .join(' ')
              .trim()}`
          : ''
      }`,
      MessageChatSquare,
    ])
    .with({ __typename: 'InteractionEvent', channel: 'EMAIL' }, () => [
      'Email sent',
      Mail01,
    ])
    .with({ __typename: 'InteractionEvent', channel: 'VOICE' }, () => [
      'Phone call',
      PhoneOutgoing02,
    ])
    .with(
      {
        __typename: 'InteractionEvent',
        channel: 'CHAT',
        externalLinks: [{ type: ExternalSystemType.Slack }],
      },
      () => ['Slack message', MessageChatSquare],
    )
    .with(
      {
        __typename: 'InteractionEvent',
        channel: 'CHAT',
        externalLinks: [{ type: ExternalSystemType.Intercom }],
      },
      () => ['Intercom message', MessageChatSquare],
    )
    .with({ __typename: 'InteractionEvent', eventType: 'meeting' }, () => [
      'Meeting',
      Calendar,
    ])
    .with({ __typename: 'Analysis' }, () => ['Analysis', () => <></>])
    .with({ __typename: 'Meeting' }, (event) => [
      `Meeting with ${event.attendedBy.length} participant${
        event.attendedBy.length === 1 ? '' : 's'
      }`,
      Calendar,
    ])
    .otherwise(() => ['', () => <></>]);

  const subLabel = label
    ? DateTimeUtils.timeAgo(lastTouchPointAt, {
        addSuffix: true,
      })
    : '';

  return (
    <Flex flexDir='column'>
      <Flex align='center'>
        <Icon boxSize='3' color='gray.700' />
        <Text color='gray.700' ml='2'>
          {label}
        </Text>
      </Flex>

      <Text color='gray.500' ml='5'>
        {subLabel}
      </Text>
    </Flex>
  );
};
