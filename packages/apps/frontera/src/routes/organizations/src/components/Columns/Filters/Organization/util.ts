import { Store } from '@store/store';
import { FilterFn } from '@tanstack/react-table';

import { Organization } from '@graphql/types';

export const filterOrganizationFn: FilterFn<Store<Organization>> = (
  row,
  id,
  filterValue,
) => {
  const value = row.getValue<Store<Organization>>(id).value;

  if (filterValue.showEmpty && !value.name) return true;

  return value.name.toLowerCase().includes(filterValue.value.toLowerCase());
};

filterOrganizationFn.autoRemove = (filterValue) => !filterValue;
