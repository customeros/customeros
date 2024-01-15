'use client';

import { useState, useCallback } from 'react';

import { useIsRestoring } from '@tanstack/react-query';
import { useFeatureIsOn } from '@growthbook/growthbook-react';

import { Organization } from '@graphql/types';
import { Table, SortingState } from '@ui/presentation/Table';

import { useOrganizationsPageData } from './src/hooks';
import { TableActions } from './src/components/Actions';
import { columns } from './src/components/Columns/Columns';
import { EmptyState } from './src/components/EmptyState/EmptyState';

export function OrganizationsTable() {
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
    <Table<Organization>
      data={data}
      columns={columns}
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
        <TableActions table={table} allOrganizationsIds={allOrganizationIds} />
      )}
    />
  );
}
