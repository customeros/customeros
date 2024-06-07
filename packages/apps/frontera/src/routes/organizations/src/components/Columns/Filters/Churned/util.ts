import { Store } from '@store/store';
import { FilterFn } from '@tanstack/react-table';

import { Organization } from '@graphql/types';

export const filterChurnedFn: FilterFn<Store<Organization>> = (row, id) => {
  const value =
    row.getValue<Store<Organization>['value']['accountDetails']>(id);
  const churnedDate = value?.churned;
  if (!churnedDate) return false;

  return true;
};
filterChurnedFn.autoRemove = (filterValue) => !filterValue;
