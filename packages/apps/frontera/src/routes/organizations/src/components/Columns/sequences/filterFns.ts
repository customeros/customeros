import { match } from 'ts-pattern';
import { FilterItem } from '@store/types';
import { FlowSequenceStore } from '@store/Sequences/FlowSequence.store';

import { Filter, ColumnViewType } from '@graphql/types';

export const getPredefinedFilterFn = (
  serverFilter: FilterItem | null | undefined,
) => {
  const noop = (_row: FlowSequenceStore) => true;

  if (!serverFilter) return noop;

  return match(serverFilter)
    .with(
      { property: ColumnViewType.FlowSequenceStatus },
      (filter) => (row: FlowSequenceStore) => {
        const filterValues = filter?.value;

        if (!filter.active) return true;

        return filterValues?.includes(row.value?.status);
      },
    )
    .otherwise(() => noop);
};

export const getSequencesFilterFns = (filters: Filter | null) => {
  if (!filters || !filters.AND) return [];

  const data = filters?.AND;

  return data.map(({ filter }) => getPredefinedFilterFn(filter));
};
