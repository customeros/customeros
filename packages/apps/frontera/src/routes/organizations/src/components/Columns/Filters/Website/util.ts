import { Store } from '@store/store';
import { FilterFn } from '@tanstack/react-table';

import { Organization } from '@graphql/types';

export const filterWebsiteFn: FilterFn<Store<Organization>> = (
  row,
  id,
  filterValue,
) => {
  const value = row.getValue<Store<Organization>['value']['website']>(id);

  if (filterValue.showEmpty && !value) return true;
  if (!value) return false;

  return value.toLowerCase().includes(filterValue.value.toLowerCase());
};

filterWebsiteFn.autoRemove = (filterValue) => !filterValue;
