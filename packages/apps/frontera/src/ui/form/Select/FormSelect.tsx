import { forwardRef } from 'react';
import { SelectInstance } from 'react-select';
import { useField } from 'react-inverted-form';

import { cn } from '@ui/utils/cn';

import { Select, SelectProps } from './Select';

interface FormSelectProps extends SelectProps {
  name: string;
  formId: string;
  label?: string;
  isLabelVisible?: boolean;
  labelProps?: React.LabelHTMLAttributes<HTMLLabelElement>;
}

// TODO: Label props are different from FormInput. They should be in sync
export const FormSelect = forwardRef<SelectInstance, FormSelectProps>(
  ({ name, formId, isLabelVisible, label, labelProps, ...rest }, ref) => {
    const { getInputProps } = useField(name, formId);
    const { id, onChange, onBlur, value } = getInputProps();

    return (
      <div className='w-full'>
        <label
          className={cn({
            absolute: !isLabelVisible,
            'top-[-999999px]': !isLabelVisible,
          })}
          {...labelProps}
        >
          {label}
        </label>

        <Select
          ref={ref}
          id={id}
          name={name}
          value={value}
          onBlur={() => onBlur(value)}
          defaultValue={value}
          onChange={onChange}
          {...rest}
        />
      </div>
    );
  },
);
