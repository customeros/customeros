'use client';
import { useMemo, useState, useEffect, RefObject, useTransition } from 'react';

import { produce } from 'immer';
import { useRecoilValue } from 'recoil';
import { Column } from '@tanstack/react-table';

import { Flex } from '@ui/layout/Flex';
import { Switch } from '@ui/form/Switch';
import { VStack } from '@ui/layout/Stack';
import { Text } from '@ui/typography/Text';
import { Organization } from '@graphql/types';
import { Checkbox, CheckboxGroup } from '@ui/form/Checkbox';

import { DebouncedSearchInput } from '../shared';
import { TouchPoint, touchpoints } from './util';
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
  const [_, startTransition] = useTransition();

  const touchpointOptions = useMemo(() => {
    if (!searchValue) return touchpoints;

    return touchpoints.filter(({ label }) =>
      label.toLowerCase().includes(searchValue.toLowerCase()),
    );
  }, [searchValue]);

  const handleSelect = (value: TouchPoint | 'ALL') => () => {
    startTransition(() => {
      setFilter((prev) =>
        produce(prev, (draft) => {
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
  }, [filterValue.value.length, filterValue.isActive]);

  return (
    <>
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
