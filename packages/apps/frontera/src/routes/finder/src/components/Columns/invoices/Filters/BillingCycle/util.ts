import { Store } from '@store/store.ts';
import { FilterFn } from '@tanstack/react-table';

import { Invoice } from '@graphql/types';
import { ContractBillingCycle } from '@graphql/types';

const billingCycleOptions: Record<ContractBillingCycle, string> = {
  [ContractBillingCycle.MonthlyBilling]: 'MONTHLY',
  [ContractBillingCycle.QuarterlyBilling]: 'QUARTERLY',
  [ContractBillingCycle.AnnualBilling]: 'ANNUALLY',
  [ContractBillingCycle.CustomBilling]: 'CUSTOM',
  [ContractBillingCycle.None]: '',
};

export const filterBillingCycleFn: FilterFn<Store<Invoice>> = (
  row,
  filterValue,
) => {
  const value = row.original?.value?.contract.metadata.id;
  const billingCycle =
    row.original?.root.contracts.value.get(value)?.value.billingDetails
      ?.billingCycle;

  if (!filterValue || filterValue.length === 0) return true;
  if (!billingCycle) return false;
  const selectedOption = billingCycleOptions?.[billingCycle];

  return filterValue.includes(selectedOption);
};

filterBillingCycleFn.autoRemove = (filterValue) => {
  return !filterValue || filterValue.length === 0;
};
