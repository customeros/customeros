import type {
  ColumnDef,
  OnChangeFn,
  SortingState,
  RowSelectionState,
  ColumnSizingState,
  ColumnFiltersState,
  Table as TableInstance,
} from '@tanstack/react-table';

import React, {
  useRef,
  useMemo,
  useState,
  useEffect,
  forwardRef,
  HTMLAttributes,
  MutableRefObject,
} from 'react';

import { twMerge } from 'tailwind-merge';
import { difference, intersection } from 'lodash';
import { Virtualizer, useVirtualizer } from '@tanstack/react-virtual';
import { useKey, useMergeRefs, useKeyBindings, useOutsideClick } from 'rooks';
import {
  createRow,
  flexRender,
  useReactTable,
  getCoreRowModel,
  getSortedRowModel,
  getFacetedRowModel,
  getFilteredRowModel,
} from '@tanstack/react-table';

import { cn } from '@ui/utils/cn';
import { TableIdType } from '@graphql/types';
import { useModKey } from '@shared/hooks/useModKey';
import { Tumbleweed } from '@ui/media/icons/Tumbleweed';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip.tsx';
import { Checkbox, CheckboxProps } from '@ui/form/Checkbox/Checkbox';

declare module '@tanstack/table-core' {
  // REASON: TData & TValue are not used in this interface but need to be defined
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  interface ColumnDefBase<TData, TValue = unknown> {
    fixWidth?: boolean;
    maxW?: number | string;
    skeleton: () => React.ReactNode;
  }
}

interface TableProps<T extends object> {
  data: T[];
  dataTest?: string;
  rowHeight?: number;
  isLoading?: boolean;
  totalItems?: number;
  borderColor?: string;
  tableId?: TableIdType;
  sorting?: SortingState;
  canFetchMore?: boolean;
  onFetchMore?: () => void;
  manualFiltering?: boolean;
  fullRowSelection?: boolean;
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  columns: ColumnDef<T, any>[];
  enableRowSelection?: boolean;
  enableTableActions?: boolean;
  selection?: RowSelectionState;
  enableColumnResizing?: boolean;
  contentHeight?: number | string;
  enableKeyboardShortcuts?: boolean;
  onFullRowSelection?: (row?: T) => void;
  onSortingChange?: OnChangeFn<SortingState>;
  getRowId?: (row: T, index: number) => string;
  onResizeColumn?: OnChangeFn<ColumnSizingState>;
  onSelectionChange?: (selectedIds: string[]) => void;
  tableRef: MutableRefObject<TableInstance<T> | null>;
  onSelectedIndexChange?: (index: number | null) => void;
  onFocusedRowChange?: (index: number | null, selectedIds: string[]) => void;
  // REASON: Typing TValue is too exhaustive and has no benefit
  renderTableActions?: (
    table: TableInstance<T>,
    focusRow: number | null,
    selectedIds: string[],
  ) => React.ReactNode;
}

export const Table = <T extends object>({
  data,
  dataTest,
  columns,
  tableRef,
  getRowId,
  isLoading,
  onFetchMore,
  canFetchMore,
  totalItems = 40,
  onSortingChange,
  sorting: _sorting,
  renderTableActions,
  enableRowSelection,
  enableTableActions,
  fullRowSelection,
  rowHeight = 33,
  contentHeight,
  borderColor,
  onSelectionChange,
  manualFiltering,
  onFocusedRowChange,
  onFullRowSelection,
  enableKeyboardShortcuts,
  enableColumnResizing = false,
  onResizeColumn,
  tableId,
}: TableProps<T>) => {
  const scrollElementRef = useRef<HTMLDivElement>(null);
  const [sorting, setSorting] = useState<SortingState>([]);
  const [selection, setSelection] = useState<RowSelectionState>({});
  const [focusedRowIndex, setFocusedRowIndex] = useState<number | null>(null);
  const [selectedIndex, setSelectedIndex] = useState<number | null>(null);
  const [isShiftPressed, setIsShiftPressed] = useState(false);

  const selectedIds = Object.keys(selection);

  useKey(
    'Shift',
    (e) => {
      setIsShiftPressed(e.type === 'keydown');
    },
    { eventTypes: ['keydown', 'keyup'] },
  );

  const handleSelectionChange: OnChangeFn<RowSelectionState> = (
    nextSelection,
  ) => {
    if (!isShiftPressed) {
      setSelection(nextSelection);

      return;
    }

    if (isShiftPressed && selectedIndex !== null && focusedRowIndex !== null) {
      setSelection((prev) => {
        const edgeIndexes = [
          Math.min(selectedIndex, focusedRowIndex),
          Math.max(selectedIndex, focusedRowIndex),
        ];

        const ids = data
          .slice(edgeIndexes[0], edgeIndexes[1] + 1)
          .map((d) => (d as unknown as { id: string }).id);

        const newSelection: Record<string, boolean> = {
          ...prev,
        };

        const prevIds = Object.keys(prev);
        const diff = difference(ids, prevIds);
        const match = intersection(ids, prevIds);
        const shouldRemove = diff.length < match.length;

        const endId = (data[edgeIndexes[1]] as unknown as { id: string }).id;

        diff.forEach((id) => {
          newSelection[id] = true;
        });
        shouldRemove &&
          [endId, ...match].forEach((id) => {
            delete newSelection[id];
          });

        return newSelection;
      });
    }
  };

  const table = useReactTable<T>({
    data,
    columns,
    state: {
      sorting: _sorting ?? sorting,
      rowSelection: selection,
    },
    enableColumnResizing,
    columnResizeMode: 'onChange',
    getRowId,
    manualFiltering,
    manualSorting: true,
    enableRowSelection: enableRowSelection || fullRowSelection,
    enableMultiRowSelection: enableRowSelection && !fullRowSelection,
    enableColumnFilters: true,
    enableSortingRemoval: false,
    getCoreRowModel: getCoreRowModel<T>(),
    getFacetedRowModel: getFacetedRowModel<T>(),
    getSortedRowModel: getSortedRowModel<T>(),
    getFilteredRowModel: getFilteredRowModel<T>(),
    onSortingChange: onSortingChange ?? setSorting,
    onRowSelectionChange: handleSelectionChange,
    columnResizeDirection: 'ltr',
    onColumnSizingChange: onResizeColumn,
  });

  const { rows } = table.getRowModel();
  const rowVirtualizer = useVirtualizer({
    count: !data.length && isLoading ? 40 : totalItems,
    overscan: 30,
    getScrollElement: () => scrollElementRef.current,
    estimateSize: () => rowHeight,
  });

  const columnSizeVars = React.useMemo(() => {
    const headers = table.getFlatHeaders();
    const colSizes: { [key: string]: number } = {};

    for (let i = 0; i < headers.length; i++) {
      const header = headers[i]!;

      colSizes[`--header-${header.id}-size`] = header.getSize();
      colSizes[`--col-${header.column.id}-size`] = header.column.getSize();
    }

    return colSizes;
  }, [table.getState().columnSizingInfo, table.getState().columnSizing, data]);

  const virtualRows = rowVirtualizer.getVirtualItems();

  useEffect(() => {
    const [lastItem] = [...virtualRows].reverse();

    if (!lastItem) {
      return;
    }

    if (lastItem.index >= data.length - 1 && canFetchMore && !isLoading) {
      onFetchMore?.();
    }
  }, [
    onFetchMore,
    data.length,
    isLoading,
    totalItems,
    virtualRows,
    canFetchMore,
  ]);

  useEffect(() => {
    if (tableRef) {
      tableRef.current = table;
    }
  }, [table]);

  useEffect(() => {
    onFocusedRowChange?.(focusedRowIndex, selectedIds);
  }, [focusedRowIndex, onFocusedRowChange, selectedIds.length]);

  useEffect(() => {
    onSelectionChange?.(selectedIds);
  }, [selectedIds, onSelectionChange]);

  useKeyBindings(
    {
      ArrowDown: () => {
        setFocusedRowIndex((prev) => {
          if (prev === null) return 0;
          if (prev === data.length - 1) return prev;

          return prev + 1;
        });
        scrollElementRef.current?.focus();
      },
      ArrowUp: () => {
        setFocusedRowIndex((prev) => {
          if (prev === null) return 0;
          if (prev === 0) return prev;

          return prev - 1;
        });
        scrollElementRef.current?.focus();

        if (!focusedRowIndex) return;
      },
      Enter: () => {
        if (focusedRowIndex === null) return;

        const row = rows[focusedRowIndex];

        setSelectedIndex(focusedRowIndex);
        row?.getToggleSelectedHandler()(true);
      },
      '/': () => {
        setFocusedRowIndex(null);
        scrollElementRef.current?.blur();
      },
    },
    {
      when: enableKeyboardShortcuts,
    },
  );

  useModKey(
    'a',
    (e) => {
      const tag = (e.target as HTMLElement).tagName;

      if (tag === 'INPUT' || tag === 'TEXTAREA') return;
      e.preventDefault();
      table.toggleAllRowsSelected();
    },
    { when: enableKeyboardShortcuts },
  );

  useEffect(() => {
    setFocusedRowIndex((prev) => (prev === null ? prev : 0));
  }, [totalItems]);

  useOutsideClick(scrollElementRef, () => {
    setFocusedRowIndex(null);
  });

  useEffect(() => {
    // If the table is not being navigated by ArrowUp or ArrowDown
    if (!scrollElementRef?.current?.hasAttribute('data-hide-cursor')) return;

    const endIndex = rowVirtualizer.range?.endIndex ?? 0;
    const startIndex = rowVirtualizer.range?.startIndex ?? 0;

    if (focusedRowIndex === null) return;

    if (endIndex - 2 < focusedRowIndex) {
      rowVirtualizer.scrollToIndex(focusedRowIndex, { align: 'end' });
    }

    if (startIndex > focusedRowIndex) {
      rowVirtualizer.scrollToIndex(focusedRowIndex);
    }
  }, [rowVirtualizer.range, focusedRowIndex]);

  const THeaderMinW =
    table.getCenterTotalSize() + (enableRowSelection ? 32 : 22);

  const includesAvatar = table
    .getHeaderGroups()?.[0]
    ?.headers?.[0]?.id?.includes('AVATAR');

  return (
    <div className={cn('flex flex-col relative w-full min-w-[300px]')}>
      <TContent
        ref={scrollElementRef}
        height={contentHeight}
        borderColor={borderColor}
        style={{
          ...columnSizeVars,
        }}
        onScrollToTop={() => {
          rowVirtualizer.scrollToIndex(0);
          setFocusedRowIndex(0);
        }}
        onScrollToBottom={() => {
          rowVirtualizer.scrollToIndex(totalItems - 1, { align: 'end' });
          setFocusedRowIndex(totalItems - 1);
        }}
      >
        <THeader
          style={{ minWidth: THeaderMinW }}
          className='top-0 sticky group/header'
        >
          {table.getHeaderGroups().map((headerGroup) => {
            return (
              <THeaderGroup key={headerGroup.id}>
                <THeaderCell
                  className={cn('p-0 min-h-8 w-4', enableRowSelection && 'w-8')}
                >
                  {!fullRowSelection && (
                    <div className={cn('items-center ml-2')}>
                      {enableRowSelection && (
                        <Tooltip
                          label={
                            table.getIsAllRowsSelected()
                              ? 'Deselect All'
                              : 'Select All'
                          }
                        >
                          <div>
                            <MemoizedCheckbox
                              key={`checkbox-header-select-all`}
                              dataTest={'all-orgs-select-all-orgs'}
                              isChecked={table.getIsAllRowsSelected()}
                              onChange={() => table.toggleAllRowsSelected()}
                              className='group-hover/header:visible group-hover/header:opacity-100'
                            />
                          </div>
                        </Tooltip>
                      )}
                    </div>
                  )}
                </THeaderCell>
                {headerGroup.headers
                  .filter((cell) => !cell.column.columnDef.enableHiding)
                  .map((header, index) => {
                    return (
                      <THeaderCell
                        key={header.id}
                        style={{
                          width: `calc(var(--header-${header?.id}-size) * 1px)`,
                          minWidth: `calc(var(--header-${header?.id}-size) * 1px)`,
                        }}
                        className={cn(`relative group/header-item`, {
                          'cursor-col-resize': header.column.getIsResizing(),
                          'pl-6': index === 1 && includesAvatar,
                          'pl-4': index === 0 && !includesAvatar,
                        })}
                      >
                        {header.isPlaceholder
                          ? null
                          : flexRender(
                              header.column.columnDef.header,
                              header.getContext(),
                            )}
                        {header.column.getCanResize() &&
                          enableColumnResizing && (
                            <div
                              onMouseDown={header.getResizeHandler()}
                              onTouchStart={header.getResizeHandler()}
                              onDoubleClick={() => header.column.resetSize()}
                              className={cn(
                                `absolute top-0 h-full px-1 cursor-col-resize right-6 opacity-0 group-hover/header-item:visible group-hover/header-item:opacity-100`,
                                {
                                  'opacity-100': header.column.getIsResizing(),
                                },
                              )}
                            >
                              <div className='h-full w-[2px]  bg-gray-300' />
                            </div>
                          )}
                      </THeaderCell>
                    );
                  })}
              </THeaderGroup>
            );
          })}
        </THeader>

        <TableBody
          table={table}
          tableId={tableId}
          dataTest={dataTest}
          isLoading={isLoading}
          totalItems={totalItems}
          rowVirtualizer={rowVirtualizer}
          includesAvatar={includesAvatar}
          focusedRowIndex={focusedRowIndex}
          fullRowSelection={fullRowSelection}
          setSelectedIndex={setSelectedIndex}
          onFullRowSelection={onFullRowSelection}
          enableRowSelection={enableRowSelection}
          setFocusedRowIndex={setFocusedRowIndex}
          enableKeyboardShortcuts={enableKeyboardShortcuts}
        />
      </TContent>

      {enableTableActions && (
        <TActions>
          {renderTableActions?.(table, focusedRowIndex, selectedIds)}
        </TActions>
      )}
    </div>
  );
};

interface TableBodyProps<T extends object> {
  dataTest?: string;
  totalItems: number;
  isLoading?: boolean;
  tableId?: TableIdType;
  table: TableInstance<T>;
  includesAvatar?: boolean;
  fullRowSelection?: boolean;
  enableRowSelection?: boolean;
  focusedRowIndex: number | null;
  enableKeyboardShortcuts?: boolean;
  onFullRowSelection?: (row: T) => void;
  setSelectedIndex: (index: number | null) => void;
  setFocusedRowIndex: (index: number | null) => void;
  rowVirtualizer: Virtualizer<HTMLDivElement, Element>;
}

const TableBody = <T extends object>({
  table,
  dataTest,
  isLoading,
  totalItems,
  rowVirtualizer,
  focusedRowIndex,
  fullRowSelection,
  setSelectedIndex,
  onFullRowSelection,
  setFocusedRowIndex,
  enableRowSelection,
  includesAvatar,
  tableId,
}: TableBodyProps<T>) => {
  const { rows } = table.getRowModel();
  const virtualRows = rowVirtualizer.getVirtualItems();
  const skeletonRow = useMemo(
    () => createRow<T>(table, 'SKELETON', {} as T, totalItems + 1, 0),
    [table, totalItems],
  );

  return (
    <TBody className='w-full' data-test={dataTest}>
      {!virtualRows.length && !isLoading && <NoResults tableId={tableId} />}
      {virtualRows.map((virtualRow) => {
        const row = rows[virtualRow.index];

        const minW = table.getCenterTotalSize() + (enableRowSelection ? 32 : 0);
        const top = `${virtualRow.start}px`;
        const hoverStyle = 'group';

        const enabledRowOpacity = enableRowSelection
          ? 'opacity-100'
          : 'opacity-0';

        const enabledRowPointer = enableRowSelection
          ? 'pointer-events-auto'
          : 'pointer-events-none';

        const focusStyle = 'data-[focused=true]:bg-grayModern-100';

        return (
          <TRow
            tabIndex={1}
            key={row?.index}
            data-index={virtualRow.index}
            ref={rowVirtualizer.measureElement}
            data-selected={row?.getIsSelected()}
            data-focused={row?.index === focusedRowIndex}
            style={{
              minWidth: minW,
              top: top,
            }}
            onFocus={() => {
              setFocusedRowIndex(row?.index);
            }}
            onMouseOver={() => {
              setFocusedRowIndex(row?.index);
            }}
            className={twMerge(
              hoverStyle,
              focusStyle,
              'group',
              row?.getIsSelected() && 'bg-gray-50',
            )}
            onClick={
              fullRowSelection
                ? (s) => {
                    row?.getToggleSelectedHandler()(s);

                    onFullRowSelection?.(row.original);
                    setFocusedRowIndex(row?.index);
                  }
                : undefined
            }
          >
            <TCell className={cn('pl-2 pr-2 max-w-fit')}>
              {!fullRowSelection && (
                <div
                  className={cn(
                    enabledRowPointer,
                    enabledRowOpacity,
                    'items-center ',
                  )}
                >
                  {enableRowSelection && (
                    <MemoizedCheckbox
                      isChecked={row?.getIsSelected()}
                      key={`checkbox-${virtualRow.index}`}
                      disabled={!row || !row?.getCanSelect()}
                      isFocused={row?.index === focusedRowIndex}
                      className='group-hover:visible group-hover:opacity-100'
                      onChange={(isChecked) => {
                        row?.getToggleSelectedHandler()(isChecked);
                        setSelectedIndex(virtualRow.index);
                      }}
                    />
                  )}
                </div>
              )}
            </TCell>
            {(isLoading ? row ?? skeletonRow : row)
              ?.getAllCells()
              .filter((cell) => !cell.column.columnDef.enableHiding)
              ?.map((cell, index) => {
                return (
                  <TCell
                    key={cell.id}
                    data-index={cell.row.index}
                    style={{
                      width: `calc((var(--col-${
                        cell.column.id
                      }-size) * 1px) - ${index > 0 ? 26 : 0}px)`,
                    }}
                    className={cn({
                      'pl-6': index === 1 && includesAvatar,
                      'pl-4': index === 0 && !includesAvatar,
                      'ml-[26px]': index > 1,
                    })}
                  >
                    {row
                      ? flexRender(
                          cell.column.columnDef.cell,
                          cell.getContext(),
                        )
                      : cell.column.columnDef?.skeleton?.()}
                  </TCell>
                );
              })}
          </TRow>
        );
      })}
    </TBody>
  );
};

interface GenericProps {
  tabIndex?: number;
  dataTest?: string;
  className?: string;
  children?: React.ReactNode;
  style?: React.CSSProperties;
  onClick?: (event: React.MouseEvent<HTMLDivElement, MouseEvent>) => void;
}

const TBody = forwardRef<HTMLDivElement, GenericProps>(
  ({ className, children, style, dataTest, ...props }, ref) => {
    return (
      <div
        ref={ref}
        style={style}
        data-test={dataTest}
        className={twMerge('flex w-full flex-1 relative', className)}
        {...props}
      >
        {children}
      </div>
    );
  },
);

const TRow = forwardRef<HTMLDivElement, React.HTMLAttributes<HTMLDivElement>>(
  ({ className, style, tabIndex, onClick, children, ...props }, ref) => {
    return (
      <div
        ref={ref}
        style={style}
        onClick={onClick}
        className={cn(
          'top-0 left-0 inline-flex items-center flex-1 w-full text-sm absolute border-b bg-white border-gray-100 transition-all animate-fadeIn',
          className,
        )}
        {...props}
      >
        {children}
      </div>
    );
  },
);

const TCell = forwardRef<HTMLDivElement, GenericProps>(
  ({ children, className, style, ...props }, ref) => {
    return (
      <div
        {...props}
        ref={ref}
        style={style}
        className={twMerge(
          'inline-block py-1 h-auto whitespace-nowrap justify-center break-keep truncate',
          className,
        )}
      >
        {children}
      </div>
    );
  },
);

interface TContentProps {
  dataTest?: string;
  className?: string;
  borderColor?: string;
  isScrolling?: boolean;
  height?: string | number;
  children?: React.ReactNode;
  onScrollToTop?: () => void;
  style?: React.CSSProperties;
  onScrollToBottom?: () => void;
  enableKeyboardShortcuts?: boolean;
}

const TContent = forwardRef<HTMLDivElement, TContentProps>(
  (
    {
      height,
      borderColor,
      children,
      className,
      style,
      isScrolling,
      onScrollToBottom,
      onScrollToTop,
      enableKeyboardShortcuts,
      ...props
    },
    ref,
  ) => {
    const _ref = useRef<HTMLDivElement>(null);
    const timeoutRef = useRef<NodeJS.Timeout>();
    const mergedRef = useMergeRefs(ref, _ref);

    const borderColorDynamic = borderColor ? borderColor : 'gray.200';
    const heightDynamic = height ? height : 'calc(100vh - 48px)';
    const scrollBarStyle =
      '[&::-webkit-scrollbar-track]:size-2 [&::-webkit-scrollbar-track]:bg-transparent [&::-webkit-scrollbar-thumb]:bg-gray-500 [&::-webkit-scrollbar-thumb]:rounded-lg [&::-webkit-scrollbar]:size-2 [&::-webkit-scrollbar]:bg-transparent';

    const hideCursor = () => {
      if (_ref?.current) {
        _ref.current.setAttribute('data-hide-cursor', '');

        if (timeoutRef.current) {
          clearTimeout(timeoutRef.current);
        }

        timeoutRef.current = setTimeout(() => {
          _ref.current?.removeAttribute('data-hide-cursor');
        }, 1000);
      }
    };

    return (
      <div
        tabIndex={-1}
        ref={mergedRef}
        style={{
          height: heightDynamic,
          borderColor: borderColorDynamic,
          ...style,
        }}
        className={twMerge(
          'flex flex-col bg-white border-t overflow-auto focus:outline-none data-[hide-cursor]:cursor-none data-[hide-cursor]:pointer-events-none',
          scrollBarStyle,
          className,
        )}
        onKeyDown={(e) => {
          if (!enableKeyboardShortcuts) return;

          if (e.code === 'Space') {
            // prevent scrolling when pressing space
            e.preventDefault();
          }

          if (e.code === 'ArrowUp') {
            // prevent scrolling when pressing arrow up
            e.preventDefault();

            if (e.metaKey) {
              onScrollToTop?.();
            }
            hideCursor();
          }

          if (e.code === 'ArrowDown') {
            // prevent scrolling when pressing arrow down
            e.preventDefault();

            if (e.metaKey) {
              onScrollToBottom?.();
            }
            hideCursor();
          }
        }}
        {...props}
      >
        {children}
      </div>
    );
  },
);

const THeader = forwardRef<HTMLDivElement, GenericProps>(
  ({ className, children, style, ...props }, ref) => {
    return (
      <div
        ref={ref}
        {...props}
        style={style}
        className={twMerge(
          'bg-white border-b border-gray-100 z-[1]',
          className,
        )}
      >
        {children}
      </div>
    );
  },
);

const THeaderGroup = forwardRef<HTMLDivElement, GenericProps>(
  ({ className, children, ...props }, ref) => {
    return (
      <div ref={ref} className='flex flex-1' {...props}>
        {children}
      </div>
    );
  },
);

const THeaderCell = forwardRef<
  HTMLDivElement,
  React.HTMLAttributes<HTMLDivElement>
>(({ className, style, children, ...props }, ref) => {
  return (
    <div
      ref={ref}
      style={style}
      className={twMerge('flex items-center py-1 whitespace-nowrap', className)}
      {...props}
    >
      {children}
    </div>
  );
});

const TActions = forwardRef<HTMLDivElement, HTMLAttributes<HTMLDivElement>>(
  (props, ref) => {
    return (
      <div
        ref={ref}
        className='flex items-center justify-center left-[50%] bottom-[32px] absolute'
        {...props}
      />
    );
  },
);

const NoResults = ({ tableId }: { tableId?: TableIdType }) => {
  return (
    <div className='pt-12 mx-auto text-gray-700 text-center'>
      <Tumbleweed className='w-12 h-12 text-gray-500' />
      <p className='text-md font-semibold'>Empty here in No Resultsville</p>
      <p className='max-w-[380px]'>
        Try using different keywords, checking for typos, or adjusting your
        filters.
        <br />
        <br />
        {tableId &&
          [TableIdType.Customers, TableIdType.Targets].includes(tableId) && (
            <>
              Alternatively, you can add more organizations here by changing
              their relationship.
            </>
          )}
      </p>
    </div>
  );
};

const MemoizedCheckbox = ({
  className,
  disabled,
  isChecked,
  isFocused,
  onChange,
  dataTest,
}: CheckboxProps & { isFocused?: boolean }) => {
  return (
    <Checkbox
      size='sm'
      iconSize='sm'
      disabled={disabled}
      onChange={onChange}
      data-test={dataTest}
      isChecked={isChecked}
      className={cn(
        className,
        isChecked || isFocused ? 'opacity-100' : 'opacity-0',
        isChecked || isFocused ? 'visible' : 'hidden',
      )}
    />
  );
};

export type {
  RowSelectionState,
  SortingState,
  TableInstance,
  ColumnFiltersState,
};

export { createColumnHelper } from '@tanstack/react-table';
