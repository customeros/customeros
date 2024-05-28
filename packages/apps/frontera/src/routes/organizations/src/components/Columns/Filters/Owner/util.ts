import { Store } from '@store/store';
import { FilterFn } from '@tanstack/react-table';

import { Organization } from '@graphql/types';

export const filterOwnerFn: FilterFn<Store<Organization>> = (
  row,
  id,
  filterValue,
) => {
  const value = row.getValue<Store<Organization>['value']['owner']>(id)?.id;

  if (filterValue?.showEmpty && !value) return true;
  if (!value) return false;

  return filterValue?.value?.includes(value);
};

filterOwnerFn.autoRemove = (filterValue) => {
  return !filterValue;
};
