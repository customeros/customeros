import { match } from 'ts-pattern';
import { FilterItem } from '@store/types.ts';
import { isAfter, isBefore } from 'date-fns';
import { OpportunityStore } from '@store/Opportunities/Opportunity.store';

import { DateTimeUtils } from '@utils/date.ts';
import { Filter, ColumnViewType, ComparisonOperator } from '@graphql/types';

const getFilterV2Fn = (filter: FilterItem | undefined | null) => {
  const noop = (_row: OpportunityStore) => true;

  if (!filter) return noop;

  return match(filter)
    .with(
      { property: ColumnViewType.OpportunitiesName },
      (filter) => (row: OpportunityStore) => {
        if (!filter.active) return true;
        const value = row.value?.name.trim();

        return filterTypeText(filter, value);
      },
    )
    .with(
      { property: ColumnViewType.OpportunitiesOrganization },
      (filter) => (row: OpportunityStore) => {
        if (!filter.active) return true;
        const value = row.organization?.value?.name.toLowerCase().trim();

        return filterTypeText(filter, value);
      },
    )
    .with(
      { property: ColumnViewType.OpportunitiesStage },
      (filter) => (row: OpportunityStore) => {
        if (!filter.active) return true;
        const value = [row.value.externalStage, row.value.internalStage];

        if (!value)
          return (
            filter.operation === ComparisonOperator.IsEmpty ||
            filter.operation === ComparisonOperator.NotContains
          );

        return filterTypeList(filter, value);
      },
    )
    .with(
      { property: ColumnViewType.OpportunitiesNextStep },
      (filter) => (row: OpportunityStore) => {
        if (!filter.active) return true;

        const value = row.value.nextSteps.toLowerCase();

        return filterTypeText(filter, value);
      },
    )
    .with(
      { property: ColumnViewType.OpportunitiesEstimatedArr },
      (filter) => (row: OpportunityStore) => {
        if (!filter.active) return true;
        const value = row.value?.maxAmount;

        return filterTypeNumber(filter, value);
      },
    )
    .with(
      { property: ColumnViewType.OpportunitiesTimeInStage },
      (filter) => (row: OpportunityStore) => {
        if (!filter.active) return true;

        const numberOfDays = DateTimeUtils.differenceInDays(
          new Date().toISOString(),
          row.value?.stageLastUpdated,
        );

        return filterTypeNumber(filter, Number(numberOfDays) + 1);
      },
    )
    .with(
      { property: ColumnViewType.OpportunitiesOwner },
      (filter) => (row: OpportunityStore) => {
        if (!filter.active) return true;
        const value = row.value.owner?.id;

        if (!value)
          return (
            filter.operation === ComparisonOperator.IsEmpty ||
            filter.operation === ComparisonOperator.NotContains
          );

        return filterTypeList(filter, [value || '']);
      },
    )

    .with(
      { property: ColumnViewType.OpportunitiesCreatedDate },
      (filter) => (row: OpportunityStore) => {
        if (!filter.active) return true;
        const value = row.value?.metadata.created?.split('T')[0];

        if (!value) return false;

        return filterTypeDate(filter, value);
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

const getFilterFn = (filter: FilterItem | undefined | null) => {
  const noop = (_row: OpportunityStore) => true;

  if (!filter) return noop;

  return match(filter)
    .with(
      { property: ColumnViewType.OpportunitiesName },
      (filter) => (row: OpportunityStore) => {
        if (!filter.active) return true;
        const filterValue = filter?.value;

        if (!filterValue && filter.active && !filter.includeEmpty) return true;
        if (!row.value?.name?.length && filter.includeEmpty) return true;
        if (!filterValue || !row.value?.name?.length) return false;

        return row.value.name.toLowerCase().includes(filterValue.toLowerCase());
      },
    )
    .with(
      { property: ColumnViewType.OpportunitiesNextStep },
      (filter) => (row: OpportunityStore) => {
        if (!filter.active) return true;
        const filterValue = filter?.value;

        if (!filterValue && filter.active && !filter.includeEmpty) return true;
        if (!row.value?.nextSteps?.length && filter.includeEmpty) return true;
        if (!filterValue || !row.value?.nextSteps?.length) return false;

        return row.value.nextSteps
          .toLowerCase()
          .includes(filterValue.toLowerCase());
      },
    )
    .with(
      { property: ColumnViewType.OpportunitiesOrganization },
      (filter) => (row: OpportunityStore) => {
        const filterValues = filter?.value;

        if (!filter.active) return true;
        const orgName = row.organization?.value?.name.toLowerCase().trim();

        if (filter.includeEmpty && !orgName?.length) {
          return true;
        }

        if (filter.includeEmpty && filterValues.length === 0) {
          return false;
        }

        return orgName?.includes(filterValues);
      },
    )
    .with(
      { property: ColumnViewType.OpportunitiesStage },
      (filter) => (row: OpportunityStore) => {
        if (!filter.active) return true;
        const filterValue = filter?.value;

        return (
          filterValue.includes(row.value.externalStage) ||
          filterValue.includes(row.value.internalStage)
        );
      },
    )
    .with(
      { property: ColumnViewType.OpportunitiesEstimatedArr },
      (filter) => (row: OpportunityStore) => {
        if (!filter.active) return true;
        const filterValue = filter?.value;
        const forecastValue = row.value?.maxAmount;

        if (!forecastValue) return false;

        return (
          forecastValue >= filterValue[0] && forecastValue <= filterValue[1]
        );
      },
    )
    .with(
      { property: ColumnViewType.OpportunitiesTimeInStage },
      (filter) => (row: OpportunityStore) => {
        if (!filter.active) return true;
        const filterValue = filter?.value;

        const operator = filter.operation;
        const numberOfDays = DateTimeUtils.differenceInDays(
          new Date().toISOString(),
          row.value?.stageLastUpdated,
        );

        if (operator === ComparisonOperator.Lt) {
          return Number(numberOfDays) < Number(filterValue);
        }

        if (operator === ComparisonOperator.Gt) {
          return Number(numberOfDays) > Number(filterValue);
        }

        if (operator === ComparisonOperator.Between) {
          const filterValue = filter?.value?.map(Number) as number[];

          return (
            numberOfDays >= Number(filterValue[0]) &&
            numberOfDays <= Number(filterValue[1])
          );
        }
      },
    )
    .with(
      { property: ColumnViewType.OpportunitiesOwner },
      (filter) => (row: OpportunityStore) => {
        if (!filter.active) return true;
        const filterValue = filter?.value;

        if (filterValue === '__EMPTY__' && !row.value.owner) {
          return true;
        }

        return filterValue.includes(row.value.owner?.id);
      },
    )

    .with(
      { property: ColumnViewType.OpportunitiesCreatedDate },
      (filter) => (row: OpportunityStore) => {
        if (!filter.active) return true;
        const filterValue = filter?.value;
        const createdAt = row.value?.metadata.created?.split('T')[0];

        if (!filterValue) return true;
        if (filterValue?.[1] === null) return filterValue?.[0] <= createdAt;
        if (filterValue?.[0] === null) return filterValue?.[1] >= createdAt;

        return filterValue[0] <= createdAt && filterValue[1] >= createdAt;
      },
    )

    .otherwise(() => noop);
};

export const getOpportunityFilterFns = (
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
