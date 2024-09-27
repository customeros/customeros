import { useMemo } from 'react';
import { useSearchParams } from 'react-router-dom';

import { FilterItem } from '@store/types';
import { observer } from 'mobx-react-lite';
import { subDays } from 'date-fns/subDays';

import { useStore } from '@shared/hooks/useStore';
import { Radio, RadioGroup } from '@ui/form/Radio/Radio';
import { ColumnViewType, ComparisonOperator } from '@graphql/types';

import { FilterHeader } from '../../../shared/Filters/abstract/FilterHeader';

const defaultFilter: FilterItem = {
  property: ColumnViewType.OrganizationsCreatedDate,
  value: new Date().toISOString().split('T')[0],
  active: false,
  caseSensitive: false,
  includeEmpty: false,
  operation: ComparisonOperator.Lte,
};

export const CreatedDateFilter = observer(() => {
  const [searchParams] = useSearchParams();
  const preset = searchParams.get('preset');

  const store = useStore();
  const tableViewDef = store.tableViewDefs.getById(preset ?? '');
  const filter =
    tableViewDef?.getFilter(defaultFilter.property) ?? defaultFilter;

  const [lastDay, last7Days, moreThan7Days] = useMemo(
    () =>
      [0, 7, 10000].map((value) => {
        return subDays(new Date(), value).toISOString().split('T')[0];
      }),
    [],
  );

  const toggle = () => {
    tableViewDef?.toggleFilter(filter);
  };

  const handleChange = (value: string) => {
    tableViewDef?.setFilter({
      ...filter,
      value,
      active: filter?.active || true,
    });
  };

  return (
    <>
      <FilterHeader
        onToggle={toggle}
        onDisplayChange={() => {}} // to be removed
        isChecked={filter?.active ?? false}
      />
      <RadioGroup
        name='created-date'
        value={filter?.value}
        onValueChange={handleChange}
      >
        <div className='flex flex-col gap-2 items-start'>
          <Radio value={lastDay}>
            <span className='text-sm'>Last day</span>
          </Radio>
          <Radio value={last7Days}>
            <span className='text-sm'>Last 7 days</span>
          </Radio>
          <Radio value={moreThan7Days}>
            <span className='text-sm'>More than 7 days ago</span>
          </Radio>
        </div>
      </RadioGroup>
    </>
  );
});
