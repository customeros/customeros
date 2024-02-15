'use client';

import { ReactNode, forwardRef } from 'react';
import { useField } from 'react-inverted-form';

import {
  Flex,
  Switch,
  FormLabel,
  FormControl,
  SwitchProps,
  VisuallyHidden,
  FormLabelProps,
} from '@chakra-ui/react';

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
      onChangeCallback,
      ...props
    }: FormSwitchProps,
    ref,
  ) => {
    const { getInputProps } = useField(name, formId);
    const { value, onChange, ...rest } = getInputProps();

    const handleChange = (value: boolean) => {
      if (!value) {
        onChange(value);

        return;
      }
      if (value) {
        if (onChangeCallback) {
          onChangeCallback(() => onChange(value));

          return;
        }
        onChange(value);
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
          <FormLabel {...labelProps} fontWeight='medium'>
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
            {...rest}
            {...props}
            isChecked={value}
            onChange={() => handleChange(!value)}
            colorScheme='primary'
          />
        </Flex>
      </FormControl>
    );
  },
);
