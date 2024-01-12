'use client';

import { useState, useCallback } from 'react';
import { useSearchParams } from 'next/navigation';

import { useIsRestoring } from '@tanstack/react-query';

import { GridItem } from '@ui/layout/Grid';
import { Table, SortingState } from '@ui/presentation/Table';
import { TableViewDef, RenewalRecord } from '@graphql/types';

import { Search } from './src/components/Search';
import { useRenewalsPageData } from './src/hooks';
import { getColumnsConfig } from './src/components/Columns/Columns';
import { EmptyState } from './src/components/EmptyState/EmptyState';
import { useFilterSetter } from './src/components/Columns/Filters/filterSetters';

const tableViewDef1: TableViewDef = {
  id: '1',
  order: 0,
  name: 'Monthly renewals',
  type: {
    id: '1',
    name: 'RENEWALS',
    createdAt: '2021-08-10T14:00:00.000Z',
    updatedAt: '2021-08-10T14:00:00.000Z',
  },
  filters: '',
  columns: [
    {
      id: '1',
      createdAt: '2021-08-10T14:00:00.000Z',
      updatedAt: '2021-08-10T14:00:00.000Z',
      isVisible: true,
      isDefaultSort: false,

      columnType: {
        id: '1',
        name: 'AVATAR',
        createdAt: '2021-08-10T14:00:00.000Z',
        updatedAt: '2021-08-10T14:00:00.000Z',
      },
    },
    {
      id: '2',
      createdAt: '2021-08-10T14:00:00.000Z',
      updatedAt: '2021-08-10T14:00:00.000Z',
      columnType: {
        id: '2',
        name: 'NAME',
        createdAt: '2021-08-10T14:00:00.000Z',
        updatedAt: '2021-08-10T14:00:00.000Z',
      },
    },
    {
      id: '7',
      createdAt: '2021-08-10T14:00:00.000Z',
      updatedAt: '2021-08-10T14:00:00.000Z',
      columnType: {
        id: '7',
        name: 'RENEWAL_DATE',
        createdAt: '2021-08-10T14:00:00.000Z',
        updatedAt: '2021-08-10T14:00:00.000Z',
      },
    },
    {
      id: '8',
      createdAt: '2021-08-10T14:00:00.000Z',
      updatedAt: '2021-08-10T14:00:00.000Z',
      columnType: {
        id: '8',
        name: 'FORECAST_ARR',
        createdAt: '2021-08-10T14:00:00.000Z',
        updatedAt: '2021-08-10T14:00:00.000Z',
      },
    },
    {
      id: '6',
      createdAt: '2021-08-10T14:00:00.000Z',
      updatedAt: '2021-08-10T14:00:00.000Z',
      columnType: {
        id: '6',
        name: 'RENEWAL_LIKELIHOOD',
        createdAt: '2021-08-10T14:00:00.000Z',
        updatedAt: '2021-08-10T14:00:00.000Z',
      },
    },
    {
      id: '9',
      createdAt: '2021-08-10T14:00:00.000Z',
      updatedAt: '2021-08-10T14:00:00.000Z',
      columnType: {
        id: '9',
        name: 'OWNER',
        createdAt: '2021-08-10T14:00:00.000Z',
        updatedAt: '2021-08-10T14:00:00.000Z',
      },
    },
    {
      id: '10',
      createdAt: '2021-08-10T14:00:00.000Z',
      updatedAt: '2021-08-10T14:00:00.000Z',
      columnType: {
        id: '10',
        name: 'LAST_TOUCHPOINT',
        createdAt: '2021-08-10T14:00:00.000Z',
        updatedAt: '2021-08-10T14:00:00.000Z',
      },
    },
  ],
  createdAt: '2021-08-10T14:00:00.000Z',
  updatedAt: '2021-08-10T14:00:00.000Z',
};

const tableViewDef2: TableViewDef = {
  id: '2',
  order: 1,
  name: 'Quarterly renewals',
  type: {
    id: '1',
    name: 'RENEWALS',
    createdAt: '2021-08-10T14:00:00.000Z',
    updatedAt: '2021-08-10T14:00:00.000Z',
  },
  filters: '',
  columns: [
    {
      id: '1',
      createdAt: '2021-08-10T14:00:00.000Z',
      updatedAt: '2021-08-10T14:00:00.000Z',
      isVisible: true,
      isDefaultSort: false,

      columnType: {
        id: '1',
        name: 'AVATAR',
        createdAt: '2021-08-10T14:00:00.000Z',
        updatedAt: '2021-08-10T14:00:00.000Z',
      },
    },
    {
      id: '2',
      createdAt: '2021-08-10T14:00:00.000Z',
      updatedAt: '2021-08-10T14:00:00.000Z',
      columnType: {
        id: '2',
        name: 'NAME',
        createdAt: '2021-08-10T14:00:00.000Z',
        updatedAt: '2021-08-10T14:00:00.000Z',
      },
    },
    {
      id: '7',
      createdAt: '2021-08-10T14:00:00.000Z',
      updatedAt: '2021-08-10T14:00:00.000Z',
      columnType: {
        id: '7',
        name: 'RENEWAL_DATE',
        createdAt: '2021-08-10T14:00:00.000Z',
        updatedAt: '2021-08-10T14:00:00.000Z',
      },
    },
    {
      id: '8',
      createdAt: '2021-08-10T14:00:00.000Z',
      updatedAt: '2021-08-10T14:00:00.000Z',
      columnType: {
        id: '8',
        name: 'FORECAST_ARR',
        createdAt: '2021-08-10T14:00:00.000Z',
        updatedAt: '2021-08-10T14:00:00.000Z',
      },
    },
    {
      id: '6',
      createdAt: '2021-08-10T14:00:00.000Z',
      updatedAt: '2021-08-10T14:00:00.000Z',
      columnType: {
        id: '6',
        name: 'RENEWAL_LIKELIHOOD',
        createdAt: '2021-08-10T14:00:00.000Z',
        updatedAt: '2021-08-10T14:00:00.000Z',
      },
    },
    {
      id: '9',
      createdAt: '2021-08-10T14:00:00.000Z',
      updatedAt: '2021-08-10T14:00:00.000Z',
      columnType: {
        id: '9',
        name: 'OWNER',
        createdAt: '2021-08-10T14:00:00.000Z',
        updatedAt: '2021-08-10T14:00:00.000Z',
      },
    },
    {
      id: '10',
      createdAt: '2021-08-10T14:00:00.000Z',
      updatedAt: '2021-08-10T14:00:00.000Z',
      columnType: {
        id: '10',
        name: 'LAST_TOUCHPOINT',
        createdAt: '2021-08-10T14:00:00.000Z',
        updatedAt: '2021-08-10T14:00:00.000Z',
      },
    },
  ],
  createdAt: '2021-08-10T14:00:00.000Z',
  updatedAt: '2021-08-10T14:00:00.000Z',
};

const tableViewDef3: TableViewDef = {
  id: '3',
  order: 2,
  name: 'Anual renewals',
  type: {
    id: '1',
    name: 'RENEWALS',
    createdAt: '2021-08-10T14:00:00.000Z',
    updatedAt: '2021-08-10T14:00:00.000Z',
  },
  filters: '',
  columns: [
    {
      id: '1',
      createdAt: '2021-08-10T14:00:00.000Z',
      updatedAt: '2021-08-10T14:00:00.000Z',
      isVisible: true,
      isDefaultSort: false,

      columnType: {
        id: '1',
        name: 'AVATAR',
        createdAt: '2021-08-10T14:00:00.000Z',
        updatedAt: '2021-08-10T14:00:00.000Z',
      },
    },
    {
      id: '2',
      createdAt: '2021-08-10T14:00:00.000Z',
      updatedAt: '2021-08-10T14:00:00.000Z',
      columnType: {
        id: '2',
        name: 'NAME',
        createdAt: '2021-08-10T14:00:00.000Z',
        updatedAt: '2021-08-10T14:00:00.000Z',
      },
    },
    {
      id: '7',
      createdAt: '2021-08-10T14:00:00.000Z',
      updatedAt: '2021-08-10T14:00:00.000Z',
      columnType: {
        id: '7',
        name: 'RENEWAL_DATE',
        createdAt: '2021-08-10T14:00:00.000Z',
        updatedAt: '2021-08-10T14:00:00.000Z',
      },
    },
    {
      id: '8',
      createdAt: '2021-08-10T14:00:00.000Z',
      updatedAt: '2021-08-10T14:00:00.000Z',
      columnType: {
        id: '8',
        name: 'FORECAST_ARR',
        createdAt: '2021-08-10T14:00:00.000Z',
        updatedAt: '2021-08-10T14:00:00.000Z',
      },
    },
    {
      id: '6',
      createdAt: '2021-08-10T14:00:00.000Z',
      updatedAt: '2021-08-10T14:00:00.000Z',
      columnType: {
        id: '6',
        name: 'RENEWAL_LIKELIHOOD',
        createdAt: '2021-08-10T14:00:00.000Z',
        updatedAt: '2021-08-10T14:00:00.000Z',
      },
    },
    {
      id: '9',
      createdAt: '2021-08-10T14:00:00.000Z',
      updatedAt: '2021-08-10T14:00:00.000Z',
      columnType: {
        id: '9',
        name: 'OWNER',
        createdAt: '2021-08-10T14:00:00.000Z',
        updatedAt: '2021-08-10T14:00:00.000Z',
      },
    },
    {
      id: '10',
      createdAt: '2021-08-10T14:00:00.000Z',
      updatedAt: '2021-08-10T14:00:00.000Z',
      columnType: {
        id: '10',
        name: 'LAST_TOUCHPOINT',
        createdAt: '2021-08-10T14:00:00.000Z',
        updatedAt: '2021-08-10T14:00:00.000Z',
      },
    },
  ],
  createdAt: '2021-08-10T14:00:00.000Z',
  updatedAt: '2021-08-10T14:00:00.000Z',
};

const mockedTableViewDefs = [tableViewDef1, tableViewDef2, tableViewDef3];

export default function RenewalsPage() {
  const isRestoring = useIsRestoring();
  const searchParams = useSearchParams();
  // const enableFeature = useFeatureIsOn('gp-dedicated-1');
  const [sorting, setSorting] = useState<SortingState>([
    { id: 'LAST_TOUCHPOINT', desc: true },
  ]);

  const preset = searchParams?.get('preset');

  useFilterSetter(tableViewDef1.filters);

  const {
    data,
    tableRef,
    isLoading,
    isFetching,
    totalCount,
    hasNextPage,
    fetchNextPage,
    totalAvailable,
  } = useRenewalsPageData({ sorting });

  const handleFetchMore = useCallback(() => {
    !isFetching && fetchNextPage();
  }, [fetchNextPage, isFetching]);

  const tableViewDef = mockedTableViewDefs.find((t) => t.id === preset);

  if (!tableViewDef || totalAvailable === 0) {
    return <EmptyState />;
  }

  const columns = getColumnsConfig(tableViewDef);

  return (
    <GridItem h='100%' area='content' overflowX='hidden' overflowY='auto'>
      <Search />

      <Table<RenewalRecord>
        data={data}
        columns={columns}
        sorting={sorting}
        tableRef={tableRef}
        canFetchMore={hasNextPage}
        onSortingChange={setSorting}
        onFetchMore={handleFetchMore}
        isLoading={isRestoring ? false : isLoading}
        totalItems={isRestoring ? 40 : totalCount || 0}
      />
    </GridItem>
  );
}
