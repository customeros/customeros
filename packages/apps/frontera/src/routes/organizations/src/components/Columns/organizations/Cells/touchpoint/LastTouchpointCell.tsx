import React, { useRef } from 'react';

import { match } from 'ts-pattern';

import { DateTimeUtils } from '@utils/date';
import { File02 } from '@ui/media/icons/File02';
import { Mail01 } from '@ui/media/icons/Mail01';
import { Calendar } from '@ui/media/icons/Calendar';
import { Ticket02 } from '@ui/media/icons/Ticket02';
import { Building07 } from '@ui/media/icons/Building07';
import { TableCellTooltip } from '@ui/presentation/Table';
import { PhoneOutgoing02 } from '@ui/media/icons/PhoneOutgoing02';
import { MessageChatSquare } from '@ui/media/icons/MessageChatSquare';
import {
  Maybe,
  ActionType,
  TimelineEvent,
  ExternalSystemType,
  LastTouchpointType,
} from '@graphql/types';

export const LastTouchpointCell = ({
  lastTouchPointAt,
  lastTouchPointTimelineEvent,
  lastTouchPointType,
}: {
  lastTouchPointAt: string;
  lastTouchPointType: Maybe<LastTouchpointType> | undefined;
  lastTouchPointTimelineEvent: Maybe<TimelineEvent> | undefined;
}) => {
  const cellRef = useRef<HTMLDivElement>(null);

  const [label, Icon] = match(lastTouchPointTimelineEvent)
    .returnType<
      [string, (props: React.SVGAttributes<SVGElement>) => JSX.Element]
    >()
    .with({ __typename: 'Action', actionType: ActionType.Created }, () => [
      'Created',
      Building07,
    ])
    .with({ __typename: 'PageView' }, () => ['Page View', () => <></>])
    .with({ __typename: 'Issue' }, (event) => {
      const issueLastUpdateType =
        DateTimeUtils.differenceInMins(event.updatedAt, event.createdAt) > 1
          ? 'updated'
          : 'created';

      return [`Issue ${issueLastUpdateType}`, Ticket02];
    })
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
      `Email ${
        lastTouchPointType === LastTouchpointType.InteractionEventEmailSent
          ? 'sent'
          : 'received'
      }`,
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
    .with({ __typename: 'Meeting' }, (event) => [
      `Meeting with ${event.attendedBy.length} participant${
        event.attendedBy.length === 1 ? '' : 's'
      }`,
      Calendar,
    ])
    .otherwise(() => ['', () => <></>]);

  const subLabel = label
    ? DateTimeUtils.timeAgo(lastTouchPointAt, {
        strict: true,
        addSuffix: true,
        includeMin: true,
      })
    : '';

  return (
    <TableCellTooltip
      hasArrow
      align='start'
      side='bottom'
      targetRef={cellRef}
      label={`${label} • ${subLabel}`}
    >
      <span ref={cellRef}>
        <Icon className='size-3 min-w-3 text-gray-700' />
        <span
          className='text-gray-700 ml-2 leading-none'
          data-test='organization-last-touchpoint-in-all-orgs-table'
        >
          {label}
        </span>
        <span className='text-gray-500 text-xs ml-1 leading-none'>•</span>
        <span className='text-gray-500 text-xs ml-1 leading-none'>
          {subLabel}
        </span>
      </span>
    </TableCellTooltip>
  );
};
