import { useMemo, useEffect } from 'react';

import { produce } from 'immer';
import { useRecoilValue } from 'recoil';
import { subDays } from 'date-fns/subDays';
import { Column } from '@tanstack/react-table';

import { Organization } from '@graphql/types';
import { Radio, RadioGroup } from '@ui/form/Radio/Radio';

import { FilterHeader, useFilterToggle } from '../shared';
import { useChurnedFilter, ChurnedFilterSelector } from './ChurnedFilter.atom';

interface TimeToRenewalProps {
  onFilterValueChange?: Column<Organization>['setFilterValue'];
}

export const ChurnedFilter = ({ onFilterValueChange }: TimeToRenewalProps) => {
  const [filter, setFilter] = useChurnedFilter();
  const filterValue = useRecoilValue(ChurnedFilterSelector);

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

  const [month, quarter, year] = useMemo(
    () =>
      [30, 90, 365].map((value) => {
        return subDays(new Date(), value).toISOString().split('T')[0];
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
        onValueChange={handleChange}
        disabled={!filter.isActive}
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
};
