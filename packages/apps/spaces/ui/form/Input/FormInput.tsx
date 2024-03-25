'use client';

import { useField } from 'react-inverted-form';
import React, { forwardRef, ForwardedRef } from 'react';

import { Text } from '@ui/typography/Text';

import { Input, InputProps } from './Input2';

export interface FormInputProps extends InputProps {
  name: string;
  formId: string;
  label?: string;
  labelProps?: React.LabelHTMLAttributes<HTMLLabelElement>;
}

//todo add visually hidden label - accessibility

export const FormInput = forwardRef(
  (
    { name, formId, label, labelProps, ...props }: FormInputProps,
    ref: ForwardedRef<HTMLInputElement>,
  ) => {
    const { getInputProps, renderError, state } = useField(name, formId);

    return (
      <div>
        <label {...labelProps}>{label}</label>

        <Input
          ref={ref}
          {...getInputProps()}
          {...props}
          onInvalid={() => state.meta?.meta?.hasError}
          autoComplete='off'
          data-1p-ignore
        />
        {renderError((error) => (
          <Text fontSize='xs' color='error.500'>
            {error}
          </Text>
        ))}
      </div>
    );
  },
);
