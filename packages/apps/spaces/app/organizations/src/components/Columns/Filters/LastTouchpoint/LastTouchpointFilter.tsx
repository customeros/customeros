'use client';
import { useMemo, useState, useEffect, RefObject } from 'react';

import { produce } from 'immer';
import { useRecoilValue } from 'recoil';
import { Column } from '@tanstack/react-table';

import { Flex } from '@ui/layout/Flex';
import { VStack } from '@ui/layout/Stack';
import { Text } from '@ui/typography/Text';
import { Organization } from '@graphql/types';
import { Checkbox, CheckboxGroup } from '@ui/form/Checkbox';

import { TouchPoint, touchpoints } from './util';
import { FilterHeader, useFilterToggle, DebouncedSearchInput } from '../shared';
import {
  LastTouchpointSelector,
  useLastTouchpointFilter,
} from './LastTouchpointFilter.atom';

interface LastTouchpointFilterProps {
  column: Column<Organization>;
  initialFocusRef: RefObject<HTMLInputElement>;
}

export const LastTouchpointFilter = ({
  column,
  initialFocusRef,
}: LastTouchpointFilterProps) => {
  const [filter, setFilter] = useLastTouchpointFilter();
  const [searchValue, setSearchValue] = useState('');
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

  const touchpointOptions = useMemo(() => {
    if (!searchValue) return touchpoints;

    return touchpoints.filter(({ label }) =>
      label.toLowerCase().includes(searchValue.toLowerCase()),
    );
  }, [searchValue]);

  const handleSelect = (value: TouchPoint | 'ALL') => () => {
    setFilter((prev) => {
      const next = produce(prev, (draft) => {
        draft.isActive = true;

        if (value === 'ALL') {
          if (
            draft.value.length === touchpointOptions.length &&
            draft.value.length > 0
          ) {
            draft.value = [];
          } else {
            draft.value = touchpointOptions.map(({ value }) => value);
          }

          return;
        }

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

  useEffect(() => {
    column.setFilterValue(filterValue.isActive ? filterValue.value : undefined);
  }, [filterValue.value.length, filterValue.isActive]);

  return (
    <>
      <FilterHeader
        isChecked={toggle.isActive}
        onToggle={toggle.handleChange}
        onDisplayChange={toggle.handleClick}
      />

      <DebouncedSearchInput
        value={searchValue}
        ref={initialFocusRef}
        onChange={(v) => setSearchValue(v)}
      />

      <VStack
        spacing={2}
        align='flex-start'
        maxH='11rem'
        mt='2'
        px='4px'
        mx='-4px'
        position='relative'
        overflowX='hidden'
        overflowY='auto'
      >
        <CheckboxGroup size='md' value={filter.value}>
          {touchpointOptions.length > 1 && (
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
                onChange={handleSelect('ALL')}
                isChecked={
                  filter.value.length > 0 &&
                  filter.value.length === touchpointOptions.length
                }
                isIndeterminate={
                  filter.value.length > 0 &&
                  filter.value.length < touchpointOptions.length
                }
              >
                <Text fontSize='sm'>
                  {filter.value.length === touchpointOptions.length &&
                  filter.value.length > 0
                    ? 'Deselect All'
                    : 'Select All'}
                </Text>
              </Checkbox>
            </Flex>
          )}
          {touchpointOptions.length > 0 ? (
            touchpointOptions.map(({ label, value }) => (
              <Checkbox
                key={value}
                value={value}
                onChange={handleSelect(value)}
              >
                <Text fontSize='sm' noOfLines={1}>
                  {label}
                </Text>
              </Checkbox>
            ))
          ) : (
            <Text fontSize='sm' color='gray.500'>
              Empty here in <b>No Results ville</b>
            </Text>
          )}
        </CheckboxGroup>
      </VStack>
    </>
  );
};
