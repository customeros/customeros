import { useMemo, useEffect } from 'react';

import { produce } from 'immer';
import addDays from 'date-fns/addDays';
import subDays from 'date-fns/subDays';
import { useRecoilValue } from 'recoil';
import { Column } from '@tanstack/react-table';

import { Organization } from '@graphql/types';
import { Radio, RadioGroup } from '@ui/form/Radio';
import { FilterHeader, useFilterToggle } from '@shared/components/Filters';

import {
  useIssueDateFilter,
  IssueDateFilterSelector,
} from './IssueDateFilter.atom';

interface IssueDateProps {
  isPast?: boolean;
  onFilterValueChange?: Column<Organization>['setFilterValue'];
}

export const IssueDateFilter = ({
  onFilterValueChange,
  isPast,
}: IssueDateProps) => {
  const [filter, setFilter] = useIssueDateFilter();
  const filterValue = useRecoilValue(IssueDateFilterSelector);

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
        const op = isPast ? subDays : addDays;

        return op(new Date(), value).toISOString();
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

  useEffect(() => {
    setFilter({
      ...filter,
      value: week,
    });
  }, []);

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
        <div className='flex flex-col gap-2 items-start'>
          <Radio value={week}>
            <p className='text-sm'>{`${
              isPast ? 'Previous' : 'Next'
            } 7 days`}</p>
          </Radio>
          <Radio value={month}>
            <p className='text-sm'>{`${
              isPast ? 'Previous' : 'Next'
            } 30 days`}</p>
          </Radio>
          <Radio value={quarter}>
            <p className='text-sm'>{`${
              isPast ? 'Previous' : 'Next'
            } 90 days`}</p>
          </Radio>
        </div>
      </RadioGroup>
    </>
  );
};
