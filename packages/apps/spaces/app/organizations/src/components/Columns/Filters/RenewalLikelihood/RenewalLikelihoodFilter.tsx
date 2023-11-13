'use client';
import { useEffect } from 'react';

import { produce } from 'immer';
import { useRecoilValue } from 'recoil';
import { Column } from '@tanstack/react-table';

import { VStack } from '@ui/layout/Stack';
import { Text } from '@ui/typography/Text';
import { Checkbox, CheckboxGroup } from '@ui/form/Checkbox';
import { RenewalLikelihoodProbability } from '@graphql/types';

import { FilterHeader, useFilterToggle } from '../shared';
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

  const handleSelect = (value: RenewalLikelihoodProbability) => () => {
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

  useEffect(() => {
    column.setFilterValue(filterValue.isActive ? filterValue.value : undefined);
  }, [filterValue.value.length, filterValue.isActive]);

  return (
    <CheckboxGroup size='md' value={filter.value}>
      <FilterHeader
        isChecked={toggle.isActive}
        onToggle={toggle.handleChange}
        onDisplayChange={toggle.handleClick}
      />
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
