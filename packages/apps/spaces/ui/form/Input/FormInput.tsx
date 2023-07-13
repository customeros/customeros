'use client';

import { forwardRef } from 'react';
import { useField } from 'react-inverted-form';

import { Input, InputProps } from './Input';

interface FormInputProps extends InputProps {
  formId: string;
  name: string;
}

export const FormInput = forwardRef((props: FormInputProps, ref) => {
  const { getInputProps } = useField(props.name, props.formId);

  return <Input ref={ref} {...getInputProps()} {...props} />;
});
