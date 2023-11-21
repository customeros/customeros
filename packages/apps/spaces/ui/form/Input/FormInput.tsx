'use client';

import { forwardRef } from 'react';
import { useField } from 'react-inverted-form';

import { FormLabel, FormControl, VisuallyHidden } from '@chakra-ui/react';

import { Input, InputProps } from './Input';

interface FormInputProps extends InputProps {
  name: string;
  formId: string;
  label?: string;
  isLabelVisible?: boolean;
}

//todo add visually hidden label - accessibility

export const FormInput = forwardRef(
  ({ name, formId, label, isLabelVisible, ...props }: FormInputProps, ref) => {
    const { getInputProps } = useField(name, formId);

    return (
      <FormControl>
        {isLabelVisible ? (
          <FormLabel>{label}</FormLabel>
        ) : (
          <VisuallyHidden>
            <FormLabel>{label}</FormLabel>
          </VisuallyHidden>
        )}

        <Input ref={ref} {...getInputProps()} {...props} autoComplete='off' />
      </FormControl>
    );
  },
);
