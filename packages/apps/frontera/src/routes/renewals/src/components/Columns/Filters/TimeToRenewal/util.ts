import { isBefore } from 'date-fns/isBefore';
import { FilterFn } from '@tanstack/react-table';

import { RenewalRecord } from '@graphql/types';

export const filterTimeToRenewalFn: FilterFn<RenewalRecord> = (
  row,
  id,
  filterValue,
) => {
  const value = row.getValue<RenewalRecord['organization']>(id);
  const renewalDate = value?.accountDetails?.renewalSummary?.nextRenewalDate;

  if (!renewalDate) return false;

  const leftDate = new Date(renewalDate.split('T')[0]);
  const rightDate = new Date(filterValue);

  return isBefore(leftDate, rightDate);
};
filterTimeToRenewalFn.autoRemove = (filterValue) => !filterValue;
