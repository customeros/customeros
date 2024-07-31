import { forwardRef } from 'react';
import { useField } from 'react-inverted-form';

import {
  AutoresizeTextarea,
  AutoresizeTextareaProps,
} from './AutoresizeTextarea';

interface FormAutoresizeTextareaProps extends AutoresizeTextareaProps {
  name: string;
  formId: string;
  label?: string;
  labelProps?: React.LabelHTMLAttributes<HTMLLabelElement>;
}

/**
 * @deprecated use `<Textarea />` instead
 */
export const FormAutoresizeTextarea = forwardRef<
  HTMLTextAreaElement,
  FormAutoresizeTextareaProps
>(({ label, formId, labelProps, ...props }, ref) => {
  const { getInputProps } = useField(props.name, formId);

  return (
    <div className='w-full'>
      <label
        {...labelProps}
        htmlFor={props.name}
        className='mb-1 text-gray-700 font-semibold text-sm'
      >
        {label}
      </label>

      <AutoresizeTextarea ref={ref} {...getInputProps()} {...props} />
    </div>
  );
});
