'use client';

import { forwardRef } from 'react';
import { useField } from 'react-inverted-form';

import { NumberInput, NumberInputProps } from './NumberInput';

interface FormNumberInputProps extends NumberInputProps {
  name: string;
  formId: string;
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
