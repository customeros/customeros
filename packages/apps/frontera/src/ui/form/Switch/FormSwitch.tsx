import { useField } from 'react-inverted-form';
import React, { ReactNode, forwardRef, ForwardedRef } from 'react';

import { cn } from '@ui/utils/cn';
import { Switch, SwitchProps } from '@ui/form/Switch/Switch';
export interface FormSwitchProps extends SwitchProps {
  name: string;
  formId: string;
  label?: ReactNode;
  leftElement?: ReactNode;
  isLabelVisible?: boolean;
  onChangeCallback?: (onCallback: () => void) => void;
  labelProps?: React.LabelHTMLAttributes<HTMLLabelElement>;
}

export const FormSwitch = forwardRef(
  (
    {
      name,
      formId,
      label,
      isLabelVisible = true,
      labelProps,
      leftElement,
      size = 'md',
      onChangeCallback,
      ...props
    }: FormSwitchProps,
    ref: ForwardedRef<HTMLButtonElement>,
  ) => {
    const { getInputProps } = useField(name, formId);
    const { value, onChange, onBlur, ...rest } = getInputProps();

    const handleChange = (newValue: boolean) => {
      if (!newValue) {
        onChange(newValue);

        return;
      }

      if (newValue) {
        if (onChangeCallback) {
          onChangeCallback(() => onChange(newValue));

          return;
        }
        onChange(newValue);
      }
    };

    return (
      <div className='flex w-full justify-between items-center'>
        {isLabelVisible ? (
          <label
            {...labelProps}
            className={cn(
              props.isDisabled ? 'text-gray-500' : 'text-gray-700',
              'text-base',
            )}
          >
            {label}
          </label>
        ) : (
          <label>{label}</label>
        )}
        <div className='items-center flex'>
          {leftElement && leftElement}
          <Switch
            ref={ref}
            size={size}
            isChecked={value}
            onChange={() => handleChange(!value)}
            {...rest}
          />
        </div>
      </div>
    );
  },
);
