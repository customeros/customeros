import { useSearchParams } from 'react-router-dom';
import { useState, RefObject, startTransition } from 'react';

import { FilterItem } from '@store/types';
import { observer } from 'mobx-react-lite';
import difference from 'lodash/difference';

import { Input } from '@ui/form/Input';
import { useStore } from '@shared/hooks/useStore';
import { Checkbox } from '@ui/form/Checkbox/Checkbox';
import { SearchSm } from '@ui/media/icons/SearchSm.tsx';
import { InputGroup, LeftElement } from '@ui/form/InputGroup';
import { ColumnViewType, ComparisonOperator } from '@graphql/types';

import { FilterHeader } from '../../../shared/Filters/abstract';

interface ContactFilterProps {
  placeholder?: string;
  property?: ColumnViewType;
  initialFocusRef: RefObject<HTMLInputElement>;
}

const defaultFilter: FilterItem = {
  property: ColumnViewType.ContactsPersona,
  value: [],
  active: false,
  caseSensitive: false,
  includeEmpty: false,
  operation: ComparisonOperator.Contains,
};

export const PersonaFilter = observer(
  ({ initialFocusRef, property, placeholder }: ContactFilterProps) => {
    const [searchParams] = useSearchParams();
    const [searchValue, setSearchValue] = useState('');
    const preset = searchParams.get('preset');

    const store = useStore();
    const tableViewDef = store.tableViewDefs.getById(preset ?? '');

    const filter = tableViewDef?.getFilter(
      property || defaultFilter.property,
    ) ?? { ...defaultFilter, property: property || defaultFilter.property };

    const toggle = () => {
      tableViewDef?.toggleFilter(filter);
    };

    const handleChange = (newValue: string) => {
      const filterValue = Array.isArray(filter.value) ? filter.value : [];
      const value = filterValue?.includes(newValue)
        ? filterValue.filter((e) => e !== newValue)
        : [...filterValue, newValue];

      startTransition(() => {
        tableViewDef?.setFilter({
          ...filter,
          value,
          active: filter.active || true,
        });
      });
      setSearchValue('');
    };

    const options = store.tags
      .toArray()
      .map((e) => e.value.name)
      .filter((e) => !!e?.length)
      .sort((a, b) => a.localeCompare(b));

    const isAllChecked = filter.value.length === options?.length;

    const handleSelectAll = () => {
      let nextValue: string[] = [];

      if (isAllChecked) {
        tableViewDef?.setFilter({
          ...filter,
          value: difference(filter.value, options),
          active: false,
        });

        return;
      }

      if (searchValue) {
        nextValue = [...options, ...difference(filter.value, options)];
      } else {
        nextValue = options;
      }

      tableViewDef?.setFilter({
        ...filter,
        value: nextValue,
        active: nextValue.length > 0,
      });
    };

    return (
      <div className='max-h-[500px] overflow-auto'>
        <FilterHeader
          onToggle={toggle}
          onDisplayChange={() => {}}
          isChecked={filter.active ?? false}
        />
        <InputGroup>
          <LeftElement>
            <SearchSm color='gray.500' />
          </LeftElement>
          <Input
            size='sm'
            value={searchValue}
            ref={initialFocusRef}
            className='border-none'
            placeholder={placeholder || 'e.g. CustomerOS'}
            onChange={(e) => setSearchValue(e.target.value)}
          />
        </InputGroup>

        <div className='pt-2 pb-2 border-b border-gray-200'>
          <Checkbox isChecked={isAllChecked} onChange={handleSelectAll}>
            <p className='text-sm'>
              {isAllChecked ? 'Deselect all' : 'Select all'}
            </p>
          </Checkbox>
        </div>

        <div className='max-h-[80vh] overflow-y-auto -mr-3 '>
          {options
            .filter((e) =>
              searchValue?.length
                ? e.trim().toLowerCase().includes(searchValue)
                : true,
            )
            ?.map((e) => (
              <Checkbox
                size='md'
                className='mt-2 '
                key={`option-${e}`}
                onChange={() => handleChange(e)}
                isChecked={filter.value.includes(e) ?? false}
                labelProps={{
                  className: 'text-sm mt-2 line-clamp-1 w-[150px]',
                }}
              >
                {e ?? 'Unnamed'}
              </Checkbox>
            ))}
        </div>
      </div>
    );
  },
);
