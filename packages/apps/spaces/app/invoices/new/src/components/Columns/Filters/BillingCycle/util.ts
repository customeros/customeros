import { FilterFn } from '@tanstack/react-table';

import { Invoice } from '@graphql/types';

export const filterBillingCycleFn: FilterFn<Invoice> = (
  row,
  id,
  filterValue,
) => {
  const value = row.getValue<Invoice['contract']>(id)?.billingDetails
    ?.billingCycle as string;

  if (filterValue.length === 0) return true;

  return filterValue.includes(value);
};

filterBillingCycleFn.autoRemove = (filterValue) => {
  return !filterValue || filterValue.length === 0;
};
