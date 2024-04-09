import { FilterFn } from '@tanstack/react-table';

import { Invoice } from '@graphql/types';

export const filterPaymentStatusFn: FilterFn<Invoice> = (
  row,
  id,
  filterValue,
) => {
  const value = row.getValue<Invoice>(id).status;

  if (filterValue.length === 0) return true;

  return filterValue.includes(value);
};

filterPaymentStatusFn.autoRemove = (filterValue) => {
  return !filterValue || filterValue.length === 0;
};
