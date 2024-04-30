import { FilterFn } from '@tanstack/react-table';

import { RenewalRecord } from '@graphql/types';

/**
 * Filter function for the ARR Forecast column
 * used to optimistically filter the table while waiting for the server to respond.
 * This needs to be kept in sync with the server-side filtering logic
 */
export const filterForecastFn: FilterFn<RenewalRecord> = (
  row,
  id,
  filterValue,
) => {
  const value = row.getValue<RenewalRecord['organization']>(id);
  const potentialValue =
    value?.accountDetails?.renewalSummary?.arrForecast || 0;

  if (!potentialValue) return false;

  return potentialValue >= filterValue[0] && potentialValue <= filterValue[1];
};

filterForecastFn.autoRemove = (filterValue) => {
  return !filterValue || filterValue.length === 0;
};
