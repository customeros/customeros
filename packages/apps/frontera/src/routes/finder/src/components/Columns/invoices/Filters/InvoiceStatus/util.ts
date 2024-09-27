import { Store } from '@store/store.ts';
import { FilterFn } from '@tanstack/react-table';

import { Invoice } from '@graphql/types';

export const filterInvoiceStatusFn: FilterFn<Store<Invoice>> = (
  row,
  filterValue,
) => {
  const data = row.original?.value?.contract.metadata.id;
  const value =
    row.original?.root.contracts.value.get(data)?.value.contractEnded;

  if (filterValue.length === 0 || filterValue.length === 2) return true;

  return (
    (filterValue[0] === 'ON_HOLD' && value) ||
    (filterValue[0] === 'SCHEDULED' && !value)
  );
};

filterInvoiceStatusFn.autoRemove = (filterValue) => {
  return !filterValue || filterValue.length === 0;
};
