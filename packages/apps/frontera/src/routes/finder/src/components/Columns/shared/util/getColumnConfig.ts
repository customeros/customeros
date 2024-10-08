import { ColumnDef } from '@tanstack/react-table';

import { TableViewDef } from '@graphql/types';

export function getColumnConfig<Datum>(
  columns: Record<string, ColumnDef<Datum>>,
  tableViewDef?: Array<TableViewDef>[0],
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
): ColumnDef<Datum, any>[] {
  if (!tableViewDef) return [];

  return (tableViewDef.columns ?? []).reduce((acc, curr) => {
    const columnTypeName = curr?.columnType;

    if (!columnTypeName) return acc;

    if (columns[columnTypeName] === undefined) return acc;
    const column = {
      ...columns[columnTypeName],
      enableHiding: !curr.visible,
      size: curr.visible ? curr.width : 0,
      minSize: curr.visible ? columns[columnTypeName].minSize : 0,
    };

    if (!column) return acc;

    return [...acc, column];
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
  }, [] as ColumnDef<Datum, any>[]);
}
