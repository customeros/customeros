import { useField } from 'react-inverted-form';
import React, { useRef, ChangeEvent, PropsWithChildren } from 'react';

import { Checkbox, CheckboxProps } from '@chakra-ui/react';

export interface FormCheckboxProps extends Omit<CheckboxProps, 'onChange'> {
  name: string;
  formId: string;
}

export const FormCheckbox = ({
  formId,
  name,
  children,
  ...props
}: PropsWithChildren<FormCheckboxProps>) => {
  const inputRef = useRef<HTMLInputElement>(null);
  const { getInputProps } = useField(name, formId);
  const { value, onChange } = getInputProps();

  const handleChange = (e: ChangeEvent<HTMLInputElement>) => {
    onChange?.(e.target.checked);
  };

  return (
    <Checkbox
      {...props}
      value={value}
      ref={inputRef}
      isChecked={value}
      onChange={handleChange}
      isDisabled={false}
    >
      {children}
    </Checkbox>
  );
};
