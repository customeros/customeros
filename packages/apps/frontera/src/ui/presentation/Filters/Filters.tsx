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
  filterSearch: string;
  filters: FilterItem[];
  onClearFilter: (filter: FilterItem) => void;
  handleFilterSearch: (value: string) => void;
  availableFilters: Partial<ColumnView & FilterType>[];
  filterTypes: Partial<Record<ColumnViewType, FilterType>>;
  onChangeOperator: (operator: string, filter: FilterItem) => void;
  onChangeFilterValue: (value: string | Date, filter: FilterItem) => void;
  onFilterSelect: (
    filter: Partial<ColumnView & FilterType>,
    getFilterOperators: (prop: string) => ComparisonOperator[],
  ) => void;
}

export const Filters = ({
  filters,
  filterTypes,
  onClearFilter,
  onChangeFilterValue,
  filterSearch,
  handleFilterSearch,
  availableFilters,
  onChangeOperator,
  onFilterSelect,
}: FiltersProps) => {
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
    <div className='flex gap-2 flex-wrap'>
      {filters.map((f) => (
        <Filter
          key={f.property}
          filterValue={f.value}
          onClearFilter={() => onClearFilter(f)}
          filterName={handleFilterName(f.property)}
          operators={getFilterOperators(f.property)}
          filterType={getFilterTypes(f.property) || ''}
          operatorValue={f.operation || ComparisonOperator.Between}
          onChangeOperator={(operator) => onChangeOperator(operator, f)}
          onChangeFilterValue={(value) => onChangeFilterValue(value, f)}
        />
      ))}
      <Menu>
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
            variant='unstyled'
            className='px-2.5'
            value={filterSearch}
            placeholder='Filter by'
            onChange={(e) => handleFilterSearch(e.target.value)}
          />
          {availableFilters.map((filter: Partial<ColumnView & FilterType>) => {
            return (
              <MenuItem
                key={filter?.columnType}
                onClick={() =>
                  onFilterSelect(filter, (property) =>
                    getFilterOperators(filter?.filterAccesor ?? property),
                  )
                }
              >
                <div className='flex items-center justify-center gap-2'>
                  {filter?.icon}
                  {filter?.filterName}
                </div>
              </MenuItem>
            );
          })}
        </MenuList>
      </Menu>
    </div>
  );
};
