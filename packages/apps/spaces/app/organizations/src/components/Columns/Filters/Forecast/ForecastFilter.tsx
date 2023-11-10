'use client';
import {
  memo,
  useRef,
  useState,
  useEffect,
  forwardRef,
  ChangeEvent,
  useCallback,
  useTransition,
} from 'react';

import { produce } from 'immer';
// import { useRecoilValue } from 'recoil';
import { Column } from '@tanstack/react-table';

import { Flex } from '@ui/layout/Flex';
import { Input } from '@ui/form/Input';
import { Switch } from '@ui/form/Switch';
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

import {
  useForecastFilter,
  // ForecastFilterSelector,
} from './ForecastFilter.atom';

interface ForecastFilterProps {
  column: Column<Organization>;
  initialFocusRef: React.RefObject<HTMLInputElement>;
}

export const ForecastFilter = memo(
  ({ column, initialFocusRef }: ForecastFilterProps) => {
    const [filter, setFilter] = useForecastFilter();
    // const filterValue = useRecoilValue(ForecastFilterSelector);
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

    const handleInputDisplayChange = useCallback(
      (index: number) => (value: number) => {
        setDisplayValue((prev) => {
          return produce(prev, (draft) => {
            draft[index] = value;
          });
        });
      },
      [setDisplayValue],
    );

    const handleInputChange = useCallback(
      (index: number) => (value: number) => {
        setFilter((prev) => {
          return produce(prev, (draft) => {
            draft.isActive = true;
            draft.value[index] = value;
          });
        });
      },
      [setFilter],
    );

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

    // investigate why this does not work

    // useEffect(() => {
    //   column.setFilterValue(
    //     filterValue.isActive ? filterValue.value : undefined,
    //   );
    // }, [filterValue.value[0], filterValue.value[1], filterValue.isActive]);

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
  },
);

interface DebouncedNumberInputProps {
  min?: number;
  max?: number;
  value: number;
  placeholder?: string;
  onFocus?: () => void;
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
          if (timeout.current) {
            clearTimeout(timeout.current);
          }
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
