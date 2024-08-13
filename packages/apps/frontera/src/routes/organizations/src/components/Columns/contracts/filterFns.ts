import { match } from 'ts-pattern';
import { Filter, FilterItem } from '@store/types.ts';
import { ContractStore } from '@store/Contracts/Contract.store.ts';

import { ColumnViewType } from '@graphql/types';

const getFilterFn = (filter: FilterItem | undefined | null) => {
  const noop = (_row: ContractStore) => true;

  if (!filter) return noop;

  return match(filter)
    .with(
      { property: ColumnViewType.ContractsName },
      (filter) => (row: ContractStore) => {
        if (!filter?.active) return true;
        const filterValue = filter?.value;

        if (!filterValue && filter.active && !filter.includeEmpty) return true;
        if (!row.value?.contractName?.length && filter.includeEmpty)
          return true;
        if (!filterValue || !row.value?.contractName?.length) return false;

        return row.value.contractName
          .toLowerCase()
          .includes(filterValue.toLowerCase());
      },
    )
    .with(
      { property: ColumnViewType.ContractsSignDate },
      (filter) => (row: ContractStore) => {
        if (!filter.active) return true;
        const filterValue = filter?.value;
        const nextRenewalDate = row.value?.contractSigned?.split('T')?.[0];

        if (!filterValue) return true;
        if (filterValue?.[1] === null)
          return filterValue?.[0] <= nextRenewalDate;
        if (filterValue?.[0] === null)
          return filterValue?.[1] >= nextRenewalDate;

        return (
          filterValue[0] <= nextRenewalDate && filterValue[1] >= nextRenewalDate
        );
      },
    )
    .with(
      { property: ColumnViewType.ContractsEnded },
      (filter) => (row: ContractStore) => {
        if (!filter.active) return true;
        const filterValue = filter?.value;
        const nextRenewalDate = row.value?.contractEnded?.split('T')?.[0];

        if (!filterValue) return true;
        if (filterValue?.[1] === null)
          return filterValue?.[0] <= nextRenewalDate;
        if (filterValue?.[0] === null)
          return filterValue?.[1] >= nextRenewalDate;

        return (
          filterValue[0] <= nextRenewalDate && filterValue[1] >= nextRenewalDate
        );
      },
    )
    .with(
      { property: ColumnViewType.ContractsCurrency },
      (filter) => (row: ContractStore) => {
        if (!filter.active) return true;
        const filterValue = filter?.value;
        const currency = row.value?.currency;

        if (!filterValue) return true;

        return filterValue.includes(currency);
      },
    )
    .with(
      { property: ColumnViewType.ContractsStatus },
      (filter) => (row: ContractStore) => {
        if (!filter.active) return true;
        const filterValue = filter?.value;
        const status = row.value?.contractStatus;

        if (!filterValue) return true;

        return filterValue.includes(status);
      },
    )
    .with(
      { property: ColumnViewType.ContractsRenewal },
      (filter) => (row: ContractStore) => {
        if (!filter.active) return true;
        const filterValue = filter?.value;
        const renewalStatus = row.value?.autoRenew;

        if (!filterValue) return true;

        return filterValue.includes(renewalStatus);
      },
    )
    .with(
      { property: ColumnViewType.ContractsLtv },
      (filter) => (row: ContractStore) => {
        if (!filter.active) return true;
        const filterValue = filter?.value;
        const forecastValue = row.value?.ltv;

        if (!forecastValue) return false;

        return (
          forecastValue >= filterValue[0] && forecastValue <= filterValue[1]
        );
      },
    )

    .otherwise(() => noop);
};

export const getContractFilterFns = (filters: Filter | null) => {
  if (!filters || !filters.AND) return [];
  const data = filters?.AND;

  return data.map(({ filter }) => getFilterFn(filter));
};
