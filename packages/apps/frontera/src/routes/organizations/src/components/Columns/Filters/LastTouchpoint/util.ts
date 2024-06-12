import { LastTouchpointType } from '@graphql/types';

export const touchpoints: { label: string; value: LastTouchpointType }[] = [
  { value: LastTouchpointType.InteractionEventEmailSent, label: 'Email sent' },
  { value: LastTouchpointType.IssueCreated, label: 'Issue created' },
  { value: LastTouchpointType.IssueUpdated, label: 'Issue updated' },
  { value: LastTouchpointType.LogEntry, label: 'Log entry' },
  { value: LastTouchpointType.Meeting, label: 'Meeting' },
  { value: LastTouchpointType.InteractionEventChat, label: 'Message received' },
  { value: LastTouchpointType.ActionCreated, label: 'Organization created' },
];

export const allTime = new Date('1970-01-01').toISOString().split('T')[0];
