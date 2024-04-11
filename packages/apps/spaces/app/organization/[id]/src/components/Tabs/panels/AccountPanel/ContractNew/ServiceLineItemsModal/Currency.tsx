import { useIMask } from 'react-imask';
import React, { useEffect } from 'react';

import { InputElement } from 'imask';

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
        onChange={handleValueChange}
        onClick={handleFocusOnClick}
        autoComplete='off'
      />
    </FormControl>
  );
};
