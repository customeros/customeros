'use client';
import { useEffect, RefObject, useCallback, useTransition } from 'react';

import { produce } from 'immer';
import { useRecoilValue } from 'recoil';
import { Column } from '@tanstack/react-table';

import { Flex } from '@ui/layout/Flex';
import { Switch } from '@ui/form/Switch';
import { Text } from '@ui/typography/Text';
import { Organization } from '@graphql/types';

import { DebouncedSearchInput } from '../shared';
import { useWebsiteFilter, WebsiteFilterSelector } from './WebsiteFilter.atom';

interface WebsiteFilterProps {
  column: Column<Organization>;
  initialFocusRef: RefObject<HTMLInputElement>;
}

export const WebsiteFilter = ({
  column,
  initialFocusRef,
}: WebsiteFilterProps) => {
  const [filter, setFilter] = useWebsiteFilter();
  const filterValue = useRecoilValue(WebsiteFilterSelector);
  const [_, startTransition] = useTransition();

  const handleChange = useCallback(
    (value: string) => {
      startTransition(() => {
        setFilter((prev) =>
          produce(prev, (draft) => {
            draft.value = value;
            if (!value) {
              draft.isActive = false;
            } else {
              draft.isActive = true;
            }
          }),
        );
      });
    },
    [setFilter],
  );

  const handleToggle = useCallback(() => {
    startTransition(() => {
      setFilter((prev) =>
        produce(prev, (draft) => {
          draft.isActive = !draft.isActive;
        }),
      );
    });
  }, [setFilter]);

  useEffect(() => {
    column.setFilterValue(filterValue.isActive ? filterValue.value : undefined);
  }, [filterValue.value, filterValue.isActive]);

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
        ref={initialFocusRef}
        value={filter.value}
        onChange={handleChange}
      />
    </>
  );
};
