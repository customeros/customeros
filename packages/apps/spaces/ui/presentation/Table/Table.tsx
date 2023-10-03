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
import { Fade } from '@ui/transitions/Fade';

const CELL_PADDING_X = 24;

declare module '@tanstack/table-core' {
  // REASON: TData & TValue are not used in this interface but need to be defined
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  interface ColumnDefBase<TData, TValue = unknown> {
    skeleton: () => React.ReactNode;
  }
}

interface TableProps<T extends object> {
  data: T[];
  // REASON: Typing TValue is too exhaustive and has no benefit
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
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
  totalItems = 40,
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

  const tableActionsCellWidth = tableActionsWidth + 2 * CELL_PADDING_X;

  const { rows } = table.getRowModel();
  const rowVirtualizer = useVirtualizer({
    count: !data.length && isLoading ? 40 : totalItems,
    overscan: 30,
    getScrollElement: () => scrollElementRef.current,
    estimateSize: () => 80,
  });

  const virtualRows = rowVirtualizer.getVirtualItems();

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
  }, [enableRowSelection, enableTableActions, _selection, data.length]);

  const skeletonRow = useMemo(
    () => createRow<T>(table, 'SKELETON', {} as T, totalItems + 1, 0),
    [table, totalItems],
  );

  return (
    <Flex w='100%' flexDir='column'>
      <TContent ref={scrollElementRef}>
        <THeader
          top='0'
          position='sticky'
          minW={
            table.getCenterTotalSize() +
            tableActionsCellWidth +
            (enableRowSelection ? 44 : 0)
          }
        >
          {table.getHeaderGroups().map((headerGroup) => (
            <THeaderGroup key={headerGroup.id}>
              {enableRowSelection && <THeaderCell w='44px' p='0' />}
              {headerGroup.headers.map((header, index) => (
                <THeaderCell
                  key={header.id}
                  flex={header.colSpan ?? '1'}
                  minWidth={`${header.getSize()}`}
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
        <TBody width='100%'>
          {!virtualRows.length && (
            <TRow justifyContent='center'>No results found</TRow>
          )}
          {virtualRows.map((virtualRow) => {
            const row = rows[virtualRow.index];
            return (
              <TRow
                key={virtualRow.key}
                data-index={virtualRow.index}
                minH={`${virtualRow.size}px`}
                minW={
                  table.getCenterTotalSize() +
                  tableActionsCellWidth +
                  (enableRowSelection ? 44 : 0)
                }
                top={`${virtualRow.start}px`}
                ref={rowVirtualizer.measureElement}
                bg={virtualRow.index % 2 === 0 ? 'gray.25' : 'white'}
              >
                {enableRowSelection && (
                  <TCell maxW='fit-content' pl='6' pr='0'>
                    <Flex align='center' flexDir='row' h='full'>
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
                    {row ? (
                      <Fade in={!!row}>
                        {flexRender(
                          cell.column.columnDef.cell,
                          cell.getContext(),
                        )}
                      </Fade>
                    ) : (
                      <Fade in={!row}>
                        {cell.column.columnDef?.skeleton?.()}
                      </Fade>
                    )}
                  </TCell>
                ))}
                {enableTableActions && (
                  <TCell
                    data-index={(row ?? skeletonRow).getAllCells().length + 1}
                    flex='0'
                    p='0'
                  >
                    <Flex flex='0' w={`${tableActionsWidth}px`} />
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
      height='inherit'
      position='relative'
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
      position='absolute'
      borderBottom='1px solid'
      borderBottomColor='gray.100'
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
      justify='center'
      ref={ref}
      {...props}
    />
  );
});

const TContent = forwardRef<HTMLDivElement, FlexProps>((props, ref) => {
  return (
    <Fade in>
      <Flex
        ref={ref}
        bg='gray.25'
        overflow='auto'
        flexDir='column'
        borderRadius='2xl'
        borderStyle='hidden'
        border='1px solid'
        borderColor='gray.200'
        height='calc(100vh - 74px)'
        sx={{
          '&::-webkit-scrollbar': {
            width: '4px',
            height: '4px',
            background: 'transparent',
          },
          '&::-webkit-scrollbar-track': {
            width: '4px',
            height: '4px',
            background: 'transparent',
          },
          '&::-webkit-scrollbar-thumb': {
            background: 'gray.500',
            borderRadius: '8px',
          },
        }}
        {...props}
      />
    </Fade>
  );
});

const THeader = forwardRef<HTMLDivElement, FlexProps>((props, ref) => {
  return (
    <Flex
      ref={ref}
      bg='white'
      width='inherit'
      borderBottom='1px solid'
      borderBottomColor='gray.100'
      zIndex='docked'
      {...props}
    />
  );
});

const THeaderGroup = forwardRef<HTMLDivElement, FlexProps>((props, ref) => {
  return <Flex ref={ref} flex='1' {...props} />;
});

const THeaderCell = forwardRef<HTMLDivElement, FlexProps>((props, ref) => {
  return (
    <Flex
      align='center'
      px='6'
      py='3'
      whiteSpace='nowrap'
      ref={ref}
      {...props}
    />
  );
});

export type { RowSelectionState, SortingState, TableInstance };
export { createColumnHelper } from '@tanstack/react-table';
