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
import { RenewalLikelihoodProbability } from '@graphql/types';

import {
  useRenewalLikelihoodFilter,
  RenewalLikelihoodFilterSelector,
} from './RenewalLikelihoodFilter.atom';

interface RenewalLikelihoodFilterProps<T> {
  column: Column<T>;
}

export const RenewalLikelihoodFilter = <T,>({
  column,
}: RenewalLikelihoodFilterProps<T>) => {
  const [filter, setFilter] = useRenewalLikelihoodFilter();
  const filterValue = useRecoilValue(RenewalLikelihoodFilterSelector);
  const [_, startTransition] = useTransition();

  const handleSelect = (value: RenewalLikelihoodProbability) => () => {
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
        <Checkbox
          value={RenewalLikelihoodProbability.High}
          onChange={handleSelect(RenewalLikelihoodProbability.High)}
        >
          <Text fontSize='sm'>High</Text>
        </Checkbox>
        <Checkbox
          value={RenewalLikelihoodProbability.Medium}
          onChange={handleSelect(RenewalLikelihoodProbability.Medium)}
        >
          <Text fontSize='sm'>Medium</Text>
        </Checkbox>
        <Checkbox
          value={RenewalLikelihoodProbability.Low}
          onChange={handleSelect(RenewalLikelihoodProbability.Low)}
        >
          <Text fontSize='sm'>Low</Text>
        </Checkbox>
        <Checkbox
          value={RenewalLikelihoodProbability.Zero}
          onChange={handleSelect(RenewalLikelihoodProbability.Zero)}
        >
          <Text fontSize='sm'>Zero</Text>
        </Checkbox>
      </VStack>
    </CheckboxGroup>
  );
};
