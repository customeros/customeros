import { useState, useEffect } from 'react';

import { Input, InputProps } from '@ui/form/Input';

interface RangeSelectorProps extends Omit<InputProps, 'onChange'> {
  filter: string;
  years?: boolean;
  placeholder: string;
  onChange: (values: [number | string, (number | string)?]) => void;
}

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
        onChange={(e) =>
          setMinValue(e.target.value ? Number(e.target.value) : '')
        }
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
        type='number'
        placeholder='Max'
        className='w-[50px]'
        value={maxValue}
        onChange={(e) =>
          setMaxValue(e.target.value ? Number(e.target.value) : '')
        }
        {...rest}
      />
      {filter === 'between' && years && <span>yrs</span>}
    </div>
  );
};
