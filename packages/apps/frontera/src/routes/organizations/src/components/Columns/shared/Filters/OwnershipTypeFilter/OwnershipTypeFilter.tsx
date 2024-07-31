import { useSearchParams } from 'react-router-dom';
import { RefObject, startTransition } from 'react';

import { FilterItem } from '@store/types';
import { observer } from 'mobx-react-lite';

import { useStore } from '@shared/hooks/useStore';
import { Checkbox } from '@ui/form/Checkbox/Checkbox';
import { FilterHeader } from '@shared/components/Filters';
import { ColumnViewType, ComparisonOperator } from '@graphql/types';

interface ContactFilterProps {
  placeholder?: string;
  property?: ColumnViewType;
  initialFocusRef: RefObject<HTMLInputElement>;
}

const defaultFilter: FilterItem = {
  property: ColumnViewType.OrganizationsIsPublic,
  value: [],
  active: false,
  caseSensitive: false,
  includeEmpty: false,
  operation: ComparisonOperator.Contains,
};

export const OwnershipTypeFilter = observer(
  ({ property }: ContactFilterProps) => {
    const [searchParams] = useSearchParams();
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
    };

    return (
      <>
        <FilterHeader
          onToggle={toggle}
          onDisplayChange={() => {}}
          isChecked={filter.active ?? false}
        />

        <div className='max-h-[80vh] overflow-y-auto -mr-3'>
          <Checkbox
            size='md'
            className='mt-2'
            onChange={() => handleChange('public')}
            labelProps={{ className: 'text-sm mt-2' }}
            isChecked={filter.value.includes('public') ?? false}
          >
            Public
          </Checkbox>
          <Checkbox
            size='md'
            className='mt-2'
            onChange={() => handleChange('private')}
            labelProps={{ className: 'text-sm mt-2' }}
            isChecked={filter.value.includes('private') ?? false}
          >
            Private
          </Checkbox>
        </div>
      </>
    );
  },
);
