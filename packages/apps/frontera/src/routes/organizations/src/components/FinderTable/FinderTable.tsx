import { useSearchParams } from 'react-router-dom';
import { useRef, useState, useEffect } from 'react';

import { useKey } from 'rooks';
import { inPlaceSort } from 'fast-sort';
import { observer } from 'mobx-react-lite';
import difference from 'lodash/difference';
import intersection from 'lodash/intersection';
import { OnChangeFn } from '@tanstack/table-core';
import { InvoiceStore } from '@store/Invoices/Invoice.store';
import { useFeatureIsOn } from '@growthbook/growthbook-react';
import { ContactStore } from '@store/Contacts/Contact.store.ts';
import { useTableActions } from '@invoices/hooks/useTableActions';
import { OrganizationStore } from '@store/Organizations/Organization.store.ts';

import { useStore } from '@shared/hooks/useStore';
import { SlashCircle01 } from '@ui/media/icons/SlashCircle01';
import { Invoice, WorkflowType, TableViewType } from '@graphql/types';
import { ConfirmDeleteDialog } from '@ui/overlay/AlertDialog/ConfirmDeleteDialog';
import {
  Table,
  SortingState,
  TableInstance,
  RowSelectionState,
} from '@ui/presentation/Table';

import { SidePanel } from '../SidePanel';
import { EmptyState } from '../EmptyState/EmptyState';
import { getColumnSortFn } from '../Columns/invoices/sortFns';
import { MergedColumnDefs } from '../Columns/shared/util/types';
import { getInvoiceFilterFns } from '../Columns/invoices/filterFns';
import { getInvoiceColumnsConfig } from '../Columns/invoices/columns';
import { getFlowFilterFns } from '../Columns/organizations/flowFilters';
import { ContactTableActions, OrganizationTableActions } from '../Actions';
import {
  getContactSortFn,
  getContactFilterFns,
  getContactColumnsConfig,
} from '../Columns/contacts';
import {
  getOrganizationSortFn,
  getOrganizationFilterFns,
  getOrganizationColumnsConfig,
} from '../Columns/organizations';

interface FinderTableProps {
  isSidePanelOpen: boolean;
}

export const FinderTable = observer(({ isSidePanelOpen }: FinderTableProps) => {
  const store = useStore();
  const [searchParams] = useSearchParams();
  const enableFeature = useFeatureIsOn('gp-dedicated-1');
  const tableRef = useRef<TableInstance<
    OrganizationStore | ContactStore | InvoiceStore
  > | null>(null);
  const [sorting, setSorting] = useState<SortingState>([
    { id: 'ORGANIZATIONS_LAST_TOUCHPOINT', desc: true },
  ]);
  const [selection, setSelection] = useState<RowSelectionState>({});
  const [isShiftPressed, setIsShiftPressed] = useState(false);
  const [focusIndex, setFocusIndex] = useState<number | null>(null);
  const [selectedIndex, setSelectedIndex] = useState<number | null>(null);
  const searchTerm = searchParams?.get('search');
  const { reset, targetId, isConfirming, onConfirm } = useTableActions();

  const preset = searchParams?.get('preset');
  const tableViewDef = store.tableViewDefs.getById(preset ?? '1');
  const tableType = tableViewDef?.value?.tableType;
  const getWorkFlow = store.workFlows
    .toArray()
    .filter((wf) => wf.value.type === WorkflowType.IdealCustomerProfile);

  const getWorkFlowId = getWorkFlow.map((wf) => wf.value.id);
  const workFlow = store.workFlows.getByType(getWorkFlowId[0]);
  const flowFiltersStatus = store.ui.isFilteringICP;

  const contactColumns = getContactColumnsConfig(tableViewDef?.value);
  const organizationColumns = getOrganizationColumnsConfig(tableViewDef?.value);
  const invoiceColumns = getInvoiceColumnsConfig(tableViewDef?.value);
  const tableColumns = (
    tableType === TableViewType.Organizations
      ? organizationColumns
      : tableType === TableViewType.Contacts
      ? contactColumns
      : invoiceColumns
  ) as MergedColumnDefs;

  const organizationsData = store.organizations?.toComputedArray((arr) => {
    if (tableType !== TableViewType.Organizations) return arr;
    const filters = getOrganizationFilterFns(tableViewDef?.getFilters());

    const flowFilters = getFlowFilterFns(workFlow?.getFilters());

    if (flowFilters.length && flowFiltersStatus) {
      arr = arr.filter((v) => !flowFilters.every((fn) => fn(v)));
    }

    if (filters) {
      arr = arr.filter((v) => filters.every((fn) => fn(v)));
    }

    if (searchTerm) {
      arr = arr.filter((entity) =>
        entity.value?.name
          ?.toLowerCase()
          .includes(searchTerm?.toLowerCase() as string),
      );
    }

    if (tableType) {
      const columnId = sorting[0]?.id;
      const isDesc = sorting[0]?.desc;

      const computed = inPlaceSort(arr)?.[isDesc ? 'desc' : 'asc'](
        getOrganizationSortFn(columnId),
      );

      return computed;
    }

    return arr;
  });

  const contactsData = store.contacts?.toComputedArray((arr) => {
    if (tableType !== TableViewType.Contacts) return arr;

    const filters = getContactFilterFns(tableViewDef?.getFilters());

    if (filters) {
      arr = arr.filter((v) => filters.every((fn) => fn(v)));
    }

    if (searchTerm) {
      arr = arr.filter((entity) =>
        entity.value?.name
          ?.toLowerCase()
          .includes(searchTerm?.toLowerCase() as string),
      );
    }

    if (tableType) {
      const columnId = sorting[0]?.id;
      const isDesc = sorting[0]?.desc;

      const computed = inPlaceSort(arr)?.[isDesc ? 'desc' : 'asc'](
        getContactSortFn(columnId),
      );

      return computed;
    }

    return arr;
  });

  const invoicesData = store.invoices.toComputedArray((arr) => {
    if (tableType !== TableViewType.Invoices) return arr;
    const filters = getInvoiceFilterFns(tableViewDef?.getFilters());

    if (filters) {
      arr = arr.filter((v) => filters.every((fn) => fn(v)));
    }

    if (searchTerm) {
      arr = arr.filter((entity) =>
        entity.contract?.contractName
          ?.toLowerCase()
          .includes(searchTerm?.toLowerCase() as string),
      );
    }

    if (tableType) {
      const columnId = sorting[0]?.id;
      const isDesc = sorting[0]?.desc;

      const computed = inPlaceSort(arr)?.[isDesc ? 'desc' : 'asc'](
        getColumnSortFn(columnId),
      );

      return computed;
    }

    return arr;
  });

  const data =
    tableType === TableViewType.Organizations
      ? organizationsData
      : tableType === TableViewType.Contacts
      ? contactsData
      : invoicesData;

  const handleSelectionChange: OnChangeFn<RowSelectionState> = (
    nextSelection,
  ) => {
    if (!isShiftPressed) {
      setSelection(nextSelection);

      return;
    }

    if (isShiftPressed && selectedIndex !== null && focusIndex !== null) {
      setSelection((prev) => {
        const edgeIndexes = [
          Math.min(selectedIndex, focusIndex),
          Math.max(selectedIndex, focusIndex),
        ];

        const ids = data
          .slice(edgeIndexes[0], edgeIndexes[1] + 1)
          .map((d) => d.id);

        const newSelection: Record<string, boolean> = {
          ...prev,
        };

        const prevIds = Object.keys(prev);
        const diff = difference(ids, prevIds);
        const match = intersection(ids, prevIds);
        const shouldRemove = diff.length < match.length;

        const endId = data[edgeIndexes[1]].id;

        diff.forEach((id) => {
          newSelection[id] = true;
        });
        shouldRemove &&
          [endId, ...match].forEach((id) => {
            delete newSelection[id];
          });

        return newSelection;
      });
    }
  };

  useEffect(() => {
    tableRef.current?.resetRowSelection();
    store.ui.commandMenu.setType('OrganizationHub');
  }, [tableViewDef?.value.id]);

  useEffect(() => {
    store.ui.setSearchCount(data.length);
    store.ui.setFilteredTable(data);
  }, [data.length]);

  useKey(
    'Shift',
    (e) => {
      setIsShiftPressed(e.type === 'keydown');
    },
    { eventTypes: ['keydown', 'keyup'] },
  );

  if (
    (tableViewDef?.value.tableType === TableViewType.Organizations &&
      store.organizations?.toArray().length === 0 &&
      !store.organizations.isLoading) ||
    (tableViewDef?.value.tableType === TableViewType.Contacts &&
      store.contacts?.toArray().length === 0 &&
      !store.contacts.isLoading) ||
    (tableViewDef?.value.tableType === TableViewType.Invoices &&
      store.invoices?.toArray().length === 0 &&
      !store.invoices.isLoading)
  ) {
    return <EmptyState />;
  }

  const isEditing = store.ui.isEditingTableCell;
  const isFiltering = store.ui.isFilteringTable;
  const isCommandMenuPrompted = store.ui.commandMenu.isOpen;
  const isSearching =
    store.ui.isSearching === tableViewDef?.value?.tableType?.toLowerCase();

  const targetInvoice: Invoice = data?.find(
    (i) => i.value.metadata?.id === targetId,
  )?.value as Invoice;
  const targetInvoiceNumber = targetInvoice?.invoiceNumber || '';
  const targetInvoiceEmail = targetInvoice?.customer?.email || '';

  const createSocial = () => {
    if (!focusIndex) return;
    store.ui.commandMenu.setType('AddContactViaLinkedInUrl');

    store.ui.commandMenu.setOpen(true);
    store.ui.commandMenu.setContext({
      entity: 'Organization',
      id: data?.[focusIndex]?.id,
    });
  };

  const focusedId =
    typeof focusIndex === 'number' ? data?.[focusIndex]?.id : null;

  const handleSetFocused = (index: number | null) => {
    if (isCommandMenuPrompted) return;

    setFocusIndex(index);

    // Todo replace with match when command k actions are available for other table types
    if (tableType === TableViewType.Organizations) {
      if (typeof index !== 'number') {
        store.ui.commandMenu.setType('OrganizationHub');
        return;
      }

      if (index) {
        store.ui.commandMenu.setType('OrganizationCommands');
        store.ui.commandMenu.setContext({
          entity: 'Organization',
          id: data?.[index]?.id,
        });
      }
    }
  };

  return (
    <div className='flex'>
      <Table<OrganizationStore | ContactStore | InvoiceStore>
        data={data}
        manualFiltering
        sorting={sorting}
        tableRef={tableRef}
        selection={selection}
        columns={tableColumns}
        getRowId={(row) => row.id}
        enableColumnResizing={true}
        onSortingChange={setSorting}
        onFocusedRowChange={handleSetFocused}
        onSelectedIndexChange={setSelectedIndex}
        isLoading={store.organizations.isLoading}
        onSelectionChange={handleSelectionChange}
        fullRowSelection={tableType === TableViewType.Invoices}
        totalItems={store.organizations.isLoading ? 40 : data.length}
        enableKeyboardShortcuts={
          !isEditing && !isFiltering && !isCommandMenuPrompted
        }
        enableTableActions={
          tableType === TableViewType.Invoices
            ? false
            : enableFeature !== null
            ? enableFeature
            : true
        }
        enableRowSelection={
          tableType === TableViewType.Invoices
            ? false
            : enableFeature !== null
            ? enableFeature
            : true
        }
        renderTableActions={(table) =>
          tableType === TableViewType.Organizations ? (
            <OrganizationTableActions
              focusedId={focusedId}
              onCreateContact={createSocial}
              onHide={store.organizations.hide}
              onMerge={store.organizations.merge}
              tableId={tableViewDef?.value.tableId}
              onUpdateStage={store.organizations.updateStage}
              table={table as TableInstance<OrganizationStore>}
              enableKeyboardShortcuts={
                !isEditing &&
                !isFiltering &&
                !isSearching &&
                !isCommandMenuPrompted
              }
            />
          ) : (
            <ContactTableActions
              onAddTags={store.contacts.updateTags}
              onHideContacts={store.contacts.archive}
              table={table as TableInstance<ContactStore>}
              enableKeyboardShortcuts={
                !isSearching &&
                !isFiltering &&
                !isEditing &&
                !isCommandMenuPrompted
              }
            />
          )
        }
      />
      {isSidePanelOpen && <SidePanel />}
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
    </div>
  );
});
