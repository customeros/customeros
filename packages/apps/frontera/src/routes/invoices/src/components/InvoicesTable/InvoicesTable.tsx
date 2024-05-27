import { useSearchParams } from 'react-router-dom';
import { useMemo, useState, useCallback } from 'react';

import { observer } from 'mobx-react-lite';

import { Invoice } from '@graphql/types';
import { useStore } from '@shared/hooks/useStore';
import { Table, SortingState } from '@ui/presentation/Table';
import { SlashCircle01 } from '@ui/media/icons/SlashCircle01';
import { ViewSettings } from '@shared/components/ViewSettings';
import { ConfirmDeleteDialog } from '@ui/overlay/AlertDialog/ConfirmDeleteDialog/ConfirmDeleteDialog';

import { Empty } from '../Empty';
import { Search } from '../Search';
import { useTableActions } from '../../hooks/useTableActions';
import { getColumnsConfig } from '../../components/Columns/Columns';
import { useInvoicesPageData } from '../../hooks/useInvoicesPageData';

export const InvoicesTable = observer(() => {
  const [searchParams] = useSearchParams();
  const [sorting, setSorting] = useState<SortingState>([
    { id: 'INVOICE_DUE_DATE', desc: true },
  ]);
  const store = useStore();

  const preset = searchParams?.get('preset');

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
  } = useInvoicesPageData({ sorting });

  const handleFetchMore = useCallback(() => {
    !isFetching && fetchNextPage();
  }, [fetchNextPage, isFetching]);

  const tableViewDef = store.tableViewDefs.getById(preset ?? '1');

  const columns = useMemo(
    () => getColumnsConfig(tableViewDef?.value),
    [tableViewDef?.value],
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
      <div className='flex items-center'>
        <Search />
        <ViewSettings type='invoices' />
      </div>
      <Table<Invoice>
        data={data}
        columns={columns}
        sorting={sorting}
        tableRef={tableRef}
        canFetchMore={hasNextPage}
        onSortingChange={setSorting}
        onFetchMore={handleFetchMore}
        rowHeight={40}
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
});
