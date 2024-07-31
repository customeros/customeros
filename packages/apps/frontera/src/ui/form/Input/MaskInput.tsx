import { useIMask } from 'react-imask';
import React, { useEffect } from 'react';

import { Input, InputProps } from '@ui/form/Input/Input';

interface MaskInputProps extends InputProps {
  name?: string;
  label?: string;
  symbol: string;
  value: string | number;
  onValueChange?: (value: string) => void;
  labelProps?: React.LabelHTMLAttributes<HTMLLabelElement>;
}

/**
 * @deprecated Use `<MaskedInput />` instead
 */
export const MaskInput = ({
  labelProps,
  label,
  name,
  value,
  symbol,
  onValueChange,

  ...props
}: MaskInputProps) => {
  const { ref, setUnmaskedValue, unmaskedValue, setValue } = useIMask({
    mask: `num${symbol}`,
    blocks: {
      num: {
        mask: Number,
        scale: 0.01,
        radix: '%',
        padFractionalZeros: true,
        mapToRadix: ['%'],
        min: 0,
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
    (ref?.current as HTMLInputElement)?.focus();

    if (unmaskedValue === '0') {
      (ref?.current as HTMLInputElement)?.setSelectionRange(1, 5);
    }
  };

  return (
    <Input
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      ref={ref as any}
      {...props}
      autoComplete='off'
      onChange={handleValueChange}
      onClick={handleFocusOnClick}
    />
  );
};
