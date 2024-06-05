import { Store } from '@store/store';
import { FilterFn } from '@tanstack/react-table';

import { Invoice } from '@graphql/types';

export const filterPaymentStatusFn: FilterFn<Store<Invoice>> = (
  row,
  filterValue,
) => {
  const value = row?.original?.value?.status;

  if (filterValue.length === 0) return true;

  return filterValue.includes(value as string);
};

filterPaymentStatusFn.autoRemove = (filterValue) => {
  return !filterValue || filterValue.length === 0;
};
