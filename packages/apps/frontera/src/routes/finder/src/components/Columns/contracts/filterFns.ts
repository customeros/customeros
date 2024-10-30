import { match } from 'ts-pattern';
import { isAfter, isBefore } from 'date-fns';
import { Filter, FilterItem } from '@store/types.ts';
import { ContractStore } from '@store/Contracts/Contract.store.ts';

import { ColumnViewType, ComparisonOperator } from '@graphql/types';

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
      { property: ColumnViewType.ContractsRenewalDate },
      (filter) => (row: ContractStore) => {
        if (!filter.active) return true;
        const filterValue = filter?.value;
        const nextRenewalDate = row?.openOpportunity?.renewedAt;

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
    .with(
      { property: ColumnViewType.ContractsOwner },
      (filter) => (row: ContractStore) => {
        if (!filter.active) return true;
        const owner = row?.openOpportunity?.owner;

        const filterValue = filter?.value;

        if (filterValue === '__EMPTY__' && !owner) {
          return true;
        }

        return filterValue.includes(owner?.id);
      },
    )
    .with(
      { property: ColumnViewType.ContractsHealth },
      (filter) => (row: ContractStore) => {
        if (!filter.active) return true;
        const renewalLikelihood = row?.openOpportunity?.renewalLikelihood;

        if (!filter.active) return true;
        const filterValue = filter?.value;

        return filterValue.includes(renewalLikelihood);
      },
    )
    .with(
      { property: ColumnViewType.ContractsForecastArr },
      (filter) => (row: ContractStore) => {
        if (!filter.active) return true;
        const filterValue = filter?.value;
        const forecastValue = row?.openOpportunity?.amount;

        if (!forecastValue) return false;

        return (
          forecastValue >= filterValue[0] && forecastValue <= filterValue[1]
        );
      },
    )
    .with(
      { property: ColumnViewType.ContractsPeriod },
      (filter) => (row: ContractStore) => {
        if (!filter.active) return true;
        const filterValue = filter?.value;
        const committedPeriodInMonths = row?.value.committedPeriodInMonths;

        return filterValue.includes(committedPeriodInMonths);
      },
    )
    .otherwise(() => noop);
};

const getFilterV2Fn = (filter: FilterItem | undefined | null) => {
  const noop = (_row: ContractStore) => true;

  if (!filter) return noop;

  return match(filter)
    .with(
      { property: ColumnViewType.ContractsName },
      (filter) => (row: ContractStore) => {
        if (!filter?.active) return true;
        const value = row.value?.contractName;

        return filterTypeText(filter, value);
      },
    )
    .with(
      { property: ColumnViewType.ContractsEnded },
      (filter) => (row: ContractStore) => {
        if (!filter.active) return true;
        const value = row.value?.contractEnded?.split('T')?.[0];

        return filterTypeDate(filter, value);
      },
    )
    .with(
      { property: ColumnViewType.ContractsRenewalDate },
      (filter) => (row: ContractStore) => {
        if (!filter.active) return true;
        const value = row?.openOpportunity?.renewedAt;

        return filterTypeDate(filter, value);
      },
    )
    .with(
      { property: ColumnViewType.ContractsCurrency },
      (filter) => (row: ContractStore) => {
        if (!filter.active) return true;
        const value = row.value?.currency;

        if (!value)
          return (
            filter.operation === ComparisonOperator.IsEmpty ||
            filter.operation === ComparisonOperator.NotContains
          );

        return filterTypeList(filter, [value]);
      },
    )
    .with(
      { property: ColumnViewType.ContractsStatus },
      (filter) => (row: ContractStore) => {
        if (!filter.active) return true;
        const value = row.value?.contractStatus;

        if (!value)
          return (
            filter.operation === ComparisonOperator.IsEmpty ||
            filter.operation === ComparisonOperator.NotContains
          );

        return filterTypeList(filter, [value]);
      },
    )
    .with(
      { property: ColumnViewType.ContractsRenewal },
      (filter) => (row: ContractStore) => {
        if (!filter.active) return true;
        const value =
          row.value?.autoRenew === true ? 'Auido-renews' : 'Not auto-renewing';

        if (!value)
          return (
            filter.operation === ComparisonOperator.IsEmpty ||
            filter.operation === ComparisonOperator.NotContains
          );

        return filterTypeList(filter, [value]);
      },
    )
    .with(
      { property: ColumnViewType.ContractsLtv },
      (filter) => (row: ContractStore) => {
        if (!filter.active) return true;
        const value = row.value?.ltv;

        if (value !== undefined || value !== null) return false;

        return filterTypeNumber(filter, value);
      },
    )
    .with(
      { property: ColumnViewType.ContractsOwner },
      (filter) => (row: ContractStore) => {
        if (!filter.active) return true;
        const value = row?.openOpportunity?.owner?.id;

        if (!value)
          return (
            filter.operation === ComparisonOperator.IsEmpty ||
            filter.operation === ComparisonOperator.NotContains
          );

        return filterTypeList(filter, [value]);
      },
    )
    .with(
      { property: ColumnViewType.ContractsHealth },
      (filter) => (row: ContractStore) => {
        if (!filter.active) return true;
        const value = row?.openOpportunity?.renewalLikelihood;

        if (!value)
          return (
            filter.operation === ComparisonOperator.IsEmpty ||
            filter.operation === ComparisonOperator.NotContains
          );

        return filterTypeList(filter, [value]);
      },
    )
    .with(
      { property: ColumnViewType.ContractsForecastArr },
      (filter) => (row: ContractStore) => {
        if (!filter.active) return true;
        const value = row?.openOpportunity?.amount;

        if (value !== undefined || value !== null) return false;

        return filterTypeNumber(filter, value);
      },
    )
    .with(
      { property: ColumnViewType.ContractsPeriod },
      (filter) => (row: ContractStore) => {
        if (!filter.active) return true;
        const value = row?.value.committedPeriodInMonths;

        if (!value)
          return (
            filter.operation === ComparisonOperator.IsEmpty ||
            filter.operation === ComparisonOperator.NotContains
          );

        return filterTypeList(filter, [value]);
      },
    )
    .otherwise(() => noop);
};

const filterTypeText = (filter: FilterItem, value: string | undefined) => {
  const filterValue = filter?.value?.toLowerCase();
  const filterOperator = filter?.operation;
  const valueLower = value?.toLowerCase();

  return match(filterOperator)
    .with(ComparisonOperator.IsEmpty, () => !value)
    .with(ComparisonOperator.IsNotEmpty, () => value)
    .with(
      ComparisonOperator.NotContains,
      () => !valueLower?.includes(filterValue),
    )
    .with(ComparisonOperator.Contains, () => valueLower?.includes(filterValue))
    .otherwise(() => false);
};

const filterTypeNumber = (filter: FilterItem, value: number | undefined) => {
  const filterValue = filter?.value;
  const filterOperator = filter?.operation;

  if (value === undefined || value === null) return false;

  return match(filterOperator)
    .with(ComparisonOperator.Lt, () => value < Number(filterValue))
    .with(ComparisonOperator.Gt, () => value > Number(filterValue))
    .with(ComparisonOperator.Eq, () => value === Number(filterValue))
    .with(ComparisonOperator.NotEqual, () => value !== Number(filterValue))
    .otherwise(() => true);
};

const filterTypeList = (filter: FilterItem, value: string[] | undefined) => {
  const filterValue = filter?.value;
  const filterOperator = filter?.operation;

  return match(filterOperator)
    .with(ComparisonOperator.IsEmpty, () => !value?.length)
    .with(ComparisonOperator.IsNotEmpty, () => value?.length)
    .with(
      ComparisonOperator.NotContains,
      () =>
        !value?.length ||
        (value?.length && !value.some((v) => filterValue?.includes(v))),
    )
    .with(
      ComparisonOperator.Contains,
      () => value?.length && value.some((v) => filterValue?.includes(v)),
    )
    .otherwise(() => false);
};

const filterTypeDate = (filter: FilterItem, value: string | undefined) => {
  const filterValue = filter?.value;
  const filterOperator = filter?.operation;

  if (!value) return false;

  return match(filterOperator)
    .with(ComparisonOperator.Lt, () =>
      isBefore(new Date(value), new Date(filterValue)),
    )
    .with(ComparisonOperator.Gt, () =>
      isAfter(new Date(value), new Date(filterValue)),
    )

    .otherwise(() => true);
};

export const getContractDefaultFilters = (
  defaultFilters: Filter | null,
  isFeatureEnabled: boolean,
) => {
  if (!defaultFilters || !defaultFilters.AND) return [];
  const data = defaultFilters?.AND;

  if (isFeatureEnabled) {
    return data.map(({ filter }) => getFilterV2Fn(filter));
  }

  return data.map(({ filter }) => getFilterFn(filter));
};

export const getContractFilterFns = (
  filters: Filter | null,
  isFeatureEnabled: boolean,
) => {
  if (!filters || !filters.AND) return [];
  const data = filters?.AND;

  if (isFeatureEnabled) {
    return data.map(({ filter }) => getFilterV2Fn(filter));
  }

  return data.map(({ filter }) => getFilterFn(filter));
};
