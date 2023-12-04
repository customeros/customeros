import { FilterFn } from '@tanstack/react-table';

import { Organization } from '@graphql/types';

export const filterForecastFn: FilterFn<Organization> = (
  row,
  id,
  filterValue,
) => {
  const value = row.getValue<Organization['accountDetails']>(id);
  const potentialValue = value?.renewalSummary?.arrForecast || 0;

  if (!potentialValue) return false;

  return potentialValue >= filterValue[0] && potentialValue <= filterValue[1];
};

filterForecastFn.autoRemove = (filterValue) => {
  return !filterValue || filterValue.length === 0;
};
