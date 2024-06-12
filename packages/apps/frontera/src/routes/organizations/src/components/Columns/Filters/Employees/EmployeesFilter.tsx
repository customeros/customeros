import { useSearchParams } from 'react-router-dom';

import { FilterItem } from '@store/types';
import { observer } from 'mobx-react-lite';

import { useStore } from '@shared/hooks/useStore';
import { Radio, RadioGroup } from '@ui/form/Radio/Radio';
import { ColumnViewType, ComparisonOperator } from '@graphql/types';

import { FilterHeader } from '../shared';

const defaultFilter: FilterItem = {
  property: ColumnViewType.OrganizationsEmployeeCount,
  value: '0-10',
  active: false,
  caseSensitive: false,
  includeEmpty: false,
  operation: ComparisonOperator.Eq,
};

export const EmployeesFilter = observer(() => {
  const [searchParams] = useSearchParams();
  const preset = searchParams.get('preset');

  const store = useStore();
  const tableViewDef = store.tableViewDefs.getById(preset ?? '');
  const filter =
    tableViewDef?.getFilter(defaultFilter.property) ?? defaultFilter;

  const toggle = () => {
    tableViewDef?.toggleFilter(filter);
  };

  const handleDateChange = (value: string) => {
    tableViewDef?.setFilter({
      ...filter,
      value,
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
        name='employees'
        value={filter.value}
        onValueChange={handleDateChange}
        disabled={!filter.active}
      >
        <div className='flex flex-col gap-2 items-start'>
          <Radio value='0-10'>
            <span className='text-sm'>0-10</span>
          </Radio>
          <Radio value='11-50'>
            <span className='text-sm'>11-50</span>
          </Radio>
          <Radio value='51-100'>
            <span className='text-sm'>51-100</span>
          </Radio>
          <Radio value='101-250'>
            <span className='text-sm'>101-250</span>
          </Radio>
          <Radio value='251-500'>
            <span className='text-sm'>251-500</span>
          </Radio>
          <Radio value='501-1000'>
            <span className='text-sm'>501-1,000</span>
          </Radio>
          <Radio value='1000'>
            <span className='text-sm'>1,000+</span>
          </Radio>
        </div>
      </RadioGroup>
    </>
  );
});
