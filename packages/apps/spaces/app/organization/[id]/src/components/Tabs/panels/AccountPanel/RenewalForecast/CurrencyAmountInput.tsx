import { useState } from 'react';

import {
  NumberInput,
  NumberInputField,
  NumberInputProps,
} from '@ui/form/NumberInput';

interface CurrencyAmountInputProps
  extends Omit<NumberInputProps, 'onChange' | 'value'> {
  value?: number | null;
  onChange?: (value: number) => void;
}

export const CurrencyAmountInput = ({
  onChange,
  value,
  ...props
}: CurrencyAmountInputProps) => {
  const [_value, setValue] = useState<string>(() => `${value}`);

  const handleChange = (valueAsString: string, valueAsNumber: number) => {
    setValue(parse(valueAsString));
    // console.log('val', valueAsNumber, valueAsString);
    onChange?.(valueAsNumber);
  };

  return (
    <NumberInput value={format(_value)} onChange={handleChange} {...props}>
      <NumberInputField />
    </NumberInput>
  );
};

function parse(val: string) {
  return val.replace(/^\$/, '');
}

function format(number: string) {
  if (!number) {
    return '$';
  }
  // Split the number into integer and decimal parts
  const parts = number.split('.');

  // Format the integer part with commas
  parts[0] = parts[0].replace(/\B(?=(\d{3})+(?!\d))/g, ',');

  // Join the integer and decimal parts back together
  if (number.endsWith('.') && parts.length === 1) {
    return number + '.';
  }

  return '$' + parts.join('.');
}
