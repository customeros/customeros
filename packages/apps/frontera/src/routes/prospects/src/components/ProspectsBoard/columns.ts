import { P, match } from 'ts-pattern';

import {
  Filter,
  ColumnView,
  FilterItem,
  Opportunity,
  TableViewDef,
} from '@graphql/types';

export const getFilterFn = (filter?: FilterItem | null) => {
  return match(filter)
    .with(
      { property: P.string },
      (filter) => (opportunity: Opportunity) =>
        opportunity[filter.property as keyof Opportunity] === filter.value,
    )
    .otherwise(() => () => true);
};

export const getColumns = (viewDef?: TableViewDef) => {
  if (!viewDef) return [];

  return viewDef?.columns
    .map((v) => {
      const parsedFilters = JSON.parse(v.filter) as Filter;
      const filterItems = parsedFilters.AND;

      const stage = getStageFromColumn(v);
      const filterFns = filterItems?.map(({ filter }) => getFilterFn(filter));

      return { ...v, stage, filterFns };
    })
    .filter((v) => v.visible);
};

export const getStageFromColumn = (column: ColumnView) => {
  const parsedFilters = JSON.parse(column.filter) as Filter;
  const filterItems = parsedFilters.AND;

  const internalStageFilter = filterItems?.find(
    (f) => f.filter?.property === 'internalStage',
  )?.filter;
  const externalStageFilter = filterItems?.find(
    (f) => f.filter?.property === 'externalStage',
  )?.filter;

  return externalStageFilter?.value ?? internalStageFilter?.value;
};
