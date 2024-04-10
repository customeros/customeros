import { forwardRef } from 'react';
import { useField } from 'react-inverted-form';

import {
  AutoresizeTextarea,
  AutoresizeTextareaProps,
} from './AutoresizeTextarea2';

interface FormAutoresizeTextareaProps extends AutoresizeTextareaProps {
  name: string;
  formId: string;
  label?: string;
  labelProps?: React.LabelHTMLAttributes<HTMLLabelElement>;
}

export const FormAutoresizeTextarea = forwardRef<
  HTMLTextAreaElement,
  FormAutoresizeTextareaProps
>(({ label, formId, labelProps, ...props }, ref) => {
  const { getInputProps } = useField(props.name, formId);

  return (
    <>
      <label
        {...labelProps}
        className='mb-1 text-gray-700 font-semibold text-sm'
        htmlFor={props.name}
      >
        {label}
      </label>

      <AutoresizeTextarea ref={ref} {...getInputProps()} {...props} />
    </>
  );
});
