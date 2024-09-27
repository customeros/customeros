import { useSearchParams } from 'react-router-dom';

import { FilterItem } from '@store/types';
import { observer } from 'mobx-react-lite';

import { useStore } from '@shared/hooks/useStore';
import { Checkbox } from '@ui/form/Checkbox/Checkbox';
import { ColumnViewType, ComparisonOperator } from '@graphql/types';

import { FilterHeader } from '../../../shared/Filters/abstract';

const defaultFilter: FilterItem = {
  property: ColumnViewType.ContractsRenewal,
  value: [],
  active: false,
  caseSensitive: false,
  includeEmpty: false,
  operation: ComparisonOperator.In,
};

export const RenewalFilter = observer(() => {
  const [searchParams] = useSearchParams();
  const preset = searchParams.get('preset');

  const store = useStore();
  const tableViewDef = store.tableViewDefs.getById(preset ?? '');
  const filter =
    tableViewDef?.getFilter(defaultFilter.property) ?? defaultFilter;

  const toggle = () => {
    tableViewDef?.toggleFilter(filter);
  };

  const handleSelect = (value: boolean) => () => {
    tableViewDef?.setFilter({
      ...filter,
      value: filter.value.includes(value)
        ? filter.value.filter((v: boolean) => v !== value)
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
        <Checkbox
          onChange={handleSelect(true)}
          isChecked={filter.value.includes(true)}
        >
          <span className='text-sm'>Auto-renews</span>
        </Checkbox>
        <Checkbox
          onChange={handleSelect(false)}
          isChecked={filter.value.includes(false)}
        >
          <span className='text-sm'>Not auto-renewing</span>
        </Checkbox>
      </div>
    </>
  );
});
