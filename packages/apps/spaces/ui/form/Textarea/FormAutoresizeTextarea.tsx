import { forwardRef } from 'react';
import { useField } from 'react-inverted-form';

import {
  AutoresizeTextarea,
  AutoresizeTextareaProps,
} from './AutoresizeTextarea';

interface FormAutoresizeTextareaProps extends AutoresizeTextareaProps {
  name: string;
  formId: string;
}

export const FormAutoresizeTextarea = forwardRef<
  HTMLTextAreaElement,
  FormAutoresizeTextareaProps
>((props, ref) => {
  const { getInputProps } = useField(props.name, props.formId);

  return <AutoresizeTextarea ref={ref} {...getInputProps()} {...props} />;
});
