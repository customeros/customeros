import { FilterFn } from '@tanstack/react-table';

import { Organization } from '@graphql/types';

export const filterRelationshipFn: FilterFn<Organization> = (
  row,
  id,
  filterValue,
) => {
  const value = row.getValue(id);

  return filterValue.includes(value);
};

filterRelationshipFn.autoRemove = (filterValue) => {
  return !filterValue;
};
