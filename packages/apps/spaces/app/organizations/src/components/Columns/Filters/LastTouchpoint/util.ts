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
  | 'ISSUE_CREATED'
  | 'ISSUE_UPDATED'
  | 'LOG_ENTRY'
  | 'EMAIL'
  | 'SLACK_MESSAGE'
  | 'INTERCOM_MESSAGE'
  | 'MEETING'
  | 'ANALYSIS';

export const touchpoints: { label: string; value: TouchPoint }[] = [
  { value: 'CREATED', label: 'Organization created' },
  { value: 'ISSUE_CREATED', label: 'Issue created' },
  { value: 'ISSUE_UPDATED', label: 'Issue updated' },
  { value: 'LOG_ENTRY', label: 'Log entry' },
  { value: 'EMAIL', label: 'Email sent' },
  { value: 'SLACK_MESSAGE', label: 'Slack received' },
  { value: 'INTERCOM_MESSAGE', label: 'Intercom received' },
  { value: 'MEETING', label: 'Meeting' },
];

const touchPointPatterns: Record<TouchPoint, Partial<TimelineEvent>> = {
  CREATED: { __typename: 'Action', actionType: ActionType.Created },
  ISSUE_CREATED: { __typename: 'Issue' },
  ISSUE_UPDATED: { __typename: 'Issue' },
  LOG_ENTRY: { __typename: 'LogEntry' },
  EMAIL: { __typename: 'InteractionEvent', channel: 'EMAIL' },
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
