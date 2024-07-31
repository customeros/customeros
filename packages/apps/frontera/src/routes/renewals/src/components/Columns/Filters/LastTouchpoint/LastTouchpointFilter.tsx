import { memo, useRef, useMemo, useState, useEffect, useCallback } from 'react';

import { produce } from 'immer';
import { useRecoilValue } from 'recoil';
import { subDays } from 'date-fns/subDays';
import { Column } from '@tanstack/react-table';

import { Checkbox } from '@ui/form/Checkbox/Checkbox';
import { Radio, RadioGroup } from '@ui/form/Radio/Radio';
import { Organization, LastTouchpointType } from '@graphql/types';

import { allTime, touchpoints } from './util';
import { FilterHeader, useFilterToggle } from '../shared';
import {
  LastTouchpointSelector,
  useLastTouchpointFilter,
} from './LastTouchpointFilter.atom';

interface LastTouchpointFilterProps {
  onFilterValueChange?: Column<Organization>['setFilterValue'];
}

export const LastTouchpointFilter = ({
  onFilterValueChange,
}: LastTouchpointFilterProps) => {
  const [filter, setFilter] = useLastTouchpointFilter();
  const filterValue = useRecoilValue(LastTouchpointSelector);

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
        return subDays(new Date(), value).toISOString().split('T')[0];
      }),
    [],
  );

  const isAllSelected =
    filter.value.length === touchpoints.length && filter.value.length > 0;

  const handleSelectAll = useCallback(() => {
    setFilter((prev) => {
      const next = produce(prev, (draft) => {
        if (isAllSelected) {
          draft.value = [];
        } else {
          draft.isActive = true;
          draft.value = touchpoints.map(({ value }) => value);
        }
      });

      toggle.setIsActive(next.isActive);

      return next;
    });
  }, [isAllSelected, setFilter, toggle.setIsActive]);

  const handleSelect = useCallback(
    (value: LastTouchpointType) => {
      setFilter((prev) => {
        const next = produce(prev, (draft) => {
          draft.isActive = true;

          if (draft.value.includes(value)) {
            draft.value = draft.value.filter((item) => item !== value);
          } else {
            draft.value.push(value);
          }
        });

        toggle.setIsActive(next.isActive);

        return next;
      });
    },
    [setFilter, toggle.setIsActive],
  );

  const handleDateChange = (value: string) => {
    setFilter((prev) => {
      const next = produce(prev, (draft) => {
        draft.isActive = true;
        draft.after = value;
      });

      toggle.setIsActive(next.isActive);

      return next;
    });
  };

  useEffect(() => {
    onFilterValueChange?.(filterValue.isActive ? filterValue : undefined);
  }, [filterValue.value, filterValue.isActive, filterValue.after]);

  return (
    <>
      <FilterHeader
        isChecked={toggle.isActive}
        onToggle={toggle.handleChange}
        onDisplayChange={toggle.handleClick}
      />

      <RadioGroup
        value={filter.after}
        disabled={!filter.isActive}
        name='last-touchpoint-before'
        onValueChange={handleDateChange}
        className='border-b pb-2 border-gray-200'
      >
        <div className='flex flex-col gap-2 items-start'>
          <Radio value={week}>
            <span className='text-sm'>Last 7 days</span>
          </Radio>
          <Radio value={month}>
            <span className='text-sm'>Last 30 days</span>
          </Radio>
          <Radio value={quarter}>
            <span className='text-sm'>Last 90 days</span>
          </Radio>
          <Radio value={allTime}>
            <span className='text-sm'>All time</span>
          </Radio>
        </div>
      </RadioGroup>

      <div className='flex flex-col space-y-2 items-start mt-2 px-[4px] mx-[-4px] relative overflow-x-hidden overflow-y-auto'>
        <Checkboxes
          value={filter.value}
          onCheck={handleSelect}
          onCheckAll={handleSelectAll}
          isAllSelected={isAllSelected}
        />
      </div>
    </>
  );
};

interface CheckboxOptionsProps {
  value: string[];
  isAllSelected: boolean;
  onCheckAll: () => void;
  onCheck: (value: LastTouchpointType) => void;
}

const makeState = (values: string[]) =>
  values.reduce((acc, curr) => ({ ...acc, [curr]: true }), {});

const allCheckedState = touchpoints.reduce(
  (acc, { value }) => ({ ...acc, [value]: true }),
  {},
);
const allUnchecked = touchpoints.reduce(
  (acc, { value }) => ({ ...acc, [value]: false }),
  {},
);

const Checkboxes = memo(
  ({
    value = [],
    onCheck,
    onCheckAll,
    isAllSelected,
  }: CheckboxOptionsProps) => {
    const timeoutRef = useRef<NodeJS.Timeout>();
    const [_isAllChecked, setIsAllChecked] = useState(() => isAllSelected);
    const [checked, setChecked] = useState<Record<string, boolean>>(
      makeState(value),
    );

    const handleCheck = (v: string) => {
      setChecked((prev) =>
        produce(prev, (draft) => {
          draft[v] = !draft[v];
        }),
      );

      timeoutRef.current && clearTimeout(timeoutRef.current);
      timeoutRef.current = setTimeout(
        () => onCheck(v as LastTouchpointType),
        250,
      );
    };

    const handleCheckAll = () => {
      setIsAllChecked((prev) => !prev);
      setChecked(_isAllChecked ? allUnchecked : allCheckedState);

      timeoutRef.current && clearTimeout(timeoutRef.current);
      timeoutRef.current = setTimeout(onCheckAll, 250);
    };

    useEffect(() => {
      setIsAllChecked(Object.values(checked).every((v) => v));
    }, [checked]);

    return (
      <>
        <div className='sticky top-0 w-full z-10 bg-white border-b border-gray-200 pb-2'>
          <Checkbox
            size='sm'
            onChange={handleCheckAll}
            isChecked={_isAllChecked}
          >
            <span className='text-sm'>
              {_isAllChecked ? 'Deselect all' : 'Select all'}
            </span>
          </Checkbox>
        </div>
        {touchpoints.map(({ label, value }) => (
          <Checkbox
            size='sm'
            key={value}
            iconSize='md'
            isChecked={checked[value]}
            onChange={() => handleCheck(value)}
            className='rounded-sm border border-gray-200'
          >
            <span className='text-sm line-clamp-1'>{label}</span>
          </Checkbox>
        ))}
      </>
    );
  },
);
