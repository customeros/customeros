import { Store } from '@store/store';
import { isBefore } from 'date-fns/isBefore';
import { FilterFn } from '@tanstack/react-table';

import { Organization } from '@graphql/types';

export const filterTimeToRenewalFn: FilterFn<Store<Organization>> = (
  row,
  id,
  filterValue,
) => {
  const value =
    row.getValue<Store<Organization>['value']['accountDetails']>(id);
  const renewalDate = value?.renewalSummary?.nextRenewalDate;

  if (!renewalDate) return false;

  const leftDate = new Date(renewalDate.split('T')[0]);
  const rightDate = new Date(filterValue);

  return isBefore(leftDate, rightDate);
};
filterTimeToRenewalFn.autoRemove = (filterValue) => !filterValue;
