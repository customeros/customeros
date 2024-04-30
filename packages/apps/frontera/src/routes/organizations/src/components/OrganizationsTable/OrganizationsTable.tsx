'use client';

import { useState, useCallback } from 'react';

import { useIsRestoring } from '@tanstack/react-query';
import { useFeatureIsOn } from '@growthbook/growthbook-react';

import { Organization } from '@graphql/types';
import { Table, SortingState } from '@ui/presentation/Table';
import { GetOrganizationsQuery } from '@organizations/graphql/getOrganizations.generated';

import { TableActions } from '../Actions';
import { columns } from '../Columns/Columns';
import { EmptyState } from '../EmptyState/EmptyState';
import { useOrganizationsPageData } from '../../hooks';

interface OrganizationsTableProps {
  initialData?: GetOrganizationsQuery;
}

export function OrganizationsTable({ initialData }: OrganizationsTableProps) {
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
  } = useOrganizationsPageData({ sorting, initialData });

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
