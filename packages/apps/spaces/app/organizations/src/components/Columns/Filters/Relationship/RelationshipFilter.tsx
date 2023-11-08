'use client';
import { useEffect, useTransition } from 'react';

import { produce } from 'immer';
import { useRecoilValue } from 'recoil';
import { Column } from '@tanstack/react-table';

import { Flex } from '@ui/layout/Flex';
import { Switch } from '@ui/form/Switch';
import { VStack } from '@ui/layout/Stack';
import { Text } from '@ui/typography/Text';
import { Checkbox, CheckboxGroup } from '@ui/form/Checkbox';

import {
  useRelationshipFilter,
  RelationshipFilterSelector,
} from './RelationshipFilter.atom';

interface RelationshipFilterProps<T> {
  column: Column<T>;
}

export const RelationshipFilter = <T,>({
  column,
}: RelationshipFilterProps<T>) => {
  const [filter, setFilter] = useRelationshipFilter();
  const filterValue = useRecoilValue(RelationshipFilterSelector);
  const [_, startTransition] = useTransition();

  const handleSelect = (value: string) => () => {
    startTransition(() => {
      setFilter((prev) =>
        produce(prev, (draft) => {
          draft.isActive = true;

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
    <CheckboxGroup size='md' value={filter.value}>
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
        <Checkbox value='customer' onChange={handleSelect('customer')}>
          <Text fontSize='sm'>Customer</Text>
        </Checkbox>
        <Checkbox value='prospect' onChange={handleSelect('prospect')}>
          <Text fontSize='sm'>Prospect</Text>
        </Checkbox>
      </VStack>
    </CheckboxGroup>
  );
};
