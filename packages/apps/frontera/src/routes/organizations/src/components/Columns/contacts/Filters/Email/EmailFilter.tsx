import { RefObject, startTransition } from 'react';
import { useSearchParams } from 'react-router-dom';

import { FilterItem } from '@store/types';
import { observer } from 'mobx-react-lite';

import { useStore } from '@shared/hooks/useStore';
import { Checkbox } from '@ui/form/Checkbox/Checkbox';
import { ColumnViewType, ComparisonOperator } from '@graphql/types';

import {
  FilterHeader,
  DebouncedSearchInput,
} from '../../../shared/Filters/abstract';

interface EmailFilterProps {
  property?: ColumnViewType;
  initialFocusRef: RefObject<HTMLInputElement>;
}

const defaultFilter: FilterItem = {
  property: ColumnViewType.ContactsEmails,
  value: '',
  active: false,
  caseSensitive: false,
  includeEmpty: true,
  operation: ComparisonOperator.Contains,
};
const defaultVerifiedFilter: FilterItem = {
  property: 'EMAIL_VERIFIED',
  value: null,
  active: false,
  caseSensitive: false,
  includeEmpty: false,
  operation: ComparisonOperator.Eq,
};

export const EmailFilter = observer(
  ({ initialFocusRef, property }: EmailFilterProps) => {
    const [searchParams] = useSearchParams();
    const preset = searchParams.get('preset');

    const store = useStore();
    const tableViewDef = store.tableViewDefs.getById(preset ?? '');

    const filter = tableViewDef?.getFilter(
      property || defaultFilter.property,
    ) ?? { ...defaultFilter, property: property || defaultFilter.property };

    const filterVerified = tableViewDef?.getFilter('EMAIL_VERIFIED') ?? {
      ...defaultVerifiedFilter,
    };
    const toggle = () => {
      if (filter.active || filterVerified.active) {
        tableViewDef?.setFilter({
          ...filter,
          active: false,
        });
        setTimeout(() => {
          tableViewDef?.setFilter({
            ...filterVerified,
            active: false,
          });
        }, 1);
      }

      if (!filter.active && !filterVerified.active) {
        tableViewDef?.setFilter({
          ...filter,
          active: true,
        });
      }
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

    const handleFilterVerified = (value: boolean) => {
      tableViewDef?.setFilter({
        ...filterVerified,
        value: value ? 'verified' : '',
        active: value,
      });
    };

    return (
      <>
        <FilterHeader
          onToggle={toggle}
          onDisplayChange={() => {}}
          isChecked={(filter.active || filterVerified.active) ?? false}
        />

        <DebouncedSearchInput
          value={filter.value}
          ref={initialFocusRef}
          onChange={handleChange}
          placeholder='e.g. john.doe@acme.com'
        />
        <div className='flex flex-col gap-2 mt-2 items-start'>
          <Checkbox
            isChecked={filterVerified.active}
            onChange={(value) => handleFilterVerified(value as boolean)}
          >
            <p className='text-sm'>Verified</p>
          </Checkbox>
        </div>
      </>
    );
  },
);
