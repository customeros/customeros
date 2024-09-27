import { useSearchParams } from 'react-router-dom';

import { observer } from 'mobx-react-lite';
import { FilterItem } from '@store/types.ts';

import { SelectOption } from '@ui/utils/types';
import { useStore } from '@shared/hooks/useStore';
import { Checkbox } from '@ui/form/Checkbox/Checkbox';
import { ColumnViewType, ComparisonOperator } from '@graphql/types';
import { FilterHeader } from '@shared/components/Filters/FilterHeader';

const options: SelectOption<number>[] = [
  { label: 'Monthly', value: 1 },
  { label: 'Quarterly', value: 3 },
  { label: 'Annually', value: 12 },
  { label: 'None', value: 0 },
];
const defaultFilter: FilterItem = {
  property: ColumnViewType.InvoicesBillingCycle,
  value: [],
  active: false,
  caseSensitive: false,
  includeEmpty: false,
  operation: ComparisonOperator.In,
};

export const BillingCycleFilter = observer(() => {
  const [searchParams] = useSearchParams();
  const preset = searchParams.get('preset');
  const store = useStore();
  const tableViewDef = store.tableViewDefs.getById(preset ?? '');
  const filter =
    tableViewDef?.getFilter(defaultFilter.property) ?? defaultFilter;

  const toggle = () => {
    tableViewDef?.toggleFilter(filter);
  };

  const handleSelect = (value: number) => () => {
    const newValue = filter.value.includes(value)
      ? filter.value.filter((v: number) => v !== value)
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
