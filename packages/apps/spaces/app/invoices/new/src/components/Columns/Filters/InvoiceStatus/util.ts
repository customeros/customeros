import { FilterFn } from '@tanstack/react-table';

import { Invoice } from '@graphql/types';

export const filterInvoiceStatusFn: FilterFn<Invoice> = (
  row,
  id,
  filterValue,
) => {
  const value = row.getValue<Invoice['contract']['contractEnded']>(
    id,
  ) as boolean;

  if (filterValue.length === 0) return true;

  return (filterValue as ('ON_HOLD' | 'SCHEDULED')[])
    .map((v) => (v === 'ON_HOLD' ? true : false))
    .includes(value);
};

filterInvoiceStatusFn.autoRemove = (filterValue) => {
  return !filterValue || filterValue.length === 0;
};
