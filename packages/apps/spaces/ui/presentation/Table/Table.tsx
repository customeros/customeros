'use client';
import type {
  ColumnDef,
  OnChangeFn,
  SortingState,
  RowSelectionState,
  ColumnFiltersState,
  Table as TableInstance,
} from '@tanstack/react-table';

import { memo, useRef, useMemo, useState, useEffect, forwardRef } from 'react';

import { useVirtualizer } from '@tanstack/react-virtual';
import {
  createRow,
  flexRender,
  useReactTable,
  getCoreRowModel,
  getSortedRowModel,
  getFacetedRowModel,
  getFilteredRowModel,
} from '@tanstack/react-table';

import { Center } from '@ui/layout/Center';
import { Checkbox } from '@ui/form/Checkbox';
import { Flex, FlexProps } from '@ui/layout/Flex';

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
  // REASON: Typing TValue is too exhaustive and has no benefit
  isLoading?: boolean;
  totalItems?: number;
  sorting?: SortingState;
  canFetchMore?: boolean;
  onFetchMore?: () => void;
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  columns: ColumnDef<T, any>[];
  enableRowSelection?: boolean;
  enableTableActions?: boolean;
  onSortingChange?: OnChangeFn<SortingState>;
  renderTableActions?: (table: TableInstance<T>) => React.ReactNode;
}

export const Table = <T extends object>({
  data,
  columns,
  isLoading,
  onFetchMore,
  canFetchMore,
  totalItems = 40,
  onSortingChange,
  sorting: _sorting,
  renderTableActions,
  enableRowSelection,
  enableTableActions,
}: TableProps<T>) => {
  const scrollElementRef = useRef<HTMLDivElement>(null);
  const [sorting, setSorting] = useState<SortingState>([]);

  const table = useReactTable<T>({
    data,
    columns,
    state: {
      sorting: _sorting ?? sorting,
    },
    manualSorting: true,
    enableRowSelection,
    enableColumnFilters: true,
    enableSortingRemoval: false,
    getCoreRowModel: getCoreRowModel<T>(),
    getSortedRowModel: getSortedRowModel<T>(),
    getFacetedRowModel: getFacetedRowModel<T>(),
    getFilteredRowModel: getFilteredRowModel<T>(),
    onSortingChange: onSortingChange ?? setSorting,
  });

  const { rows } = table.getRowModel();
  const rowVirtualizer = useVirtualizer({
    count: !data.length && isLoading ? 40 : totalItems,
    overscan: 30,
    getScrollElement: () => scrollElementRef.current,
    estimateSize: () => 64,
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

  const skeletonRow = useMemo(
    () => createRow<T>(table, 'SKELETON', {} as T, totalItems + 1, 0),
    [table, totalItems],
  );

  return (
    <Flex w='100%' flexDir='column' position='relative'>
      <TContent ref={scrollElementRef}>
        <THeader
          top='0'
          position='sticky'
          minW={table.getCenterTotalSize() + (enableRowSelection ? 28 : 0)}
        >
          {table.getHeaderGroups().map((headerGroup) => (
            <THeaderGroup key={headerGroup.id}>
              {enableRowSelection && <THeaderCell w='28px' p='0' />}
              {headerGroup.headers.map((header, index) => (
                <THeaderCell
                  key={header.id}
                  flex={header.colSpan ?? '1'}
                  minWidth={`${header.getSize()}px`}
                  maxWidth={
                    header.column.columnDef.fixWidth
                      ? `${header.getSize()}px`
                      : 'none'
                  }
                  pr={index === 0 ? '0' : undefined}
                  pl={(() => {
                    switch (index) {
                      case 0:
                        return '2';
                      case 1:
                        return '0';
                      default:
                        return '6';
                    }
                  })()}
                >
                  {header.isPlaceholder
                    ? null
                    : flexRender(
                        header.column.columnDef.header,
                        header.getContext(),
                      )}
                </THeaderCell>
              ))}
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
                  table.getCenterTotalSize() + (enableRowSelection ? 28 : 0)
                }
                top={`${virtualRow.start}px`}
                ref={rowVirtualizer.measureElement}
                bg={virtualRow.index % 2 === 0 ? 'gray.25' : 'white'}
                _hover={{
                  '& .row-select-checkbox': {
                    opacity: '1',
                    visibility: 'visible',
                  },
                }}
              >
                {enableRowSelection && (
                  <TCell maxW='fit-content' pl='2' pr='0'>
                    <Flex align='center' flexDir='row' h='full'>
                      <MemoizedCheckbox
                        key={`checkbox-${virtualRow.index}`}
                        isSelected={row?.getIsSelected()}
                        isDisabled={!row || !row?.getCanSelect()}
                        onChange={row?.getToggleSelectedHandler()}
                      />
                    </Flex>
                  </TCell>
                )}
                {(row ?? skeletonRow).getAllCells()?.map((cell, index) => {
                  return (
                    <TCell
                      key={cell.id}
                      data-index={cell.row.index}
                      pr={index === 0 ? '0' : undefined}
                      pl={(() => {
                        switch (index) {
                          case 0:
                            return '2';
                          case 1:
                            return '0';
                          default:
                            return '6';
                        }
                      })()}
                      minW={`${cell.column.getSize()}px`}
                      maxW={
                        cell.column.columnDef.fixWidth
                          ? `${cell.column.getSize()}px`
                          : 'none'
                      }
                      flex={
                        table
                          .getFlatHeaders()
                          .find((h) => h.id === cell.column.columnDef.id)
                          ?.colSpan ?? '1'
                      }
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
      </TContent>

      {enableTableActions && <TActions>{renderTableActions?.(table)}</TActions>}
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
      py='2'
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
      py='1'
      whiteSpace='nowrap'
      ref={ref}
      {...props}
    />
  );
});

const TActions = forwardRef<HTMLDivElement, FlexProps>((props, ref) => {
  return (
    <Center left='50%' position='absolute' bottom='32px' ref={ref} {...props} />
  );
});

interface MemoizedCheckboxProps {
  isSelected: boolean;
  isDisabled: boolean;
  onChange: (event: unknown) => void;
}

const MemoizedCheckbox = memo<MemoizedCheckboxProps>(
  ({ isDisabled, isSelected, onChange }) => {
    return (
      <Checkbox
        size='lg'
        className='row-select-checkbox'
        isChecked={isSelected}
        disabled={isDisabled}
        onChange={onChange}
        opacity={isSelected ? '1' : '0'}
        visibility={isSelected ? 'visible' : 'hidden'}
      />
    );
  },
);

export type {
  RowSelectionState,
  SortingState,
  TableInstance,
  ColumnFiltersState,
};
export { createColumnHelper } from '@tanstack/react-table';
