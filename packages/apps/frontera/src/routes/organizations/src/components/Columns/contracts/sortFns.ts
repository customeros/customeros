import { match } from 'ts-pattern';
import { ContractStore } from '@store/Contracts/Contract.store';

import { ColumnViewType } from '@graphql/types';

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

    .otherwise(() => (_row: ContractStore) => false);
