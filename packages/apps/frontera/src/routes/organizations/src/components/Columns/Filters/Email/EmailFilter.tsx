import { RefObject, startTransition } from 'react';
import { useSearchParams } from 'react-router-dom';

import { FilterItem } from '@store/types';
import { observer } from 'mobx-react-lite';

import { useStore } from '@shared/hooks/useStore';
import { Radio, RadioGroup } from '@ui/form/Radio';
import { ColumnViewType, ComparisonOperator } from '@graphql/types';

import { FilterHeader, DebouncedSearchInput } from '../shared';

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
      property: property || defaultVerifiedFilter.property,
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

    const handleFilterVerified = (newState: string) => {
      tableViewDef?.setFilter({
        ...filterVerified,
        value: newState,
        active: filterVerified.active || true,
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
          placeholder='e.g. john.doe@acme.com'
        />

        <RadioGroup
          name='emailVerified'
          value={filterVerified.value}
          onValueChange={(value) => handleFilterVerified(value)}
          disabled={!filterVerified.active}
        >
          <div className='flex flex-col gap-2 mt-2 items-start'>
            <Radio value={'verified'}>
              <p className='text-sm'>Verified</p>
            </Radio>
            <Radio value={'not-verified'}>
              <p className='text-sm'>Not Verified</p>
            </Radio>
          </div>
        </RadioGroup>
      </>
    );
  },
);
