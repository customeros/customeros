import { useField } from 'react-inverted-form';
import React, { forwardRef, ForwardedRef } from 'react';

import { Input, InputProps } from './Input';

export interface FormInputProps extends InputProps {
  name: string;
  formId: string;
  label?: string;
  isLabelHidden?: boolean;
  labelProps?: React.LabelHTMLAttributes<HTMLLabelElement>;
}

//todo add visually hidden label - accessibility

export const FormInput = forwardRef(
  (
    {
      name,
      formId,
      label,
      labelProps,
      isLabelHidden,
      ...props
    }: FormInputProps,
    ref: ForwardedRef<HTMLInputElement>,
  ) => {
    const { getInputProps, renderError, state } = useField(name, formId);

    return (
      <div className='w-full'>
        <label className={isLabelHidden ? 'sr-only' : ''} {...labelProps}>
          {label}
        </label>

        <Input
          ref={ref}
          {...getInputProps()}
          {...props}
          onInvalid={() => state.meta?.meta?.hasError}
          autoComplete='off'
          data-1p-ignore
        />
        {renderError((error) => (
          <span className='text-xs text-error-500'>{error}</span>
        ))}
      </div>
    );
  },
);
