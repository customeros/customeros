import { useRef, useState, useEffect } from 'react';

import { FilterType } from '@finder/components/Columns/organizations/filtersType';

import { Input } from '@ui/form/Input';
import { Button } from '@ui/form/Button/Button';
import { IconButton } from '@ui/form/IconButton';
import { FilterLines } from '@ui/media/icons/FilterLines';
import { Menu, MenuList, MenuItem, MenuButton } from '@ui/overlay/Menu/Menu';
import {
  FilterItem,
  ColumnView,
  ColumnViewType,
  ComparisonOperator,
} from '@shared/types/__generated__/graphql.types';

import { Filter } from '../Filter/Filter';

interface FiltersProps {
  filters: FilterItem[];
  columns: ColumnView[];
  onClearFilter: (filter: FilterItem, idx: number) => void;
  filterTypes: Partial<Record<ColumnViewType, FilterType>>;
  setFilters: (
    filter: FilterItem & { active?: boolean },
    index: number,
  ) => void;
  onFilterSelect: (
    filter: Partial<ColumnView & FilterType>,
    getFilterOperators: (prop: string) => ComparisonOperator[],
  ) => void;
}

export const Filters = ({
  filters,
  filterTypes,
  onClearFilter,
  columns,
  setFilters,
  onFilterSelect,
}: FiltersProps) => {
  const [isOpen, setIsOpen] = useState(false);
  const [search, setSearch] = useState<string>('');
  const inputRef = useRef<HTMLInputElement>(null);

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

  const getFilterOptions = (property: string) => {
    const filterType = Object.values(filterTypes).find(
      (type) => type.filterAccesor === property,
    );

    return filterType?.options;
  };

  useEffect(() => {
    setTimeout(() => {
      if (isOpen) {
        inputRef.current?.focus();
      }
    }, 100);
  }, [isOpen]);

  const handleChangeOperator = (
    operation: string,
    filter: FilterItem,
    index: number,
  ) => {
    const selectedOperation =
      operation === ComparisonOperator.IsEmpty ||
      operation === ComparisonOperator.IsNotEmpty ||
      filter.value
        ? true
        : false;

    setFilters(
      {
        ...filter,
        operation: (operation as ComparisonOperator) || '',
        property: filter.property,
        active: selectedOperation,
        includeEmpty: operation === ComparisonOperator.IsEmpty ? true : false,
      },
      index,
    );

    if (ComparisonOperator.Lt === operation) {
      setFilters(
        {
          ...filter,
          value: [null, filter.value[0]],
          property: filter.property,
          operation: (operation as ComparisonOperator) || '',
        },
        index,
      );
    } else {
      if (ComparisonOperator.Gt === operation) {
        setFilters(
          {
            ...filter,
            value: [filter.value[1], null],
            property: filter.property,
            operation: (operation as ComparisonOperator) || '',
          },
          index,
        );
      }
    }
  };

  const handleChangeFilterValue = (
    value: string | Date | string[],
    filter: FilterItem,
    index: number,
  ) => {
    if (Array.isArray(value) && value.length === 0) {
      setFilters(
        {
          ...filter,
          property: filter.property,
          active: false,
          operation: filter.operation,
          value: value,
        },
        index,
      );
    } else {
      setFilters(
        {
          ...filter,
          value: value,
          property: filter.property,
          active: true,
        },
        index,
      );
    }
  };

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
      f?.filterName.toLowerCase().includes(search || ''.toLowerCase()),
    );

  return (
    <div className='flex gap-2 flex-wrap'>
      {filters.map((f, idx) => (
        <Filter
          filterValue={f.value}
          key={`${f.property}-${idx}`}
          filterName={handleFilterName(f.property)}
          operators={getFilterOperators(f.property)}
          onClearFilter={() => onClearFilter(f, idx)}
          filterType={getFilterTypes(f.property) || ''}
          listFilterOptions={getFilterOptions(f.property) || []}
          operatorValue={f.operation || ComparisonOperator.Between}
          onChangeFilterValue={(value) =>
            handleChangeFilterValue(value, f, idx)
          }
          onChangeOperator={(operator) =>
            handleChangeOperator(operator, f, idx)
          }
        />
      ))}
      <Menu open={isOpen} onOpenChange={(v) => setIsOpen(v)}>
        <MenuButton asChild>
          {filters.length ? (
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
              variant='ghost'
              colorScheme='grayModern'
              leftIcon={<FilterLines />}
            >
              Filters
            </Button>
          )}
        </MenuButton>
        <MenuList align='start' side='bottom'>
          <Input
            size='sm'
            ref={inputRef}
            value={search}
            variant='unstyled'
            className='px-2.5'
            placeholder='Filter by'
            onChange={(e) => setSearch(e.target.value)}
          />
          {availableFilters.map((filter) => {
            if (filter === null) return null;

            return (
              <MenuItem
                className='group'
                key={filter.columnType}
                onClick={() =>
                  onFilterSelect(filter, (property) =>
                    getFilterOperators(filter.filterAccesor ?? property),
                  )
                }
              >
                <div className='flex items-center justify-center gap-2 '>
                  <span className='group-hover:text-gray-700 text-gray-500'>
                    {filter.icon}
                  </span>
                  {filter.filterName}
                </div>
              </MenuItem>
            );
          })}
        </MenuList>
      </Menu>
    </div>
  );
};
