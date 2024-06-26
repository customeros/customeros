import { useSearchParams } from 'react-router-dom';
import { useMemo, useState, useCallback } from 'react';

import { observer } from 'mobx-react-lite';

import { useStore } from '@shared/hooks/useStore';
import { Table, SortingState } from '@ui/presentation/Table';
import { RenewalRecord, TableViewType } from '@graphql/types';
import { ViewSettings } from '@shared/components/ViewSettings';

import { Search } from '../Search';
import { useRenewalsPageData } from '../../hooks';
import { getColumnsConfig } from '../../components/Columns/Columns';
import { EmptyState } from '../../components/EmptyState/EmptyState';

export const RenewalsTable = observer(() => {
  const [searchParams] = useSearchParams();
  const [sorting, setSorting] = useState<SortingState>([
    { id: 'LAST_TOUCHPOINT', desc: true },
  ]);
  const store = useStore();
  const preset = searchParams?.get('preset');

  const tableViewDef = store.tableViewDefs.getById(preset ?? '1');

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
  } = useRenewalsPageData({ sorting });

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
      <div className='flex items-center w-full'>
        <Search />
        <ViewSettings type={TableViewType.Renewals} />
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
