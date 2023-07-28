'use client';

import { forwardRef } from 'react';
import { useField } from 'react-inverted-form';

import { Input, InputProps } from './Input';

interface FormInputProps extends InputProps {
  formId: string;
  name: string;
}

export const FormInput = forwardRef(
  ({ name, formId, ...props }: FormInputProps, ref) => {
    const { getInputProps } = useField(name, formId);

    return (
      <Input ref={ref} {...getInputProps()} {...props} autoComplete='off' />
    );
  },
);
