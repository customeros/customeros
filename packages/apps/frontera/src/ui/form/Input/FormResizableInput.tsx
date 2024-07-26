import { useField } from 'react-inverted-form';
import { forwardRef, ForwardedRef } from 'react';

import { InputProps } from './Input';
import { ResizableInput } from './ResizableInput';

interface FormInputProps extends InputProps {
  name: string;
  formId: string;
  label?: string;
  error?: string | null;
  rightElement?: React.ReactNode;
  labelProps?: React.LabelHTMLAttributes<HTMLLabelElement>;
}

/**
 * @deprecated Use `<ResizableInput />` instead
 */
export const FormResizableInput = forwardRef<HTMLInputElement, FormInputProps>(
  (
    {
      name,
      size = 'md',
      formId,
      label,
      labelProps,
      rightElement,
      ...props
    }: FormInputProps,
    ref: ForwardedRef<HTMLInputElement>,
  ) => {
    const { getInputProps, renderError, state } = useField(name, formId);

    return (
      <div>
        <label {...labelProps}>{label}</label>

        <div className='flex items-center'>
          <ResizableInput
            ref={ref}
            {...getInputProps()}
            {...props}
            onInvalid={() => state.meta?.meta?.hasError}
            autoComplete='off'
          />
          {rightElement && rightElement}
        </div>

        {renderError((error) => (
          <span className='text-xs text-error-500'>{error}</span>
        ))}
      </div>
    );
  },
);
