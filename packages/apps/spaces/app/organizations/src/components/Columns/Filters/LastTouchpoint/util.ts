import isAfter from 'date-fns/isAfter';
import { FilterFn } from '@tanstack/react-table';

import { Organization, LastTouchpointType } from '@graphql/types';

export const touchpoints: { label: string; value: LastTouchpointType }[] = [
  { value: LastTouchpointType.ActionCreated, label: 'Organization created' },
  { value: LastTouchpointType.IssueCreated, label: 'Issue created' },
  { value: LastTouchpointType.IssueUpdated, label: 'Issue updated' },
  { value: LastTouchpointType.LogEntry, label: 'Log entry' },
  { value: LastTouchpointType.InteractionEventEmailSent, label: 'Email sent' },
  { value: LastTouchpointType.InteractionEventChat, label: 'Message received' },
  { value: LastTouchpointType.Meeting, label: 'Meeting' },
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
