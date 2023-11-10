import { FilterFn } from '@tanstack/react-table';

import { Organization } from '@graphql/types';

export const filterOrganizationFn: FilterFn<Organization> = (
  row,
  id,
  filterValue,
) => {
  const value = row.getValue<Organization>(id);

  return value.name.toLowerCase().includes(filterValue.toLowerCase());
};

filterOrganizationFn.autoRemove = (filterValue) => !filterValue;
