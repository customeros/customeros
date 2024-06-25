import { noop } from 'lodash';
import { Filter, FilterItem } from '@store/types.ts';

import { TableViewDef, TableViewType } from '@graphql/types';
import { contactColumns } from '@organizations/components/Columns/Dictionaries/contactColumns.tsx';
import { organizationColumns } from '@organizations/components/Columns/Dictionaries/organizationColumns.tsx';
import {
  getContactColumnSortFn,
  getOrganizationColumnSortFn,
} from '@organizations/components/Columns/Dictionaries/SortAndFilterDictionary';

const allColumns = { ...organizationColumns, ...contactColumns };

export const getColumnsConfig = (tableViewDef?: Array<TableViewDef>[0]) => {
  if (!tableViewDef) return [];

  // @ts-expect-error fixme
  return (tableViewDef.columns ?? []).reduce((acc, curr) => {
    const columnTypeName = curr?.columnType;
    if (!columnTypeName) return acc;

    if (allColumns[columnTypeName] === undefined) return acc;
    const column = {
      ...allColumns[columnTypeName],
      enableHiding: !curr.visible,
    };

    if (!column) return acc;

    return [...acc, column];
  }, []);
};

export const getColumnSortFn = (columnId: string, type: TableViewType) => {
  switch (type) {
    case TableViewType.Contacts:
      return getContactColumnSortFn(columnId);
    case TableViewType.Organizations:
      return getOrganizationColumnSortFn(columnId);
    default:
      return noop;
  }
};

export const getAllFilterFns = (
  filters: Filter | null,
  filterFunc: (f: FilterItem | undefined | null) => void = noop,
) => {
  if (!filters || !filters.AND) return [];

  const data = filters?.AND;

  return data.map(({ filter }) => filterFunc(filter));
};
