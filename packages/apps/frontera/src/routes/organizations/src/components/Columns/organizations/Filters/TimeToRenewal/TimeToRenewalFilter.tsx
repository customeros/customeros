import { useMemo } from 'react';
import { useSearchParams } from 'react-router-dom';

import { FilterItem } from '@store/types';
import { addDays } from 'date-fns/addDays';
import { observer } from 'mobx-react-lite';

import { useStore } from '@shared/hooks/useStore';
import { Radio, RadioGroup } from '@ui/form/Radio/Radio';
import { ColumnViewType, ComparisonOperator } from '@graphql/types';

import { FilterHeader } from '../../../shared/Filters/abstract';

const defaultFilter: FilterItem = {
  property: ColumnViewType.OrganizationsRenewalDate,
  value: '',
  active: false,
  caseSensitive: false,
  includeEmpty: false,
  operation: ComparisonOperator.Lte,
};

export const TimeToRenewalFilter = observer(() => {
  const [searchParams] = useSearchParams();
  const preset = searchParams.get('preset');

  const store = useStore();
  const tableViewDef = store.tableViewDefs.getById(preset ?? '');
  const filter =
    tableViewDef?.getFilter(defaultFilter.property) ?? defaultFilter;

  const [week, month, quarter] = useMemo(
    () =>
      [7, 30, 90].map((value) => {
        return addDays(new Date(), value).toISOString().split('T')[0];
      }),
    [],
  );

  const toggle = () => {
    tableViewDef?.toggleFilter({
      ...filter,
      value: filter.value || week,
    });
  };

  const handleChange = (value: string) => {
    tableViewDef?.setFilter({
      ...filter,
      value,
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
      <RadioGroup
        name='timeToRenewal'
        value={filter.value}
        onValueChange={handleChange}
      >
        <div className='gap-2 flex flex-col items-start'>
          <Radio value={week}>
            <label className='text-sm'>Next 7 days</label>
          </Radio>
          <Radio value={month}>
            <label className='text-sm'>Next 30 days</label>
          </Radio>
          <Radio value={quarter}>
            <label className='text-sm'>Next 90 days</label>
          </Radio>
        </div>
      </RadioGroup>
    </>
  );
});
