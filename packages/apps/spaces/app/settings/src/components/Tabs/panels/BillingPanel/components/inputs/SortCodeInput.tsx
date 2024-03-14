import { useIMask } from 'react-imask';
import React, { useEffect } from 'react';
import { useField } from 'react-inverted-form';

import { VisuallyHidden } from '@ui/presentation/VisuallyHidden';
import { Input, FormLabel, FormControl, FormInputProps } from '@ui/form/Input';

const opts = {
  mask: '00-00-00',
  definitions: {
    '0': /[0-9]/,
  },
};
export const SortCodeInput = ({
  isLabelVisible,
  labelProps,
  label,
  formId,
  name,
  ...props
}: FormInputProps) => {
  const { ref, setUnmaskedValue } = useIMask(
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
    <FormControl maxW='80px'>
      {isLabelVisible ? (
        <FormLabel {...labelProps}>{label}</FormLabel>
      ) : (
        <VisuallyHidden>
          <FormLabel>{label}</FormLabel>
        </VisuallyHidden>
      )}

      <Input ref={ref} {...props} onChange={onChange} autoComplete='off' />
    </FormControl>
  );
};
