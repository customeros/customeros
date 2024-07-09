import { useSearchParams } from 'react-router-dom';
import { useRef, useState, useEffect } from 'react';

import { useKey } from 'rooks';
import { inPlaceSort } from 'fast-sort';
import { observer } from 'mobx-react-lite';
import difference from 'lodash/difference';
import intersection from 'lodash/intersection';
import { OnChangeFn } from '@tanstack/table-core';
import { useFeatureIsOn } from '@growthbook/growthbook-react';
import { ContactStore } from '@store/Contacts/Contact.store.ts';
import { OrganizationStore } from '@store/Organizations/Organization.store.ts';

import { useStore } from '@shared/hooks/useStore';
import { WorkflowType, TableViewType } from '@graphql/types';
import {
  Table,
  SortingState,
  TableInstance,
  RowSelectionState,
} from '@ui/presentation/Table';

import { SidePanel } from '../SidePanel';
import { EmptyState } from '../EmptyState/EmptyState';
import { MergedColumnDefs } from '../Columns/shared/util/types';
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
    OrganizationStore | ContactStore
  > | null>(null);
  const [sorting, setSorting] = useState<SortingState>([
    { id: 'ORGANIZATIONS_LAST_TOUCHPOINT', desc: true },
  ]);
  const [selection, setSelection] = useState<RowSelectionState>({});
  const [isShiftPressed, setIsShiftPressed] = useState(false);
  const [focusIndex, setFocusIndex] = useState<number | null>(null);
  const [selectedIndex, setSelectedIndex] = useState<number | null>(null);
  const searchTerm = searchParams?.get('search');

  const preset = searchParams?.get('preset');
  const tableViewDef = store.tableViewDefs.getById(preset ?? '1');
  const tableType = tableViewDef?.value?.tableType;
  const getWorkFlow = store.workFlows
    .toArray()
    .filter((wf) => wf.value.type === WorkflowType.IdealCustomerProfile);

  const getWorkFlowId = getWorkFlow.map((wf) => wf.value.id);
  // TODO take care of flow filters
  const workFlow = store.workFlows.getByType(getWorkFlowId[0]);
  const flowFiltersStatus = store.ui.isFilteringICP;

  const contactColumns = getContactColumnsConfig(tableViewDef?.value);
  const organizationColumns = getOrganizationColumnsConfig(tableViewDef?.value);
  const tableColumns = (
    tableType === TableViewType.Organizations
      ? organizationColumns
      : contactColumns
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

  const data =
    tableType === TableViewType.Organizations
      ? organizationsData
      : contactsData;

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
      !store.contacts.isLoading)
  ) {
    return <EmptyState />;
  }

  const isEditing = store.ui.isEditingTableCell;
  const isFiltering = store.ui.isFilteringTable;
  const isSearching =
    store.ui.isSearching === tableViewDef?.value?.tableType?.toLowerCase();

  const createSocial = (data: {
    socialUrl: string;
    organizationId: string;
    options?: {
      onSuccess?: (serverId: string) => void;
    };
  }) => {
    store.contacts.createWithSocial(data);
  };

  return (
    <div className='flex'>
      <Table<OrganizationStore | ContactStore>
        data={data}
        manualFiltering
        sorting={sorting}
        tableRef={tableRef}
        columns={tableColumns}
        enableTableActions={enableFeature !== null ? enableFeature : true}
        enableRowSelection={enableFeature !== null ? enableFeature : true}
        onSortingChange={setSorting}
        getRowId={(row) => row.id}
        isSidePanelOpen={isSidePanelOpen}
        isLoading={store.organizations.isLoading}
        totalItems={store.organizations.isLoading ? 40 : data.length}
        selection={selection}
        onFocusedRowChange={setFocusIndex}
        onSelectedIndexChange={setSelectedIndex}
        onSelectionChange={handleSelectionChange}
        enableKeyboardShortcuts={!isEditing && !isFiltering}
        renderTableActions={(table) =>
          tableType === TableViewType.Organizations ? (
            <OrganizationTableActions
              table={table as TableInstance<OrganizationStore>}
              onHide={store.organizations.hide}
              onMerge={store.organizations.merge}
              tableId={tableViewDef?.value.tableId}
              onUpdateStage={store.organizations.updateStage}
              onCreateContact={createSocial}
              enableKeyboardShortcuts={!isEditing && !isFiltering}
              focusedId={focusIndex ? data?.[focusIndex]?.id : null}
            />
          ) : (
            <ContactTableActions
              table={table as TableInstance<ContactStore>}
              enableKeyboardShortcuts={!isSearching || !isFiltering}
              onAddTags={store.contacts.updateTags}
              onHideContacts={store.contacts.archive}
            />
          )
        }
      />
      {isSidePanelOpen && <SidePanel />}
    </div>
  );
});
