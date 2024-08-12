import { useSearchParams } from 'react-router-dom';
import { useRef, useMemo, useState, useEffect } from 'react';

import { useKey } from 'rooks';
import { inPlaceSort } from 'fast-sort';
import { observer } from 'mobx-react-lite';
import difference from 'lodash/difference';
import intersection from 'lodash/intersection';
import { ColumnDef } from '@tanstack/react-table';
import { OnChangeFn } from '@tanstack/table-core';
import { InvoiceStore } from '@store/Invoices/Invoice.store';
import { useFeatureIsOn } from '@growthbook/growthbook-react';
import { ContactStore } from '@store/Contacts/Contact.store.ts';
import { CommandMenuType } from '@store/UI/CommandMenu.store.ts';
import { useTableActions } from '@invoices/hooks/useTableActions';
import { ContractStore } from '@store/Contracts/Contract.store.ts';
import { OrganizationStore } from '@store/Organizations/Organization.store.ts';

import { useStore } from '@shared/hooks/useStore';
import { Invoice, WorkflowType, TableViewType } from '@graphql/types';
import { ConfirmDeleteDialog } from '@ui/overlay/AlertDialog/ConfirmDeleteDialog';
import {
  Table,
  SortingState,
  TableInstance,
  RowSelectionState,
} from '@ui/presentation/Table';
import {
  getContractSortFn,
  getContractColumnsConfig,
} from '@organizations/components/Columns/contracts';

import { SidePanel } from '../SidePanel';
import { EmptyState } from '../EmptyState/EmptyState';
import { getColumnSortFn } from '../Columns/invoices/sortFns';
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
    OrganizationStore | ContactStore | InvoiceStore | ContractStore
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
  const contractColumns = getContractColumnsConfig(tableViewDef?.value);
  const organizationColumns = getOrganizationColumnsConfig(tableViewDef?.value);
  const invoiceColumns = getInvoiceColumnsConfig(tableViewDef?.value);
  const tableColumns = useMemo(() => {
    switch (tableType) {
      case TableViewType.Organizations:
        return organizationColumns;
      case TableViewType.Contacts:
        return contactColumns;
      case TableViewType.Contracts:
        return contractColumns;
      default:
        return invoiceColumns;
    }
  }, [tableType]) as ColumnDef<
    OrganizationStore | ContactStore | InvoiceStore | ContractStore
  >[];
  const isCommandMenuPrompted = store.ui.commandMenu.isOpen;

  const removeAccents = (str: string) => {
    return str
      .toLowerCase()
      .normalize('NFD')
      .replace(/[\u0300-\u036f]/g, '');
  };

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
      const normalizedSearchTerm = removeAccents(searchTerm);

      arr = arr.filter((entity) => {
        const name = entity.value?.name || '';

        return removeAccents(name).includes(normalizedSearchTerm);
      });
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
      const normalizedSearchTerm = removeAccents(searchTerm);

      arr = arr.filter((entity) => {
        const name = entity.value?.name || '';

        return removeAccents(name).includes(normalizedSearchTerm);
      });
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

  const contractsData = store.contracts?.toComputedArray((arr) => {
    if (tableType !== TableViewType.Contracts) return arr;

    // todo uncommment when filters are added
    // const filters = getContactFilterFns(tableViewDef?.getFilters());
    // if (filters) {
    //   arr = arr.filter((v) => filters.every((fn) => fn(v)));
    // }
    if (searchTerm) {
      const normalizedSearchTerm = removeAccents(searchTerm);

      arr = arr.filter((entity) => {
        const name = entity.value?.contractName || '';

        return removeAccents(name).includes(normalizedSearchTerm);
      });
    }

    if (tableType) {
      const columnId = sorting[0]?.id;
      const isDesc = sorting[0]?.desc;

      const computed = inPlaceSort(arr)?.[isDesc ? 'desc' : 'asc'](
        getContractSortFn(columnId),
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
      : tableType === TableViewType.Contracts
      ? contractsData
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

  const selectedIds = Object.keys(selection);

  useEffect(() => {
    if (selectedIds.length > 0 && !isCommandMenuPrompted) {
      if (tableType === TableViewType.Organizations) {
        if (selectedIds.length === 1) {
          store.ui.commandMenu.setType('OrganizationCommands');
          store.ui.commandMenu.setContext({
            entity: 'Organization',
            ids: selectedIds,
          });
        }

        if (selectedIds.length > 1) {
          store.ui.commandMenu.setType('OrganizationBulkCommands');
          store.ui.commandMenu.setContext({
            entity: 'Organizations',
            ids: selectedIds,
          });
        }
      } else {
        if (selectedIds.length === 1) {
          store.ui.commandMenu.setType('ContactCommands');
          store.ui.commandMenu.setContext({
            entity: 'Contact',
            ids: selectedIds,
          });
        }

        if (selectedIds.length > 1) {
          store.ui.commandMenu.setType('ContactBulkCommands');
          store.ui.commandMenu.setContext({
            entity: 'Contact',
            ids: selectedIds,
          });
        }
      }
    }
  }, [isCommandMenuPrompted, selectedIds.length]);

  useEffect(() => {
    tableRef.current?.resetRowSelection();

    if (tableType === TableViewType.Organizations) {
      store.ui.commandMenu.setType('OrganizationHub');
    } else {
      store.ui.commandMenu.setType('ContactHub');
    }
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
  const isSearching =
    store.ui.isSearching === tableViewDef?.value?.tableType?.toLowerCase();

  const focusedId =
    typeof focusIndex === 'number' ? data?.[focusIndex]?.id : null;
  const targetInvoice: Invoice = data?.find(
    (i) => i.value.metadata?.id === targetId,
  )?.value as Invoice;
  const targetInvoiceNumber = targetInvoice?.invoiceNumber || '';
  const targetInvoiceEmail = targetInvoice?.customer?.email || '';

  const createSocial = () => {
    if (!focusedId) return;
    store.ui.commandMenu.setType('AddContactViaLinkedInUrl');

    store.ui.commandMenu.setOpen(true);
    store.ui.commandMenu.setContext({
      entity: 'Organization',
      ids: [focusedId],
    });
  };

  const handleSetFocused = (index: number | null) => {
    if (isCommandMenuPrompted) return;

    setFocusIndex(index);

    if (tableType === TableViewType.Organizations) {
      if (typeof index !== 'number') {
        store.ui.commandMenu.setType('OrganizationHub');

        return;
      }

      if (index > -1 && Object.keys(selection).length === 0) {
        store.ui.commandMenu.setType('OrganizationCommands');
        store.ui.commandMenu.setContext({
          entity: 'Organization',
          ids: [data?.[index]?.id],
        });
      }
    }

    if (tableType === TableViewType.Contacts) {
      if (typeof index !== 'number') {
        store.ui.commandMenu.setType('ContactHub');

        return;
      }

      if (index > -1 && Object.keys(selection).length === 0) {
        store.ui.commandMenu.setType('ContactCommands');
        store.ui.commandMenu.setContext({
          entity: 'Contact',
          ids: [data?.[index]?.id],
        });
      }
    }
  };

  const handleOpenCommandKMenu = () => {
    const selectedIds = Object.keys(selection);

    if (tableType === TableViewType.Organizations) {
      if (selectedIds.length === 1) {
        store.ui.commandMenu.setType('OrganizationCommands');
        store.ui.commandMenu.setContext({
          entity: 'Organization',
          ids: selectedIds,
        });
      }

      if (selectedIds.length > 1) {
        store.ui.commandMenu.setType('OrganizationBulkCommands');
        store.ui.commandMenu.setContext({
          entity: 'Organizations',
          ids: selectedIds,
        });
      }
    } else {
      if (selectedIds.length === 1) {
        store.ui.commandMenu.setType('ContactCommands');
        store.ui.commandMenu.setContext({
          entity: 'Contact',
          ids: selectedIds,
        });
      }

      if (selectedIds.length > 1) {
        store.ui.commandMenu.setType('ContactBulkCommands');
        store.ui.commandMenu.setContext({
          entity: 'Contact',
          ids: selectedIds,
        });
      }
    }

    store.ui.commandMenu.setOpen(true);
  };

  return (
    <div className='flex'>
      <Table<OrganizationStore | ContactStore | InvoiceStore | ContractStore>
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
        enableRowSelection={
          tableType &&
          [TableViewType.Invoices, TableViewType.Contracts].includes(tableType)
            ? false
            : enableFeature !== null
            ? enableFeature
            : true
        }
        enableTableActions={
          tableType &&
          [TableViewType.Invoices, TableViewType.Contracts].includes(tableType)
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
              tableId={tableViewDef?.value.tableId}
              onOpenCommandK={handleOpenCommandKMenu}
              isCommandMenuOpen={isCommandMenuPrompted}
              onUpdateStage={store.organizations.updateStage}
              table={table as TableInstance<OrganizationStore>}
              handleOpen={(type: CommandMenuType) => {
                handleOpenCommandKMenu();
                store.ui.commandMenu.setType(type);
              }}
              enableKeyboardShortcuts={
                !isEditing &&
                !isFiltering &&
                !isSearching &&
                !isCommandMenuPrompted
              }
              onHide={() => {
                store.ui.commandMenu.setCallback(() =>
                  table.resetRowSelection(),
                );
                handleOpenCommandKMenu();
                store.ui.commandMenu.setType('DeleteConfirmationModal');
              }}
            />
          ) : (
            <ContactTableActions
              focusedId={focusedId}
              onAddTags={store.contacts.updateTags}
              onOpenCommandK={handleOpenCommandKMenu}
              table={table as TableInstance<ContactStore>}
              handleOpen={(type: CommandMenuType) => {
                handleOpenCommandKMenu();
                store.ui.commandMenu.setType(type);
              }}
              enableKeyboardShortcuts={
                !isSearching &&
                !isFiltering &&
                !isEditing &&
                !isCommandMenuPrompted
              }
              onHideContacts={() => {
                store.ui.commandMenu.setCallback(() =>
                  table.resetRowSelection(),
                );
                handleOpenCommandKMenu();
                store.ui.commandMenu.setType('DeleteConfirmationModal');
              }}
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
        confirmButtonLabel='Void invoice'
        label={`Void invoice ${targetInvoiceNumber}`}
        description={`Voiding this invoice will send an email notification to ${targetInvoiceEmail}`}
      />
    </div>
  );
});
