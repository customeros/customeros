import { FilterFn } from '@tanstack/react-table';

import { Invoice } from '@graphql/types';
import { ContractBillingCycle } from '@graphql/types';

const billingCycleOptions: Record<ContractBillingCycle, string> = {
  [ContractBillingCycle.MonthlyBilling]: 'MONTHLY',
  [ContractBillingCycle.QuarterlyBilling]: 'QUARTERLY',
  [ContractBillingCycle.AnnualBilling]: 'ANNUALY',
  [ContractBillingCycle.None]: '',
};

export const filterBillingCycleFn: FilterFn<Invoice> = (
  row,
  id,
  filterValue,
) => {
  const value = row.getValue<Invoice['contract']>(id)?.billingDetails
    ?.billingCycle as ContractBillingCycle;

  if (filterValue.length === 0) return true;

  return (filterValue as ContractBillingCycle[])
    .map((v) => billingCycleOptions?.[v])
    .includes(value);
};

filterBillingCycleFn.autoRemove = (filterValue) => {
  return !filterValue || filterValue.length === 0;
};
