'use client';
import { useMemo, useState, useEffect, RefObject } from 'react';

import { produce } from 'immer';
import { useRecoilValue } from 'recoil';
import difference from 'lodash/difference';
import intersection from 'lodash/intersection';
import { Column } from '@tanstack/react-table';

import { Flex } from '@ui/layout/Flex';
import { VStack } from '@ui/layout/Stack';
import { Text } from '@ui/typography/Text';
import { Organization } from '@graphql/types';
import { Tumbleweed } from '@ui/media/icons/Tumbleweed';
import { Checkbox, CheckboxGroup } from '@ui/form/Checkbox';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useGetUsersQuery } from '@organizations/graphql/getUsers.generated';

import { useOwnerFilter, OwnerFilterSelector } from './OwnerFilter.atom';
import { FilterHeader, useFilterToggle, DebouncedSearchInput } from '../shared';

interface OwnerFilterProps {
  initialFocusRef: RefObject<HTMLInputElement>;
  onFilterValueChange?: Column<Organization>['setFilterValue'];
}

export const OwnerFilter = ({
  initialFocusRef,
  onFilterValueChange,
}: OwnerFilterProps) => {
  const client = getGraphQLClient();
  const [filter, setFilter] = useOwnerFilter();
  const [searchValue, setSearchValue] = useState('');
  const filterValue = useRecoilValue(OwnerFilterSelector);

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

  const { data } = useGetUsersQuery(client, {
    pagination: { limit: 100, page: 1 },
  });

  const users = useMemo(() => {
    const items = [
      { value: '__EMPTY__', label: 'Unknown' },
      ...(data?.users.content.map(({ id, name, firstName, lastName }) => ({
        value: id,
        label: name ? name : [firstName, lastName].filter(Boolean).join(' '),
      })) ?? []),
    ];

    if (!searchValue) return items;

    return items.filter(({ label }) =>
      label?.toLowerCase().includes(searchValue.toLowerCase()),
    );
  }, [data?.users.content, searchValue]);

  const userIds = users.map(({ value }) => value);
  const isAllSelected =
    intersection(filter.value, userIds).length === users.length &&
    users.length > 0;

  const handleSelectAll = () => {
    setFilter((prev) => {
      const next = produce(prev, (draft) => {
        draft.isActive = true;

        if (isAllSelected) {
          draft.value = draft.value.filter((item) => !userIds.includes(item));

          if (draft.value.length === 0) {
            draft.isActive = false;
          }

          return;
        }

        if (searchValue) {
          draft.value = [...userIds, ...difference(draft.value, userIds)];

          return;
        }

        draft.value = userIds;
      });

      toggle.setIsActive(next.isActive);

      return next;
    });
  };

  const handleSelect = (value: string) => () => {
    setFilter((prev) => {
      const next = produce(prev, (draft) => {
        draft.isActive = true;

        if (draft.value.includes(value)) {
          draft.value = draft.value.filter((item) => item !== value);
          if (draft.value.length === 0) {
            draft.isActive = false;
          }
        } else {
          draft.value.push(value);
        }
      });

      toggle.setIsActive(next.isActive);

      return next;
    });
  };

  useEffect(() => {
    onFilterValueChange?.(filterValue.isActive ? filterValue : undefined);
  }, [filterValue.value.length, filterValue.isActive, filterValue.showEmpty]);

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
        onDisplayChange={(v) => setSearchValue(v)}
      />

      <VStack
        spacing={2}
        align='flex-start'
        maxH='13rem'
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
              gap='2'
              flexDir='column'
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
                  {isAllSelected
                    ? 'Deselect all'
                    : 'Select all' +
                      (searchValue && users.length > 2
                        ? ` ${users.length}`
                        : '')}
                </Text>
              </Checkbox>
            </Flex>
          )}

          {users.length > 0 ? (
            users.map(({ value, label }) => (
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
            <Flex w='full' justify='center' align='center' flexDir='column'>
              <Tumbleweed
                mr='10'
                boxSize='8'
                color='gray.400'
                alignSelf='flex-end'
              />
              <Text fontSize='sm' color='gray.500'>
                Empty here in <b>No Resultsville</b>
              </Text>
            </Flex>
          )}
        </CheckboxGroup>
      </VStack>
    </>
  );
};
