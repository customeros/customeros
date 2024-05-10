import { useEffect } from 'react';

import { produce } from 'immer';
import { useRecoilValue } from 'recoil';
import { Column } from '@tanstack/react-table';

import { SelectOption } from '@ui/utils/types';
import { Checkbox } from '@ui/form/Checkbox/Checkbox';
import {
  FilterHeader,
  useFilterToggle,
} from '@shared/components/Filters/FilterHeader';

import {
  useBillingCycleFilter,
  BillingCycleFilterSelector,
} from './BillingCycleFilter.atom';

const options: SelectOption<string>[] = [
  { label: 'Monthly', value: 'MONTHLY' },
  { label: 'Quarterly', value: 'QUARTERLY' },
  { label: 'Annually', value: 'ANNUALLY' },
  { label: 'None', value: 'NONE' },
];

interface BillingCycleFilterProps<T> {
  column: Column<T>;
}

export const BillingCycleFilter = <T,>({
  column,
}: BillingCycleFilterProps<T>) => {
  const [filter, setFilter] = useBillingCycleFilter();
  const filterValue = useRecoilValue(BillingCycleFilterSelector);

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

  const handleSelect = (value: string) => () => {
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
  };

  const handleSelectAll = () => {
    setFilter((prev) => {
      const next = produce(prev, (draft) => {
        draft.isActive = true;

        if (draft.value.length === options.length) {
          draft.value = [];
        } else {
          draft.value = options.map((option) => option.value);
        }
      });

      toggle.setIsActive(next.isActive);

      return next;
    });
  };

  useEffect(() => {
    column.setFilterValue?.(
      filterValue.isActive ? filterValue.value : undefined,
    );
  }, [filterValue.value.length, filterValue.isActive]);

  const isAllChecked = filterValue.value.length === options.length;

  return (
    <>
      <FilterHeader
        isChecked={toggle.isActive}
        onToggle={toggle.handleChange}
        onDisplayChange={toggle.handleClick}
      />
      <div className='flex flex-col gap-2 items-start'>
        <Checkbox isChecked={isAllChecked} onChange={handleSelectAll}>
          <p className='text-sm'>
            {isAllChecked ? 'Deselect all' : 'Select all'}
          </p>
        </Checkbox>
        {options.map((option) => (
          <Checkbox
            key={option.label}
            isChecked={filter.value.includes(option.value)}
            onChange={handleSelect(option.value)}
          >
            <p className='text-sm'>{option.label}</p>
          </Checkbox>
        ))}
      </div>
    </>
  );
};
