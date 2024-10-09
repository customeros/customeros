import { useRef, useState, useEffect } from 'react';

import { FilterType } from '@finder/components/Columns/organizations/filtersType';

import { Combobox } from '@ui/form/Combobox';
import { Button } from '@ui/form/Button/Button';
import { IconButton } from '@ui/form/IconButton';
import { FilterLines } from '@ui/media/icons/FilterLines';
import { components, OptionProps } from '@ui/form/Select/Select';
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from '@ui/overlay/Popover/Popover';
import {
  FilterItem,
  ColumnView,
  ColumnViewType,
  ComparisonOperator,
} from '@graphql/types';

import { Filter } from '../Filter/Filter';
import { categorySelected } from '../Filter/utils/utils';

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

  const getFilterGroupOptions = (property: string) => {
    const filterType = Object.values(filterTypes).find(
      (type) => type.filterAccesor === property,
    );

    return filterType?.groupOptions;
  };

  useEffect(() => {
    if (isOpen) {
      setTimeout(() => {
        inputRef.current?.focus();
      }, 100);
    }
  }, [isOpen]);

  useEffect(() => {
    const handleKeyDown = (event: KeyboardEvent) => {
      if (event.key === 'f') {
        setIsOpen(true);
        event.preventDefault();
      }
    };

    window.addEventListener('keydown', handleKeyDown);

    return () => {
      window.removeEventListener('keydown', handleKeyDown);
    };
  }, []);

  const handleChangeOperator = (
    operation: string,
    filter: FilterItem,
    index: number,
  ) => {
    const newValue = filter.value;
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
        value: newValue,
        includeEmpty: operation === ComparisonOperator.IsEmpty ? true : false,
      },
      index,
    );
  };

  const handleChangeFilterValue = (
    value: string | Date | string[],
    filter: FilterItem,
    index: number,
  ) => {
    if (
      filter.property === 'EMAIL_VERIFICATION_WORK_EMAIL' ||
      filter.property === 'EMAIL_VERIFICATION_PERSONAL_EMAIL'
    ) {
      const emailVerificationValues = Array.isArray(value) ? value : [value];

      const updatedValue = emailVerificationValues.map((val) => ({
        category: categorySelected(val as string),
        value: val,
      }));

      setFilters(
        {
          ...filter,
          value: updatedValue,
          property: filter.property,
          active: true,
        },
        index,
      );

      if ((value as [])?.length === 0) {
        setFilters(
          {
            ...filter,
            active: false,
          },
          index,
        );
      }
    }

    if (
      filter.property !== 'EMAIL_VERIFICATION_WORK_EMAIL' &&
      filter.property !== 'EMAIL_VERIFICATION_PERSONAL_EMAIL'
    ) {
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
    .filter(Boolean);

  const Option = ({ children, ...props }: OptionProps) => {
    const data = props?.data as {
      label: string;
      value: string;
      icon: JSX.Element;
    };

    return (
      <components.Option {...props} className='group'>
        <div className='flex justify-start items-center gap-2'>
          <span className='h-6 align-middle'>{data.icon}</span>
          <span className='align-middle text-sm'>{data.label}</span>
        </div>
      </components.Option>
    );
  };

  return (
    <div className='flex gap-2 flex-wrap'>
      {filters.map((f, idx) => {
        const value =
          f.property === 'EMAIL_VERIFICATION' ? f?.value?.[0]?.value : f?.value;

        return (
          <Filter
            filterValue={value}
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
            icon={
              availableFilters.find(
                (filter) => filter?.columnType === f.property,
              )?.icon || <></>
            }
            groupOptions={
              getFilterGroupOptions(f.property) as unknown as {
                label: string;
                options: { id: string; label: string }[];
              }[]
            }
          />
        );
      })}

      <Popover open={isOpen} onOpenChange={(open) => setIsOpen(open)}>
        <PopoverTrigger asChild>
          {filters.length ? (
            <IconButton
              size='xs'
              variant='ghost'
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
        </PopoverTrigger>
        <PopoverContent
          side='bottom'
          align='start'
          className='py-1 min-w-[254px]'
        >
          <Combobox
            size='xs'
            escapeClearsValue
            maxHeight='600px'
            closeMenuOnSelect={false}
            placeholder='Select filter...'
            noOptionsMessage={() => 'Nothing in sight...'}
            components={{
              Option,
            }}
            onKeyDown={(e) => {
              if (e.key === 'Escape') setIsOpen(false);
            }}
            options={availableFilters.map((filter) => ({
              label: filter?.filterName,
              value: filter?.columnType,
              icon: filter?.icon,
            }))}
            onChange={(selectedOption) => {
              const selectedFilter = availableFilters.find(
                (filter) => filter?.columnType === selectedOption?.value,
              );

              if (selectedFilter) {
                onFilterSelect(selectedFilter, (property) =>
                  getFilterOperators(selectedFilter.filterAccesor || property),
                );
              }
              setIsOpen(false);
            }}
          />
        </PopoverContent>
      </Popover>
    </div>
  );
};
