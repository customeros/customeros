import { useState } from 'react';
import { useSearchParams } from 'react-router-dom';

import { match } from 'ts-pattern';
import { observer } from 'mobx-react-lite';

import { Input } from '@ui/form/Input';
import { Button } from '@ui/form/Button/Button';
import { IconButton } from '@ui/form/IconButton';
import { Filter } from '@ui/presentation/Filter';
import { useStore } from '@shared/hooks/useStore';
import { TableViewType } from '@shared/types/tableDef';
import { FilterLines } from '@ui/media/icons/FilterLines';
import { Menu, MenuItem, MenuList, MenuButton } from '@ui/overlay/Menu/Menu';
import { useTableColumnOptionsMap } from '@organizations/hooks/useTableColumnOptionsMap';
import {
  FilterItem,
  TableIdType,
  ColumnViewType,
  ComparisonOperator,
} from '@shared/types/__generated__/graphql.types';

import { filterTypes } from '../Columns/organizations/filtersType';

export const Filters = observer(
  ({ tableId, type }: { type: TableViewType; tableId: TableIdType }) => {
    const store = useStore();
    const [search, setSearch] = useState<string>('');
    const [searchParams] = useSearchParams();
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    const [optionsMap, helperTextMap] = useTableColumnOptionsMap(type as any);

    const preset = match(tableId)
      .with(
        TableIdType.Opportunities,
        () => store.tableViewDefs.opportunitiesPreset,
      )
      .otherwise(() => searchParams?.get('preset'));

    const tableViewDef = store.tableViewDefs.getById(preset ?? '0');

    const columns =
      tableViewDef?.value?.columns
        .filter(
          (c) =>
            ![
              ColumnViewType.FlowSequenceContactCount,
              ColumnViewType.FlowName,
              ColumnViewType.ContactsFlows,
            ].includes(c.columnType),
        )
        .map((c) => ({
          ...c,
          label: optionsMap[c.columnType],
          helperText: helperTextMap[c.columnType],
        })) ?? [];

    const availableFilters = columns
      .map((column) => {
        const filterType = filterTypes[column.columnType];

        if (filterType) {
          return {
            ...filterType,
            columnType: column.columnType,
          };
        }

        return null;
      })
      .filter(Boolean)
      .filter((f) =>
        f?.filterName.toLowerCase().includes(search.toLowerCase()),
      );

    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    const filters = tableViewDef?.getFilters()?.AND as any | undefined;

    const flattenedFilters: FilterItem[] =
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      filters?.map((f: any[]) => ({ ...f.filter })) ?? [];

    const handleFilterName = (property: string) => {
      const filterType = Object.values(filterTypes).find(
        (type) => type.filterAccesor === property,
      );

      return filterType ? filterType.filterName : property;
    };

    const getFilterOperators = (property: string) => {
      const filterType = Object.values(filterTypes).find(
        (type) => type.filterAccesor === property,
      );

      return filterType?.filterOperators ?? [];
    };

    const getFilterTypes = (property: string) => {
      const filterType = Object.values(filterTypes).find(
        (type) => type.filterAccesor === property,
      );

      return filterType?.filterType;
    };

    return (
      <div className='flex gap-2'>
        {flattenedFilters.map((f) => (
          <Filter
            key={f.property}
            filterValue={f.value}
            filterName={handleFilterName(f.property)}
            operators={getFilterOperators(f.property)}
            filterType={getFilterTypes(f.property) || ''}
            operatorValue={f.operation || ComparisonOperator.Between}
            onClearFilter={() => {
              tableViewDef?.removeFilter(f.property);
            }}
            onChangeOperator={(operation: string) => {
              tableViewDef?.setFilter({
                ...f,
                operation: (operation as ComparisonOperator) || '',
                property: f.property,
              });
            }}
            onChangeFilterValue={(value: string) => {
              if (value.length === 0) {
                tableViewDef?.setFilter({
                  ...f,
                  property: f.property,
                  active: false,
                  operation: f.operation,
                  value: null,
                });
              } else {
                tableViewDef?.setFilter({
                  ...f,
                  value: value,
                  property: f.property,
                  active: true,
                });
              }
            }}
          />
        ))}
        <Menu>
          <MenuButton>
            {flattenedFilters.length ? (
              <IconButton
                size='xs'
                variant='outline'
                aria-label='filters'
                icon={<FilterLines />}
                colorScheme='grayModern'
                className='border-transparent'
              />
            ) : (
              <Button
                size='xs'
                colorScheme='grayModern'
                leftIcon={<FilterLines />}
              >
                Filters
              </Button>
            )}
          </MenuButton>
          <MenuList>
            <Input
              size='sm'
              value={search}
              variant='unstyled'
              className='px-2.5'
              placeholder='Filter by'
              onChange={(e) => setSearch(e.target.value)}
            />
            {availableFilters.map((filter) => {
              return (
                <>
                  <MenuItem
                    key={filter?.columnType}
                    onClick={() =>
                      tableViewDef?.appendFilter({
                        property: filter?.filterAccesor || '',
                        value: undefined,
                        active: false,
                        operation:
                          getFilterOperators(filter?.filterAccesor ?? '')[0] ||
                          '',
                      })
                    }
                  >
                    <div className='flex items-center justify-center gap-2'>
                      {filter?.icon}
                      {filter?.filterName}
                    </div>
                  </MenuItem>
                </>
              );
            })}
          </MenuList>
        </Menu>
      </div>
    );
  },
);
