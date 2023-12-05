import { QueryClient } from '@tanstack/react-query';

import { SelectOption } from '@shared/types/SelectOptions';
import { useOrganizationAccountDetailsQuery } from '@organization/src/graphql/getAccountPanelDetails.generated';
import {
  Maybe,
  BilledType,
  ContractRenewalCycle,
  OpportunityRenewalLikelihood,
} from '@graphql/types';

export const invalidateAccountDetailsQuery = (
  queryClient: QueryClient,
  id: string,
) =>
  queryClient.invalidateQueries(
    useOrganizationAccountDetailsQuery.getKey({ id }),
  );

export function getRenewalLikelihoodColor(
  renewalLikelihood?: Maybe<OpportunityRenewalLikelihood> | undefined,
) {
  switch (renewalLikelihood) {
    case OpportunityRenewalLikelihood.HighRenewal:
      return 'success';
    case OpportunityRenewalLikelihood.MediumRenewal:
      return 'warning';
    case OpportunityRenewalLikelihood.LowRenewal:
      return 'error';
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
export const billingFrequencyOptions: SelectOption<ContractRenewalCycle>[] = [
  { label: 'Monthly', value: ContractRenewalCycle.MonthlyRenewal },
  { label: 'Annually', value: ContractRenewalCycle.AnnualRenewal },
];

export const billedTypeOptions: SelectOption<BilledType>[] = [
  { label: 'Once', value: BilledType.Once },
  { label: 'Usage', value: BilledType.Usage },
  { label: 'Monthly', value: BilledType.Monthly },
  { label: 'Annually', value: BilledType.Annually },
];
