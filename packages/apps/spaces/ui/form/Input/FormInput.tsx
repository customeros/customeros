'use client';

import { forwardRef } from 'react';
import { useField } from 'react-inverted-form';

import { Input, InputProps } from './Input';
import { FormControl, FormLabel, VisuallyHidden } from '@chakra-ui/react';

interface FormInputProps extends InputProps {
  formId: string;
  name: string;
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
