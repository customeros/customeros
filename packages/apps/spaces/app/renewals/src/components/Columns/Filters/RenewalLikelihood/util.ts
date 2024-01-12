import { FilterFn } from '@tanstack/react-table';

import { RenewalRecord } from '@graphql/types';

export const filterRenewalLikelihoodFn: FilterFn<RenewalRecord> = (
  row,
  id,
  filterValue,
) => {
  const value = row.getValue<RenewalRecord['organization']>(id);

  if (filterValue.length === 0) {
    return value?.accountDetails?.renewalSummary?.renewalLikelihood === null;
  }

  return filterValue.includes(
    value?.accountDetails?.renewalSummary?.renewalLikelihood,
  );
};
filterRenewalLikelihoodFn.autoRemove = (filterValue) => {
  return !filterValue;
};
