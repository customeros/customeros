'use client';
import { useEffect } from 'react';

import { produce } from 'immer';
import { useRecoilValue } from 'recoil';
import { Column } from '@tanstack/react-table';

import { VStack } from '@ui/layout/Stack';
import { Text } from '@ui/typography/Text';
import { SelectOption } from '@ui/utils/types';
import { InvoiceStatus } from '@graphql/types';
import { Checkbox, CheckboxGroup } from '@ui/form/Checkbox';
import {
  FilterHeader,
  useFilterToggle,
} from '@shared/components/Filters/FilterHeader';

import {
  usePaymentStatusFilter,
  PaymentStatusFilterSelector,
} from './PaymentStatusFilter.atom';

const options: SelectOption<InvoiceStatus>[] = [
  { label: 'Due', value: InvoiceStatus.Due },
  { label: 'Paid', value: InvoiceStatus.Paid },
  { label: 'Void', value: InvoiceStatus.Void },
];

interface PaymentStatusFilterProps<T> {
  column: Column<T>;
}

export const PaymentStatusFilter = <T,>({
  column,
}: PaymentStatusFilterProps<T>) => {
  const [filter, setFilter] = usePaymentStatusFilter();
  const filterValue = useRecoilValue(PaymentStatusFilterSelector);

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

  const handleSelect = (value: InvoiceStatus) => () => {
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
