import type {
  ColumnDef,
  OnChangeFn,
  SortingState,
  RowSelectionState,
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
import { useMergeRefs, useKeyBindings, useOutsideClick } from 'rooks';
import { Virtualizer, useVirtualizer } from '@tanstack/react-virtual';
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
  rowHeight?: number;
  isLoading?: boolean;
  totalItems?: number;
  borderColor?: string;
  sorting?: SortingState;
  canFetchMore?: boolean;
  onFetchMore?: () => void;
  isSidePanelOpen?: boolean;
  manualFiltering?: boolean;
  fullRowSelection?: boolean;
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  columns: ColumnDef<T, any>[];
  enableRowSelection?: boolean;
  enableTableActions?: boolean;
  selection?: RowSelectionState;
  contentHeight?: number | string;
  enableKeyboardShortcuts?: boolean;
  onFullRowSelection?: (id?: string) => void;
  onSortingChange?: OnChangeFn<SortingState>;
  getRowId?: (row: T, index: number) => string;
  onSelectionChange?: OnChangeFn<RowSelectionState>;
  tableRef: MutableRefObject<TableInstance<T> | null>;
  onFocusedRowChange?: (index: number | null) => void;
  onSelectedIndexChange?: (index: number | null) => void;
  // REASON: Typing TValue is too exhaustive and has no benefit
  renderTableActions?: (table: TableInstance<T>) => React.ReactNode;
}

export const Table = <T extends object>({
  data,
  columns,
  tableRef,
  getRowId,
  isLoading,
  onFetchMore,
  canFetchMore,
  totalItems = 40,
  onSortingChange,
  sorting: _sorting,
  selection: _selection,
  renderTableActions,
  enableRowSelection,
  enableTableActions,
  isSidePanelOpen,
  fullRowSelection,
  rowHeight = 33,
  contentHeight,
  borderColor,
  manualFiltering,
  onSelectionChange,
  onFocusedRowChange,
  onFullRowSelection,
  onSelectedIndexChange,
  enableKeyboardShortcuts,
}: TableProps<T>) => {
  const scrollElementRef = useRef<HTMLDivElement>(null);
  const [sorting, setSorting] = useState<SortingState>([]);
  const [selection, setSelection] = useState<RowSelectionState>({});
  const [focusedRowIndex, setFocusedRowIndex] = useState<number | null>(null);
  const [selectedIndex, setSelectedIndex] = useState<number | null>(null);

  const table = useReactTable<T>({
    data,
    columns,
    state: {
      sorting: _sorting ?? sorting,
      rowSelection: _selection ?? selection,
    },
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
    onRowSelectionChange: onSelectionChange ?? setSelection,
  });

  const { rows } = table.getRowModel();
  const rowVirtualizer = useVirtualizer({
    count: !data.length && isLoading ? 40 : totalItems,
    overscan: 30,
    getScrollElement: () => scrollElementRef.current,
    estimateSize: () => rowHeight,
  });

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
    onFocusedRowChange?.(focusedRowIndex);
  }, [focusedRowIndex, onFocusedRowChange]);

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
      Space: () => {
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

  useEffect(() => {
    if (selectedIndex === -1) return;
    onSelectedIndexChange?.(selectedIndex);
  }, [selectedIndex]);

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
    table.getCenterTotalSize() + (enableRowSelection ? 32 : 0);

  return (
    <div
      className={cn('flex flex-col relative w-full')}
      style={{
        minWidth: isSidePanelOpen ? '300px' : '100%',
      }}
    >
      <TContent
        ref={scrollElementRef}
        height={contentHeight}
        borderColor={borderColor}
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
          className='top-0 sticky group/header'
          style={{ minWidth: THeaderMinW }}
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
                              isChecked={table.getIsAllRowsSelected()}
                              onChange={() => table.toggleAllRowsSelected()}
                              key={`checkbox-header-select-all`}
                              className='group-hover/header:visible group-hover/header:opacity-100'
                            />
                          </div>
                        </Tooltip>
                      )}
                    </div>
                  )}
                </THeaderCell>
                {headerGroup.headers.map((header, index) => {
                  const isHidden = header.column.columnDef.enableHiding;

                  if (isHidden) return null;

                  return (
                    <THeaderCell
                      key={header.id}
                      className={cn('relative', index === 1 && 'pl-6')}
                      style={{
                        width: header.getSize(),
                      }}
                    >
                      {header.isPlaceholder
                        ? null
                        : flexRender(
                            header.column.columnDef.header,
                            header.getContext(),
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
          isLoading={isLoading}
          totalItems={totalItems}
          rowVirtualizer={rowVirtualizer}
          focusedRowIndex={focusedRowIndex}
          onFullRowSelection={onFullRowSelection}
          fullRowSelection={fullRowSelection}
          enableRowSelection={enableRowSelection}
          setFocusedRowIndex={setFocusedRowIndex}
          setSelectedIndex={setSelectedIndex}
          enableKeyboardShortcuts={enableKeyboardShortcuts}
        />
      </TContent>

      {enableTableActions && <TActions>{renderTableActions?.(table)}</TActions>}
    </div>
  );
};

interface TableBodyProps<T extends object> {
  totalItems: number;
  isLoading?: boolean;
  table: TableInstance<T>;
  fullRowSelection?: boolean;
  enableRowSelection?: boolean;
  focusedRowIndex: number | null;
  enableKeyboardShortcuts?: boolean;
  onFullRowSelection?: (id?: string) => void;
  setSelectedIndex: (index: number | null) => void;
  setFocusedRowIndex: (index: number | null) => void;
  rowVirtualizer: Virtualizer<HTMLDivElement, Element>;
}

const TableBody = <T extends object>({
  table,
  isLoading,
  totalItems,
  rowVirtualizer,
  focusedRowIndex,
  fullRowSelection,
  setSelectedIndex,
  onFullRowSelection,
  setFocusedRowIndex,
  enableRowSelection,
}: TableBodyProps<T>) => {
  const { rows } = table.getRowModel();
  const virtualRows = rowVirtualizer.getVirtualItems();

  const skeletonRow = useMemo(
    () => createRow<T>(table, 'SKELETON', {} as T, totalItems + 1, 0),
    [table, totalItems],
  );

  return (
    <TBody className='w-full'>
      {!virtualRows.length && !isLoading && <NoResults />}
      {virtualRows.map((virtualRow) => {
        const row = rows[virtualRow.index];

        const minW = table.getCenterTotalSize() + (enableRowSelection ? 32 : 0);
        const top = `${virtualRow.start}px`;
        const hoverStyle = fullRowSelection ? 'hover:cursor-pointer' : 'group';

        const enabledRowOpacity = enableRowSelection
          ? 'opacity-100'
          : 'opacity-0';

        const enabledRowPointer = enableRowSelection
          ? 'pointer-events-auto'
          : 'pointer-events-none';

        const fullRowSelectionStyleDynamic = cn(
          virtualRow.index === 0
            ? 'hover:before:top-[-1px]'
            : 'hover:before:top-[-2px]',
          `hover:after:contents-[""] hover:after:h-[2px] hover:after:w-full hover:after:bg-gray-200 hover:after:bottom-[-1px] hover:after:absolute
           hover:before:contents-[""] hover:before:w-full hover:before:bg-gray-200 hover:before:h-[2px] hover:before:absolute`,
        );

        const rowHoverStyle = fullRowSelection
          ? fullRowSelectionStyleDynamic
          : undefined;

        // this might need to be removed
        const selectedStyle =
          fullRowSelection &&
          cn(
            'data-[selected=true]:before:contents-[""] data-[selected=true]:before:h-[2px] data-[selected=true]:before:w-full data-[selected=true]:before:bg-gray-200 data-[selected=true]:before:absolute',
            'data-[selected=true]:after:contents-[""] data-[selected=true]:after:w-full data-[selected=true]:after:bottom-[-1px] data-[selected=true]:after:bg-gray-200 data-[selected=true]:after:h-[2px] data-[selected=true]:after:absolute',
            virtualRow.index === 0
              ? 'data-[selected=true]:before:top[-1px]'
              : 'data-[selected=true]:before:top-[-2px]',
          );

        const focusStyle = 'data-[focused=true]:bg-grayModern-100';

        return (
          <TRow
            className={twMerge(
              hoverStyle,
              rowHoverStyle,
              selectedStyle,
              focusStyle,
              'group',
              row?.getIsSelected() && 'bg-gray-50',
            )}
            style={{
              minWidth: minW,
              top: top,
            }}
            key={row?.index}
            data-index={virtualRow.index}
            data-selected={row?.getIsSelected()}
            data-focused={row?.index === focusedRowIndex}
            ref={rowVirtualizer.measureElement}
            tabIndex={1}
            onMouseOver={() => {
              setFocusedRowIndex(row?.index);
            }}
            onFocus={() => {
              setFocusedRowIndex(row?.index);
            }}
            onClick={
              fullRowSelection
                ? (s) => {
                    row?.getToggleSelectedHandler()(s);
                    /// @ts-expect-error improve this later
                    const rowId = (row.original as unknown)?.id;
                    onFullRowSelection?.(rowId);
                    setFocusedRowIndex(row?.index);
                  }
                : undefined
            }
          >
            <TCell className='pl-2 pr-2 max-w-fit'>
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
                      isFocused={row?.index === focusedRowIndex}
                      key={`checkbox-${virtualRow.index}`}
                      disabled={!row || !row?.getCanSelect()}
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
              ?.map((cell, index) => {
                const isHidden = cell.column.columnDef.enableHiding;

                if (isHidden) return null;

                return (
                  <TCell
                    key={cell.id}
                    className={cn(
                      index === 1 && 'pl-6',
                      index > 1 && 'ml-[24px]',
                    )}
                    style={{
                      width:
                        (cell.column.columnDef.size ?? cell.column.getSize()) -
                        (index > 0 ? 24 : 0),
                    }}
                    data-index={cell.row.index}
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
  className?: string;
  children?: React.ReactNode;
  style?: React.CSSProperties;
  onClick?: (event: React.MouseEvent<HTMLDivElement, MouseEvent>) => void;
}

const TBody = forwardRef<HTMLDivElement, GenericProps>(
  ({ className, children, style, ...props }, ref) => {
    return (
      <div
        ref={ref}
        className={twMerge('flex w-full flex-1 relative', className)}
        style={style}
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
        className={cn(
          'top-0 left-0 inline-flex items-center flex-1 w-full text-sm absolute border-b bg-white border-gray-100 transition-all animate-fadeIn',
          className,
        )}
        ref={ref}
        style={style}
        onClick={onClick}
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
        className={twMerge(
          'inline-block py-1 h-auto whitespace-nowrap justify-center break-keep truncate',
          className,
        )}
        style={style}
        ref={ref}
      >
        {children}
      </div>
    );
  },
);

interface TContentProps {
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
        ref={mergedRef}
        tabIndex={-1}
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
        className={twMerge(
          'flex flex-col bg-white border-t overflow-auto focus:outline-none data-[hide-cursor]:cursor-none data-[hide-cursor]:pointer-events-none',
          scrollBarStyle,
          className,
        )}
        style={{
          height: heightDynamic,
          borderColor: borderColorDynamic,
          ...style,
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
        className={twMerge(
          'bg-white border-b border-gray-100 z-[1]',
          className,
        )}
        style={style}
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
      className={twMerge('flex items-center py-1 whitespace-nowrap', className)}
      style={style}
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
        className='flex items-center justify-center left-[50%] bottom-[32px] absolute'
        ref={ref}
        {...props}
      />
    );
  },
);

const NoResults = () => {
  return (
    <div className='pt-12 mx-auto text-gray-700 text-center'>
      <Tumbleweed className='w-12 h-12 text-gray-500' />
      <p className='text-md font-semibold'>Empty here in No Resultsville</p>
      <p className='max-w-[380px]'>
        Try using different keywords, checking for typos, or adjusting your
        filters.
        <br />
        <br /> Alternatively, you can add more organizations here by changing
        their relationship.
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
}: CheckboxProps & { isFocused?: boolean }) => {
  return (
    <Checkbox
      className={cn(
        className,
        isChecked || isFocused ? 'opacity-100' : 'opacity-0',
        isChecked || isFocused ? 'visible' : 'hidden',
      )}
      size='sm'
      iconSize='sm'
      isChecked={isChecked}
      disabled={disabled}
      onChange={onChange}
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
