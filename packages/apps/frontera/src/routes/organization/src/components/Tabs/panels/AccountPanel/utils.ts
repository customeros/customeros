import { QueryClient } from '@tanstack/react-query';

import { SelectOption } from '@shared/types/SelectOptions';
import { useOrganizationAccountDetailsQuery } from '@organization/graphql/getAccountPanelDetails.generated';
import {
  Maybe,
  BilledType,
  ContractBillingCycle,
  ContractRenewalCycle,
  OpportunityRenewalLikelihood,
} from '@graphql/types';

export const invalidateAccountDetailsQuery = (
  queryClient: QueryClient,
  id: string,
) =>
  queryClient.invalidateQueries({
    queryKey: useOrganizationAccountDetailsQuery.getKey({ id }),
  });

export function getRenewalLikelihoodColor(
  renewalLikelihood?: Maybe<OpportunityRenewalLikelihood> | undefined,
) {
  switch (renewalLikelihood) {
    case OpportunityRenewalLikelihood.HighRenewal:
      return 'greenLight';
    case OpportunityRenewalLikelihood.MediumRenewal:
      return 'yellow';
    case OpportunityRenewalLikelihood.LowRenewal:
      return 'orangeDark';
    case OpportunityRenewalLikelihood.ZeroRenewal:
      return 'gray';
    default:
      return 'gray';
  }
}
export function getRenewalLikelihoodLabel(
  renewalLikelihood?: Maybe<OpportunityRenewalLikelihood> | undefined,
) {
  switch (renewalLikelihood) {
    case OpportunityRenewalLikelihood.HighRenewal:
      return 'High';
    case OpportunityRenewalLikelihood.MediumRenewal:
      return 'Medium';
    case OpportunityRenewalLikelihood.LowRenewal:
      return 'Low';
    case OpportunityRenewalLikelihood.ZeroRenewal:
      return 'Zero';
    default:
      return '';
  }
}
export const contractRenewalCycle: SelectOption<
  ContractRenewalCycle | 'MULTI_YEAR'
>[] = [
  { label: 'Monthly', value: ContractRenewalCycle.MonthlyRenewal },
  { label: 'Quarterly', value: ContractRenewalCycle.QuarterlyRenewal },
  { label: 'Annual', value: ContractRenewalCycle.AnnualRenewal },
  { label: 'Multi-year', value: 'MULTI_YEAR' },
];
export const billingFrequencyOptions: SelectOption<ContractBillingCycle>[] = [
  { label: 'Monthly', value: ContractBillingCycle.MonthlyBilling },
  { label: 'Quarterly', value: ContractBillingCycle.QuarterlyBilling },
  { label: 'Annual', value: ContractBillingCycle.AnnualBilling },
];
export const contractBillingCycleOptions: SelectOption<ContractBillingCycle>[] =
  [
    { label: 'monthly', value: ContractBillingCycle.MonthlyBilling },
    { label: 'quarterly', value: ContractBillingCycle.QuarterlyBilling },
    { label: 'annually', value: ContractBillingCycle.AnnualBilling },
  ];

export const billedTypeOptions: SelectOption<BilledType>[] = [
  { label: 'once', value: BilledType.Once },
  { label: 'month', value: BilledType.Monthly },
  { label: 'quarter', value: BilledType.Quarterly },
  { label: 'year', value: BilledType.Annually },
];

export const currencyOptions: SelectOption<string>[] = [
  { label: 'United States Dollar', value: 'USD' },
];

export const paymentDueOptions: SelectOption<number>[] = [
  { label: '0 days', value: 0 },
  { label: '15 days', value: 15 },
  { label: '30 days', value: 30 },
  { label: '45 days', value: 45 },
  { label: '60 days', value: 60 },
  { label: '90 days', value: 90 },
];
