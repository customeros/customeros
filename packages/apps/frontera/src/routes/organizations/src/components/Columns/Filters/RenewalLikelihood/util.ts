import { FilterFn } from '@tanstack/react-table';

import { Organization } from '@graphql/types';

export const filterRenewalLikelihoodFn: FilterFn<Organization> = (
  row,
  id,
  filterValue,
) => {
  const value = row.getValue<Organization['accountDetails']>(id);

  if (filterValue.length === 0) {
    return value?.renewalSummary?.renewalLikelihood === null;
  }

  return filterValue.includes(value?.renewalSummary?.renewalLikelihood);
};
filterRenewalLikelihoodFn.autoRemove = (filterValue) => {
  return !filterValue;
};
