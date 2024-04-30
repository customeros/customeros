import { FilterFn } from '@tanstack/react-table';

import { Organization } from '@graphql/types';

export const filterOnboardingFn: FilterFn<Organization> = (
  row,
  id,
  filterValue,
) => {
  const value = row.getValue<Organization['accountDetails']>(id)?.onboarding
    ?.status as string;

  if (filterValue.length === 0) return true;

  return filterValue.includes(value);
};

filterOnboardingFn.autoRemove = (filterValue) => {
  return !filterValue || filterValue.length === 0;
};
