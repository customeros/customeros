import { useSearchParams } from 'react-router-dom';

import { FilterItem } from '@store/types';
import { observer } from 'mobx-react-lite';

import { useStore } from '@shared/hooks/useStore';
import { Checkbox } from '@ui/form/Checkbox/Checkbox';
import { flowOptions } from '@organizations/components/Columns/flows/utils';
import { FlowStatus, ColumnViewType, ComparisonOperator } from '@graphql/types';

import { FilterHeader } from '../../../shared/Filters/abstract';

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

  const handleSelect = (value: FlowStatus) => () => {
    tableViewDef?.setFilter({
      ...filter,
      value: filter.value.includes(value)
        ? filter.value.filter((v: FlowStatus) => v !== value)
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
        {flowOptions.map((option) => {
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
