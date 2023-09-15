import React from 'react';
import { ExternalSystemType, Maybe, TimelineEvent } from '@graphql/types';
import { Icons } from '@ui/media/Icon';
import { Flex } from '@ui/layout/Flex';
import { Text } from '@ui/typography/Text';

import { DateTimeUtils } from '../../../utils';

export const LastTouchpointTableCell = ({
  lastTouchPointAt,
  lastTouchPointTimelineEvent,
}: {
  lastTouchPointAt: any;
  lastTouchPointTimelineEvent: Maybe<TimelineEvent> | undefined;
}) => {
  let icon = undefined;
  let label = '';

  if (lastTouchPointTimelineEvent) {
    switch (lastTouchPointTimelineEvent.__typename) {
      case 'Action':
        label = '';
        switch (lastTouchPointTimelineEvent.actionType) {
          case 'CREATED':
            label += 'Created';
            icon = <Icons.Building7 boxSize='3' color='gray.700' />;
            break;
          default:
            break;
        }
        break;
      case 'PageView':
        label = 'Page View';
        break;
      case 'Issue':
        label = 'Issue';
        icon = <Icons.Ticket2 boxSize='3' color='gray.700' />;
        break;
      case 'Note': {
        label =
          'Note by ' +
          lastTouchPointTimelineEvent.createdBy?.firstName +
          ' ' +
          lastTouchPointTimelineEvent.createdBy?.lastName;
        icon = <Icons.File2 boxSize='3' color='gray.700' />;
        break;
      }
      case 'InteractionEvent': {
        if (lastTouchPointTimelineEvent.channel === 'EMAIL') {
          label = 'Email sent';
          icon = <Icons.Mail1 boxSize='3' color='gray.700' />;
        } else if (lastTouchPointTimelineEvent.channel === 'VOICE') {
          label = 'Phone call';
          icon = <Icons.PhoneOutgoing2 boxSize='3' color='greay.700' />;
        } else if (
          !lastTouchPointTimelineEvent.channel &&
          lastTouchPointTimelineEvent.eventType === 'meeting'
        ) {
          label = 'Meeting';
          icon = <Icons.Calendar boxSize='3' color='gray.700' />;
        } else if (lastTouchPointTimelineEvent.channel === 'CHAT') {
          if (
            lastTouchPointTimelineEvent.externalLinks?.[0].type ===
            ExternalSystemType.Slack
          ) {
            label = 'Slack message';
            icon = <Icons.MessageTextSquare1 boxSize='3' color='gray.700' />;
          } else if (
            lastTouchPointTimelineEvent.externalLinks?.[0].type ===
            ExternalSystemType.Intercom
          ) {
            label = 'Intercom message';
            icon = <Icons.MessageTextSquare1 boxSize='3' color='gray.700' />;
          }
        } else {
          label = 'InteractionEvent';
        }
        break;
      }
      case 'Analysis':
        label = 'Analysis';
        break;
      case 'Meeting': {
        label =
          'Meeting with ' +
          lastTouchPointTimelineEvent.attendedBy.length +
          ' participant' +
          (lastTouchPointTimelineEvent.attendedBy.length === 1 ? '' : 's');
        icon = <Icons.Calendar boxSize='3' color='gray.700' />;
        break;
      }
      default:
        label = 'Unknown';
        icon = <Icons.Ticket2 boxSize='3' color='gray.700' />;
        break;
    }
  }

  const subLabel = label
    ? DateTimeUtils.timeAgo(lastTouchPointAt, {
        addSuffix: true,
      })
    : '';

  return (
    <Flex flexDir='column'>
      <Flex align='center'>
        {icon}
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
