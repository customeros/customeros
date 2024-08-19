import { useMemo, useCallback } from 'react';

import { ContactStore } from '@store/Contacts/Contact.store.ts';
import { InvoiceStore } from '@store/Invoices/Invoice.store.ts';
import { ContractStore } from '@store/Contracts/Contract.store.ts';
import { ColumnDef, ColumnSizingState } from '@tanstack/react-table';
import { OrganizationStore } from '@store/Organizations/Organization.store.ts';
import { TableViewDefStore } from '@store/TableViewDefs/TableViewDef.store.ts';

export const useColumnSizing = (
  tableColumns: ColumnDef<
    OrganizationStore | ContactStore | InvoiceStore | ContractStore
  >[],
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
