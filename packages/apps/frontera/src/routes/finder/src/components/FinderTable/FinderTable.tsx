import { useRef, useState, useEffect } from 'react';
import { useParams, useSearchParams } from 'react-router-dom';

import { match } from 'ts-pattern';
import { useKeyBindings } from 'rooks';
import { observer } from 'mobx-react-lite';
import { useFeatureIsOn } from '@growthbook/growthbook-react';
import { useColumnSizing } from '@finder/hooks/useColumnSizing';
import { useTableActions } from '@invoices/hooks/useTableActions';
import { OpportunitiesTableActions } from '@finder/components/Actions/OpportunityActions';

import { useStore } from '@shared/hooks/useStore';
import { Invoice, TableViewType, ColumnViewType } from '@graphql/types';
import { Table, SortingState, TableInstance } from '@ui/presentation/Table';
import { ConfirmDeleteDialog } from '@ui/overlay/AlertDialog/ConfirmDeleteDialog';

import { SidePanel } from '../SidePanel';
import { EmptyState } from '../EmptyState/EmptyState';
import { computeFinderData } from './computeFinderData';
import { computeFinderColumns } from './computeFinderColumns';
import { ContactPreviewCard } from '../ContactPreviewCard/ContactPreviewCard';
import {
  ContactTableActions,
  OrganizationTableActions,
  FlowSequencesTableActions,
} from '../Actions';

interface FinderTableProps {
  isSidePanelOpen: boolean;
}

export const FinderTable = observer(({ isSidePanelOpen }: FinderTableProps) => {
  const store = useStore();
  const params = useParams();
  const [searchParams] = useSearchParams();
  const enableFeature = useFeatureIsOn('gp-dedicated-1');
  const tableRef = useRef<TableInstance<object> | null>(null);
  const preset = searchParams?.get('preset');
  const tableViewDef = store.tableViewDefs.getById(preset ?? '1');

  const contactsPreset = store.tableViewDefs.contactsPreset;

  const [sorting, setSorting] = useState<SortingState>([
    preset === contactsPreset
      ? { id: ColumnViewType.ContactsCreatedAt, desc: true }
      : { id: ColumnViewType.OrganizationsLastTouchpoint, desc: true },
  ]);
  const filtersV2 = useFeatureIsOn('filters-v2');

  const searchTerm = searchParams?.get('search');
  const { reset, targetId, isConfirming, onConfirm } = useTableActions();

  const tableType =
    tableViewDef?.value?.tableType || TableViewType.Organizations;

  const columns = computeFinderColumns(store, {
    tableType,
    currentPreset: preset,
  });

  const data = computeFinderData(store, filtersV2, {
    sorting,
    tableViewDef,
    urlParams: params,
    searchTerm: searchTerm ?? '',
  });

  const isCommandMenuPrompted = store.ui.commandMenu.isOpen;
  const handleColumnSizing = useColumnSizing(columns, tableViewDef);

  const onSelectionChange = (selectedIds: string[]) => {
    if (selectedIds.length > 0 && !isCommandMenuPrompted) {
      store.ui.commandMenu.setCallback(() => {
        tableRef?.current?.resetRowSelection();
      });

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
      } else if (tableType === TableViewType.Opportunities) {
        if (selectedIds.length === 1) {
          store.ui.commandMenu.setType('OpportunityCommands');
          store.ui.commandMenu.setContext({
            entity: 'Opportunity',
            ids: selectedIds,
          });
        }

        if (selectedIds.length > 1) {
          store.ui.commandMenu.setType('OpportunityBulkCommands');
          store.ui.commandMenu.setContext({
            entity: 'Opportunities',
            ids: selectedIds,
          });
        }
      } else if (tableType === TableViewType.Flow) {
        if (selectedIds.length === 1) {
          store.ui.commandMenu.setType('FlowCommands');
          store.ui.commandMenu.setContext({
            entity: 'Flow',
            ids: selectedIds,
          });
        }

        if (selectedIds.length > 1) {
          store.ui.commandMenu.setType('FlowsBulkCommands');
          store.ui.commandMenu.setContext({
            entity: 'Flows',
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
  };

  useEffect(() => {
    tableRef.current?.resetRowSelection();

    if (tableType === TableViewType.Organizations) {
      store.ui.commandMenu.setType('OrganizationHub');
    } else if (tableType === TableViewType.Opportunities) {
      store.ui.commandMenu.setType('OpportunityHub');
    } else if (tableType === TableViewType.Flow) {
      store.ui.commandMenu.setType('FlowHub');
    } else {
      store.ui.commandMenu.setType('ContactHub');
    }
  }, [tableViewDef?.value.id]);

  useEffect(() => {
    store.ui.setSearchCount(data.length);
    store.ui.setFilteredTable(data);
  }, [data.length]);

  const isEditing = store.ui.isEditingTableCell;
  const isFiltering = store.ui.isFilteringTable;
  const isSearching =
    store.ui.isSearching === tableViewDef?.value?.tableType?.toLowerCase();

  const [targetInvoiceNumber, targetInvoiceEmail] = match(tableType)
    .with(TableViewType.Invoices, () => {
      const invoice = data?.find((i) => i.value.metadata?.id === targetId)
        ?.value as Invoice;

      const targetInvoiceNumber = invoice?.invoiceNumber || '';
      const targetInvoiceEmail = invoice?.customer?.email || '';

      return [targetInvoiceNumber, targetInvoiceEmail];
    })
    .otherwise(() => ['', '']);

  const handleSetFocused = (index: number | null, selectedIds: string[]) => {
    if (isCommandMenuPrompted) return;

    if (selectedIds.length > 0) return;

    if (tableType === TableViewType.Organizations) {
      if (typeof index !== 'number') {
        store.ui.commandMenu.setType('OrganizationHub');

        return;
      }

      if (index > -1 && selectedIds.length === 0) {
        store.ui.commandMenu.setType('OrganizationCommands');
        store.ui.commandMenu.setContext({
          entity: 'Organization',
          ids: [data?.[index]?.id],
        });
      }
    }

    if (tableType === TableViewType.Contacts) {
      if (!store.ui.contactPreviewCardOpen) {
        if (index !== null) {
          store.ui.setFocusRow(data?.[index]?.id);
        }
      }

      if (typeof index !== 'number') {
        store.ui.commandMenu.setType('ContactHub');

        return;
      }

      if (index > -1 && selectedIds.length === 0) {
        store.ui.commandMenu.setType('ContactCommands');
        store.ui.commandMenu.setContext({
          entity: 'Contact',
          ids: [data?.[index]?.id],
        });
      }
    }

    if (tableType === TableViewType.Flow) {
      if (typeof index !== 'number') {
        store.ui.commandMenu.setType('FlowHub');

        return;
      }

      if (index > -1 && selectedIds.length === 0) {
        store.ui.commandMenu.setType('FlowCommands');
        store.ui.commandMenu.setContext({
          entity: 'Flow',
          ids: [data?.[index]?.id],
        });
      }
    }

    if (tableType === TableViewType.Opportunities) {
      if (typeof index !== 'number') {
        store.ui.commandMenu.setType('OpportunityHub');

        return;
      }

      if (index > -1 && selectedIds.length === 0) {
        store.ui.commandMenu.setType('OpportunityCommands');
        store.ui.commandMenu.setContext({
          entity: 'Opportunity',
          ids: [data?.[index]?.id],
        });
      }
    }
  };

  useEffect(() => {
    return () => {
      store.ui.setContactPreviewCardOpen(false);
    };
  }, [preset]);

  useKeyBindings(
    {
      Escape: () => {
        store.ui.setContactPreviewCardOpen(false);
      },
      Space: (e) => {
        e.preventDefault();
        store.ui.setContactPreviewCardOpen(false);
      },
    },
    {
      when: store.ui.contactPreviewCardOpen,
    },
  );

  const checkIfEmpty = () => {
    return match(tableType)
      .with(
        TableViewType.Organizations,
        () => store.organizations?.totalElements === 0,
      )
      .with(TableViewType.Contacts, () => store.contacts?.totalElements === 0)
      .with(TableViewType.Invoices, () => store.invoices?.totalElements === 0)
      .with(TableViewType.Contracts, () => store.contracts?.totalElements === 0)
      .with(TableViewType.Flow, () => store.flows?.totalElements === 0)
      .otherwise(() => false);
  };

  const checkIfLoading = () => {
    return match(tableType)
      .with(TableViewType.Organizations, () => store.organizations?.isLoading)
      .with(TableViewType.Contacts, () => store.contacts?.isLoading)
      .with(TableViewType.Invoices, () => store.invoices?.isLoading)
      .with(TableViewType.Contracts, () => store.contracts?.isLoading)
      .with(TableViewType.Flow, () => store.flows?.isLoading)
      .otherwise(() => false);
  };

  if (checkIfEmpty() && checkIfLoading()) {
    return <EmptyState />;
  }

  return (
    <div className='flex'>
      {/* eslint-disable-next-line @typescript-eslint/no-explicit-any */}
      <Table<any>
        data={data}
        manualFiltering
        sorting={sorting}
        columns={columns}
        tableRef={tableRef}
        getRowId={(row) => row.id}
        enableColumnResizing={true}
        onSortingChange={setSorting}
        onResizeColumn={handleColumnSizing}
        onSelectionChange={onSelectionChange}
        onFocusedRowChange={handleSetFocused}
        tableId={tableViewDef?.value?.tableId}
        dataTest={`finder-table-${tableType}`}
        isLoading={store.organizations.isLoading}
        fullRowSelection={tableType === TableViewType.Invoices}
        totalItems={store.organizations.isLoading ? 40 : data.length}
        enableKeyboardShortcuts={
          !isEditing && !isFiltering && !isCommandMenuPrompted
        }
        enableTableActions={
          tableType &&
          [TableViewType.Invoices, TableViewType.Contracts].includes(tableType)
            ? false
            : enableFeature !== null
            ? enableFeature
            : true
        }
        enableRowSelection={
          tableType &&
          [TableViewType.Invoices, TableViewType.Contracts].includes(tableType)
            ? false
            : enableFeature !== null
            ? enableFeature
            : true
        }
        renderTableActions={(table, focusRow, selectedIds) => {
          if (tableType === TableViewType.Organizations) {
            return (
              <OrganizationTableActions
                table={table}
                selection={selectedIds}
                isCommandMenuOpen={isCommandMenuPrompted}
                focusedId={focusRow !== null ? data?.[focusRow]?.id : null}
                enableKeyboardShortcuts={
                  !isEditing &&
                  !isFiltering &&
                  !isSearching &&
                  !isCommandMenuPrompted
                }
              />
            );
          }

          if (tableType === TableViewType.Contacts) {
            return (
              <ContactTableActions
                table={table}
                selection={selectedIds}
                isCommandMenuOpen={isCommandMenuPrompted}
                focusedId={focusRow !== null ? data?.[focusRow]?.id : null}
                enableKeyboardShortcuts={
                  !isSearching &&
                  !isFiltering &&
                  !isEditing &&
                  !isCommandMenuPrompted
                }
              />
            );
          }

          if (tableType === TableViewType.Opportunities) {
            return (
              <OpportunitiesTableActions
                table={table}
                selection={selectedIds}
                focusedId={focusRow !== null ? data?.[focusRow]?.id : null}
                enableKeyboardShortcuts={
                  !isSearching &&
                  !isFiltering &&
                  !isEditing &&
                  !isCommandMenuPrompted
                }
              />
            );
          }

          if (tableType === TableViewType.Flow) {
            return (
              <FlowSequencesTableActions
                table={table}
                selection={selectedIds}
                focusedId={focusRow !== null ? data?.[focusRow]?.id : null}
                enableKeyboardShortcuts={
                  !isSearching &&
                  !isFiltering &&
                  !isEditing &&
                  !isCommandMenuPrompted
                }
              />
            );
          }

          return <></>;
        }}
      />
      {isSidePanelOpen && <SidePanel />}
      {store.ui.contactPreviewCardOpen && !store.ui.isSearching && (
        <ContactPreviewCard />
      )}
      {tableType === TableViewType.Invoices && (
        <ConfirmDeleteDialog
          onClose={reset}
          hideCloseButton
          isOpen={isConfirming}
          onConfirm={onConfirm}
          confirmButtonLabel='Void invoice'
          label={`Void invoice ${targetInvoiceNumber}`}
          description={`Voiding this invoice will send an email notification to ${targetInvoiceEmail}`}
        />
      )}
    </div>
  );
});
