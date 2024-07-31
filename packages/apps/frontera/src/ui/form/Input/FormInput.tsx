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

/**
 * @deprecated Use `<Input />` instead
 */
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
          data-1p-ignore
          autoComplete='off'
          onInvalid={() => state.meta?.meta?.hasError}
        />
        {renderError((error) => (
          <span className='text-xs text-error-500'>{error}</span>
        ))}
      </div>
    );
  },
);
