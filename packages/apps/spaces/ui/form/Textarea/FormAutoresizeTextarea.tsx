import { forwardRef } from 'react';
import { useField } from 'react-inverted-form';

import { FormLabel, FormControl, VisuallyHidden } from '@chakra-ui/react';

import {
  AutoresizeTextarea,
  AutoresizeTextareaProps,
} from './AutoresizeTextarea';

interface FormAutoresizeTextareaProps extends AutoresizeTextareaProps {
  name: string;
  formId: string;
  label?: string;
  isLabelVisible?: boolean;
}

export const FormAutoresizeTextarea = forwardRef<
  HTMLTextAreaElement,
  FormAutoresizeTextareaProps
>(({ isLabelVisible, label, formId, ...props }, ref) => {
  const { getInputProps } = useField(props.name, formId);

  return (
    <FormControl>
      {isLabelVisible ? (
        <FormLabel
          fontWeight={600}
          color='gray.700'
          fontSize='sm'
          mb={-1}
          htmlFor={props.name}
        >
          {label}
        </FormLabel>
      ) : (
        <VisuallyHidden>
          <FormLabel>{label}</FormLabel>
        </VisuallyHidden>
      )}
      <AutoresizeTextarea ref={ref} {...getInputProps()} {...props} />
    </FormControl>
  );
});
