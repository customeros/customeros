'use client';

import { forwardRef } from 'react';
import { useField } from 'react-inverted-form';

import { NumberInputProps, NumberInput } from './NumberInput';

interface FormNumberInputProps extends NumberInputProps {
  formId: string;
  name: string;
}

export const FormNumberInput = forwardRef(
  ({ name, formId, ...props }: FormNumberInputProps, ref) => {
    const { getInputProps } = useField(name, formId);

    return (
      <NumberInput
        ref={ref}
        {...getInputProps()}
        {...props}
        autoComplete='off'
      />
    );
  },
);
