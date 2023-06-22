import React from 'react';
import { Maybe, TimelineEvent } from '@spaces/graphql';
import { DateTimeUtils } from '../../../utils';
import { TableCell } from '@spaces/atoms/table';
import {
  Meeting,
  Notes,
  OutgoingEmail,
  UpdateOnIssue,
  OutgoingVoice,
} from '@spaces/atoms/icons';

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
      case 'PageView':
        label = 'Page View';
        break;
      case 'Issue':
        label = 'Issue';
        icon = <UpdateOnIssue width='24' height='24' viewBox='0 0 24 24' />;
        break;
      case 'Note': {
        label =
          'Note by ' +
          lastTouchPointTimelineEvent.createdBy?.firstName +
          ' ' +
          lastTouchPointTimelineEvent.createdBy?.lastName;
        icon = <Notes width='24' height='24' viewBox='0 0 24 24' />;
        break;
      }
      case 'InteractionEvent': {
        if (lastTouchPointTimelineEvent.channel === 'EMAIL') {
          label = 'Email sent';
          icon = <OutgoingEmail width='24' height='24' viewBox='0 0 24 24' />;
        } else if (lastTouchPointTimelineEvent.channel === 'VOICE') {
          label = 'Phone call';
          icon = <OutgoingVoice width='24' height='24' viewBox='0 0 24 24' />;
        } else if (
          !lastTouchPointTimelineEvent.channel &&
          lastTouchPointTimelineEvent.eventType === 'meeting'
        ) {
          label = '';
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
        icon = <Meeting width='24' height='24' viewBox='0 0 24 24' />;
        break;
      }
      default:
        console.log(
          'not able to print: ' + lastTouchPointTimelineEvent.__typename,
        );
        break;
    }
  }

  const subLabel = label
    ? DateTimeUtils.timeAgo(lastTouchPointAt, {
        addSuffix: true,
      })
    : '';
  return (
    <TableCell
      label={label}
      subLabel={subLabel}
      customStyleLabel={{ marginLeft: 4 }}
      customStyleSubLabel={{ marginLeft: 4 }}
    >
      {icon}
    </TableCell>
  );
};
