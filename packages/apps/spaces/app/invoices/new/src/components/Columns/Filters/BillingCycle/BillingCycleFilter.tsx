'use client';
import { useEffect } from 'react';

import { produce } from 'immer';
import { useRecoilValue } from 'recoil';
import { Column } from '@tanstack/react-table';

import { VStack } from '@ui/layout/Stack';
import { Text } from '@ui/typography/Text';
import { SelectOption } from '@ui/utils/types';
import { ContractBillingCycle } from '@graphql/types';
import { Checkbox, CheckboxGroup } from '@ui/form/Checkbox';
import {
  FilterHeader,
  useFilterToggle,
} from '@shared/components/Filters/FilterHeader';

import {
  useBillingCycleFilter,
  BillingCycleFilterSelector,
} from './BillingCycleFilter.atom';

const options: SelectOption<ContractBillingCycle>[] = [
  { label: 'Monthly', value: ContractBillingCycle.MonthlyBilling },
  { label: 'Quarterly', value: ContractBillingCycle.QuarterlyBilling },
  { label: 'Annualy', value: ContractBillingCycle.AnnualBilling },
  { label: 'None', value: ContractBillingCycle.None },
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

  const handleSelect = (value: ContractBillingCycle) => () => {
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
    <CheckboxGroup size='md' value={filterValue.value}>
      <FilterHeader
        isChecked={toggle.isActive}
        onToggle={toggle.handleChange}
        onDisplayChange={toggle.handleClick}
      />
      <VStack spacing={2} align='flex-start'>
        <Checkbox isChecked={isAllChecked} onChange={handleSelectAll}>
          <Text fontSize='sm'>
            {isAllChecked ? 'Deselect all' : 'Select all'}
          </Text>
        </Checkbox>
        {options.map((option) => (
          <Checkbox
            key={option.label}
            value={option.value}
            onChange={handleSelect(option.value)}
          >
            <Text fontSize='sm'>{option.label}</Text>
          </Checkbox>
        ))}
      </VStack>
    </CheckboxGroup>
  );
};
