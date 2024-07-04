import { useMemo } from 'react';
import { useSearchParams } from 'react-router-dom';

import { FilterItem } from '@store/types';
import { subDays } from 'date-fns/subDays';
import { observer } from 'mobx-react-lite';

import { useStore } from '@shared/hooks/useStore';
import { Radio, RadioGroup } from '@ui/form/Radio/Radio';
import { ColumnViewType, ComparisonOperator } from '@graphql/types';

import { FilterHeader } from '../../../shared/Filters/abstract';

const defaultFilter: FilterItem = {
  property: ColumnViewType.OrganizationsChurnDate,
  value: subDays(new Date(), 30).toISOString().split('T')[0],
  active: false,
  caseSensitive: false,
  includeEmpty: false,
  operation: ComparisonOperator.Lte,
};

export const ChurnedFilter = observer(() => {
  const [searchParams] = useSearchParams();
  const preset = searchParams.get('preset');

  const store = useStore();
  const tableViewDef = store.tableViewDefs.getById(preset ?? '');
  const filter =
    tableViewDef?.getFilter(defaultFilter.property) ?? defaultFilter;

  const toggle = () => {
    tableViewDef?.toggleFilter(filter);
  };

  const [month, quarter, year] = useMemo(
    () =>
      [30, 90, 365].map((value) => {
        return subDays(new Date(), value).toISOString().split('T')[0];
      }),
    [],
  );

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
        disabled={!filter.active}
      >
        <div className='gap-2 flex flex-col items-start'>
          <Radio value={month}>
            <label className='text-sm'>Last 30 days</label>
          </Radio>
          <Radio value={quarter}>
            <label className='text-sm'>Last quarter</label>
          </Radio>
          <Radio value={year}>
            <label className='text-sm'>Last year</label>
          </Radio>
        </div>
      </RadioGroup>
    </>
  );
});
