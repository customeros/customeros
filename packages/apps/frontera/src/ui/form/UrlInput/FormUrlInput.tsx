import React from 'react';
import { useField } from 'react-inverted-form';

import { UrlInput } from './UrlInput';
import { FormInputProps } from '../Input/FormInput';

interface FormUrlInputProps extends FormInputProps {
  name: string;
  formId: string;
  label?: string;
  isLabelVisible?: boolean;
  leftElement?: React.ReactNode;
  rightElement?: React.ReactNode;
}

export const FormUrlInput = ({ name, formId, ...props }: FormUrlInputProps) => {
  const { getInputProps } = useField(name, formId);

  return <UrlInput {...getInputProps()} {...props} />;
};
