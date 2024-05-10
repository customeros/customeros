import { FilterFn } from '@tanstack/react-table';

import { Organization } from '@graphql/types';

export const filterRelationshipFn: FilterFn<Organization> = (
  row,
  id,
  filterValue,
) => {
  const value = row.getValue(id);

  if (filterValue.length === 0) return true;

  return filterValue.includes(value);
};

filterRelationshipFn.autoRemove = (filterValue) => {
  return !filterValue;
};
