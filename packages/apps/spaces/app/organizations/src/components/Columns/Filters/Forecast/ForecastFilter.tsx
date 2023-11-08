'use client';
import { useState, useEffect, useTransition } from 'react';

import { produce } from 'immer';
import { useRecoilValue } from 'recoil';
import { Column } from '@tanstack/react-table';

import { Flex } from '@ui/layout/Flex';
import { Switch } from '@ui/form/Switch';
import { Text } from '@ui/typography/Text';
import { CurrencyDollar } from '@ui/media/icons/CurrencyDollar';
import {
  RangeSlider,
  RangeSliderTrack,
  RangeSliderThumb,
  RangeSliderFilledTrack,
} from '@ui/form/RangeSlider';

import {
  useForecastFilter,
  ForecastFilterSelector,
} from './ForecastFilter.atom';

interface ForecastFilterProps<T> {
  column: Column<T>;
}

export const ForecastFilter = <T,>({ column }: ForecastFilterProps<T>) => {
  const [filter, setFilter] = useForecastFilter();
  const filterValue = useRecoilValue(ForecastFilterSelector);
  const [_, startTransition] = useTransition();
  const [displayValue, setDisplayValue] = useState<[number, number]>(
    () => filter.value,
  );

  const handleChange = (value: [number, number]) => {
    startTransition(() => {
      setFilter((prev) =>
        produce(prev, (draft) => {
          draft.isActive = true;
          draft.value = value;
        }),
      );
    });
  };

  const handleDragChange = (value: [number, number]) => {
    setDisplayValue(value);
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
  }, [filterValue.value[0], filterValue.value[1], filterValue.isActive]);

  return (
    <>
      <Flex
        mb='3'
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

      <Flex justify='space-between' mb='8'>
        <Flex flexDir='column' flex='1' gap='1'>
          <Text fontSize='sm' fontWeight='medium'>
            Min
          </Text>
          <Flex align='center' gap='3'>
            <CurrencyDollar color='gray.500' />
            <Text fontSize='sm'>{displayValue[0]}</Text>
          </Flex>
        </Flex>
        <Flex flexDir='column' flex='1' gap='1'>
          <Text fontSize='sm' fontWeight='medium'>
            Max
          </Text>
          <Flex align='center' gap='3'>
            <CurrencyDollar color='gray.500' />
            <Text fontSize='sm'>{displayValue[1]}</Text>
          </Flex>
        </Flex>
      </Flex>

      <Flex px='1'>
        <RangeSlider
          min={0}
          max={1000 * 10}
          colorScheme='gray'
          onChangeEnd={handleChange}
          onChange={handleDragChange}
          defaultValue={filter.value}
        >
          <RangeSliderTrack bg='gray.200' h='2px'>
            <RangeSliderFilledTrack h='2px' bg='gray.400' />
          </RangeSliderTrack>
          <RangeSliderThumb
            boxSize='5'
            border='2px solid'
            borderColor='gray.400'
            index={0}
          />
          <RangeSliderThumb
            boxSize='5'
            border='2px solid'
            borderColor='gray.400'
            index={1}
          />
        </RangeSlider>
      </Flex>
    </>
  );
};
