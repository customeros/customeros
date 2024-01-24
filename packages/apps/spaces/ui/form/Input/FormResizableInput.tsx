'use client';

import { forwardRef } from 'react';
import { useField } from 'react-inverted-form';

import {
  FormLabel,
  FormControl,
  VisuallyHidden,
  FormLabelProps,
} from '@chakra-ui/react';

import { InputProps } from './Input';
import { ResizableInput } from './ResizableInput';

interface FormInputProps extends InputProps {
  name: string;
  formId: string;
  label?: string;
  isLabelVisible?: boolean;
  labelProps?: FormLabelProps;
}

//todo add visually hidden label - accessibility

export const FormResizableInput = forwardRef<HTMLInputElement, FormInputProps>(
  (
    {
      name,
      formId,
      label,
      isLabelVisible,
      labelProps,
      ...props
    }: FormInputProps,
    ref,
  ) => {
    const { getInputProps } = useField(name, formId);

    return (
      <FormControl>
        {isLabelVisible ? (
          <FormLabel {...labelProps}>{label}</FormLabel>
        ) : (
          <VisuallyHidden>
            <FormLabel>{label}</FormLabel>
          </VisuallyHidden>
        )}

        <ResizableInput
          ref={ref}
          {...getInputProps()}
          {...props}
          autoComplete='off'
        />
      </FormControl>
    );
  },
);
