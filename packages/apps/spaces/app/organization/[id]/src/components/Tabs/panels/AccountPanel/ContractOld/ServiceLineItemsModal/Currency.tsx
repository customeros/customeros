import { useIMask } from 'react-imask';
import React, { useEffect } from 'react';

import { InputElement } from 'imask';

import { Input, InputProps } from '@ui/form/Input/Input';

interface CurrencyProps extends InputProps {
  label?: string;
  isLabelVisible?: boolean;
  currency?: string | null;
  onValueChange?: (value: string) => void;
  labelProps?: React.HTMLProps<HTMLLabelElement>;
}

export const Currency = ({
  isLabelVisible,
  label,
  value = '',
  labelProps,
  onValueChange,
  currency,
  ...props
}: CurrencyProps) => {
  const { ref, setUnmaskedValue, unmaskedValue, setValue } = useIMask({
    mask: `${currency}num`,

    blocks: {
      num: {
        mask: Number,
        thousandsSeparator: ',',
        radix: '.',
        mapToRadix: ['.'],
        normalizeZeros: true,
        padFractionalZeros: true,
      },
    },
  });

  useEffect(() => {
    if (unmaskedValue) {
      onValueChange?.(`${unmaskedValue ?? ''}`);
    }
  }, [unmaskedValue]);

  useEffect(() => {
    if (value) {
      setUnmaskedValue(`${value}`);
    }
  }, []);

  const handleValueChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const isNegative = e.target.value?.indexOf('-') !== -1;
    if (isNegative) {
      setValue('0');

      return;
    }
    setValue(e.target.value);
  };

  const handleFocusOnClick = () => {
    // fixes cos-2594 - focus masked input on single click
    (ref?.current as InputElement)?.focus();
    if (unmaskedValue === '0') {
      (ref?.current as InputElement)?.setSelectionRange(1, 5);
    }
  };

  return (
    <div>
      <label {...labelProps}>{label}</label>

      <Input
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        ref={ref as any}
        {...props}
        onChange={handleValueChange}
        onClick={handleFocusOnClick}
        autoComplete='off'
      />
    </div>
  );
};
