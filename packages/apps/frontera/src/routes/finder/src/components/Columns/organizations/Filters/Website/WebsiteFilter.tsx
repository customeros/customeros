import { useSearchParams } from 'react-router-dom';
import { RefObject, startTransition } from 'react';

import { FilterItem } from '@store/types';
import { observer } from 'mobx-react-lite';

import { useStore } from '@shared/hooks/useStore';
import { Checkbox } from '@ui/form/Checkbox/Checkbox';
import { ColumnViewType, ComparisonOperator } from '@graphql/types';

import {
  FilterHeader,
  DebouncedSearchInput,
} from '../../../shared/Filters/abstract';

interface WebsiteFilterProps {
  initialFocusRef: RefObject<HTMLInputElement>;
}

const defaultFilter: FilterItem = {
  property: ColumnViewType.OrganizationsWebsite,
  value: '',
  active: false,
  caseSensitive: false,
  includeEmpty: false,
  operation: ComparisonOperator.Contains,
};

export const WebsiteFilter = observer(
  ({ initialFocusRef }: WebsiteFilterProps) => {
    const [searchParams] = useSearchParams();
    const preset = searchParams.get('preset');

    const store = useStore();
    const tableViewDef = store.tableViewDefs.getById(preset ?? '');
    const filter =
      tableViewDef?.getFilter(defaultFilter.property) ?? defaultFilter;

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
          placeholder='e.g. www.customeros.ai'
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