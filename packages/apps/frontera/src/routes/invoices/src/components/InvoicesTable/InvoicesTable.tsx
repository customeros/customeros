import { useRef, useMemo, useState } from 'react';
import { useSearchParams } from 'react-router-dom';

import { Store } from '@store/store.ts';
import { inPlaceSort } from 'fast-sort';
import { observer } from 'mobx-react-lite';

import { Invoice } from '@graphql/types';
import { useStore } from '@shared/hooks/useStore';
import { SlashCircle01 } from '@ui/media/icons/SlashCircle01';
import { ViewSettings } from '@shared/components/ViewSettings';
import { Table, SortingState, TableInstance } from '@ui/presentation/Table';
import { ConfirmDeleteDialog } from '@ui/overlay/AlertDialog/ConfirmDeleteDialog/ConfirmDeleteDialog';

import { Empty } from '../Empty';
import { Search } from '../Search';
import { useTableActions } from '../../hooks/useTableActions';
import {
  getColumnSortFn,
  getColumnsConfig,
  getPredefinedFilterFn,
} from '../../components/Columns/Columns';

export const InvoicesTable = observer(() => {
  const [searchParams] = useSearchParams();
  const [sorting, setSorting] = useState<SortingState>([
    { id: 'INVOICE_DUE_DATE', desc: true },
  ]);
  const tableRef = useRef<TableInstance<Store<Invoice>> | null>(null);

  const store = useStore();

  const preset = searchParams?.get('preset');
  const searchTerm = searchParams?.get('search');

  const tableViewDef = store.tableViewDefs.getById(preset ?? '1');

  const columns = useMemo(
    () => getColumnsConfig(tableViewDef?.value),
    [tableViewDef?.value],
  );

  const data = store.invoices.toComputedArray((arr: Store<Invoice>[]) => {
    const predefinedFilter = getPredefinedFilterFn(tableViewDef?.getFilters());
    if (predefinedFilter) {
      arr = arr.filter(predefinedFilter);
    }
    if (searchTerm) {
      arr = arr.filter((invoiceStore: Store<Invoice>) => {
        const invoice = invoiceStore.value?.organization?.metadata?.id;

        return store.organizations.value
          ?.get(invoice)
          ?.value?.name?.toLowerCase()
          .includes(searchTerm?.toLowerCase());
      });
    }
    const columnId = sorting[0]?.id;
    const isDesc = sorting[0]?.desc;

    return inPlaceSort<Store<Invoice>>(arr)?.[isDesc ? 'desc' : 'asc'](
      getColumnSortFn(columnId),
    );
  });
  const { reset, targetId, isConfirming, onConfirm } = useTableActions();

  const targetInvoice = data?.find(
    (i) => i.value.metadata?.id === targetId,
  )?.value;
  const targetInvoiceNumber = targetInvoice?.invoiceNumber || '';
  const targetInvoiceEmail = targetInvoice?.customer?.email || '';
  if (!columns.length || data?.length === 0) {
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
      <Table<Store<Invoice>>
        data={data}
        columns={columns}
        sorting={sorting}
        tableRef={tableRef}
        onSortingChange={setSorting}
        rowHeight={45}
        isLoading={store.invoices.isLoading}
        totalItems={store.invoices.isLoading ? 40 : data.length}
      />
      <ConfirmDeleteDialog
        onClose={reset}
        hideCloseButton
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
