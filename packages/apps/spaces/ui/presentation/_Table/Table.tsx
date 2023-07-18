'use client';
import { useRef, useState, useEffect, useMemo, forwardRef } from 'react';
import {
  flexRender,
  useReactTable,
  getCoreRowModel,
  getSortedRowModel,
  createRow,
} from '@tanstack/react-table';
import type {
  Table as TableInstance,
  ColumnDef,
  SortingState,
  RowSelectionState,
  OnChangeFn,
} from '@tanstack/react-table';
import { useVirtualizer } from '@tanstack/react-virtual';

import { Flex, FlexProps } from '@ui/layout/Flex';
import { Checkbox } from '@ui/form/Checkbox';

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
  renderTableActions?: (table: TableInstance<T>) => React.ReactNode;
}

export const Table = <T extends object>({
  data,
  columns,
  isLoading,
  onFetchMore,
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
    getCoreRowModel: getCoreRowModel<T>(),
    getSortedRowModel: getSortedRowModel<T>(),
    onSortingChange: onSortingChange ?? setSorting,
    onRowSelectionChange: onSelectionChange ?? setSelection,
  });

  const { rows } = table.getRowModel();
  const rowVirtualizer = useVirtualizer({
    count: !data.length && isLoading ? 5 : totalItems,
    overscan: 25,
    getScrollElement: () => scrollElementRef.current,
    estimateSize: () => 21,
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

  const skeletonRow = useMemo(
    () => createRow<T>(table, 'SKELETON', {} as T, totalItems + 1, 0),
    [table, totalItems],
  );
  return (
    <Flex w='100%' flexDir='column'>
      <Flex fontSize='sm' alignSelf='flex-end'>
        Total items: {totalItems}
      </Flex>
      <TContent minW={table.getCenterTotalSize()}>
        <THeader>
          {table.getHeaderGroups().map((headerGroup) => (
            <THeaderGroup key={headerGroup.id}>
              {enableRowSelection && <THeaderCell w='44px' p='0' />}
              {headerGroup.headers.map((header, index) => (
                <THeaderCell
                  key={header.id}
                  flex={header.colSpan}
                  minWidth={header.getSize()}
                  pl={index === 0 ? '3' : '6'}
                >
                  {header.isPlaceholder
                    ? null
                    : flexRender(
                        header.column.columnDef.header,
                        header.getContext(),
                      )}
                </THeaderCell>
              ))}
              {enableTableActions && (
                <THeaderCell
                  align='center'
                  ref={tableActionsRef}
                  justifyContent='flex-end'
                >
                  {renderTableActions?.(table)}
                </THeaderCell>
              )}
            </THeaderGroup>
          ))}
        </THeader>
        <TBody ref={scrollElementRef}>
          {!virtualRows.length && <TRow justifyContent='center'>No data</TRow>}
          {virtualRows.map((virtualRow) => {
            const row = rows[virtualRow.index];
            return (
              <TRow
                key={virtualRow.key}
                data-index={virtualRow.index}
                minH={`${virtualRow.size}px`}
                top={`${virtualRow.start}px`}
                ref={rowVirtualizer.measureElement}
                bg={virtualRow.index % 2 === 0 ? 'gray.50' : 'white'}
              >
                {enableRowSelection && (
                  <TCell maxW='fit-content' pl='6' pr='0'>
                    <Flex align='center' flexDir='row'>
                      <Checkbox
                        size='lg'
                        checked={row?.getIsSelected()}
                        disabled={!row || !row?.getCanSelect()}
                        onChange={row?.getToggleSelectedHandler()}
                      />
                    </Flex>
                  </TCell>
                )}
                {(row ?? skeletonRow).getAllCells().map((cell, index) => (
                  <TCell
                    key={cell.id}
                    data-index={cell.row.index}
                    pl={index === 0 ? '3' : '6'}
                    minW={`${cell.column.getSize()}px`}
                    flex={
                      table
                        .getFlatHeaders()
                        .find((h) => h.id === cell.column.columnDef.id)
                        ?.colSpan ?? '1'
                    }
                  >
                    {row && !isLoading
                      ? flexRender(
                          cell.column.columnDef.cell,
                          cell.getContext(),
                        )
                      : cell.column.columnDef?.skeleton?.()}
                  </TCell>
                ))}
                {enableTableActions && (
                  <TCell
                    flex='0'
                    // w={`${tableActionsWidth}px`}
                    // maxW={`${tableActionsWidth}px`}
                    minW={`${tableActionsWidth - 4}px`}
                  >
                    <Flex flex='0' w={`${tableActionsWidth - 4}px`} />
                  </TCell>
                )}
              </TRow>
            );
          })}
        </TBody>
      </TContent>
    </Flex>
  );
};

const TBody = forwardRef<HTMLDivElement, FlexProps>((props, ref) => {
  return (
    <Flex
      ref={ref}
      flex='1'
      w='100%'
      overflowY='auto'
      height='inherit'
      overflowX='hidden'
      position='relative'
      sx={{
        '&::-webkit-scrollbar': {
          width: '6px',
          background: 'transparent',
        },
        '&::-webkit-scrollbar-track': {
          width: '6px',
          background: 'white',
        },
        '&::-webkit-scrollbar-thumb': {
          background: 'gray.200',
          borderRadius: '24px',
        },
      }}
      {...props}
    />
  );
});

const TRow = forwardRef<HTMLDivElement, FlexProps>((props, ref) => {
  return (
    <Flex
      top='0'
      left='0'
      ref={ref}
      flex='1'
      width='100%'
      fontSize='sm'
      overflow='visible'
      position='absolute'
      borderBottom='1px solid'
      borderBottomColor='gray.200'
      {...props}
    />
  );
});

const TCell = forwardRef<HTMLDivElement, FlexProps>((props, ref) => {
  return (
    <Flex
      px='6'
      py='4'
      flex='1'
      flexDir='column'
      whiteSpace='nowrap'
      wordBreak='keep-all'
      ref={ref}
      {...props}
    />
  );
});

const TContent = forwardRef<HTMLDivElement, FlexProps>((props, ref) => {
  return (
    <Flex
      ref={ref}
      bg='white'
      boxShadow='sm'
      overflowX='auto'
      flexDir='column'
      borderRadius='2xl'
      borderStyle='hidden'
      height='calc(100vh - 122px)'
      {...props}
    />
  );
});

const THeader = forwardRef<HTMLDivElement, FlexProps>((props, ref) => {
  return (
    <Flex
      ref={ref}
      bg='white'
      borderBottom='1px solid'
      borderBottomColor='gray.200'
      {...props}
    />
  );
});

const THeaderGroup = forwardRef<HTMLDivElement, FlexProps>((props, ref) => {
  return <Flex ref={ref} flex='1' {...props} />;
});

const THeaderCell = forwardRef<HTMLDivElement, FlexProps>((props, ref) => {
  return <Flex px='6' py='3' whiteSpace='nowrap' ref={ref} {...props} />;
});

export type { RowSelectionState, SortingState, TableInstance };
export { createColumnHelper } from '@tanstack/react-table';
