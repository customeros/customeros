import { FilterFn } from '@tanstack/react-table';

import { Organization } from '@graphql/types';

export const filterRelationshipFn: FilterFn<Organization> = (
  row,
  id,
  filterValue,
) => {
  const value = row.getValue(id);

  if (filterValue.length === 2) {
    return filterValue.includes(value);
  }

  return filterValue[0] === value;
};
filterRelationshipFn.autoRemove = (filterValue) => {
  return !filterValue || filterValue.length === 0;
};
