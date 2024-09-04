import { useSearchParams } from 'react-router-dom';

import { FilterItem } from '@store/types';
import { observer } from 'mobx-react-lite';

import { useStore } from '@shared/hooks/useStore';
import { Checkbox } from '@ui/form/Checkbox/Checkbox';
import {
  ColumnViewType,
  ComparisonOperator,
  FlowSequenceStatus,
} from '@graphql/types';

import { FilterHeader } from '../../../shared/Filters/abstract';

const flowSequencesOptions = [
  { label: 'Active', value: FlowSequenceStatus.Active },
  { label: 'Archived', value: FlowSequenceStatus.Archived },
  { label: 'Inactive', value: FlowSequenceStatus.Inactive },
  { label: 'Paused', value: FlowSequenceStatus.Paused },
];
const defaultFilter: FilterItem = {
  property: ColumnViewType.FlowSequenceStatus,
  value: [],
  active: false,
  caseSensitive: false,
  includeEmpty: false,
  operation: ComparisonOperator.In,
};

export const StatusFilter = observer(() => {
  const [searchParams] = useSearchParams();
  const preset = searchParams.get('preset');

  const store = useStore();
  const tableViewDef = store.tableViewDefs.getById(preset ?? '');
  const filter =
    tableViewDef?.getFilter(defaultFilter.property) ?? defaultFilter;

  const toggle = () => {
    tableViewDef?.toggleFilter(filter);
  };

  const handleSelect = (value: FlowSequenceStatus) => () => {
    tableViewDef?.setFilter({
      ...filter,
      value: filter.value.includes(value)
        ? filter.value.filter((v: FlowSequenceStatus) => v !== value)
        : [...filter.value, value],
      active: true,
    });
  };

  return (
    <>
      <FilterHeader
        onToggle={toggle}
        onDisplayChange={() => {}}
        isChecked={filter.active ?? false}
      />
      <div className='flex flex-col space-y-2 items-start'>
        {flowSequencesOptions.map((option) => {
          return (
            <Checkbox
              key={option.value}
              onChange={handleSelect(option.value)}
              isChecked={filter.value.includes(option.value)}
            >
              <span className='text-sm'>{option.label}</span>
            </Checkbox>
          );
        })}
      </div>
    </>
  );
});
