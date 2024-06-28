import { useSearchParams } from 'react-router-dom';
import { useRef, useMemo, useState, useEffect } from 'react';

import { useKey } from 'rooks';
import { Store } from '@store/store';
import { inPlaceSort } from 'fast-sort';
import { observer } from 'mobx-react-lite';
import difference from 'lodash/difference';
import intersection from 'lodash/intersection';
import { OnChangeFn } from '@tanstack/table-core';
import { useFeatureIsOn } from '@growthbook/growthbook-react';

import { useStore } from '@shared/hooks/useStore';
import { Contact, Organization, TableViewType } from '@graphql/types';
import {
  Table,
  SortingState,
  TableInstance,
  RowSelectionState,
} from '@ui/presentation/Table';
import {
  getContactFilterFn,
  getOrganizationFilterFn,
} from '@organizations/components/Columns/Dictionaries/SortAndFilterDictionary';

import { EmptyState } from '../EmptyState/EmptyState';
import { ContactTableActions, OrganizationTableActions } from '../Actions';
import {
  getAllFilterFns,
  getColumnSortFn,
  getColumnsConfig,
} from '../Columns/Dictionaries/columnsDictionary.tsx';

export const FinderTable = observer(() => {
  const store = useStore();
  const [searchParams] = useSearchParams();
  const enableFeature = useFeatureIsOn('gp-dedicated-1');
  const tableRef = useRef<TableInstance<Store<unknown>> | null>(null);
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
  const tableColumns = getColumnsConfig(tableViewDef?.value);

  const tableType = tableViewDef?.value?.tableType;

  const dataSet = useMemo(() => {
    if (tableType === TableViewType.Organizations) {
      return store.organizations;
    }
    if (tableType === TableViewType.Contacts) {
      return store.contacts;
    }

    return store.organizations;
  }, [tableType]);

  const filterFunction = useMemo(() => {
    if (tableType === TableViewType.Organizations) {
      return getOrganizationFilterFn;
    }
    if (tableType === TableViewType.Contacts) {
      return getContactFilterFn;
    }

    return getOrganizationFilterFn;
  }, [tableType]);

  // @ts-expect-error fixme
  const data = dataSet?.toComputedArray((arr) => {
    const filters = getAllFilterFns(tableViewDef?.getFilters(), filterFunction);
    if (filters) {
      // @ts-expect-error fixme

      arr = arr.filter((v) => filters.every((fn) => fn(v)));
    }

    if (searchTerm) {
      arr = arr.filter((entity) =>
        entity.value?.name
          ?.toLowerCase()
          .includes(searchTerm?.toLowerCase() as string),
      ) as Store<Contact>[] | Store<Organization>[];
    }
    if (tableType) {
      const columnId = sorting[0]?.id;
      const isDesc = sorting[0]?.desc;
      // @ts-expect-error fixme
      const computed = inPlaceSort(arr)?.[isDesc ? 'desc' : 'asc'](
        getColumnSortFn(columnId, tableType),
      );

      return computed;
    }

    return arr;
  });

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

  useKey(
    'Shift',
    (e) => {
      setIsShiftPressed(e.type === 'keydown');
    },
    { eventTypes: ['keydown', 'keyup'] },
  );
  if (dataSet.totalElements === 0 && !dataSet.isLoading) {
    return <EmptyState />;
  }

  const isEditing = store.ui.isEditingTableCell;
  const isFiltering = store.ui.isFilteringTable;
  const isSearching =
    store.ui.isSearching === tableViewDef?.value?.tableType?.toLowerCase();

  return (
    <Table<Store<unknown>>
      data={data as Store<Organization>[] | Store<Contact>[]}
      manualFiltering
      sorting={sorting}
      tableRef={tableRef}
      // @ts-expect-error fixme
      columns={tableColumns}
      enableTableActions={enableFeature !== null ? enableFeature : true}
      enableRowSelection={enableFeature !== null ? enableFeature : true}
      onSortingChange={setSorting}
      getRowId={(row) => row.id}
      isLoading={store.organizations.isLoading}
      totalItems={store.organizations.isLoading ? 40 : data.length}
      selection={selection}
      onFocusedRowChange={setFocusIndex}
      onSelectedIndexChange={setSelectedIndex}
      onSelectionChange={handleSelectionChange}
      enableKeyboardShortcuts={!isEditing || !isFiltering}
      renderTableActions={(table) =>
        tableType === TableViewType.Organizations ? (
          <OrganizationTableActions
            table={table as TableInstance<Store<Organization>>}
            onHide={store.organizations.hide}
            onMerge={store.organizations.merge}
            tableId={tableViewDef?.value.tableId}
            onUpdateStage={store.organizations.updateStage}
            enableKeyboardShortcuts={!isSearching || !isFiltering}
          />
        ) : (
          <ContactTableActions table={table as TableInstance<Store<Contact>>} />
        )
      }
    />
  );
});
