import { useSearchParams } from 'react-router-dom';
import { RefObject, startTransition } from 'react';

import { FilterItem } from '@store/types';
import { observer } from 'mobx-react-lite';

import { useStore } from '@shared/hooks/useStore';
import { Checkbox } from '@ui/form/Checkbox/Checkbox';
import { ColumnViewType, ComparisonOperator } from '@graphql/types';

import { FilterHeader, DebouncedSearchInput } from '../abstract';

interface SocialsFilterProps {
  placeholder?: string;
  property?: ColumnViewType;
  initialFocusRef: RefObject<HTMLInputElement>;
}

const defaultFilter: FilterItem = {
  property: ColumnViewType.OrganizationsSocials,
  value: '',
  active: false,
  caseSensitive: false,
  includeEmpty: false,
  operation: ComparisonOperator.Contains,
};

export const SocialsFilter = observer(
  ({ initialFocusRef, property, placeholder }: SocialsFilterProps) => {
    const [searchParams] = useSearchParams();
    const preset = searchParams.get('preset');

    const store = useStore();
    const tableViewDef = store.tableViewDefs.getById(preset ?? '');
    const filter = tableViewDef?.getFilter(
      property || defaultFilter.property,
    ) ?? {
      ...defaultFilter,
      property: property || defaultFilter.property,
    };

    const toggle = () => {
      tableViewDef?.toggleFilter(filter);
    };

    const handleChange = (value: string) => {
      startTransition(() => {
        tableViewDef?.setFilter({
          ...filter,
          value,
          active: filter.active || true,
        });
      });
    };

    const handleShowEmpty = (isChecked: boolean) => {
      tableViewDef?.setFilter({
        ...filter,
        includeEmpty: isChecked,
        active: filter.active || true,
      });
    };

    return (
      <>
        <FilterHeader
          onToggle={toggle}
          onDisplayChange={() => {}}
          isChecked={filter.active ?? false}
        />

        <DebouncedSearchInput
          value={filter.value}
          ref={initialFocusRef}
          onChange={handleChange}
          placeholder={placeholder || 'linkedin.com/company/customer-os'}
        />

        <Checkbox
          size='sm'
          className='mt-2'
          isChecked={filter.includeEmpty ?? false}
          labelProps={{ className: 'text-sm mt-2' }}
          onChange={(isChecked) => handleShowEmpty(isChecked as boolean)}
        >
          Unknown
        </Checkbox>
      </>
    );
  },
);
