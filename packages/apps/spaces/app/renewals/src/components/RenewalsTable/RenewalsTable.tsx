'use client';

import { useSearchParams } from 'next/navigation';
import { useMemo, useState, useCallback } from 'react';

import { observer } from 'mobx-react-lite';

import { useStore } from '@shared/hooks/useStore';
import { TableViewDef, RenewalRecord } from '@graphql/types';
import { Table, SortingState } from '@ui/presentation/Table';
import { ViewSettings } from '@shared/components/ViewSettings';

import type { GetRenewalsQuery } from '../../graphql/getRenewals.generated';

import { Search } from '../Search';
import { useRenewalsPageData } from '../../hooks';
import { getColumnsConfig } from '../../components/Columns/Columns';
import { EmptyState } from '../../components/EmptyState/EmptyState';

interface RenewalsTableProps {
  bootstrap: {
    renewals: GetRenewalsQuery;
    tableViewDefs: TableViewDef[];
  };
}

export const RenewalsTable = observer(({ bootstrap }: RenewalsTableProps) => {
  const searchParams = useSearchParams();
  const [sorting, setSorting] = useState<SortingState>([
    { id: 'LAST_TOUCHPOINT', desc: true },
  ]);
  const { tableViewDefsStore } = useStore();
  const preset = searchParams?.get('preset');

  const tableViewDef = tableViewDefsStore.getById(preset ?? '1');

  const {
    data,
    tableRef,
    isLoading,
    isFetching,
    totalCount,
    hasNextPage,
    isRefetching,
    fetchNextPage,
    totalAvailable,
  } = useRenewalsPageData({ sorting, initialData: bootstrap.renewals });

  const handleFetchMore = useCallback(() => {
    !isFetching && fetchNextPage();
  }, [fetchNextPage, isFetching]);

  const columns = useMemo(
    () => getColumnsConfig(tableViewDef?.value),
    [tableViewDef?.value],
  );

  if (!columns.length || totalAvailable === 0) {
    return <EmptyState />;
  }

  return (
    <>
      <div className='flex items-center'>
        <Search />
        <ViewSettings type='renewals' />
      </div>
      <Table<RenewalRecord>
        data={data}
        columns={columns}
        sorting={sorting}
        tableRef={tableRef}
        canFetchMore={hasNextPage}
        onSortingChange={setSorting}
        onFetchMore={handleFetchMore}
        isLoading={isLoading && !isRefetching}
        totalItems={isLoading ? 40 : totalCount || 0}
      />
    </>
  );
});
