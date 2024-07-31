import { useField } from 'react-inverted-form';
import { useRef, PropsWithChildren } from 'react';

import { Checkbox, CheckboxProps } from './Checkbox';

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
  const inputRef = useRef<HTMLButtonElement>(null);
  const { getInputProps } = useField(name, formId);
  const { value, onChange } = getInputProps();

  const handleChange = (isChecked: boolean) => {
    onChange?.(isChecked);
  };

  return (
    <Checkbox
      {...props}
      ref={inputRef}
      disabled={false}
      isChecked={value}
      onChange={(isChecked) => handleChange(isChecked as boolean)}
    >
      {children}
    </Checkbox>
  );
};
