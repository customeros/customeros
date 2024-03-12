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
            {...rest}
            colorScheme='primary'
            {...props}
            isChecked={value}
            onChange={() => handleChange(!value)}
          />
        </Flex>
      </FormControl>
    );
  },
);
