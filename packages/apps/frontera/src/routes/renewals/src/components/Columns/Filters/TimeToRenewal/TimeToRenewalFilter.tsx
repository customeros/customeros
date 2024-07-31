import { useMemo, useEffect } from 'react';

import { produce } from 'immer';
import { useRecoilValue } from 'recoil';
import { addDays } from 'date-fns/addDays';
import { Column } from '@tanstack/react-table';

import { Organization } from '@graphql/types';
import { Radio, RadioGroup } from '@ui/form/Radio/Radio';

import { FilterHeader, useFilterToggle } from '../shared';
import {
  useTimeToRenewalFilter,
  TimeToRenewalFilterSelector,
} from './TimeToRenewalFilter.atom';

interface TimeToRenewalProps {
  onFilterValueChange?: Column<Organization>['setFilterValue'];
}

export const TimeToRenewalFilter = ({
  onFilterValueChange,
}: TimeToRenewalProps) => {
  const [filter, setFilter] = useTimeToRenewalFilter();
  const filterValue = useRecoilValue(TimeToRenewalFilterSelector);

  const toggle = useFilterToggle({
    defaultValue: filter.isActive,
    onToggle: (setIsActive) => {
      setFilter((prev) => {
        const next = produce(prev, (draft) => {
          draft.isActive = !draft.isActive;
        });

        setIsActive(next.isActive);

        return next;
      });
    },
  });

  const [week, month, quarter] = useMemo(
    () =>
      [7, 30, 90].map((value) => {
        return addDays(new Date(), value).toISOString().split('T')[0];
      }),
    [],
  );

  const handleChange = (value: string) => {
    setFilter((prev) => {
      const next = produce(prev, (draft) => {
        draft.isActive = true;
        draft.value = value;
      });

      toggle.setIsActive(next.isActive);

      return next;
    });
  };

  useEffect(() => {
    onFilterValueChange?.(filterValue.isActive ? filterValue.value : undefined);
  }, [filterValue.value, filterValue.isActive]);

  return (
    <>
      <FilterHeader
        isChecked={toggle.isActive}
        onToggle={toggle.handleChange}
        onDisplayChange={toggle.handleClick}
      />
      <RadioGroup
        name='timeToRenewal'
        value={filter.value}
        disabled={!filter.isActive}
        onValueChange={handleChange}
      >
        <div className='flex flex-col space-y-2 items-start'>
          <Radio value={week}>
            <span className='text-sm'>Next 7 days</span>
          </Radio>
          <Radio value={month}>
            <span className='text-sm'>Next 30 days</span>
          </Radio>
          <Radio value={quarter}>
            <span className='text-sm'>Next 90 days</span>
          </Radio>
        </div>
      </RadioGroup>
    </>
  );
};
