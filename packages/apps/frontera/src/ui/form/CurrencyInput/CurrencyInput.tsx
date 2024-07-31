import React, { ChangeEvent } from 'react';

import { Input } from '../Input/Input';
import {
  InputGroup,
  LeftElement,
  RightElement,
} from '../InputGroup/InputGroup';

export interface CurrencyInputProps
  extends Omit<React.InputHTMLAttributes<HTMLInputElement>, 'onChange'> {
  value: string;
  label?: string;
  placeholder?: string;
  isLabelVisible?: boolean;
  leftElement?: React.ReactNode;
  rightElement?: React.ReactNode;
  onChange?: (value: string) => void;
  parseValue?: (val: string) => string;
  formatValue?: (val: string) => string;
  labelProps?: React.LabelHTMLAttributes<HTMLLabelElement>;
}

export const CurrencyInput = React.forwardRef<
  HTMLInputElement,
  CurrencyInputProps
>(
  (
    {
      label,
      leftElement,
      rightElement,
      value,
      labelProps,
      onChange,
      formatValue,
      placeholder,
      parseValue,
    },
    ref,
  ) => {
    const handleValueChange = (e: ChangeEvent<HTMLInputElement>) => {
      const valueString = e.target.value;

      // handle weird case of blurring the field with an empty value
      if (valueString === '-9007199254740991') {
        onChange?.('');

        return;
      }

      if (parseValue) {
        onChange?.(parseValue(valueString));

        return;
      }
      onChange?.(valueString);
    };

    return (
      <div>
        <label className='font-semibold text-sm mb-[-4px]' {...labelProps}>
          {label}
        </label>

        <InputGroup>
          {leftElement && (
            <LeftElement className='size-4'>{leftElement}</LeftElement>
          )}
          <Input
            ref={ref}
            type='number'
            variant='flushed'
            placeholder={placeholder}
            onChange={handleValueChange}
            value={formatValue ? formatValue(value) : value}
            className='border-transparent focus:border-0 hover:border-transparent'
          />

          {rightElement && <RightElement>{rightElement}</RightElement>}
        </InputGroup>
      </div>
    );
  },
);
