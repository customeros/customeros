import { useMemo } from 'react';
import { useSearchParams } from 'react-router-dom';

import { observer } from 'mobx-react-lite';
import { subDays } from 'date-fns/subDays';
import { FilterItem } from '@store/types.ts';

import { useStore } from '@shared/hooks/useStore';
import { Radio, RadioGroup } from '@ui/form/Radio';
import { FilterHeader } from '@shared/components/Filters';
import { ColumnViewType, ComparisonOperator } from '@graphql/types';

const allTime = new Date('1970-01-01').toISOString().split('T')[0];

const defaultFilter: FilterItem = {
  property: ColumnViewType.InvoicesIssueDate,
  value: allTime,
  active: false,
  caseSensitive: false,
  includeEmpty: false,
  operation: ComparisonOperator.Eq,
};

export const IssueDateFilter = observer(
  ({ property }: { property?: string }) => {
    const [searchParams] = useSearchParams();
    const preset = searchParams.get('preset');

    const store = useStore();
    const tableViewDef = store.tableViewDefs.getById(preset ?? '');
    const filter = tableViewDef?.getFilter(
      property ?? defaultFilter.property,
    ) ?? {
      ...defaultFilter,

      property: property || defaultFilter.property,
    };

    const toggle = () => {
      tableViewDef?.toggleFilter(filter);
    };

    const [day, days3, week] = useMemo(
      () =>
        [1, 3, 7].map((value) => {
          return subDays(new Date(), value).toISOString().split('T')[0];
        }),
      [],
    );

    const handleDateChange = (value: string) => {
      tableViewDef?.setFilter({
        ...filter,
        active: filter?.active || true,
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
          name='last-touchpoint-date-before'
          value={filter.value}
          onValueChange={handleDateChange}
        >
          <div className='flex flex-col gap-2 items-start'>
            <Radio value={day}>
              <span className='text-sm'>Last day</span>
            </Radio>
            <Radio value={days3}>
              <span className='text-sm'>Last 3 days</span>
            </Radio>
            <Radio value={week}>
              <span className='text-sm'>Last 7 days</span>
            </Radio>
            <Radio value={allTime}>
              <span className='text-sm'>More than 7 days ago</span>
            </Radio>
          </div>
        </RadioGroup>
      </>
    );
  },
);
