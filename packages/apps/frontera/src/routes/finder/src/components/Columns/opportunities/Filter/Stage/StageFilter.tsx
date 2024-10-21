import { useSearchParams } from 'react-router-dom';

import { FilterItem } from '@store/types';
import { observer } from 'mobx-react-lite';

import { useStore } from '@shared/hooks/useStore';
import { Checkbox } from '@ui/form/Checkbox/Checkbox';
import {
  InternalStage,
  ColumnViewType,
  ComparisonOperator,
} from '@graphql/types';

import { makeStageLabels } from '../../Cells';
import { FilterHeader } from '../../../shared/Filters/abstract/FilterHeader';

const defaultFilter: FilterItem = {
  property: ColumnViewType.OpportunitiesStage,
  value: [],
  active: false,
  caseSensitive: false,
  includeEmpty: false,
  operation: ComparisonOperator.In,
};

const internalStageOptions = [
  { label: 'Closed Lost', value: InternalStage.ClosedLost },
  { label: 'Closed Won', value: InternalStage.ClosedWon },
];

export const StageFilter = observer(() => {
  const [searchParams] = useSearchParams();
  const preset = searchParams.get('preset') ?? undefined;

  const store = useStore();
  const tableViewDef = store.tableViewDefs.getById(preset ?? '');
  const filter =
    tableViewDef?.getFilter(defaultFilter.property) ?? defaultFilter;

  const toggle = () => {
    tableViewDef?.toggleFilter(filter);
  };
  const stageLabels = makeStageLabels(store, preset) ?? {};
  const statuses = Object.entries(stageLabels).map(([value, label]) => ({
    label,
    value,
  }));

  const handleSelect = (value: string) => () => {
    const newValue = filter.value.includes(value)
      ? filter.value.filter((v: string) => v !== value)
      : [...filter.value, value];

    tableViewDef?.setFilter({
      ...filter,
      value: newValue,
      active: newValue.length > 0,
    });
  };

  return (
    <>
      <FilterHeader
        onToggle={toggle}
        onDisplayChange={() => {}}
        isChecked={filter.active ?? false}
      />
      <div className='flex flex-col gap-2 items-start'>
        {[...statuses, ...internalStageOptions].map((option) => (
          <Checkbox
            key={option.value.toString()}
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
