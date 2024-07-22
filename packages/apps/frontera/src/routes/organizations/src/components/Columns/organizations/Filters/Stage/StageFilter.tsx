import { useSearchParams } from 'react-router-dom';

import { FilterItem } from '@store/types';
import { observer } from 'mobx-react-lite';

import { useStore } from '@shared/hooks/useStore';
import { Checkbox } from '@ui/form/Checkbox/Checkbox';
import { stageOptions } from '@organization/components/Tabs/panels/AboutPanel/util.ts';
import {
  ColumnViewType,
  OrganizationStage,
  ComparisonOperator,
} from '@graphql/types';

import { FilterHeader } from '../../../shared/Filters/abstract/FilterHeader';

const defaultFilter: FilterItem = {
  property: ColumnViewType.OrganizationsStage,
  value: [],
  active: false,
  caseSensitive: false,
  includeEmpty: false,
  operation: ComparisonOperator.In,
};

export const StageFilter = observer(() => {
  const [searchParams] = useSearchParams();
  const preset = searchParams.get('preset');

  const store = useStore();
  const tableViewDef = store.tableViewDefs.getById(preset ?? '');
  const filter =
    tableViewDef?.getFilter(defaultFilter.property) ?? defaultFilter;

  const toggle = () => {
    tableViewDef?.toggleFilter(filter);
  };

  const handleSelect = (value: OrganizationStage) => () => {
    const newValue = filter.value.includes(value)
      ? filter.value.filter((v: OrganizationStage) => v !== value)
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
        {stageOptions.map((option) => (
          <Checkbox
            key={option.value.toString()}
            isChecked={filter.value.includes(option.value)}
            onChange={handleSelect(option.value)}
          >
            <p className='text-sm'>{option.label}</p>
          </Checkbox>
        ))}
      </div>
    </>
  );
});
