'use client';

import { useSearchParams } from 'next/navigation';
import { useMemo, useState, useCallback } from 'react';

import { useIsRestoring } from '@tanstack/react-query';

import { RenewalRecord } from '@graphql/types';
import { Table, SortingState } from '@ui/presentation/Table';
import { mockedTableDefs } from '@shared/util/tableDefs.mock';

import { useRenewalsPageData } from '../../hooks';
import { getColumnsConfig } from '../../components/Columns/Columns';
import { EmptyState } from '../../components/EmptyState/EmptyState';

export function RenewalsTable() {
  const isRestoring = useIsRestoring();
  const searchParams = useSearchParams();
  const [sorting, setSorting] = useState<SortingState>([
    { id: 'LAST_TOUCHPOINT', desc: true },
  ]);
  const tableViewDefsData = mockedTableDefs;

  const preset = searchParams?.get('preset');

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

  const tableViewDef = tableViewDefsData?.find((t) => t.id === preset);
  const columns = useMemo(
    () => getColumnsConfig(tableViewDef),
    [tableViewDef?.id],
  );

  if (!columns.length || totalAvailable === 0) {
    return <EmptyState />;
  }

  return (
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
  );
}
