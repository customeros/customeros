'use client';
import { useMemo, useEffect } from 'react';

import { produce } from 'immer';
import addDays from 'date-fns/addDays';
import { useRecoilValue } from 'recoil';
import { Column } from '@tanstack/react-table';

import { VStack } from '@ui/layout/Stack';
import { Text } from '@ui/typography/Text';
import { Organization } from '@graphql/types';
import { Radio, RadioGroup } from '@ui/form/Radio';

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
        colorScheme='primary'
        value={filter.value}
        onChange={handleChange}
        isDisabled={!filter.isActive}
      >
        <VStack spacing={2} align='flex-start'>
          <Radio value={week}>
            <Text fontSize='sm'>Next 7 days</Text>
          </Radio>
          <Radio value={month}>
            <Text fontSize='sm'>Next 30 days</Text>
          </Radio>
          <Radio value={quarter}>
            <Text fontSize='sm'>Next 90 days</Text>
          </Radio>
        </VStack>
      </RadioGroup>
    </>
  );
};
