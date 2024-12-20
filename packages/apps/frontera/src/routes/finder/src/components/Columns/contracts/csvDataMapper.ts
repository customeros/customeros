import { ContractStore } from '@store/Contracts/Contract.store';
import { getRenewalLikelihoodLabel } from '@finder/components/Columns/contracts/Cells/health/utils.ts';

import { DateTimeUtils } from '@utils/date';
import { ColumnViewType, ContractStatus } from '@graphql/types';
import { formatCurrency } from '@utils/getFormattedCurrencyNumber.ts';

export const csvDataMapper = {
  [ColumnViewType.ContractsName]: (d: ContractStore) => {
    return d?.value?.contractName || d.organization?.value?.name;
  },

  [ColumnViewType.ContractsEnded]: (d: ContractStore) =>
    d?.value?.contractEnded
      ? DateTimeUtils.format(d.value.contractEnded, DateTimeUtils.iso8601)
      : 'No date yet',

  [ColumnViewType.ContractsPeriod]: (d: ContractStore) =>
    `${d?.value?.committedPeriodInMonths} months`,

  [ColumnViewType.ContractsCurrency]: (d: ContractStore) => d?.value?.currency,

  [ColumnViewType.ContractsStatus]: (d: ContractStore) =>
    d?.value?.contractStatus,

  [ColumnViewType.ContractsRenewal]: (d: ContractStore) =>
    d?.value?.autoRenew ? 'Auto-renews' : 'Not auto-renewing',

  [ColumnViewType.ContractsLtv]: (d: ContractStore) => d?.value?.ltv,

  [ColumnViewType.ContractsOwner]: (d: ContractStore) => {
    const opportunityId =
      d?.value?.opportunities?.find((e) => e.internalStage === 'OPEN')?.id ||
      d?.value?.opportunities?.[0]?.id;
    const opportunity = d.root.opportunities.value.get(opportunityId ?? '');

    const owner = opportunity?.value?.owner;

    return owner ? owner.name : 'No owner';
  },

  [ColumnViewType.ContractsRenewalDate]: (d: ContractStore) => {
    const renewsAt = d?.value?.opportunities?.find(
      (e) => e.internalStage === 'OPEN',
    )?.renewedAt;

    return renewsAt
      ? DateTimeUtils.format(renewsAt, DateTimeUtils.dateWithAbreviatedMonth)
      : 'No date yet';
  },

  [ColumnViewType.ContractsHealth]: (d: ContractStore) => {
    const opportunityId =
      d?.value?.opportunities?.find((e) => e.internalStage === 'OPEN')?.id ||
      d?.value?.opportunities?.[0]?.id;

    const opportunity = d.root.opportunities.value.get(opportunityId ?? '');
    const value = opportunity?.value?.renewalLikelihood;

    return value ? getRenewalLikelihoodLabel(value) : 'No set';
  },

  [ColumnViewType.ContractsForecastArr]: (d: ContractStore) => {
    const opportunityId =
      d?.value?.opportunities?.find((e) => e.internalStage === 'OPEN')?.id ||
      d?.value?.opportunities?.[0]?.id;

    const opportunity = d.root.opportunities.value.get(opportunityId ?? '');
    const hasEnded = d?.value?.contractStatus === ContractStatus.Ended;

    const formattedAmount = formatCurrency(
      hasEnded ? 0 : opportunity?.value?.amount ?? 0,
      2,
      'USD',
    );

    return formattedAmount;
  },
};
