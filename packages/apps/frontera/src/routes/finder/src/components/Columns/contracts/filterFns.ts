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
          row.value?.autoRenew === true ? 'Auto-renews' : 'Not auto-renewing';

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

        if (value === undefined || value === null) return false;

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
        const value = row?.value.opportunities?.find(
          (opp) => opp,
        )?.renewalLikelihood;

        if (value?.[0] === null || value?.[0] === undefined)
          return (
            filter.operation === ComparisonOperator.IsEmpty ||
            filter.operation === ComparisonOperator.NotContains
          );

        return filterTypeList(filter, [String(value)]);
      },
    )
    .with(
      { property: ColumnViewType.ContractsForecastArr },
      (filter) => (row: ContractStore) => {
        if (!filter.active) return true;
        const value = row?.openOpportunity?.amount;

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

        return filterTypeList(filter, [String(value)]);
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

  const val = value ?? 0;

  return match(filterOperator)
    .with(ComparisonOperator.Lt, () => val < Number(filterValue))
    .with(ComparisonOperator.Gt, () => val > Number(filterValue))
    .with(ComparisonOperator.Eq, () => val === Number(filterValue))
    .with(ComparisonOperator.NotEqual, () => val !== Number(filterValue))
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

export const getContractDefaultFilters = (defaultFilters: Filter | null) => {
  if (!defaultFilters || !defaultFilters.AND) return [];
  const data = defaultFilters?.AND;

  return data.map(({ filter }) => getFilterFn(filter));
};

export const getContractFilterFns = (filters: Filter | null) => {
  if (!filters || !filters.AND) return [];
  const data = filters?.AND;

  return data.map(({ filter }) => getFilterFn(filter));
};
