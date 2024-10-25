import { startTransition } from 'react';
import { useSearchParams } from 'react-router-dom';

import { FilterItem } from '@store/types';
import { observer } from 'mobx-react-lite';
import difference from 'lodash/difference';

import { useStore } from '@shared/hooks/useStore';
import { Checkbox } from '@ui/form/Checkbox/Checkbox';
import {
  ColumnViewType,
  ComparisonOperator,
  FlowParticipantStatus,
} from '@graphql/types';

import { FilterHeader } from '../../../shared/Filters/abstract';

const optionsDict = {
  // Temporary - Should be replaced with correct enum
  [FlowParticipantStatus.OnHold]: 'On Hold',
  [FlowParticipantStatus.Ready]: 'Ready',
  [FlowParticipantStatus.Scheduled]: 'Scheduled',
  [FlowParticipantStatus.InProgress]: 'In Progress',
  [FlowParticipantStatus.Completed]: 'Completed',
  [FlowParticipantStatus.GoalAchieved]: 'Goal achieved',
};
const options = Object.entries(optionsDict);

const defaultFilter: FilterItem = {
  property: ColumnViewType.ContactsFlowStatus,
  value: [],
  active: false,
  caseSensitive: false,
  includeEmpty: false,
  operation: ComparisonOperator.Contains,
};

export const ContactFlowStatusFilter = observer(() => {
  const [searchParams] = useSearchParams();
  const preset = searchParams.get('preset');

  const store = useStore();
  const tableViewDef = store.tableViewDefs.getById(preset ?? '');

  const filter =
    tableViewDef?.getFilter(defaultFilter.property) ?? defaultFilter;

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

  const isAllChecked = filter.value.length === options?.length;

  const handleSelectAll = () => {
    let nextValue: string[] = [];

    if (isAllChecked) {
      tableViewDef?.setFilter({
        ...filter,
        value: difference(
          filter.value,
          options.map(([k]) => k),
        ),
        active: false,
      });

      return;
    }

    nextValue = options.map(([k]) => k);

    tableViewDef?.setFilter({
      ...filter,
      value: nextValue,
      active: nextValue.length > 0,
    });
  };

  return (
    <div className='max-h-[500px] overflow-y-auto overflow-x-hidden '>
      <FilterHeader
        onToggle={toggle}
        onDisplayChange={() => {}}
        isChecked={filter.active ?? false}
      />

      <div className='pt-2 pb-2 border-b border-gray-200'>
        <Checkbox isChecked={isAllChecked} onChange={handleSelectAll}>
          <p className='text-sm'>
            {isAllChecked ? 'Deselect all' : 'Select all'}
          </p>
        </Checkbox>
      </div>

      <div className='max-h-[80vh] overflow-y-auto overflow-x-hidden -mr-3'>
        {options.map(([k, v]) => (
          <Checkbox
            size='md'
            className='mt-2'
            key={`option-${k}`}
            onChange={() => handleChange(k)}
            isChecked={filter.value.includes(k) ?? false}
            labelProps={{
              className:
                'text-sm mt-2 whitespace-nowrap overflow-hidden overflow-ellipsis',
            }}
          >
            {v ?? 'Unnamed'}
          </Checkbox>
        ))}
      </div>
    </div>
  );
});
