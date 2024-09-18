import { useSearchParams } from 'react-router-dom';
import { useRef, useState, useEffect } from 'react';

import { useKeyBindings } from 'rooks';
import { inPlaceSort } from 'fast-sort';
import { observer } from 'mobx-react-lite';
import { ColumnDef } from '@tanstack/react-table';
import { FlowStore } from '@store/Flows/Flow.store.ts';
import { InvoiceStore } from '@store/Invoices/Invoice.store';
import { ContactStore } from '@store/Contacts/Contact.store';
import { useFeatureIsOn } from '@growthbook/growthbook-react';
import { ContractStore } from '@store/Contracts/Contract.store';
import { useTableActions } from '@invoices/hooks/useTableActions';
import { OpportunityStore } from '@store/Opportunities/Opportunity.store';
import { OrganizationStore } from '@store/Organizations/Organization.store';

import { useStore } from '@shared/hooks/useStore';
import { Invoice, WorkflowType, TableViewType } from '@graphql/types';
import { useColumnSizing } from '@organizations/hooks/useColumnSizing';
import { Table, SortingState, TableInstance } from '@ui/presentation/Table';
import { ConfirmDeleteDialog } from '@ui/overlay/AlertDialog/ConfirmDeleteDialog';
import { getFlowsFilterFns } from '@organizations/components/Columns/flows/filterFns.ts';
import { getFlowsColumnSortFn } from '@organizations/components/Columns/flows/sortFns.ts';
import { getFlowColumnsConfig } from '@organizations/components/Columns/flows/columns.tsx';
import { getOpportunitiesSortFn } from '@organizations/components/Columns/opportunities/sortFns';
import { OpportunitiesTableActions } from '@organizations/components/Actions/OpportunityActions';
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
import { ContactPreviewCard } from '../ContactPreviewCard/ContactPreviewCard';
import {
  getContactSortFn,
  getContactFilterFns,
  getContactColumnsConfig,
} from '../Columns/contacts';
import {
  ContactTableActions,
  OrganizationTableActions,
  FlowSequencesTableActions,
} from '../Actions';
import {
  getOrganizationSortFn,
  getOrganizationFilterFns,
  getOrganizationColumnsConfig,
} from '../Columns/organizations';

export type FinderTableEntityTypes =
  | OrganizationStore
  | ContactStore
  | InvoiceStore
  | ContractStore
  | OpportunityStore
  | FlowStore;

interface FinderTableProps {
  isSidePanelOpen: boolean;
}

export const FinderTable = observer(({ isSidePanelOpen }: FinderTableProps) => {
  const store = useStore();
  const [searchParams] = useSearchParams();
  const enableFeature = useFeatureIsOn('gp-dedicated-1');
  const tableRef = useRef<TableInstance<FinderTableEntityTypes> | null>(null);
  const [sorting, setSorting] = useState<SortingState>([
    { id: 'ORGANIZATIONS_LAST_TOUCHPOINT', desc: true },
  ]);

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
  const flowSequenceColumns = getFlowColumnsConfig(tableViewDef?.value);

  const tableColumns = (
    tableType === TableViewType.Organizations
      ? organizationColumns
      : tableType === TableViewType.Contacts
      ? contactColumns
      : tableType === TableViewType.Contracts
      ? contractColumns
      : tableType === TableViewType.Opportunities
      ? opportunityColumns
      : tableType === TableViewType.Flow
      ? flowSequenceColumns
      : invoiceColumns
  ) as ColumnDef<FinderTableEntityTypes>[];
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

  const flowsData = store.flows.toComputedArray((arr) => {
    if (tableType !== TableViewType.Flow) return arr;

    const filters = getFlowsFilterFns(tableViewDef?.getFilters());

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
        getFlowsColumnSortFn(columnId),
      );

      return computed;
    }

    return arr.filter((e) => e.value.status !== 'ARCHIVED');
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
    [TableViewType.Flow]: flowsData,
  };

  const data = dataMap[tableType];

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

  const targetInvoice: Invoice = data?.find(
    (i) => i.value.metadata?.id === targetId,
  )?.value as Invoice;
  const targetInvoiceNumber = targetInvoice?.invoiceNumber || '';
  const targetInvoiceEmail = targetInvoice?.customer?.email || '';

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
  }, []);
  useKeyBindings(
    {
      Escape: () => {
        store.ui.setContactPreviewCardOpen(false);
      },
      Space: () => {
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
      !store.invoices.isLoading) ||
    (tableViewDef?.value.tableType === TableViewType.Flow &&
      store.flows?.toArray().length === 0 &&
      !store.flows.isLoading) ||
    (tableViewDef?.value.tableType === TableViewType.Contracts &&
      store.contracts?.toArray().length === 0 &&
      !store.contracts.isLoading)
  ) {
    return <EmptyState />;
  }

  return (
    <div className='flex'>
      <Table<FinderTableEntityTypes>
        data={data}
        manualFiltering
        sorting={sorting}
        tableRef={tableRef}
        columns={tableColumns}
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
                selection={selectedIds}
                isCommandMenuOpen={isCommandMenuPrompted}
                table={table as TableInstance<OrganizationStore>}
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
                selection={selectedIds}
                isCommandMenuOpen={isCommandMenuPrompted}
                table={table as TableInstance<ContactStore>}
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
                selection={selectedIds}
                table={table as TableInstance<ContactStore>}
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
                selection={selectedIds}
                table={table as TableInstance<ContactStore>}
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
      {store.ui.contactPreviewCardOpen && <ContactPreviewCard />}
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
