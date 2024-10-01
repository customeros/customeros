import { match } from 'ts-pattern';
import { FilterItem } from '@store/types';
import { FlowStore } from '@store/Flows/Flow.store.ts';

import { Filter, ColumnViewType } from '@graphql/types';

export const getPredefinedFilterFn = (
  serverFilter: FilterItem | null | undefined,
) => {
  const noop = (_row: FlowStore) => true;

  if (!serverFilter) return noop;

  return match(serverFilter)
    .with(
      { property: ColumnViewType.FlowActionName },
      (filter) => (row: FlowStore) => {
        const filterValues = filter?.value;

        if (!filter.active || !filterValues.length) return true;

        return filterValues?.includes(row.value?.status);
      },
    )
    .otherwise(() => noop);
};

export const getFlowsFilterFns = (filters: Filter | null) => {
  if (!filters || !filters.AND) return [];

  const data = filters?.AND;

  return data.map(({ filter }) => getPredefinedFilterFn(filter));
};
