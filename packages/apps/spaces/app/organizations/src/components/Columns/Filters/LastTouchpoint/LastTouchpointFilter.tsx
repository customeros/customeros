'use client';
import { memo, useRef, useMemo, useState, useEffect, useCallback } from 'react';

import { produce } from 'immer';
import subDays from 'date-fns/subDays';
import { useRecoilValue } from 'recoil';
import { Column } from '@tanstack/react-table';

import { Flex } from '@ui/layout/Flex';
import { VStack } from '@ui/layout/Stack';
import { Text } from '@ui/typography/Text';
import { Radio, RadioGroup } from '@ui/form/Radio';
import { CustomCheckbox } from '@ui/form/Checkbox';
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
          <Radio value={allTime}>
            <Text fontSize='sm'>All time</Text>
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
        <Checkboxes
          value={filter.value}
          onCheck={handleSelect}
          onCheckAll={handleSelectAll}
          isAllSelected={isAllSelected}
        />
      </VStack>
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
    const [checked, setChecked] = useState<Record<string, boolean>>(() =>
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
          <CustomCheckbox
            top='0'
            size='md'
            zIndex='10'
            onChange={handleCheckAll}
            isChecked={_isAllChecked}
          >
            <Text fontSize='sm'>
              {_isAllChecked ? 'Deselect all' : 'Select all'}
            </Text>
          </CustomCheckbox>
        </Flex>
        {touchpoints.map(({ label, value }) => (
          <CustomCheckbox
            key={value}
            size='md'
            value={value}
            onChange={handleCheck}
            isChecked={checked[value]}
          >
            <Text fontSize='sm' noOfLines={1}>
              {label}
            </Text>
          </CustomCheckbox>
        ))}
      </>
    );
  },
);
