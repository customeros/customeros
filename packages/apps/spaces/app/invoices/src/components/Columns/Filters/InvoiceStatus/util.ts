import { FilterFn } from '@tanstack/react-table';

import { Invoice } from '@graphql/types';

export const filterInvoiceStatusFn: FilterFn<Invoice> = (
  row,
  id,
  filterValue,
) => {
  const value = !!row.getValue<Invoice['contract']['contractEnded']>(
    id,
  ) as boolean;

  if (filterValue.length === 0 || filterValue.length === 2) return true;

  return (
    (filterValue[0] === 'ON_HOLD' && value) ||
    (filterValue[0] === 'SCHEDULED' && !value)
  );
};

filterInvoiceStatusFn.autoRemove = (filterValue) => {
  return !filterValue || filterValue.length === 0;
};
