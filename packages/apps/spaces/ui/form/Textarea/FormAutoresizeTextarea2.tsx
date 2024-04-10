import { forwardRef } from 'react';
import { useField } from 'react-inverted-form';

import { Textarea, TextareaProps } from './Textarea';

interface FormAutoresizeTextareaProps extends TextareaProps {
  name: string;
  formId: string;
  label?: string;
  isLabelVisible?: boolean;
  labelProps?: React.LabelHTMLAttributes<HTMLLabelElement>;
}

export const FormAutoresizeTextarea = forwardRef<
  HTMLTextAreaElement,
  FormAutoresizeTextareaProps
>(({ isLabelVisible, label, formId, labelProps, ...props }, ref) => {
  const { getInputProps } = useField(props.name, formId);

  return (
    <div className='w-full'>
      <label
        {...labelProps}
        className='mb-1 text-gray-700 font-semibold text-sm'
        htmlFor={props.name}
      >
        {label}
      </label>

      <Textarea ref={ref} {...getInputProps()} {...props} />
    </div>
  );
});
