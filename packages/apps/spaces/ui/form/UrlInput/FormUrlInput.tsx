import React from 'react';
import { useField } from 'react-inverted-form';

import { FormLabelProps } from '@chakra-ui/react';

import { InputProps } from '@ui/form/Input';

import { UrlInput } from './UrlInput';

interface FormUrlInputProps extends InputProps {
  name: string;
  formId: string;
  label?: string;
  isLabelVisible?: boolean;
  labelProps?: FormLabelProps;
  leftElement?: React.ReactNode;
  rightElement?: React.ReactNode;
}

export const FormUrlInput = ({ name, formId, ...props }: FormUrlInputProps) => {
  const { getInputProps } = useField(name, formId);

  return <UrlInput {...getInputProps()} {...props} />;
};
