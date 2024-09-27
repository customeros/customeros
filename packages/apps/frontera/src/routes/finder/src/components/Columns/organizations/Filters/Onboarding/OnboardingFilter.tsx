import { useSearchParams } from 'react-router-dom';

import { FilterItem } from '@store/types';
import { observer } from 'mobx-react-lite';

import { SelectOption } from '@ui/utils/types';
import { useStore } from '@shared/hooks/useStore';
import { Checkbox } from '@ui/form/Checkbox/Checkbox';
import {
  ColumnViewType,
  OnboardingStatus,
  ComparisonOperator,
} from '@graphql/types';

import { FilterHeader } from '../../../shared/Filters/abstract/FilterHeader';

const options: SelectOption<OnboardingStatus>[] = [
  { label: 'Not started', value: OnboardingStatus.NotStarted },
  { label: 'Late', value: OnboardingStatus.Late },
  { label: 'Stuck', value: OnboardingStatus.Stuck },
  { label: 'On track', value: OnboardingStatus.OnTrack },
  { label: 'Done', value: OnboardingStatus.Done },
  { label: 'Not applicable', value: OnboardingStatus.NotApplicable },
  { label: 'Successful', value: OnboardingStatus.Successful },
];

const defaultFilter: FilterItem = {
  property: ColumnViewType.OrganizationsOnboardingStatus,
  value: [],
  active: false,
  caseSensitive: false,
  includeEmpty: false,
  operation: ComparisonOperator.In,
};

export const OnboardingFilter = observer(() => {
  const [searchParams] = useSearchParams();
  const preset = searchParams.get('preset');

  const store = useStore();
  const tableViewDef = store.tableViewDefs.getById(preset ?? '');
  const filter =
    tableViewDef?.getFilter(defaultFilter.property) ?? defaultFilter;

  const toggle = () => {
    tableViewDef?.toggleFilter(filter);
  };

  const handleSelect = (value: OnboardingStatus) => () => {
    tableViewDef?.setFilter({
      ...filter,
      value: filter.value.includes(value)
        ? filter.value.filter((v: OnboardingStatus) => v !== value)
        : [...filter.value, value],
      active: true,
    });
  };

  const handleSelectAll = () => {
    tableViewDef?.setFilter({
      ...filter,
      value:
        filter.value.length === options.length
          ? []
          : options.map((o) => o.value),
      active: true,
    });
  };

  const isAllChecked = filter.value.length === options.length;

  return (
    <>
      <FilterHeader
        onToggle={toggle}
        onDisplayChange={() => {}}
        isChecked={filter.active ?? false}
      />
      <div className='flex flex-col gap-2 items-start'>
        <Checkbox isChecked={isAllChecked} onChange={handleSelectAll}>
          <p className='text-sm'>
            {isAllChecked ? 'Deselect all' : 'Select all'}
          </p>
        </Checkbox>
        {options.map((option) => (
          <Checkbox
            key={option.label}
            onChange={handleSelect(option.value)}
            isChecked={filter.value.includes(option.value)}
          >
            <p className='text-sm'>{option.label}</p>
          </Checkbox>
        ))}
      </div>
    </>
  );
});
