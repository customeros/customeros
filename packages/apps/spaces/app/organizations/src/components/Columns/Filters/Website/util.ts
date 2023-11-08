import { FilterFn } from '@tanstack/react-table';

import { Organization } from '@graphql/types';

export const filterWebsiteFn: FilterFn<Organization> = (
  row,
  id,
  filterValue,
) => {
  const value = row.getValue<Organization['website']>(id);

  if (!value) return false;

  return value.toLowerCase().includes(filterValue.toLowerCase());
};

filterWebsiteFn.autoRemove = (filterValue) => !filterValue;
