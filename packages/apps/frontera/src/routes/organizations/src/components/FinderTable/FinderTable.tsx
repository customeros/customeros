import { useSearchParams } from 'react-router-dom';
import { useRef, useState, useEffect } from 'react';

import { inPlaceSort } from 'fast-sort';
import { observer } from 'mobx-react-lite';
import difference from 'lodash/difference';
import { useKey, useKeyBindings } from 'rooks';
import intersection from 'lodash/intersection';
import { OnChangeFn } from '@tanstack/table-core';
import { ColumnDef } from '@tanstack/react-table';
import { InvoiceStore } from '@store/Invoices/Invoice.store';
import { useFeatureIsOn } from '@growthbook/growthbook-react';
import { ContactStore } from '@store/Contacts/Contact.store.ts';
import { CommandMenuType } from '@store/UI/CommandMenu.store.ts';
import { useTableActions } from '@invoices/hooks/useTableActions';
import { ContractStore } from '@store/Contracts/Contract.store.ts';
import { OpportunityStore } from '@store/Opportunities/Opportunity.store.ts';
import { OrganizationStore } from '@store/Organizations/Organization.store.ts';

import { useStore } from '@shared/hooks/useStore';
import { Invoice, WorkflowType, TableViewType } from '@graphql/types';
import { useColumnSizing } from '@organizations/hooks/useColumnSizing.ts';
import { ConfirmDeleteDialog } from '@ui/overlay/AlertDialog/ConfirmDeleteDialog';
import { getOpportunitiesSortFn } from '@organizations/components/Columns/opportunities/sortFns';
import { OpportunitiesTableActions } from '@organizations/components/Actions/OpportunityActions';
import {
  Table,
  SortingState,
  TableInstance,
  RowSelectionState,
} from '@ui/presentation/Table';
import {
  getOpportunityFilterFns,
  getOpportunityColumnsConfig,
} from '@organizations/components/Columns/opportunities';
import {
  getContractSortFn,
  getContractFilterFns,
  getContractColumnsConfig,
} from '@organizations/components/Columns/contracts';

import { SidePanel } from '../SidePanel';
import { EmptyState } from '../EmptyState/EmptyState';
import { getColumnSortFn } from '../Columns/invoices/sortFns';
import { getInvoiceFilterFns } from '../Columns/invoices/filterFns';
import { getInvoiceColumnsConfig } from '../Columns/invoices/columns';
import { getFlowFilterFns } from '../Columns/organizations/flowFilters';
import { ContactTableActions, OrganizationTableActions } from '../Actions';
import { ContactPreviewCard } from '../ContactPreviewCard/ContactPreviewCard';
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
    | OrganizationStore
    | ContactStore
    | InvoiceStore
    | ContractStore
    | OpportunityStore
  > | null>(null);
  const [sorting, setSorting] = useState<SortingState>([
    { id: 'ORGANIZATIONS_LAST_TOUCHPOINT', desc: true },
  ]);
  const [selection, setSelection] = useState<RowSelectionState>({});
  const [isShiftPressed, setIsShiftPressed] = useState(false);
  const [focusIndex, setFocusIndex] = useState<number | null>(null);
  const [selectedIndex, setSelectedIndex] = useState<number | null>(null);
  const [lastFocusedId, setLastFocusedId] = useState<string | null>(null);
  const searchTerm = searchParams?.get('search');
  const { reset, targetId, isConfirming, onConfirm } = useTableActions();

  const preset = searchParams?.get('preset');
  const tableViewDef = store.tableViewDefs.getById(preset ?? '1');

  const tableType =
    tableViewDef?.value?.tableType || TableViewType.Organizations;
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
  const opportunityColumns = getOpportunityColumnsConfig(tableViewDef?.value);

  const tableColumns = (
    tableType === TableViewType.Organizations
      ? organizationColumns
      : tableType === TableViewType.Contacts
      ? contactColumns
      : tableType === TableViewType.Contracts
      ? contractColumns
      : tableType === TableViewType.Opportunities
      ? opportunityColumns
      : invoiceColumns
  ) as ColumnDef<
    | OrganizationStore
    | ContactStore
    | InvoiceStore
    | ContractStore
    | OpportunityStore
  >[];
  const isCommandMenuPrompted = store.ui.commandMenu.isOpen;

  const removeAccents = (str: string) => {
    return str
      .toLowerCase()
      .normalize('NFD')
      .replace(/[\u0300-\u036f]/g, '');
  };

  const handleColumnSizing = useColumnSizing(tableColumns, tableViewDef);

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
        const org = entity.value?.organizations.content?.[0]?.name || '';
        const email = entity.value?.emails?.[0]?.email || '';

        return (
          removeAccents(name).includes(normalizedSearchTerm) ||
          removeAccents(org).includes(normalizedSearchTerm) ||
          removeAccents(email).includes(normalizedSearchTerm)
        );
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

    const filters = getContractFilterFns(tableViewDef?.getFilters());

    if (filters) {
      arr = arr.filter((v) => filters.every((fn) => fn(v)));
    }

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

  const opportunityData = store.opportunities.toComputedArray((arr) => {
    if (tableType !== TableViewType.Opportunities) return arr;
    arr = arr.filter((opp) => opp.value.internalType === 'NBO');

    const filters = getOpportunityFilterFns(tableViewDef?.getFilters());

    if (filters) {
      arr = arr.filter((v) => filters.every((fn) => fn(v)));
    }

    if (searchTerm) {
      const normalizedSearchTerm = removeAccents(searchTerm);

      arr = arr.filter((entity) => {
        const name = entity.value?.name || '';
        const org = entity.organization?.value.name || '';
        const email = entity.owner?.name || '';

        return (
          removeAccents(name).includes(normalizedSearchTerm) ||
          removeAccents(org).includes(normalizedSearchTerm) ||
          removeAccents(email).includes(normalizedSearchTerm)
        );
      });
    }

    if (tableType) {
      const columnId = sorting[0]?.id;
      const isDesc = sorting[0]?.desc;

      const computed = inPlaceSort(arr)?.[isDesc ? 'desc' : 'asc'](
        getOpportunitiesSortFn(columnId),
      );

      return computed;
    }

    return arr;
  });

  const dataMap = {
    [TableViewType.Organizations]: organizationsData,
    [TableViewType.Contacts]: contactsData,
    [TableViewType.Contracts]: contractsData,
    [TableViewType.Opportunities]: opportunityData,
    [TableViewType.Invoices]: invoicesData,
  };

  const data = dataMap[tableType];

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
      store.ui.commandMenu.setCallback((id?: string) => {
        tableRef?.current?.resetRowSelection();

        if (id) {
          setSelection(() => ({ [id]: true }));
        }
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
          reset();

          store.ui.commandMenu.setType('OpportunityCommands');
          store.ui.commandMenu.setContext({
            entity: 'Opportunity',
            ids: selectedIds,
          });
        }

        if (selectedIds.length > 1) {
          reset();

          store.ui.commandMenu.setType('OpportunityBulkCommands');
          store.ui.commandMenu.setContext({
            entity: 'Opportunities',
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
    } else if (tableType === TableViewType.Opportunities) {
      store.ui.commandMenu.setType('OpportunityHub');
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

    if (selectedIds.length > 0) return;

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

    if (tableType === TableViewType.Opportunities) {
      if (typeof index !== 'number') {
        store.ui.commandMenu.setType('OpportunityHub');

        return;
      }

      if (index > -1 && Object.keys(selection).length === 0) {
        store.ui.commandMenu.setType('OpportunityCommands');
        store.ui.commandMenu.setContext({
          entity: 'Opportunity',
          ids: [data?.[index]?.id],
        });
      }
    }
  };

  const handleOpenCommandKMenu = () => {
    const selectedIds = Object.keys(selection);
    const reset = () =>
      store.ui.commandMenu.setCallback((id?: string) => {
        tableRef?.current?.resetRowSelection();

        if (id) {
          setSelection(() => ({ [id]: true }));
        }
      });

    if (tableType === TableViewType.Organizations) {
      if (selectedIds.length === 1) {
        reset();
        store.ui.commandMenu.setType('OrganizationCommands');
        store.ui.commandMenu.setContext({
          entity: 'Organization',
          ids: selectedIds,
        });
      }

      if (selectedIds.length > 1) {
        reset();

        store.ui.commandMenu.setType('OrganizationBulkCommands');
        store.ui.commandMenu.setContext({
          entity: 'Organizations',
          ids: selectedIds,
        });
      }
    } else if (tableType === TableViewType.Opportunities) {
      if (selectedIds.length === 1) {
        reset();

        store.ui.commandMenu.setType('OpportunityCommands');
        store.ui.commandMenu.setContext({
          entity: 'Opportunity',
          ids: selectedIds,
        });
      }

      if (selectedIds.length > 1) {
        reset();

        store.ui.commandMenu.setType('OpportunityBulkCommands');
        store.ui.commandMenu.setContext({
          entity: 'Opportunities',
          ids: selectedIds,
        });
      }
    } else {
      if (selectedIds.length === 1) {
        reset();

        store.ui.commandMenu.setType('ContactCommands');
        store.ui.commandMenu.setContext({
          entity: 'Contact',
          ids: selectedIds,
        });
      }

      if (selectedIds.length > 1) {
        reset();

        store.ui.commandMenu.setType('ContactBulkCommands');
        store.ui.commandMenu.setContext({
          entity: 'Contact',
          ids: selectedIds,
        });
      }
    }

    store.ui.commandMenu.setOpen(true);
  };

  useEffect(() => {
    if (focusedId && !store.ui.contactPreviewCardOpen) {
      setLastFocusedId(focusedId);
    }
  }, [focusedId]);

  useEffect(() => {
    return () => {
      store.ui.setContactPreviewCardOpen(false);
    };
  }, []);

  useKeyBindings(
    {
      Escape: () => {
        store.ui.setContactPreviewCardOpen(false);
      },
    },
    {
      when: store.ui.contactPreviewCardOpen,
    },
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

  return (
    <div className='flex'>
      <Table<
        | OrganizationStore
        | ContactStore
        | InvoiceStore
        | ContractStore
        | OpportunityStore
      >
        data={data}
        manualFiltering
        sorting={sorting}
        tableRef={tableRef}
        selection={selection}
        columns={tableColumns}
        getRowId={(row) => row.id}
        enableColumnResizing={true}
        onSortingChange={setSorting}
        onResizeColumn={handleColumnSizing}
        onFocusedRowChange={handleSetFocused}
        dataTest={`finder-table-${tableType}`}
        onSelectedIndexChange={setSelectedIndex}
        isLoading={store.organizations.isLoading}
        onSelectionChange={handleSelectionChange}
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
        renderTableActions={(table) => {
          if (tableType === TableViewType.Organizations) {
            return (
              <OrganizationTableActions
                focusedId={focusedId}
                onCreateContact={createSocial}
                onOpenCommandK={handleOpenCommandKMenu}
                isCommandMenuOpen={isCommandMenuPrompted}
                onUpdateStage={store.organizations.updateStage}
                table={table as TableInstance<OrganizationStore>}
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
                handleOpen={(type: CommandMenuType, context) => {
                  handleOpenCommandKMenu();
                  store.ui.commandMenu.setType(type);

                  if (context) {
                    store.ui.commandMenu.setContext({
                      ...store.ui.commandMenu.context,
                      ...context,
                    });
                  }
                }}
              />
            );
          }

          if (tableType === TableViewType.Contacts) {
            return (
              <ContactTableActions
                focusedId={focusedId}
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
            );
          }

          if (tableType === TableViewType.Opportunities) {
            return (
              <OpportunitiesTableActions
                focusedId={focusedId}
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
              />
            );
          }

          return null;
        }}
      />
      {isSidePanelOpen && <SidePanel />}
      {store.ui.contactPreviewCardOpen && (
        <ContactPreviewCard contactId={lastFocusedId || ''} />
      )}
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
