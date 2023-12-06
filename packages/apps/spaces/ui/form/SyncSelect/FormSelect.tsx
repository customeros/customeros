import { forwardRef } from 'react';
import { useField } from 'react-inverted-form';

import { FormLabel, FormControl, VisuallyHidden } from '@chakra-ui/react';

import { Select, SelectProps, SelectInstance } from './Select';

interface FormSelectProps extends SelectProps {
  name: string;
  formId: string;
  label?: string;
  isLabelVisible?: boolean;
}

// TODO: Label props are different from FormInput. They should be in sync
export const FormSelect = forwardRef<SelectInstance, FormSelectProps>(
  ({ name, formId, isLabelVisible, label, ...rest }, ref) => {
    const { getInputProps } = useField(name, formId);
    const { id, onChange, onBlur, value, defaultValue } = getInputProps();

    return (
      <FormControl>
        {isLabelVisible ? (
          <FormLabel fontWeight={600} color='gray.700' fontSize='sm' mb={-1}>
            {label}
          </FormLabel>
        ) : (
          <VisuallyHidden>
            <FormLabel>{label}</FormLabel>
          </VisuallyHidden>
        )}
        <Select
          ref={ref}
          id={id}
          name={name}
          value={value}
          onBlur={() => onBlur(value)}
          defaultValue={defaultValue}
          onChange={onChange}
          {...rest}
        />
      </FormControl>
    );
  },
);
