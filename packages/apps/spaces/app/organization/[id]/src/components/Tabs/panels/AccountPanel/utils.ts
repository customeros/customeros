import { QueryClient } from '@tanstack/react-query';

import { DateTimeUtils } from '@spaces/utils/date';
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
export const billingFrequencyOptions: SelectOption<ContractRenewalCycle>[] = [
  { label: 'Monthly', value: ContractRenewalCycle.MonthlyRenewal },
  { label: 'Annually', value: ContractRenewalCycle.AnnualRenewal },
];

export const billedTypeOptions: SelectOption<BilledType>[] = [
  { label: 'Once', value: BilledType.Once },

  { label: 'Monthly', value: BilledType.Monthly },
  { label: 'Annually', value: BilledType.Annually },
];

export function calculateNextRenewalDate(
  serviceStartedAt: string,
  renewalCycle: ContractRenewalCycle,
): string {
  const difference = DateTimeUtils.differenceInMonths(
    new Date().toISOString(),
    serviceStartedAt,
  );
  switch (renewalCycle) {
    case ContractRenewalCycle.AnnualRenewal: {
      const differenceInYears = Math.ceil(difference / 12);

      let nextRenewal = DateTimeUtils.addYears(
        serviceStartedAt,
        differenceInYears,
      ).toISOString();
      if (!DateTimeUtils.isBeforeNow(nextRenewal)) {
        nextRenewal = DateTimeUtils.addYears(nextRenewal, 1).toISOString();
      }

      return nextRenewal;
    }
    case ContractRenewalCycle.MonthlyRenewal: {
      const difference = DateTimeUtils.differenceInMonths(
        new Date().toISOString(),
        serviceStartedAt,
      );
      let nextRenewal = DateTimeUtils.addMonth(serviceStartedAt, difference);

      if (!DateTimeUtils.isBeforeNow(nextRenewal.toISOString())) {
        nextRenewal = DateTimeUtils.addMonth(nextRenewal.toISOString(), 1);
      }

      return nextRenewal.toISOString();
    }
    default:
      return '';
  }
}
