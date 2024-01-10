'use client';

import { useState, useCallback } from 'react';

import { useIsRestoring } from '@tanstack/react-query';
import { useFeatureIsOn } from '@growthbook/growthbook-react';

import { GridItem } from '@ui/layout/Grid';
import { Organization, TableViewDef } from '@graphql/types';
import { Table, SortingState } from '@ui/presentation/Table';

import { KMenu } from './src/components/KMenu';
import { Search } from './src/components/Search';
import { useOrganizationsPageData } from './src/hooks';
import { TableActions } from './src/components/Actions';
import { EmptyState } from './src/components/EmptyState/EmptyState';
import { getColumnConfig } from './src/components/Columns/columnsDictionary';

const tableViewDef: TableViewDef = {
  id: '1',
  order: 0,
  name: 'Organizations',
  type: {
    id: '1',
    name: 'Organization',
    createdAt: '2021-08-10T14:00:00.000Z',
    updatedAt: '2021-08-10T14:00:00.000Z',
  },
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
      id: '3',
      createdAt: '2021-08-10T14:00:00.000Z',
      updatedAt: '2021-08-10T14:00:00.000Z',
      columnType: {
        id: '3',
        name: 'WEBSITE',
        createdAt: '2021-08-10T14:00:00.000Z',
        updatedAt: '2021-08-10T14:00:00.000Z',
      },
    },
    {
      id: '4',
      createdAt: '2021-08-10T14:00:00.000Z',
      updatedAt: '2021-08-10T14:00:00.000Z',
      columnType: {
        id: '4',
        name: 'RELATIONSHIP',
        createdAt: '2021-08-10T14:00:00.000Z',
        updatedAt: '2021-08-10T14:00:00.000Z',
      },
    },
    {
      id: '5',
      createdAt: '2021-08-10T14:00:00.000Z',
      updatedAt: '2021-08-10T14:00:00.000Z',
      columnType: {
        id: '5',
        name: 'ONBOARDING_STATUS',
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

const columnConfig = getColumnConfig(tableViewDef);

export default function OrganizationsPage() {
  const isRestoring = useIsRestoring();
  const enableFeature = useFeatureIsOn('gp-dedicated-1');
  const [sorting, setSorting] = useState<SortingState>([
    { id: 'LAST_TOUCHPOINT', desc: true },
  ]);

  const {
    data,
    tableRef,
    isLoading,
    isFetching,
    totalCount,
    hasNextPage,
    fetchNextPage,
    totalAvailable,
    allOrganizationIds,
  } = useOrganizationsPageData({ sorting });

  const handleFetchMore = useCallback(() => {
    !isFetching && fetchNextPage();
  }, [fetchNextPage, isFetching]);

  if (totalAvailable === 0) {
    return <EmptyState />;
  }

  return (
    <GridItem h='100%' area='content' overflowX='hidden' overflowY='auto'>
      <Search />

      <Table<Organization>
        data={data}
        columns={columnConfig}
        sorting={sorting}
        tableRef={tableRef}
        enableTableActions={enableFeature !== null ? enableFeature : true}
        enableRowSelection={enableFeature !== null ? enableFeature : true}
        canFetchMore={hasNextPage}
        onSortingChange={setSorting}
        onFetchMore={handleFetchMore}
        isLoading={isRestoring ? false : isLoading}
        totalItems={isRestoring ? 40 : totalCount || 0}
        renderTableActions={(table) => (
          <TableActions
            table={table}
            allOrganizationsIds={allOrganizationIds}
          />
        )}
      />

      <KMenu />
    </GridItem>
  );
}
