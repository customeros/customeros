import { useParams, useLocation, useSearchParams } from 'react-router-dom';
import {
  useRef,
  useMemo,
  useEffect,
  type Dispatch,
  type SetStateAction,
} from 'react';

import { match } from 'ts-pattern';
import { useKeyBindings } from 'rooks';
import { observer } from 'mobx-react-lite';
import { ColumnSort } from '@tanstack/table-core';
import { useFeatureIsOn } from '@growthbook/growthbook-react';
import { useColumnSizing } from '@finder/hooks/useColumnSizing';
import { useTableActions } from '@invoices/hooks/useTableActions';
import { OpportunitiesTableActions } from '@finder/components/Actions/OpportunityActions';

import { useStore } from '@shared/hooks/useStore';
import { Table, SortingState, TableInstance } from '@ui/presentation/Table';
import { ConfirmDeleteDialog } from '@ui/overlay/AlertDialog/ConfirmDeleteDialog';
import {
  Invoice,
  TableIdType,
  TableViewType,
  ColumnViewType,
} from '@graphql/types';

import { EmptyState } from '../EmptyState/EmptyState';
import { computeFinderData } from './computeFinderData';
import { computeFinderColumns } from './computeFinderColumns';
import {
  ContactTableActions,
  OrganizationTableActions,
  FlowSequencesTableActions,
} from '../Actions';

export const FinderTable = observer(() => {
  const store = useStore();
  const params = useParams();
  const [searchParams] = useSearchParams();
  const location = useLocation();

  const enableFeature = useFeatureIsOn('gp-dedicated-1');
  const tableRef = useRef<TableInstance<object> | null>(null);
  const preset = searchParams?.get('preset');
  const tableViewDef = store.tableViewDefs.getById(preset ?? '1');
  const contactsPreset = store.tableViewDefs.contactsPreset;

  const sortingData = tableViewDef?.getSorting();
  const defaultSorting =
    preset === contactsPreset
      ? [{ id: ColumnViewType.ContactsCreatedAt, desc: true }]
      : [{ id: ColumnViewType.OrganizationsLastTouchpoint, desc: true }];

  const sorting: ColumnSort[] = !sortingData?.id
    ? defaultSorting
    : [
        {
          id: sortingData.id,
          desc: sortingData.desc,
        },
      ];

  const searchTerm = searchParams?.get('search');
  const { reset, targetId, isConfirming, onConfirm } = useTableActions();

  const tableType =
    tableViewDef?.value?.tableType || TableViewType.Organizations;
  const tableId = tableViewDef?.value?.tableId || TableIdType.Organizations;

  const columns = computeFinderColumns(store, {
    tableType,
    currentPreset: preset,
  });

  const handleSortChange: Dispatch<SetStateAction<SortingState>> = (
    updaterOrValue,
  ) => {
    const next =
      typeof updaterOrValue === 'function'
        ? updaterOrValue(sorting)
        : updaterOrValue;

    tableViewDef?.setSorting(next[0]?.id, next[0]?.desc);
  };

  const data = computeFinderData(store, {
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
            meta: {
              tableId: tableId,
              id: params.id,
            },
          });
        }

        if (selectedIds.length > 1) {
          store.ui.commandMenu.setType('ContactBulkCommands');
          store.ui.commandMenu.setContext({
            entity: 'Contact',
            ids: selectedIds,
            meta: {
              tableId: tableId,
              id: params.id,
            },
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

      if (location?.state?.fromOnboarding) {
        store.ui.commandMenu.setType('CreateNewFlow');
        store.ui.commandMenu.setOpen(true);
      }
    } else {
      store.ui.commandMenu.setType('ContactHub');
    }
  }, [tableViewDef?.value.id]);

  useEffect(() => {
    if (location?.state?.fromOnboarding && tableType === TableViewType.Flow) {
      setTimeout(() => {
        store.ui.commandMenu.setType('CreateNewFlow');
        store.ui.commandMenu.setOpen(true);
      }, 100);
    }

    if (
      location?.state?.fromOnboarding &&
      tableType === TableViewType.Contacts
    ) {
      setTimeout(() => {
        store.ui.commandMenu.setType('AddContactsBulk');
        store.ui.commandMenu.setOpen(true);
      }, 100);
    }
  }, [location.state?.fromOnboarding]);

  const totalItems = match(tableType)
    .returnType<number>()
    .with(
      TableViewType.Organizations,
      () => store.organizations?.availableCounts.get(preset ?? '') ?? 0,
    )
    .otherwise(() => data.length ?? 50);

  useEffect(() => {
    store.ui.setSearchCount(totalItems);
    store.ui.setFilteredTable(data);
  }, [data.length, store.organizations?.totalElements]);

  const isEditing = store.ui.isEditingTableCell;
  const isFiltering = store.ui.isFilteringTable;
  const isSearching =
    store.ui.isSearching === tableViewDef?.value?.tableType?.toLowerCase();

  const [targetInvoiceNumber, targetInvoiceEmail] = match(tableType)
    .with(TableViewType.Invoices, () => {
      const invoice = data?.find((i) => {
        if ('metadata' in i.value) {
          return i.value!.metadata.id === targetId;
        } else {
          return i.value.id === targetId;
        }
      })?.value as Invoice;

      const targetInvoiceNumber = invoice?.invoiceNumber || '';
      const targetInvoiceEmail = invoice?.customer?.email || '';

      return [targetInvoiceNumber, targetInvoiceEmail];
    })
    .otherwise(() => ['', '']);

  const handleSetFocused = (index: number | null, selectedIds: string[]) => {
    if (isCommandMenuPrompted) return;

    if (selectedIds.length > 0) return;

    if (tableType === TableViewType.Organizations) {
      if (!store.ui.showPreviewCard) {
        if (index !== null) {
          store.ui.setFocusRow(data?.[index]?.id);
        }
      }

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
      if (!store.ui.showPreviewCard) {
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
          meta: {
            tableId: tableId,
            id: params.id,
          },
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
      store.ui.setShowPreviewCard(false);
      store.ui.setShortcutsPanel(false);
    };
  }, [preset]);

  useKeyBindings(
    {
      Escape: () => {
        store.ui.setShowPreviewCard(false);
      },
    },
    {
      when: store.ui.showPreviewCard,
    },
  );

  const rowHeight = useMemo(() => {
    return match(tableType)
      .with(TableViewType.Organizations, () => 33)
      .with(TableViewType.Opportunities, () => 33)
      .with(TableViewType.Contacts, () => 33)
      .with(TableViewType.Invoices, () => 31)
      .with(TableViewType.Contracts, () => 29)
      .with(TableViewType.Flow, () => 29)
      .otherwise(() => 33);
  }, [tableType]);

  const checkIfEmpty = () => {
    return match(tableType)
      .with(TableViewType.Organizations, () =>
        preset ? store.organizations?.totalElements === 0 : true,
      )
      .with(TableViewType.Contacts, () => {
        return preset ? store.contacts.totalElements === 0 : true;
      })
      .with(TableViewType.Invoices, () => store.invoices?.totalElements === 0)
      .with(TableViewType.Contracts, () => store.contracts?.totalElements === 0)
      .with(TableViewType.Flow, () => store.flows?.totalElements === 0)
      .otherwise(() => false);
  };

  const enableRowSelection = match(tableType)
    .with(TableViewType.Organizations, () => enableFeature || true)
    .with(TableViewType.Contacts, () => true)
    .with(TableViewType.Opportunities, () => true)
    .with(TableViewType.Flow, () => true)
    .otherwise(() => false);

  if (checkIfEmpty()) {
    return <EmptyState />;
  }

  const canFetchMore = match(tableType)
    .with(
      TableViewType.Organizations,
      () => !!preset && store.organizations.canLoadNext(preset),
    )
    .with(
      TableViewType.Contacts,
      () => !!preset && store.contacts.canLoadNext(preset),
    )
    .otherwise(() => false);

  return (
    <div className='flex w-full'>
      {/* eslint-disable-next-line @typescript-eslint/no-explicit-any */}
      <Table<any>
        data={data}
        manualFiltering
        sorting={sorting}
        columns={columns}
        tableRef={tableRef}
        rowHeight={rowHeight}
        id={tableViewDef?.id}
        totalItems={totalItems}
        getRowId={(row) => row.id}
        enableColumnResizing={true}
        canFetchMore={canFetchMore}
        onSortingChange={handleSortChange}
        onResizeColumn={handleColumnSizing}
        onSelectionChange={onSelectionChange}
        onFocusedRowChange={handleSetFocused}
        dataTest={`finder-table-${tableType}`}
        enableRowSelection={enableRowSelection}
        tableType={tableViewDef?.value?.tableId}
        isLoading={store.organizations.isLoading}
        fullRowSelection={tableType === TableViewType.Invoices}
        enableKeyboardShortcuts={
          !isEditing && !isFiltering && !isCommandMenuPrompted
        }
        onFetchMore={() => {
          store.organizations.loadNext(preset!);
          store.contacts.loadNext(preset!);
        }}
        enableTableActions={
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
