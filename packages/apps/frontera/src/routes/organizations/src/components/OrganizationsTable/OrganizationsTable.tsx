import { useRef, useState } from 'react';
import { useSearchParams } from 'react-router-dom';

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

  const computedOrganizationIds = data.map((row) => row.value.metadata.id);

  if (
    store.organizations.totalElements === 0 &&
    !store.organizations.isLoading
  ) {
    return <EmptyState />;
  }

  return (
    <Table<Store<Organization>>
      data={data}
      columns={tableColumns}
      sorting={sorting}
      tableRef={tableRef}
      enableTableActions={enableFeature !== null ? enableFeature : true}
      enableRowSelection={enableFeature !== null ? enableFeature : true}
      onSortingChange={setSorting}
      isLoading={store.organizations.isLoading}
      totalItems={store.organizations.isLoading ? 40 : data.length}
      renderTableActions={(table) => (
        <TableActions
          table={table}
          onHide={store.organizations.hide}
          onMerge={store.organizations.merge}
          organizationIds={computedOrganizationIds}
          tableName={tableViewDef?.value?.name}
        />
      )}
    />
  );
});
