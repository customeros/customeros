import type { FormLabelProps, FormControlProps } from '@chakra-ui/react';

import { useField } from 'react-inverted-form';
import { cloneElement, isValidElement } from 'react';

import { FormLabel, FormControl } from '@chakra-ui/react';

export interface FormElementProps extends FormControlProps {
  name: string;
  formId: string;
  label?: string;
}

export const FormElement = ({
  name,
  label,
  formId,
  children,
}: FormElementProps) => {
  const { getInputProps } = useField(name, formId);

  if (!isValidElement(children)) return children;

  const inputProps = getInputProps();
  const clonedChildren = cloneElement(children, inputProps);

  return (
    <FormControl id={inputProps.id}>
      {label && <FormLabel htmlFor={inputProps.id}>{label}</FormLabel>}
      {clonedChildren}
    </FormControl>
  );
};

export type { FormControlProps, FormLabelProps };
export { FormControl, FormLabel };
