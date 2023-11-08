import { FilterFn } from '@tanstack/react-table';

import { Organization } from '@graphql/types';

export const filterRenewalLikelihoodFn: FilterFn<Organization> = (
  row,
  id,
  filterValue,
) => {
  const value = row.getValue<Organization['accountDetails']>(id);

  if (filterValue.length === 0) {
    return value?.renewalLikelihood?.probability === null;
  }

  return filterValue.includes(value?.renewalLikelihood?.probability);
};
filterRenewalLikelihoodFn.autoRemove = (filterValue) => {
  return !filterValue;
};
