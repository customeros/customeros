'use client';

import { useSearchParams } from 'next/navigation';
import { useMemo, useState, useCallback } from 'react';

import { Invoice } from '@graphql/types';
import { Table, SortingState } from '@ui/presentation/Table';
import { mockedTableDefs } from '@shared/util/tableDefs.mock';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useTableViewDefsQuery } from '@shared/graphql/tableViewDefs.generated';

import { Empty } from '../Empty';
import { Search } from '../Search';
import { getColumnsConfig } from '../../components/Columns/Columns';
import { useInvoicesPageData } from '../../hooks/useInvoicesPageData';

export const InvoicesTable = () => {
  const client = getGraphQLClient();
  const searchParams = useSearchParams();
  const [sorting, setSorting] = useState<SortingState>([
    { id: 'INVOICE_DUE_DATE', desc: true },
  ]);

  const preset = searchParams?.get('preset');

  const { data: tableViewDefsData } = useTableViewDefsQuery(
    client,
    {
      pagination: { limit: 100, page: 1 },
    },
    {
      enabled: false,
      placeholderData: { tableViewDefs: { content: mockedTableDefs } },
    },
  );

  const {
    data,
    tableRef,
    isLoading,
    isFetching,
    totalCount,
    hasNextPage,
    fetchNextPage,
    totalAvailable,
  } = useInvoicesPageData({ sorting });

  const handleFetchMore = useCallback(() => {
    !isFetching && fetchNextPage();
  }, [fetchNextPage, isFetching]);

  const tableViewDef = tableViewDefsData?.tableViewDefs?.content?.find(
    (t) => t.id === preset,
  );
  const columns = useMemo(
    () => getColumnsConfig(tableViewDef),
    [tableViewDef?.id],
  );

  if (!columns.length || totalAvailable === 0) {
    return (
      <div className='flex justify-center'>
        <Empty />
      </div>
    );
  }

  return (
    <>
      <Search />
      <Table<Invoice>
        data={data}
        columns={columns}
        sorting={sorting}
        tableRef={tableRef}
        canFetchMore={hasNextPage}
        onSortingChange={setSorting}
        onFetchMore={handleFetchMore}
        isLoading={isFetching || isLoading}
        totalItems={isFetching ? 40 : totalCount || 0}
      />
    </>
  );
};
