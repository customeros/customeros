import { useState, useEffect } from 'react';

import { Input, InputProps } from '@ui/form/Input';

interface RangeSelectorProps extends Omit<InputProps, 'onChange'> {
  filter: string;
  years?: boolean;
  placeholder: string;
  onChange: (values: [number | string, (number | string)?]) => void;
}

// Helper function to format numbers with commas
const formatNumberWithCommas = (value: string | number | undefined): string => {
  if (value === undefined || value === '') return '';
  const numString = value.toString();

  return numString.replace(/\B(?=(\d{3})+(?!\d))/g, ',');
};

export const RangeSelector = ({
  filter,
  placeholder,
  onChange,
  years = false,
  ...rest
}: RangeSelectorProps) => {
  const [minValue, setMinValue] = useState<number | string | undefined>(
    undefined,
  );
  const [maxValue, setMaxValue] = useState<number | string | undefined>(
    undefined,
  );

  useEffect(() => {
    if (filter === 'between') {
      onChange([minValue as number, maxValue as number]);
    } else {
      onChange([minValue as number]);
    }
  }, [minValue, maxValue, filter]);

  const handleMinChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value.replace(/,/g, '');
    setMinValue(value ? Number(value) : undefined);
  };

  const handleMaxChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value.replace(/,/g, '');
    setMaxValue(value ? Number(value) : undefined);
  };

  return (
    <div className='flex-1 flex items-center'>
      <Input
        variant='unstyled'
        type='text'
        value={years ? minValue : formatNumberWithCommas(minValue)}
        placeholder={filter === 'between' ? 'Min' : `${placeholder}`}
        style={{
          width: filter !== 'between' && !years ? '100%' : '50px',
        }}
        onChange={handleMinChange}
        {...rest}
      />
      {years && (
        <>
          <span>yrs</span>
          <span
            className='mx-4 '
            style={{
              display: filter === 'between' ? 'block' : 'none',
            }}
          >
            -{' '}
          </span>
        </>
      )}
      {!years && (
        <span
          className='mr-[30px]'
          style={{
            display: filter === 'between' ? 'block' : 'none',
          }}
        >
          -{' '}
        </span>
      )}
      <Input
        style={{
          display: filter === 'between' ? 'block' : 'none',
        }}
        variant='unstyled'
        type='text'
        placeholder='Max'
        className='w-[50px]'
        value={years ? maxValue : formatNumberWithCommas(maxValue)}
        onChange={handleMaxChange}
        {...rest}
      />
      {filter === 'between' && years && <span>yrs</span>}
    </div>
  );
};
