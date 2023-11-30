import isAfter from 'date-fns/isAfter';
import { FilterFn } from '@tanstack/react-table';

import { Organization, LastTouchpointType } from '@graphql/types';

export const touchpoints: { label: string; value: LastTouchpointType }[] = [
  { value: LastTouchpointType.InteractionEventEmailSent, label: 'Email sent' },
  { value: LastTouchpointType.IssueCreated, label: 'Issue created' },
  { value: LastTouchpointType.IssueUpdated, label: 'Issue updated' },
  { value: LastTouchpointType.LogEntry, label: 'Log entry' },
  { value: LastTouchpointType.Meeting, label: 'Meeting' },
  { value: LastTouchpointType.InteractionEventChat, label: 'Message received' },
  { value: LastTouchpointType.ActionCreated, label: 'Organization created' },
];

export const filterLastTouchpointFn: FilterFn<Organization> = (
  row,
  id,
  filterValue,
) => {
  const value = row.getValue<Organization>(id);
  const lastTouchpoint = value?.lastTouchPointType;
  const lastTouchpointAt = value?.lastTouchPointAt;

  const isIncluded = filterValue.value.length
    ? filterValue.value.includes(lastTouchpoint)
    : true;
  const isAfterDate = isAfter(
    new Date(lastTouchpointAt),
    new Date(filterValue.after),
  );

  return isIncluded && isAfterDate;
};

filterLastTouchpointFn.autoRemove = (filterValue) => {
  return !filterValue;
};

export const allTime = new Date('1970-01-01').toISOString().split('T')[0];
