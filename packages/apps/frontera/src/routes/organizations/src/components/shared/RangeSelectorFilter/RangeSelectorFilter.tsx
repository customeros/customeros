import { useState, useEffect } from 'react';

import { Input, InputProps } from '@ui/form/Input';

interface RangeSelectorProps extends Omit<InputProps, 'onChange'> {
  filter: string;
  years?: boolean;
  placeholder: string;
  onChange: (values: [number?, number?]) => void;
}

export const RangeSelector = ({
  filter,
  placeholder,
  onChange,
  years,
  ...rest
}: RangeSelectorProps) => {
  const [minValue, setMinValue] = useState<string>('');
  const [maxValue, setMaxValue] = useState<string>('');

  useEffect(() => {
    let parsedMinValue: number | undefined = undefined;
    let parsedMaxValue: number | undefined = undefined;

    if (minValue !== '') {
      parsedMinValue = years
        ? new Date().getFullYear() - Number(minValue)
        : Number(minValue);
    }

    if (maxValue !== '') {
      parsedMaxValue = years
        ? new Date().getFullYear() - Number(maxValue)
        : Number(maxValue);
    }

    if (filter === 'between') {
      onChange([parsedMinValue, parsedMaxValue]);
    } else {
      onChange([parsedMinValue]);
    }
  }, [minValue, maxValue, filter, onChange, years]);

  return (
    <div className='flex-1 flex items-center'>
      <Input
        variant='unstyled'
        type='number'
        value={minValue}
        placeholder={filter === 'between' ? 'Min' : `${placeholder}`}
        style={{
          width: filter !== 'between' && !years ? '100%' : '50px',
        }}
        onChange={(e) => setMinValue(e.target.value)}
        {...rest}
      />
      {years && <span>yrs</span>}
      <span
        className={years ? 'mx-4' : 'mr-[30px]'}
        style={{
          display: filter === 'between' ? 'block' : 'none',
        }}
      >
        -{' '}
      </span>
      <Input
        style={{
          display: filter === 'between' ? 'block' : 'none',
        }}
        variant='unstyled'
        type='number'
        placeholder='Max'
        className='w-[50px]'
        value={maxValue}
        onChange={(e) => setMaxValue(e.target.value)}
        {...rest}
      />
      {filter === 'between' && <span>yrs</span>}
    </div>
  );
};
