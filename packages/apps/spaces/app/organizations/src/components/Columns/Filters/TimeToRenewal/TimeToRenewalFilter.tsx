'use client';
import { useMemo, useEffect, useTransition } from 'react';

import { produce } from 'immer';
import addDays from 'date-fns/addDays';
import { useRecoilValue } from 'recoil';
import { Column } from '@tanstack/react-table';

import { Flex } from '@ui/layout/Flex';
import { Switch } from '@ui/form/Switch';
import { VStack } from '@ui/layout/Stack';
import { Text } from '@ui/typography/Text';
import { Radio, RadioGroup } from '@ui/form/Radio';

import {
  useTimeToRenewalFilter,
  TimeToRenewalFilterSelector,
} from './TimeToRenewalFilter.atom';

interface TimeToRenewalProps<T> {
  column: Column<T>;
}

export const TimeToRenewalFilter = <T,>({ column }: TimeToRenewalProps<T>) => {
  const [filter, setFilter] = useTimeToRenewalFilter();
  const filterValue = useRecoilValue(TimeToRenewalFilterSelector);
  const [_, startTransition] = useTransition();

  const [week, month, quarter] = useMemo(
    () =>
      [7, 30, 90].map((value) => {
        return addDays(new Date(), value).toISOString().split('T')[0];
      }),
    [],
  );

  const handleChange = (value: string) => {
    startTransition(() => {
      setFilter((prev) =>
        produce(prev, (draft) => {
          draft.isActive = true;
          draft.value = value;
        }),
      );
    });
  };

  const handleToggle = () => {
    startTransition(() => {
      setFilter((prev) =>
        produce(prev, (draft) => {
          draft.isActive = !draft.isActive;
        }),
      );
    });
  };

  useEffect(() => {
    column.setFilterValue(filterValue.isActive ? filterValue.value : undefined);
  }, [filterValue.value, filterValue.isActive]);

  return (
    <RadioGroup
      name='timeToRenewal'
      colorScheme='primary'
      value={filter.value}
      onChange={handleChange}
      isDisabled={!filter.isActive}
    >
      <Flex
        mb='2'
        flexDir='row'
        alignItems='center'
        justifyContent='space-between'
      >
        <Text fontSize='sm' fontWeight='medium'>
          Filter
        </Text>
        <Switch
          size='sm'
          colorScheme='primary'
          onChange={handleToggle}
          isChecked={filter.isActive}
        />
      </Flex>
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
  );
};
