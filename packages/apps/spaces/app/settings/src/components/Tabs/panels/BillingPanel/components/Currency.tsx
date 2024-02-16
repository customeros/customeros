import { useIMask } from 'react-imask';
import React, { useEffect } from 'react';

import { VisuallyHidden } from '@ui/presentation/VisuallyHidden';
import { Input, FormLabel, InputProps, FormControl } from '@ui/form/Input';

interface CurrencyProps extends InputProps {
  label?: string;
  isLabelVisible?: boolean;
  currency?: string | null;
  onValueChange?: (value: string) => void;
}

export const Currency = ({
  isLabelVisible,
  label,
  value = '',
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
        min: 0.01,
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

  return (
    <FormControl>
      {isLabelVisible ? (
        <FormLabel>{label}</FormLabel>
      ) : (
        <VisuallyHidden>
          <FormLabel>{label}</FormLabel>
        </VisuallyHidden>
      )}

      <Input
        ref={ref}
        {...props}
        onChange={(e) => setValue(e.target.value)}
        autoComplete='off'
      />
    </FormControl>
  );
};
