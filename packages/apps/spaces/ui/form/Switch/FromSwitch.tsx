'use client';

import { useField } from 'react-inverted-form';
import { ReactNode, forwardRef, ForwardedRef } from 'react';

import {
  Flex,
  FormLabel,
  FormControl,
  VisuallyHidden,
  FormLabelProps,
} from '@chakra-ui/react';

import { Switch, SwitchProps } from '@ui/form/Switch/Switch2';
export interface FormSwitchProps extends SwitchProps {
  name: string;
  formId: string;
  label?: ReactNode;
  leftElement?: ReactNode;
  isLabelVisible?: boolean;
  labelProps?: FormLabelProps;
  onChangeCallback?: (onCallback: () => void) => void;
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
      <FormControl
        as={Flex}
        w='full'
        justifyContent='space-between'
        alignItems='center'
      >
        {isLabelVisible ? (
          <FormLabel
            {...labelProps}
            fontWeight='medium'
            color={props.isDisabled ? 'gray.500' : 'gray.700'}
          >
            {label}
          </FormLabel>
        ) : (
          <VisuallyHidden>
            <FormLabel>{label}</FormLabel>
          </VisuallyHidden>
        )}
        <Flex alignItems='center'>
          {leftElement && leftElement}
          <Switch
            ref={ref}
            size={size}
            isChecked={value}
            onChange={() => handleChange(!value)}
            {...rest}
          />
        </Flex>
      </FormControl>
    );
  },
);
