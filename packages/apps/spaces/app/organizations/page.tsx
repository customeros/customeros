'use client';

import { useMemo, useState, useCallback } from 'react';

import { useLocalStorage } from 'usehooks-ts';
import { useIsRestoring } from '@tanstack/react-query';

import { GridItem } from '@ui/layout/Grid';
import { Organization } from '@graphql/types';
import { Table, SortingState } from '@ui/presentation/Table';

import { Search } from './src/components/Search';
import { TableActions } from './src/components/Actions';
import { getColumns } from './src/components/Columns/Columns';
import { EmptyState } from './src/components/EmptyState/EmptyState';
import {
  useOrganizationsPageData,
  useOrganizationsPageMethods,
} from './src/hooks';

export default function OrganizationsPage() {
  const isRestoring = useIsRestoring();
  const [sorting, setSorting] = useState<SortingState>([]);
  const { createOrganization } = useOrganizationsPageMethods();
  const [tabs] = useLocalStorage<{
    [key: string]: string;
  }>(`customeros-player-last-position`, { root: 'organization' });

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

  const handleCreateOrganization = useCallback(() => {
    createOrganization.mutate({ input: { name: '' } });
  }, [createOrganization]);

  const handleFetchMore = useCallback(() => {
    !isFetching && fetchNextPage();
  }, [fetchNextPage, isFetching]);

  const columns = useMemo(
    () =>
      getColumns({
        tabs,
        createIsLoading: createOrganization.isLoading,
        onCreateOrganization: handleCreateOrganization,
      }),
    [tabs, handleCreateOrganization, createOrganization, getColumns],
  );

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
