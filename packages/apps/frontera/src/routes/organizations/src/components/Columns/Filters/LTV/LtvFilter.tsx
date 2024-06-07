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

import { Input } from '@ui/form/Input/Input';
import { Organization } from '@graphql/types';
import { CurrencyDollar } from '@ui/media/icons/CurrencyDollar';
import { InputGroup, LeftElement } from '@ui/form/InputGroup/InputGroup';
import {
  RangeSlider,
  RangeSliderTrack,
  RangeSliderThumb,
  RangeSliderFilledTrack,
} from '@ui/form/RangeSlider/RangeSlider';

import { useForecastFilter } from './LtvFilter.atom';
import { FilterHeader, useFilterToggle } from '../shared/FilterHeader';

interface ForecastFilterProps {
  initialFocusRef: React.RefObject<HTMLInputElement>;
  onFilterValueChange?: Column<Organization>['setFilterValue'];
}

export const LtvFilter = memo(
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

      setDisplayValue(filter.value);
    }, [filter.isActive, filter.value[0], filter.value[1]]);

    return (
      <>
        <FilterHeader
          isChecked={toggle.isActive}
          onToggle={toggle.handleChange}
          onDisplayChange={toggle.handleClick}
        />

        <div className='flex justify-between mb-8 gap-2'>
          <div className='flex flex-col flex-1'>
            <span className='text-sm font-medium'>Min</span>
            <DebouncedNumberInput
              min={0}
              max={displayValue[1]}
              ref={initialFocusRef}
              value={displayValue[0]}
              onChange={handleInputChange(0)}
              onDisplayChange={handleInputDisplayChange(0)}
            />
          </div>
          <div className='flex flex-col flex-1'>
            <span className='text-sm font-medium'>Max</span>
            <DebouncedNumberInput
              min={displayValue[0]}
              value={displayValue[1]}
              onChange={handleInputChange(1)}
              onDisplayChange={handleInputDisplayChange(1)}
            />
          </div>
        </div>

        <div className='flex px-1'>
          <RangeSlider
            min={0}
            step={10}
            max={1000}
            value={displayValue}
            onValueCommit={handleChange}
            onValueChange={handleDragChange}
          >
            <RangeSliderTrack className='bg-gray-200 h-[2px]'>
              <RangeSliderFilledTrack className='h-[2px] bg-gray-400' />
            </RangeSliderTrack>
            <RangeSliderThumb className='border-2 border-gray-400' />
            <RangeSliderThumb className='border-2 border-gray-400' />
          </RangeSlider>
        </div>
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
          <LeftElement className='mb-1'>
            <CurrencyDollar className='text-gray-500' />
          </LeftElement>
          <Input
            className='border-transparent focus:border-0 hover:border-transparent'
            ref={ref}
            min={min}
            max={max}
            type='number'
            value={_value}
            variant='flushed'
            onChange={handleChange}
            placeholder={placeholder}
            defaultValue={defaultValue}
          />
        </InputGroup>
      );
    },
  ),
);
