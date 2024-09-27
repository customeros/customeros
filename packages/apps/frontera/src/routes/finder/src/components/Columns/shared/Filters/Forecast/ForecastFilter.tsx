import { useSearchParams } from 'react-router-dom';
import { useRef, useState, useEffect, forwardRef, ChangeEvent } from 'react';

import { observer } from 'mobx-react-lite';
import { FilterItem } from '@store/types.ts';

import { Input } from '@ui/form/Input/Input.tsx';
import { useStore } from '@shared/hooks/useStore';
import { CurrencyDollar } from '@ui/media/icons/CurrencyDollar.tsx';
import { ColumnViewType, ComparisonOperator } from '@graphql/types';
import { InputGroup, LeftElement } from '@ui/form/InputGroup/InputGroup.tsx';
import {
  RangeSlider,
  RangeSliderTrack,
  RangeSliderThumb,
  RangeSliderFilledTrack,
} from '@ui/form/RangeSlider/RangeSlider.tsx';

import { FilterHeader } from '../abstract/FilterHeader';

interface ForecastFilterProps {
  property?: ColumnViewType;
  hideCurrencySymbol?: boolean;
  initialFocusRef: React.RefObject<HTMLInputElement>;
}

const defaultFilter: FilterItem = {
  property: ColumnViewType.OrganizationsForecastArr,
  value: [0, 10000],
  active: false,
  caseSensitive: false,
  includeEmpty: false,
  operation: ComparisonOperator.Contains,
};

export const ForecastFilter = observer(
  ({ initialFocusRef, property, hideCurrencySymbol }: ForecastFilterProps) => {
    const [searchParams] = useSearchParams();
    const preset = searchParams.get('preset');

    const store = useStore();
    const tableViewDef = store.tableViewDefs.getById(preset ?? '');
    const filter = tableViewDef?.getFilter(
      property || defaultFilter.property,
    ) ?? { ...defaultFilter, property: property || defaultFilter.property };

    const toggle = () => {
      tableViewDef?.toggleFilter(filter);
    };

    const [displayValue, setDisplayValue] = useState<[number, number]>(
      filter.value,
    );

    const handleChange = (value: [number, number]) => {
      tableViewDef?.setFilter({
        ...filter,
        value,
        active: filter.active || true,
      });
    };

    const handleInputChange = (index: number) => (value: number) => {
      const nextValue: [number, number] = [...displayValue];

      nextValue[index] = value;

      setDisplayValue(nextValue);

      tableViewDef?.setFilter({
        ...filter,
        value: nextValue,
        active: filter.active || true,
      });
    };

    const handleDragChange = (value: [number, number]) => {
      setDisplayValue(value);
    };

    return (
      <>
        <FilterHeader
          onToggle={toggle}
          onDisplayChange={() => {}}
          isChecked={filter.active ?? false}
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
              hideCurrencySymbol={hideCurrencySymbol}
            />
          </div>
          <div className='flex flex-col flex-1'>
            <span className='text-sm font-medium'>Max</span>
            <DebouncedNumberInput
              min={displayValue[0]}
              value={displayValue[1]}
              onChange={handleInputChange(1)}
              hideCurrencySymbol={hideCurrencySymbol}
            />
          </div>
        </div>

        <div className='flex px-1'>
          <RangeSlider
            min={0}
            step={10}
            max={1000 * 10}
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
  hideCurrencySymbol?: boolean;
  onChange: (value: number) => void;
}

export const DebouncedNumberInput = forwardRef<
  HTMLInputElement,
  DebouncedNumberInputProps
>(
  (
    {
      min,
      max,
      onChange,
      placeholder,
      value,
      hideCurrencySymbol,
      defaultValue = 0,
    },
    ref,
  ) => {
    const [displayValue, setDisplayValue] = useState(value);
    const timeout = useRef<NodeJS.Timeout>();

    const handleChange = (e: ChangeEvent<HTMLInputElement>) => {
      const value = e.target.valueAsNumber;

      setDisplayValue(value);

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

        onChange(value);
      }, 250);
    };

    useEffect(() => {
      return () => {
        timeout.current && clearTimeout(timeout.current);
      };
    }, []);

    useEffect(() => {
      setDisplayValue(value);
    }, [value]);

    return (
      <InputGroup>
        {!hideCurrencySymbol && (
          <LeftElement className='mb-1'>
            <CurrencyDollar className='text-gray-500' />
          </LeftElement>
        )}

        <Input
          ref={ref}
          min={min}
          max={max}
          type='number'
          variant='flushed'
          value={displayValue}
          onChange={handleChange}
          placeholder={placeholder}
          defaultValue={defaultValue}
          className='border-transparent focus:border-transparent focus:hover:border-transparent hover:border-transparent'
        />
      </InputGroup>
    );
  },
);
