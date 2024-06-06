import { Store } from '@store/store.ts';
import { isBefore } from 'date-fns/isBefore';
import { FilterFn } from '@tanstack/react-table';

import { Invoice } from '@graphql/types';

export const filterIssueDateFn: FilterFn<Store<Invoice>> = (
  row,
  id,
  filterValue,
) => {
  const value = row.getValue<Store<Invoice>>(id)?.value?.issued;

  if (!value) return false;

  const leftDate = new Date(value.split('T')[0]);
  const rightDate = new Date(filterValue);

  return isBefore(leftDate, rightDate);
};
filterIssueDateFn.autoRemove = (filterValue) => !filterValue;

export const filterIssueDatePastFn: FilterFn<Store<Invoice>> = (
  row,
  id,
  filterValue,
) => {
  const value = row.getValue<Store<Invoice>>(id)?.value?.issued;

  if (!value) return false;

  const leftDate = new Date(filterValue);
  const rightDate = new Date(value.split('T')[0]);

  return isBefore(leftDate, rightDate);
};
filterIssueDatePastFn.autoRemove = (filterValue) => !filterValue;
