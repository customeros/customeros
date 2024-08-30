import { match } from 'ts-pattern';
import { FilterItem } from '@store/types.ts';
import { OpportunityStore } from '@store/Opportunities/Opportunity.store';

import { DateTimeUtils } from '@utils/date.ts';
import { Filter, ColumnViewType, ComparisonOperator } from '@graphql/types';

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

export const getOpportunityFilterFns = (filters: Filter | null) => {
  if (!filters || !filters.AND) return [];

  const data = filters?.AND;

  return data.map(({ filter }) => getFilterFn(filter));
};
