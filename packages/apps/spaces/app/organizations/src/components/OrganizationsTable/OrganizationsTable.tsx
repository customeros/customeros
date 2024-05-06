'use client';

import { useState, useCallback } from 'react';
import { useSearchParams } from 'next/navigation';

import { observer } from 'mobx-react-lite';
import { useFeatureIsOn } from '@growthbook/growthbook-react';

import { Organization } from '@graphql/types';
import { useStore } from '@shared/hooks/useStore';
import { Table, SortingState } from '@ui/presentation/Table';

import { useOrganizationsPageData } from '../../hooks';
import { TableActions } from '../../components/Actions';
import { getColumnsConfig } from '../Columns/columnsDictionary';
import { EmptyState } from '../../components/EmptyState/EmptyState';

export const OrganizationsTable = observer(() => {
  const searchParams = useSearchParams();

  const enableFeature = useFeatureIsOn('gp-dedicated-1');
  const [sorting, setSorting] = useState<SortingState>([
    { id: 'LAST_TOUCHPOINT', desc: true },
  ]);

  const { tableViewDefsStore } = useStore();
  const preset = searchParams?.get('preset');
  const tableViewDef =
    tableViewDefsStore.getById(preset ?? '1') ||
    tableViewDefsStore.getDefault();

  const {
    data,
    tableRef,
    isFetching,
    totalCount,
    isLoading,
    hasNextPage,
    fetchNextPage,
    totalAvailable,
    allOrganizationIds,
  } = useOrganizationsPageData({
    sorting,
  });

  const handleFetchMore = useCallback(() => {
    !isFetching && fetchNextPage();
  }, [fetchNextPage, isFetching]);

  const columns = getColumnsConfig(tableViewDef?.value);

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
      isLoading={isLoading}
      totalItems={isLoading ? 40 : totalCount || 0}
      renderTableActions={(table) => (
        <TableActions table={table} allOrganizationsIds={allOrganizationIds} />
      )}
    />
  );
});
