import { forwardRef } from 'react';
import { useField } from 'react-inverted-form';

import { Select, SelectInstance, SelectProps } from './Select';

interface FormSelectProps extends SelectProps {
  name: string;
  formId: string;
}

export const FormSelect = forwardRef<SelectInstance, FormSelectProps>(
  ({ name, formId, ...rest }, ref) => {
    const { getInputProps } = useField(name, formId);
    const { id, onChange, onBlur, value } = getInputProps();

    return (
      <Select
        ref={ref}
        id={id}
        name={name}
        value={value}
        onBlur={() => onBlur(value)}
        onChange={onChange}
        {...rest}
      />
    );
  },
);
