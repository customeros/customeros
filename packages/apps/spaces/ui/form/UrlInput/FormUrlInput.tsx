import React from 'react';
import { useField } from 'react-inverted-form';

import { UrlInput } from './UrlInput';
import { FormInputProps } from '../Input/FormInput2';

interface FormUrlInputProps extends FormInputProps {
  name: string;
  formId: string;
  label?: string;
  isLabelVisible?: boolean;
  leftElement?: React.ReactNode;
  rightElement?: React.ReactNode;
}

export const FormUrlInput: React.FC<FormUrlInputProps> = ({
  name,
  formId,
  ...props
}) => {
  const { getInputProps } = useField(name, formId);

  return <UrlInput {...getInputProps()} {...props} />;
};
