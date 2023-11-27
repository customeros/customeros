import { QueryClient } from '@tanstack/react-query';

import { SelectOption } from '@shared/types/SelectOptions';
import { useOrganizationAccountDetailsQuery } from '@organization/src/graphql/getAccountPanelDetails.generated';
import {
  Maybe,
  BilledType,
  RenewalCycle,
  ContractRenewalCycle,
  RenewalLikelihoodProbability,
} from '@graphql/types';

export const invalidateAccountDetailsQuery = (
  queryClient: QueryClient,
  id: string,
) =>
  queryClient.invalidateQueries(
    useOrganizationAccountDetailsQuery.getKey({ id }),
  );

export function getARRColor(
  renewalLikelihood?: Maybe<RenewalLikelihoodProbability> | undefined,
) {
  switch (renewalLikelihood) {
    case 'HIGH':
      return 'success';
    case 'MEDIUM':
      return 'warning';
    case 'LOW':
      return 'error';
    case 'ZERO':
      return 'gray';
    default:
      return 'gray';
  }
}

export const frequencyOptions: SelectOption<RenewalCycle>[] = [
  { label: 'Weekly', value: RenewalCycle.Weekly },
  { label: 'Biweekly', value: RenewalCycle.Biweekly },
  { label: 'Monthly', value: RenewalCycle.Monthly },
  { label: 'Quarterly', value: RenewalCycle.Quarterly },
  { label: 'Biannually', value: RenewalCycle.Biannually },
  { label: 'Annually', value: RenewalCycle.Annually },
];
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
