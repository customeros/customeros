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
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useGetUsersQuery } from '@organizations/graphql/getUsers.generated';

import { DebouncedSearchInput } from '../shared';
import { useOwnerFilter, OwnerFilterSelector } from './OwnerFilter.atom';

interface OwnerFilterProps {
  column: Column<Organization>;
  initialFocusRef: RefObject<HTMLInputElement>;
}

const client = getGraphQLClient();

export const OwnerFilter = ({ column, initialFocusRef }: OwnerFilterProps) => {
  const [filter, setFilter] = useOwnerFilter();
  const [searchValue, setSearchValue] = useState('');
  const filterValue = useRecoilValue(OwnerFilterSelector);
  const [_, startTransition] = useTransition();

  const { data } = useGetUsersQuery(client, {
    pagination: { limit: 100, page: 1 },
  });

  const users = useMemo(() => {
    const items = data?.users.content ?? [];
    if (!searchValue) return items;

    return items.filter(
      ({ name, firstName, lastName }) =>
        name?.toLowerCase().includes(searchValue.toLowerCase()) ||
        firstName?.toLowerCase().includes(searchValue.toLowerCase()) ||
        lastName?.toLowerCase().includes(searchValue.toLowerCase()),
    );
  }, [data?.users.content, searchValue]);

  const handleSelect = (value: string) => () => {
    startTransition(() => {
      setFilter((prev) =>
        produce(prev, (draft) => {
          draft.isActive = true;

          if (value === 'ALL') {
            if (draft.value.length === users.length && draft.value.length > 0) {
              draft.value = [];
            } else {
              draft.value = users.map(({ id }) => id);
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
          {users.length > 1 && (
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
                  filter.value.length === users.length
                }
                isIndeterminate={
                  filter.value.length > 0 && filter.value.length < users.length
                }
              >
                <Text fontSize='sm'>
                  {filter.value.length === users.length &&
                  filter.value.length > 0
                    ? 'Deselect All'
                    : 'Select All'}
                </Text>
              </Checkbox>
            </Flex>
          )}
          {users.length > 0 ? (
            users.map(({ id, firstName, lastName, name }) => (
              <Checkbox key={id} value={id} onChange={handleSelect(id)}>
                <Text fontSize='sm' noOfLines={1}>
                  {name
                    ? name
                    : [firstName, lastName].filter(Boolean).join(' ')}
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
