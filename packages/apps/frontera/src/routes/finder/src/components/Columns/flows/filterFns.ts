import { match } from 'ts-pattern';
import { FilterItem } from '@store/types';
import { FlowStore } from '@store/Flows/Flow.store.ts';

import { Filter, ColumnViewType, ComparisonOperator } from '@graphql/types';

export const getFlowFiltersFns = (
  serverFilter: FilterItem | null | undefined,
) => {
  const noop = (_row: FlowStore) => true;

  if (!serverFilter) return noop;

  return match(serverFilter)
    .with(
      { property: ColumnViewType.FlowActionName },
      (filter) => (row: FlowStore) => {
        if (!filter.active) return true;
        const value = row.value?.name;

        return filterTypeText(filter, value);
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

export const getFlowsFilterFns = (filters: Filter | null) => {
  if (!filters || !filters.AND) return [];

  const data = filters?.AND;

  return data.map(({ filter }) => getFlowFiltersFns(filter));
};
