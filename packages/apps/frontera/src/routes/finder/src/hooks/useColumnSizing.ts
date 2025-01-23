import { useMemo, useCallback } from 'react';

import { ColumnDef, ColumnSizingState } from '@tanstack/react-table';
import { TableViewDef } from '@store/TableViewDefs/TableViewDef.dto';

export const useColumnSizing = (
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  tableColumns: ColumnDef<any, any>[],
  tableViewDef: TableViewDef | null,
) => {
  const columnCache = useMemo(
    () => new Map<string, { minSize?: number; maxSize?: number }>(),
    [],
  );

  const handleColumnSizing = useCallback(
    (
      updater:
        | ColumnSizingState
        | ((prev: ColumnSizingState) => ColumnSizingState),
    ) => {
      const update = typeof updater === 'function' ? updater({}) : updater;
      const entries = Object.entries(update);

      if (!entries.length) return;
      const [columnId, width] = entries[0];

      let columnSettings = columnCache.get(columnId);

      if (!columnSettings) {
        columnSettings = tableColumns?.find((e) => e.id === columnId) || {};
        columnCache.set(columnId, columnSettings);
      }

      const { minSize, maxSize } = columnSettings;
      const newWidth = Math.min(
        Math.max(width, minSize || 0),
        maxSize || Infinity,
      );

      tableViewDef?.setColumnSize(columnId, newWidth);
    },
    [tableColumns, columnCache, tableViewDef],
  );

  return handleColumnSizing;
};
