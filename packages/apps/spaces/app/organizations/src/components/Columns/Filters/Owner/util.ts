import { FilterFn } from '@tanstack/react-table';

import { Organization } from '@graphql/types';

export const filterOwnerFn: FilterFn<Organization> = (row, id, filterValue) => {
  const value = row.getValue<Organization['owner']>(id)?.id;

  if (!value) return false;

  return filterValue.includes(value);
};

filterOwnerFn.autoRemove = (filterValue) => {
  return !filterValue || filterValue.length === 0;
};
