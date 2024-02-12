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
  isLabelVisible?: boolean;
  labelProps?: FormLabelProps;
}

export const FormSwitch = forwardRef(
  (
    {
      name,
      formId,
      label,
      isLabelVisible = true,
      labelProps,
      ...props
    }: FormSwitchProps,
    ref,
  ) => {
    const { getInputProps } = useField(name, formId);
    const { value, onChange, ...rest } = getInputProps();

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

        <Switch
          ref={ref}
          {...rest}
          {...props}
          isChecked={value}
          onChange={() => onChange(!value)}
          colorScheme='primary'
        />
      </FormControl>
    );
  },
);
