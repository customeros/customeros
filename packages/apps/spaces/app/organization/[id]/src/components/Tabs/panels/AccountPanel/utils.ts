import { QueryClient } from '@tanstack/react-query';

import { SelectOption } from '@shared/types/SelectOptions';
import {
  Maybe,
  RenewalCycle,
  RenewalLikelihoodProbability,
} from '@graphql/types';
import { useOrganizationAccountDetailsQuery } from '@organization/src/graphql/getAccountPanelDetails.generated';

export const invalidateAccountDetailsQuery = (
  queryClient: QueryClient,
  id: string,
) =>
  queryClient.invalidateQueries(
    useOrganizationAccountDetailsQuery.getKey({ id }),
  );

export function getFeatureIconColor(
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
