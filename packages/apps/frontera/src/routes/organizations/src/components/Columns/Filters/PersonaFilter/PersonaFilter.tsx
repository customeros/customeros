import { useSearchParams } from 'react-router-dom';
import { useState, RefObject, startTransition } from 'react';

import { FilterItem } from '@store/types';
import { observer } from 'mobx-react-lite';

import { Input } from '@ui/form/Input';
import { useStore } from '@shared/hooks/useStore';
import { Checkbox } from '@ui/form/Checkbox/Checkbox';
import { ColumnViewType, ComparisonOperator } from '@graphql/types';

import { FilterHeader } from '../shared';

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

    return (
      <>
        <FilterHeader
          onToggle={toggle}
          onDisplayChange={() => {}}
          isChecked={filter.active ?? false}
        />

        <Input
          value={searchValue}
          ref={initialFocusRef}
          onChange={(e) => setSearchValue(e.target.value)}
          placeholder={placeholder || 'e.g. CustomerOS'}
        />

        <div className='max-h-[80vh] overflow-y-auto -mr-3'>
          {store.tags
            .toArray()
            .filter((e) =>
              searchValue?.length
                ? e.value.name.trim().toLowerCase().includes(searchValue)
                : true,
            )
            ?.map((e) => (
              <Checkbox
                key={e.value.id}
                className='mt-2'
                size='md'
                isChecked={filter.value.includes(e.value.name) ?? false}
                labelProps={{ className: 'text-sm mt-2' }}
                onChange={() => handleChange(e.value.name)}
              >
                {e.value.name ?? 'Unnamed'}
              </Checkbox>
            ))}
        </div>
      </>
    );
  },
);
