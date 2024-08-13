import { match } from 'ts-pattern';
import { ContractStore } from '@store/Contracts/Contract.store';

import {
  Opportunity,
  ColumnViewType,
  OpportunityRenewalLikelihood,
} from '@graphql/types';

export const getContractSortFn = (columnId: string) =>
  match(columnId)
    .with(ColumnViewType.ContractsName, () => (row: ContractStore) => {
      return row.value?.contractName?.trim().toLowerCase() || null;
    })
    .with(
      ColumnViewType.ContractsPeriod,
      () => (row: ContractStore) => row.value?.committedPeriodInMonths,
    )
    .with(
      ColumnViewType.ContractsSignDate,
      () => (row: ContractStore) => row.value?.signedAt || null,
    )
    .with(ColumnViewType.ContractsCurrency, () => (row: ContractStore) => {
      return row.value?.currency?.toLowerCase() || null;
    })
    .with(
      ColumnViewType.ContractsEnded,
      () => (row: ContractStore) => row.value.contractEnded || null,
    )
    .with(ColumnViewType.ContractsStatus, () => (row: ContractStore) => {
      return row.value.contractStatus?.toLowerCase();
    })
    .with(
      ColumnViewType.ContractsRenewal,
      () => (row: ContractStore) => row.value.autoRenew,
    )
    .with(
      ColumnViewType.ContractsLtv,
      () => (row: ContractStore) => row.value.ltv,
    )
    .with(ColumnViewType.ContractsRenewalDate, () => (row: ContractStore) => {
      const renewsAt = row?.value?.opportunities?.find(
        (e: Opportunity) => e.internalStage === 'OPEN',
      )?.renewedAt;

      return renewsAt || null;
    })
    .with(ColumnViewType.ContractsRenewalDate, () => (row: ContractStore) => {
      const renewsAt = row?.value?.opportunities?.find(
        (e: Opportunity) => e.internalStage === 'OPEN',
      )?.renewedAt;

      return renewsAt || null;
    })
    .with(ColumnViewType.ContractsHealth, () => (row: ContractStore) => {
      const renewalLikelihood = row?.value?.opportunities?.find(
        (e: Opportunity) => e.internalStage === 'OPEN',
      )?.renewalLikelihood;

      return match(renewalLikelihood)
        .with(OpportunityRenewalLikelihood.HighRenewal, () => 3)
        .with(OpportunityRenewalLikelihood.MediumRenewal, () => 2)
        .with(OpportunityRenewalLikelihood.LowRenewal, () => 1)
        .otherwise(() => null);
    })
    .with(ColumnViewType.ContractsOwner, () => (row: ContractStore) => {
      const owner = row?.value?.opportunities?.find(
        (e: Opportunity) => e.internalStage === 'OPEN',
      )?.owner;

      const name = owner?.name ?? '';
      const firstName = owner?.firstName ?? '';
      const lastName = owner?.lastName ?? '';

      const fullName = (name ?? `${firstName} ${lastName}`).trim();

      return fullName.length ? fullName.toLocaleLowerCase() : null;
    })

    .with(ColumnViewType.ContractsForecastArr, () => (row: ContractStore) => {
      const amount = row?.value?.opportunities?.find(
        (e: Opportunity) => e.internalStage === 'OPEN',
      )?.amount;

      return amount || null;
    })

    .otherwise(() => (_row: ContractStore) => false);
