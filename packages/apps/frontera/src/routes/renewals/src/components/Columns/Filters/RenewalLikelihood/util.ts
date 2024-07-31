import { FilterFn } from '@tanstack/react-table';

import { RenewalRecord, RenewalSummary } from '@graphql/types';

export const filterRenewalLikelihoodFn: FilterFn<RenewalRecord> = (
  row,
  id,
  filterValue,
) => {
  const value = row.getValue<{ renewalSummary: RenewalSummary }>(id);

  if (filterValue.length === 0) {
    return value?.renewalSummary?.renewalLikelihood === null;
  }

  return filterValue.includes(value?.renewalSummary?.renewalLikelihood);
};

filterRenewalLikelihoodFn.autoRemove = (filterValue) => {
  return !filterValue;
};
