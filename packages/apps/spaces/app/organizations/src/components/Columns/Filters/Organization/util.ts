import { FilterFn } from '@tanstack/react-table';

import { Organization } from '@graphql/types';

export const filterOrganizationFn: FilterFn<Organization> = (
  row,
  id,
  filterValue,
) => {
  const value = row.getValue<Organization>(id);

  if (filterValue.showEmpty && !value.name) return true;

  return value.name.toLowerCase().includes(filterValue.value.toLowerCase());
};

filterOrganizationFn.autoRemove = (filterValue) => !filterValue;
