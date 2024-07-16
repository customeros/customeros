import { useRef, useState } from 'react';
import { useSearchParams } from 'react-router-dom';

import { inPlaceSort } from 'fast-sort';
import { observer } from 'mobx-react-lite';
import { InvoiceStore } from '@store/Invoices/Invoice.store.ts';

import { TableViewType } from '@graphql/types';
import { useStore } from '@shared/hooks/useStore';
import { SlashCircle01 } from '@ui/media/icons/SlashCircle01';
import { ViewSettings } from '@shared/components/ViewSettings';
import { Table, SortingState, TableInstance } from '@ui/presentation/Table';
import { ConfirmDeleteDialog } from '@ui/overlay/AlertDialog/ConfirmDeleteDialog/ConfirmDeleteDialog';

import { Search } from '../Search';
import { getColumnSortFn } from '../Columns/sortFns.ts';
import { getInvoiceColumnsConfig } from '../Columns/Columns';
import { getInvoiceFilterFns } from '../Columns/filterFns.ts';
import { useTableActions } from '../../hooks/useTableActions';

export const InvoicesTable = observer(() => {
  const [searchParams] = useSearchParams();
  const [sorting, setSorting] = useState<SortingState>([
    { id: 'INVOICE_DUE_DATE', desc: true },
  ]);
  const tableRef = useRef<TableInstance<InvoiceStore> | null>(null);

  const store = useStore();

  const preset = searchParams?.get('preset');
  const searchTerm = searchParams?.get('search');

  const tableViewDef = store.tableViewDefs.getById(preset ?? '1');

  const columns = getInvoiceColumnsConfig(tableViewDef?.value);

  const data = store.invoices.toComputedArray((arr) => {
    const predefinedFilter = getInvoiceFilterFns(tableViewDef?.getFilters());
    if (predefinedFilter) {
      arr = arr.filter((v) => predefinedFilter.every((fn) => fn(v)));
    }
    if (searchTerm) {
      arr = arr.filter((invoiceStore) => {
        const invoice = invoiceStore.value?.organization?.metadata?.id;

        return store.organizations.value
          ?.get(invoice)
          ?.value?.name?.toLowerCase()
          .includes(searchTerm?.toLowerCase());
      });
    }
    const columnId = sorting[0]?.id;
    const isDesc = sorting[0]?.desc;

    return inPlaceSort(arr)?.[isDesc ? 'desc' : 'asc'](
      getColumnSortFn(columnId),
    );
  });
  const { reset, targetId, isConfirming, onConfirm } = useTableActions();

  const targetInvoice = data?.find(
    (i) => i.value.metadata?.id === targetId,
  )?.value;
  const targetInvoiceNumber = targetInvoice?.invoiceNumber || '';
  const targetInvoiceEmail = targetInvoice?.customer?.email || '';

  return (
    <>
      <div className='flex items-center'>
        <Search />
        <ViewSettings type={TableViewType.Invoices} />
      </div>
      <Table<InvoiceStore>
        data={data}
        manualFiltering
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
