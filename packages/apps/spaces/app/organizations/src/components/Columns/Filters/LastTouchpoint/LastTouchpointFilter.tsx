'use client';
import { useMemo, useEffect } from 'react';

import { produce } from 'immer';
import subDays from 'date-fns/subDays';
import { useRecoilValue } from 'recoil';
import { Column } from '@tanstack/react-table';

import { Flex } from '@ui/layout/Flex';
import { VStack } from '@ui/layout/Stack';
import { Text } from '@ui/typography/Text';
import { Radio, RadioGroup } from '@ui/form/Radio';
import { Checkbox, CheckboxGroup } from '@ui/form/Checkbox';
import { Organization, LastTouchpointType } from '@graphql/types';

import { touchpoints } from './util';
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

  const handleSelectAll = () => {
    setFilter((prev) => {
      const next = produce(prev, (draft) => {
        if (isAllSelected) {
          draft.isActive = false;
          draft.value = [];
        } else {
          draft.isActive = true;
          draft.value = touchpoints.map(({ value }) => value);
        }
      });

      toggle.setIsActive(next.isActive);

      return next;
    });
  };

  const handleSelect = (value: LastTouchpointType) => () => {
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
  }, [filterValue.value.length, filterValue.isActive, filterValue.after]);

  return (
    <>
      <FilterHeader
        isChecked={toggle.isActive}
        onToggle={toggle.handleChange}
        onDisplayChange={toggle.handleClick}
      />

      <RadioGroup
        pb='2'
        name='last-touchpoint-before'
        colorScheme='primary'
        value={filter.after}
        borderBottom='1px solid'
        borderBottomColor='gray.200'
        onChange={handleDateChange}
        isDisabled={!filter.isActive}
      >
        <VStack spacing={2} align='flex-start'>
          <Radio value={week}>
            <Text fontSize='sm'>Last 7 days</Text>
          </Radio>
          <Radio value={month}>
            <Text fontSize='sm'>Last 30 days</Text>
          </Radio>
          <Radio value={quarter}>
            <Text fontSize='sm'>Last 90 days</Text>
          </Radio>
        </VStack>
      </RadioGroup>

      <VStack
        spacing={2}
        align='flex-start'
        mt='2'
        px='4px'
        mx='-4px'
        position='relative'
        overflowX='hidden'
        overflowY='auto'
      >
        <CheckboxGroup size='md' value={filter.value}>
          <Flex
            top='0'
            w='full'
            zIndex='10'
            bg='white'
            position='sticky'
            borderBottom='1px solid'
            borderColor='gray.200'
            pb='2'
          >
            <Checkbox
              top='0'
              zIndex='10'
              isChecked={isAllSelected}
              onChange={handleSelectAll}
            >
              <Text fontSize='sm'>
                {isAllSelected ? 'Deselect all' : 'Select all'}
              </Text>
            </Checkbox>
          </Flex>
          {touchpoints.map(({ label, value }) => (
            <Checkbox key={value} value={value} onChange={handleSelect(value)}>
              <Text fontSize='sm' noOfLines={1}>
                {label}
              </Text>
            </Checkbox>
          ))}
        </CheckboxGroup>
      </VStack>
    </>
  );
};
