import React, { useRef } from 'react';

import { match } from 'ts-pattern';

import { DateTimeUtils } from '@utils/date';
import { Mail01 } from '@ui/media/icons/Mail01';
import { Calendar } from '@ui/media/icons/Calendar';
import { Ticket02 } from '@ui/media/icons/Ticket02';
import { Building07 } from '@ui/media/icons/Building07';
import { TableCellTooltip } from '@ui/presentation/Table';
import { Maybe, LastTouchpointType } from '@graphql/types';
import { PhoneOutgoing02 } from '@ui/media/icons/PhoneOutgoing02';
import { MessageChatSquare } from '@ui/media/icons/MessageChatSquare';

export const LastTouchpointCell = ({
  lastTouchPointAt,
  lastTouchPointType,
}: {
  lastTouchPointAt: string;
  lastTouchPointType: Maybe<LastTouchpointType> | undefined;
}) => {
  const cellRef = useRef<HTMLDivElement>(null);

  const [label, Icon] = match(lastTouchPointType)
    .returnType<
      [string, (props: React.SVGAttributes<SVGElement>) => JSX.Element]
    >()
    .with(LastTouchpointType.ActionCreated, () => ['Created', Building07])
    .with(LastTouchpointType.PageView, () => ['Page View', () => <></>])
    .with(LastTouchpointType.IssueCreated, () => {
      return [`Issue created`, Ticket02];
    })
    .with(LastTouchpointType.IssueUpdated, () => {
      return [`Issue updated`, Ticket02];
    })
    .with(LastTouchpointType.LogEntry, () => [`Log entry`, MessageChatSquare])
    .with(LastTouchpointType.InteractionEventEmailSent, () => [
      `Email sent`,
      Mail01,
    ])
    .with(LastTouchpointType.InteractionEventEmailReceived, () => [
      `Email received`,
      Mail01,
    ])
    .with(LastTouchpointType.InteractionEventPhoneCall, () => [
      'Phone call',
      PhoneOutgoing02,
    ])
    .with(LastTouchpointType.InteractionEventChat, () => [
      'Message',
      MessageChatSquare,
    ])

    .with(LastTouchpointType.Meeting, () => ['Meeting', Calendar])

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
