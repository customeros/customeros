import { useIMask } from 'react-imask';
import React, { useEffect } from 'react';
import { useField } from 'react-inverted-form';

import { Input } from '@ui/form/Input/Input2';
import { VisuallyHidden } from '@ui/presentation/VisuallyHidden';
import { FormLabel, FormControl, FormInputProps } from '@ui/form/Input';

const opts = {
  mask: 'AA 000 000 000',
  definitions: {
    A: /[A-Za-z]/,
    '0': /[0-9]/,
  },
  prepare: function (value: string, mask: { _value: string }) {
    if (mask._value.length < 2) {
      return value.toUpperCase();
    }

    return value;
  },
  format: function (value: string) {
    return value.toUpperCase();
  },
  parse: function (value: string) {
    return value.toUpperCase();
  },
};

export const VatInput = ({
  isLabelVisible,
  labelProps,
  label,
  formId,
  name,
  ...props
}: FormInputProps) => {
  const { ref, setUnmaskedValue } = useIMask<HTMLInputElement>(
    opts /* { onAccept, onComplete } */,
  );
  const { getInputProps } = useField(name, formId);
  const { value, onChange } = getInputProps();

  useEffect(() => {
    if (value) {
      setUnmaskedValue(value);
    }
  }, [value]);

  return (
    <FormControl>
      {isLabelVisible ? (
        <FormLabel {...labelProps}>{label}</FormLabel>
      ) : (
        <VisuallyHidden>
          <FormLabel>{label}</FormLabel>
        </VisuallyHidden>
      )}

      <Input ref={ref} onChange={onChange} autoComplete='off' {...props} />
    </FormControl>
  );
};
