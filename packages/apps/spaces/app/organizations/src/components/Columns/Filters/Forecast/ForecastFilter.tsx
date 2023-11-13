'use client';
import {
  memo,
  useRef,
  useState,
  useEffect,
  forwardRef,
  ChangeEvent,
} from 'react';

import { produce } from 'immer';
import { Column } from '@tanstack/react-table';

import { Flex } from '@ui/layout/Flex';
import { Input } from '@ui/form/Input';
import { Text } from '@ui/typography/Text';
import { Organization } from '@graphql/types';
import { CurrencyDollar } from '@ui/media/icons/CurrencyDollar';
import { InputGroup, InputLeftElement } from '@ui/form/InputGroup';
import {
  RangeSlider,
  RangeSliderTrack,
  RangeSliderThumb,
  RangeSliderFilledTrack,
} from '@ui/form/RangeSlider';

import { useForecastFilter } from './ForecastFilter.atom';
import { FilterHeader, useFilterToggle } from '../shared/FilterHeader';

interface ForecastFilterProps {
  initialFocusRef: React.RefObject<HTMLInputElement>;
  onFilterValueChange?: Column<Organization>['setFilterValue'];
}

export const ForecastFilter = memo(
  ({ initialFocusRef, onFilterValueChange }: ForecastFilterProps) => {
    const [filter, setFilter] = useForecastFilter();
    const [displayValue, setDisplayValue] = useState<[number, number]>(
      () => filter.value,
    );
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

    const handleChange = (value: [number, number]) => {
      setFilter((prev) => {
        const next = produce(prev, (draft) => {
          draft.isActive = true;
          draft.value = value;
        });

        toggle.setIsActive(next.isActive);

        return next;
      });
    };

    const handleInputDisplayChange = (index: number) => (value: number) => {
      setDisplayValue((prev) =>
        produce(prev, (draft) => {
          draft[index] = value;
        }),
      );
    };

    const handleInputChange = (index: number) => (value: number) => {
      setFilter((prev) => {
        const next = produce(prev, (draft) => {
          draft.isActive = true;
          draft.value[index] = value;
        });

        toggle.setIsActive(next.isActive);

        return next;
      });
    };

    const handleDragChange = (value: [number, number]) => {
      setDisplayValue(value);
    };

    useEffect(() => {
      onFilterValueChange?.(filter.isActive ? filter.value : undefined);
    }, [filter.isActive, filter.value[0], filter.value[1]]);

    return (
      <>
        <FilterHeader
          isChecked={toggle.isActive}
          onToggle={toggle.handleChange}
          onDisplayChange={toggle.handleClick}
        />

        <Flex justify='space-between' mb='8' gap='2'>
          <Flex flexDir='column' flex='1'>
            <Text fontSize='sm' fontWeight='medium'>
              Min
            </Text>
            <DebouncedNumberInput
              min={0}
              max={displayValue[1]}
              ref={initialFocusRef}
              value={displayValue[0]}
              onChange={handleInputChange(0)}
              onDisplayChange={handleInputDisplayChange(0)}
            />
          </Flex>
          <Flex flexDir='column' flex='1'>
            <Text fontSize='sm' fontWeight='medium'>
              Max
            </Text>
            <DebouncedNumberInput
              min={displayValue[0]}
              value={displayValue[1]}
              onChange={handleInputChange(1)}
              onDisplayChange={handleInputDisplayChange(1)}
            />
          </Flex>
        </Flex>

        <Flex px='1'>
          <RangeSlider
            min={0}
            step={10}
            max={1000 * 10}
            colorScheme='gray'
            value={displayValue}
            onChangeEnd={handleChange}
            onChange={handleDragChange}
          >
            <RangeSliderTrack bg='gray.200' h='2px'>
              <RangeSliderFilledTrack h='2px' bg='gray.400' />
            </RangeSliderTrack>
            <RangeSliderThumb
              index={0}
              boxSize='5'
              border='2px solid'
              borderColor='gray.400'
            />
            <RangeSliderThumb
              index={1}
              boxSize='5'
              border='2px solid'
              borderColor='gray.400'
            />
          </RangeSlider>
        </Flex>
      </>
    );
  },
);

interface DebouncedNumberInputProps {
  min?: number;
  max?: number;
  value: number;
  placeholder?: string;
  defaultValue?: number;
  onChange: (value: number) => void;
  onDisplayChange: (value: number) => void;
}

export const DebouncedNumberInput = memo(
  forwardRef<HTMLInputElement, DebouncedNumberInputProps>(
    (
      {
        min,
        max,
        onChange,
        placeholder,
        value: _value,
        onDisplayChange,
        defaultValue = 0,
      },
      ref,
    ) => {
      const timeout = useRef<NodeJS.Timeout>();

      const handleChange = (e: ChangeEvent<HTMLInputElement>) => {
        const value = e.target.valueAsNumber;
        onDisplayChange(value);

        if (timeout.current) {
          clearTimeout(timeout.current);
        }

        timeout.current = setTimeout(() => {
          if (max && value > max) {
            onChange(max);

            return;
          }

          if (min && value < min) {
            onChange(min);

            return;
          }

          onChange(e.target.valueAsNumber);
        }, 250);
      };

      useEffect(() => {
        return () => {
          timeout.current && clearTimeout(timeout.current);
        };
      }, []);

      return (
        <InputGroup>
          <InputLeftElement w='fit-content'>
            <CurrencyDollar color='gray.500' />
          </InputLeftElement>
          <Input
            pl='6'
            ref={ref}
            min={min}
            max={max}
            type='number'
            value={_value}
            variant='flushed'
            onChange={handleChange}
            placeholder={placeholder}
          />
        </InputGroup>
      );
    },
  ),
);
