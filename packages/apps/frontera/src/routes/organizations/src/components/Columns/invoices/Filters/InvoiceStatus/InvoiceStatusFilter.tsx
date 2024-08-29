import { useSearchParams } from 'react-router-dom';

import { observer } from 'mobx-react-lite';
import { FilterItem } from '@store/types.ts';

import { SelectOption } from '@ui/utils/types';
import { useStore } from '@shared/hooks/useStore';
import { Checkbox } from '@ui/form/Checkbox/Checkbox';
import { FilterHeader } from '@shared/components/Filters/FilterHeader';
import {
  InvoiceStatus,
  ColumnViewType,
  ComparisonOperator,
} from '@graphql/types';

const options: SelectOption<InvoiceStatus>[] = [
  { label: 'Out of contract', value: InvoiceStatus.OnHold },
  { label: 'Scheduled', value: InvoiceStatus.Scheduled },
  { label: 'Void', value: InvoiceStatus.Void },
  { label: 'Overdue', value: InvoiceStatus.Overdue },
  { label: 'Paid', value: InvoiceStatus.Paid },
];

const defaultFilter: FilterItem = {
  property: ColumnViewType.InvoicesInvoiceStatus,
  value: [],
  active: false,
  caseSensitive: false,
  includeEmpty: false,
  operation: ComparisonOperator.In,
};

export const InvoiceStatusFilter = observer(() => {
  const [searchParams] = useSearchParams();
  const preset = searchParams.get('preset');
  const store = useStore();
  const tableViewDef = store.tableViewDefs.getById(preset ?? '');
  const filter =
    tableViewDef?.getFilter(defaultFilter.property) ?? defaultFilter;

  const toggle = () => {
    tableViewDef?.toggleFilter(filter);
  };

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
