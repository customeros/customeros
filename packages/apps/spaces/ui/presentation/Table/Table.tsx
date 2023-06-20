import { useRef, useState, useEffect, RefObject } from 'react';
import {
  flexRender,
  useReactTable,
  getCoreRowModel,
  getSortedRowModel,
  createRow,
} from '@tanstack/react-table';
import type { MenuProps } from 'primereact';
import type {
  Table as TableInstance,
  ColumnDef,
  SortingState,
  RowSelectionState,
  OnChangeFn,
} from '@tanstack/react-table';
import { useVirtualizer } from '@tanstack/react-virtual';

import { TActions } from './TActions';
import styles from './Table.module.scss';
import classNames from 'classnames';

declare module '@tanstack/table-core' {
  interface ColumnDefBase<TData, TValue = unknown> {
    skeleton: () => React.ReactNode;
  }
}

interface TableProps<T extends object> {
  data: T[];
  columns: ColumnDef<T, any>[];
  isLoading?: boolean;
  totalItems?: number;
  sorting?: SortingState;
  selection?: RowSelectionState;
  onFetchMore?: () => void;
  enableRowSelection?: boolean;
  enableTableActions?: boolean;
  onSortingChange?: OnChangeFn<SortingState>;
  onSelectionChange?: OnChangeFn<RowSelectionState>;
  tableActions?: (table: TableInstance<T>) => MenuProps['model'];
  renderTableActions?: (
    ref: RefObject<HTMLDivElement>,
    table: TableInstance<T>,
  ) => React.ReactNode;
}

export const Table = <T extends object>({
  data,
  columns,
  isLoading,
  onFetchMore,
  tableActions,
  totalItems = 50,
  onSortingChange,
  onSelectionChange,
  renderTableActions,
  enableRowSelection,
  enableTableActions,
  sorting: _sorting,
  selection: _selection,
}: TableProps<T>) => {
  const tableActionsRef = useRef<HTMLDivElement>(null);
  const scrollElementRef = useRef<HTMLDivElement>(null);
  const [sorting, setSorting] = useState<SortingState>([]);
  const [selection, setSelection] = useState<RowSelectionState>({});
  const [tableActionsWidth, setTableActionsWidth] = useState<number>(0);

  const table = useReactTable<T>({
    data,
    columns,
    state: {
      sorting: _sorting ?? sorting,
      rowSelection: _selection ?? selection,
    },
    manualSorting: true,
    enableRowSelection,
    meta: {
      isLoading: false,
    },
    getCoreRowModel: getCoreRowModel<T>(),
    getSortedRowModel: getSortedRowModel<T>(),
    onSortingChange: onSortingChange ?? setSorting,
    onRowSelectionChange: onSelectionChange ?? setSelection,
  });

  const { rows } = table.getRowModel();
  const rowVirtualizer = useVirtualizer({
    count: totalItems,
    overscan: 5,
    getScrollElement: () => scrollElementRef.current,
    estimateSize: () => 69,
  });

  const { getVirtualItems } = rowVirtualizer;
  const virtualRows = getVirtualItems();

  useEffect(() => {
    const [lastItem] = [...virtualRows].reverse();

    if (!lastItem) {
      return;
    }

    if (
      lastItem.index >= data.length - 1 &&
      data.length < totalItems &&
      !isLoading
    ) {
      onFetchMore?.();
    }
  }, [onFetchMore, data.length, isLoading, totalItems, virtualRows]);

  useEffect(() => {
    setTableActionsWidth(tableActionsRef.current?.clientWidth ?? 0);
  }, [enableRowSelection, enableTableActions, _selection]);

  const skeletonRow = createRow<T>(
    table,
    'SKELETON',
    {} as T,
    totalItems + 1,
    0,
  );

  return (
    <div className={styles.container}>
      <span className={styles.totalItems}>Total items: {totalItems}</span>
      <div
        className={styles.table}
        style={{ minWidth: table.getCenterTotalSize() }}
      >
        <div className={classNames(styles.thead)}>
          {table.getHeaderGroups().map((headerGroup) => (
            <div key={headerGroup.id} className={styles.tr}>
              {enableRowSelection && (
                <div
                  className={classNames(styles.th, styles.selectCell)}
                  style={{
                    width: 24,
                    padding: 0,
                    borderRightColor: 'transparent',
                  }}
                />
              )}
              {headerGroup.headers.map((header, index) => (
                <div
                  key={header.id}
                  className={styles.th}
                  style={{
                    flex: header.colSpan,
                    minWidth: header.getSize(),
                    borderLeft:
                      index === 0 && enableRowSelection
                        ? '1px solid transparent'
                        : undefined,
                  }}
                >
                  {header.isPlaceholder
                    ? null
                    : flexRender(
                        header.column.columnDef.header,
                        header.getContext(),
                      )}
                </div>
              ))}
              {enableTableActions && (
                <div className={classNames(styles.th, styles.actionCell)}>
                  {renderTableActions ? (
                    renderTableActions(tableActionsRef, table)
                  ) : (
                    <TActions
                      ref={tableActionsRef}
                      model={tableActions?.(table)}
                    />
                  )}
                </div>
              )}
            </div>
          ))}
        </div>
        <div className={styles.tbody} ref={scrollElementRef}>
          {virtualRows.map((virtualRow) => {
            const row = rows[virtualRow.index];
            return (
              <div
                key={virtualRow.key}
                data-index={virtualRow.index}
                ref={rowVirtualizer.measureElement}
                className={classNames(styles.row, {
                  [styles.even]: virtualRow.index % 2 !== 0,
                })}
                style={{
                  minHeight: `${virtualRow.size}px`,
                  top: `${virtualRow.start}px`,
                }}
              >
                {enableRowSelection && (
                  <div
                    className={classNames(styles.rowCell, styles.selectCell)}
                  >
                    <div className={styles.selectCheckboxWrapper}>
                      <input
                        type='checkbox'
                        className={styles.selectCheckbox}
                        checked={row.getIsSelected()}
                        disabled={!row.getCanSelect()}
                        onChange={row.getToggleSelectedHandler()}
                      />
                    </div>
                  </div>
                )}
                {(row ?? skeletonRow).getVisibleCells().map((cell) => (
                  <div
                    key={cell.id}
                    data-index={cell.row.index}
                    className={classNames(styles.rowCell)}
                    style={{
                      minWidth: cell.column.getSize(),
                      flex: table
                        .getFlatHeaders()
                        .find((h) => h.id === cell.column.columnDef.id)
                        ?.colSpan,
                    }}
                  >
                    {row
                      ? flexRender(
                          cell.column.columnDef.cell,
                          cell.getContext(),
                        )
                      : cell.column.columnDef?.skeleton?.()}
                  </div>
                ))}
                {enableTableActions && (
                  <div
                    className={styles.rowCell}
                    style={{
                      width: tableActionsWidth + 21,
                      padding: 0,
                    }}
                  />
                )}
              </div>
            );
          })}
        </div>
      </div>
    </div>
  );
};

export type { RowSelectionState, SortingState, TableInstance };
export { createColumnHelper } from '@tanstack/react-table';
