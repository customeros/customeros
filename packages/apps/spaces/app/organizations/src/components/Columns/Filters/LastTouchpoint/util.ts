import { match } from 'ts-pattern';
import { FilterFn } from '@tanstack/react-table';

import {
  ActionType,
  Organization,
  TimelineEvent,
  ExternalSystemType,
} from '@graphql/types';

export type TouchPoint =
  | 'CREATED'
  | 'PAGE_VIEW'
  | 'ISSUE_CREATED'
  | 'ISSUE_UPDATED'
  | 'NOTE'
  | 'LOG_ENTRY'
  | 'EMAIL'
  | 'PHONE_CALL'
  | 'SLACK_MESSAGE'
  | 'INTERCOM_MESSAGE'
  | 'MEETING'
  | 'ANALYSIS';

export const touchpoints: { label: string; value: TouchPoint }[] = [
  { value: 'CREATED', label: 'Created' },
  { value: 'PAGE_VIEW', label: 'Page View' },
  { value: 'ISSUE_CREATED', label: 'Issue Created' },
  { value: 'ISSUE_UPDATED', label: 'Issue Updated' },
  { value: 'NOTE', label: 'Note' },
  { value: 'LOG_ENTRY', label: 'Log Entry' },
  { value: 'EMAIL', label: 'Email Sent' },
  { value: 'PHONE_CALL', label: 'Phone Call' },
  { value: 'SLACK_MESSAGE', label: 'Slack Message' },
  { value: 'INTERCOM_MESSAGE', label: 'Intercom Message' },
  { value: 'MEETING', label: 'Meeting' },
];

const touchPointPatterns: Record<TouchPoint, Partial<TimelineEvent>> = {
  CREATED: { __typename: 'Action', actionType: ActionType.Created },
  PAGE_VIEW: { __typename: 'PageView' },
  ISSUE_CREATED: { __typename: 'Issue' },
  ISSUE_UPDATED: { __typename: 'Issue' },
  NOTE: { __typename: 'Note' },
  LOG_ENTRY: { __typename: 'LogEntry' },
  EMAIL: { __typename: 'InteractionEvent', channel: 'EMAIL' },
  PHONE_CALL: { __typename: 'InteractionEvent', channel: 'VOICE' },
  SLACK_MESSAGE: {
    __typename: 'InteractionEvent',
    channel: 'CHAT',
    externalLinks: [{ type: ExternalSystemType.Slack }],
  },
  INTERCOM_MESSAGE: {
    __typename: 'InteractionEvent',
    channel: 'CHAT',
    externalLinks: [{ type: ExternalSystemType.Intercom }],
  },
  MEETING: { __typename: 'InteractionEvent', eventType: 'meeting' },
  ANALYSIS: { __typename: 'Analysis' },
};

const testPattern = (data: TimelineEvent, touchpoint: TouchPoint) => {
  return match(data)
    .returnType<boolean>()
    .with(touchPointPatterns[touchpoint] as object, () => true)
    .otherwise(() => false);
};

export const filterLastTouchpointFn: FilterFn<Organization> = (
  row,
  id,
  filterValue,
) => {
  const value = row.getValue<Organization>(id)?.lastTouchPointTimelineEvent;

  if (!value) return false;

  return (filterValue as TouchPoint[]).some((touchpoint) =>
    testPattern(value, touchpoint),
  );
};

filterLastTouchpointFn.autoRemove = (filterValue) => {
  return !filterValue || filterValue.length === 0;
};
