import { useMemo, useCallback } from 'react';

import { ColumnDef, ColumnSizingState } from '@tanstack/react-table';
import { TableViewDefStore } from '@store/TableViewDefs/TableViewDef.store.ts';

import { FinderTableEntityTypes } from '@organizations/components/FinderTable';

export const useColumnSizing = (
  tableColumns: ColumnDef<FinderTableEntityTypes>[],
  tableViewDef?: TableViewDefStore,
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
      const [columnId, width] = Object.entries(update)[0];

      let columnSettings = columnCache.get(columnId);

      if (!columnSettings) {
        columnSettings = tableColumns.find((e) => e.id === columnId) || {};
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
