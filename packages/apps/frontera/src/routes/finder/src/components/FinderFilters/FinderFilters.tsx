import { useSearchParams } from 'react-router-dom';

import { match } from 'ts-pattern';
import { observer } from 'mobx-react-lite';
import { useTableColumnOptionsMap } from '@finder/hooks/useTableColumnOptionsMap';

import { useStore } from '@shared/hooks/useStore';
import { Filters } from '@ui/presentation/Filters';
import { TableViewType } from '@shared/types/tableDef';
import {
  FilterItem,
  TableIdType,
  ColumnViewType,
} from '@shared/types/__generated__/graphql.types';

import { getFilterTypes as getFilterTypesForContacts } from '../Columns/contacts/filterTypes';
import { getFilterTypes as getFilterTypesForOrganizations } from '../Columns/organizations/filterTypes';

export const FinderFilters = observer(
  ({ tableId, type }: { type: TableViewType; tableId: TableIdType }) => {
    const store = useStore();
    const getFilterTypes = match(tableId)
      .with(TableIdType.Contacts, () => getFilterTypesForContacts)
      .with(TableIdType.Organizations, () => getFilterTypesForOrganizations)
      .otherwise(() => getFilterTypesForOrganizations);

    const [searchParams] = useSearchParams();
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    const [optionsMap, helperTextMap] = useTableColumnOptionsMap(type as any);

    const filterTypes = getFilterTypes(store);

    const preset = match(tableId)
      .with(
        TableIdType.Opportunities,
        () => store.tableViewDefs.opportunitiesPreset,
      )
      .otherwise(() => searchParams?.get('preset'));

    const tableViewDef = store.tableViewDefs.getById(preset ?? '0');

    const columns =
      tableViewDef?.value?.columns
        .filter((c) => ![ColumnViewType.FlowName].includes(c.columnType))
        .map((c) => ({
          ...c,
          label: optionsMap[c.columnType],
          helperText: helperTextMap[c.columnType],
        })) ?? [];

    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    const filters = tableViewDef?.getFilters()?.AND as any | undefined;

    if (tableId === TableIdType.Contacts) {
      columns.push({
        // @ts-expect-error fix later
        columnType: 'EMAIL_VERIFICATION_PRIMARY_EMAIL',
        label: 'Primary email status',
        helperText: 'Email Verification',
      });
    }

    const flattenedFilters: FilterItem[] =
      filters
        ?.map((f: FilterItem[]) => ({ ...f.filter }))
        .filter((filter: FilterItem) => {
          if (
            tableId === TableIdType.Customers &&
            filter.property === 'RELATIONSHIP'
          ) {
            return false;
          }

          if (
            (tableId === TableIdType.Targets &&
              filter.property === 'RELATIONSHIP') ||
            filter.property === 'STAGE'
          ) {
            return false;
          }

          return true;
        }) ?? [];

    return (
      <Filters
        columns={columns}
        filters={flattenedFilters}
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        filterTypes={filterTypes as any}
        onClearFilter={(filter, idx) =>
          tableViewDef?.removeFilter(filter.property, idx)
        }
        setFilters={(filter: FilterItem, index: number) => {
          tableViewDef?.setFilterv2(filter, index);
        }}
        onFilterSelect={(filter, getFilterOperators) => {
          tableViewDef?.appendFilter({
            property: filter?.filterAccesor || '',
            value: undefined,
            active: false,
            operation: getFilterOperators(filter?.filterAccesor ?? '')[0] || '',
          });
        }}
      />
    );
  },
);
