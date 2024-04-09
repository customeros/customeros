import isBefore from 'date-fns/isBefore';
import { FilterFn } from '@tanstack/react-table';

import { Invoice } from '@graphql/types';

export const filterIssueDateFn: FilterFn<Invoice> = (row, id, filterValue) => {
  const value = row.getValue<Invoice['metadata']['created']>(id);

  if (!value) return false;

  const leftDate = new Date(value.split('T')[0]);
  const rightDate = new Date(filterValue);

  return isBefore(leftDate, rightDate);
};
filterIssueDateFn.autoRemove = (filterValue) => !filterValue;

export const filterIssueDatePastFn: FilterFn<Invoice> = (
  row,
  id,
  filterValue,
) => {
  const value = row.getValue<Invoice['metadata']['created']>(id);

  if (!value) return false;

  const leftDate = new Date(filterValue);
  const rightDate = new Date(value.split('T')[0]);

  return isBefore(leftDate, rightDate);
};
filterIssueDatePastFn.autoRemove = (filterValue) => !filterValue;
