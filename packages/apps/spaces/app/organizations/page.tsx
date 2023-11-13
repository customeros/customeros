'use client';

import { useState, useCallback } from 'react';

import { useIsRestoring } from '@tanstack/react-query';

import { GridItem } from '@ui/layout/Grid';
import { Organization } from '@graphql/types';
import { Table, SortingState } from '@ui/presentation/Table';

import { Search } from './src/components/Search';
import { useOrganizationsPageData } from './src/hooks';
import { TableActions } from './src/components/Actions';
import { columns } from './src/components/Columns/Columns';
import { EmptyState } from './src/components/EmptyState/EmptyState';

export default function OrganizationsPage() {
  const isRestoring = useIsRestoring();
  const [sorting, setSorting] = useState<SortingState>([
    { id: 'LAST_TOUCHPOINT', desc: false },
  ]);

  const {
    data,
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
        columns={columns}
        sorting={sorting}
        enableTableActions
        enableRowSelection
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
    </GridItem>
  );
}
