import { FilterFn } from '@tanstack/react-table';

import { RenewalRecord } from '@graphql/types';

export const filterOwnerFn: FilterFn<RenewalRecord> = (
  row,
  id,
  filterValue,
) => {
  const value = row.getValue<RenewalRecord['contract']>(id)?.owner?.id;

  if (filterValue?.showEmpty && !value) return true;
  if (!value) return false;

  return filterValue?.value?.includes(value);
};

filterOwnerFn.autoRemove = (filterValue) => {
  return !filterValue;
};
