import { useSearchParams } from 'react-router-dom';
import { useRef, useState, useEffect } from 'react';

import { useKey } from 'rooks';
import { Store } from '@store/store';
import { inPlaceSort } from 'fast-sort';
import { observer } from 'mobx-react-lite';
import { useFeatureIsOn } from '@growthbook/growthbook-react';

import { Organization } from '@graphql/types';
import { useStore } from '@shared/hooks/useStore';
import { Table, SortingState, TableInstance } from '@ui/presentation/Table';

import { TableActions } from '../Actions';
import { EmptyState } from '../EmptyState/EmptyState';
import {
  getColumnSortFn,
  getColumnsConfig,
  getPredefinedFilterFn,
} from '../Columns/columnsDictionary';

export const OrganizationsTable = observer(() => {
  const store = useStore();
  const [searchParams] = useSearchParams();
  const enableFeature = useFeatureIsOn('gp-dedicated-1');
  const tableRef = useRef<TableInstance<Store<Organization>> | null>(null);
  const [sorting, setSorting] = useState<SortingState>([
    { id: 'ORGANIZATIONS_LAST_TOUCHPOINT', desc: true },
  ]);
  const [isShiftPressed, setIsShiftPressed] = useState(false);

  const preset = searchParams?.get('preset');
  const searchTerm = searchParams?.get('search');
  const tableViewDef = store.tableViewDefs.getById(preset ?? '1');
  const tableColumns = getColumnsConfig(tableViewDef?.value);

  const data = store.organizations.toComputedArray((arr) => {
    const predefinedFilter = getPredefinedFilterFn(tableViewDef?.getFilters());
    if (predefinedFilter) {
      arr = arr.filter(predefinedFilter);
    }
    if (searchTerm) {
      arr = arr.filter((org) =>
        org.value.name.toLowerCase().includes(searchTerm.toLowerCase()),
      );
    }
    const columnId = sorting[0]?.id;
    const isDesc = sorting[0]?.desc;
    const computed = inPlaceSort(arr)?.[isDesc ? 'desc' : 'asc'](
      getColumnSortFn(columnId),
    );

    return computed;
  });

  if (
    store.organizations.totalElements === 0 &&
    !store.organizations.isLoading
  ) {
    return <EmptyState />;
  }

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

  const isCurrentlySearching = store.ui.isSearching === 'organizations';

  return (
    <Table<Store<Organization>>
      data={data}
      columns={tableColumns}
      sorting={sorting}
      tableRef={tableRef}
      enableTableActions={enableFeature !== null ? enableFeature : true}
      enableRowSelection={enableFeature !== null ? enableFeature : true}
      onSortingChange={setSorting}
      getRowId={(row) => row.value?.metadata?.id}
      isLoading={store.organizations.isLoading}
      totalItems={store.organizations.isLoading ? 40 : data.length}
      renderTableActions={(table) => (
        <TableActions
          table={table}
          onHide={store.organizations.hide}
          onMerge={store.organizations.merge}
          tableId={tableViewDef?.value.tableId}
          onUpdateStage={store.organizations.updateStage}
          isCurrentlySearching={isCurrentlySearching}
        />
      )}
      onSelectionChange={(selection) => {
        if (!isShiftPressed) return;

        const selectedIds = Object.keys(selection);
        const indexes = selectedIds.map((id) =>
          data.findIndex((d) => d.value.metadata.id === id),
        );

        const edgeIndexes = [Math.min(...indexes), Math.max(...indexes)];
        const targetIds = data
          .slice(edgeIndexes[0], edgeIndexes[1] + 1)
          .map((d) => d.value.metadata.id);

        tableRef.current?.setRowSelection((prev) => ({
          ...prev,
          ...targetIds.reduce((acc, id) => ({ ...acc, [id]: true }), {}),
        }));
      }}
    />
  );
});
