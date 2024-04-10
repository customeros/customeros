'use client';

import { useSearchParams } from 'next/navigation';
import { useMemo, useState, useCallback } from 'react';

import { Invoice } from '@graphql/types';
import { Table, SortingState } from '@ui/presentation/Table';
import { mockedTableDefs } from '@shared/util/tableDefs.mock';
import { SlashCircle01 } from '@ui/media/icons/SlashCircle01';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { GetInvoicesQuery } from '@shared/graphql/getInvoices.generated';
import { useTableViewDefsQuery } from '@shared/graphql/tableViewDefs.generated';
import { ConfirmDeleteDialog } from '@ui/overlay/AlertDialog/ConfirmDeleteDialog/ConfirmDeleteDialog2';

import { Empty } from '../Empty';
import { Search } from '../Search';
import { useTableActions } from '../../hooks/useTableActions';
import { getColumnsConfig } from '../../components/Columns/Columns';
import { useInvoicesPageData } from '../../hooks/useInvoicesPageData';

interface InvoicesTableProps {
  initialData?: GetInvoicesQuery;
}

export const InvoicesTable = ({ initialData }: InvoicesTableProps) => {
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
    isRefetching,
    fetchNextPage,
    totalAvailable,
  } = useInvoicesPageData({ sorting, initialData });

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

  const { reset, targetId, isConfirming, onConfirm, isPending } =
    useTableActions();

  const targetInvoice = data?.find((i) => i.metadata.id === targetId);
  const targetInvoiceNumber = targetInvoice?.invoiceNumber || '';
  const targetInvoiceEmail = targetInvoice?.customer?.email || '';

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
        isLoading={isLoading && !isRefetching}
        totalItems={isLoading ? 40 : totalCount || 0}
      />
      <ConfirmDeleteDialog
        onClose={reset}
        hideCloseButton
        isLoading={isPending}
        isOpen={isConfirming}
        onConfirm={onConfirm}
        icon={<SlashCircle01 />}
        confirmButtonLabel='Void invoice'
        label={`Void invoice ${targetInvoiceNumber}`}
        description={`Voiding this invoice will send an email notification to ${targetInvoiceEmail}`}
      />
    </>
  );
};
