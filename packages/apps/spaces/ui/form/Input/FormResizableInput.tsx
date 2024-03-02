'use client';

import { forwardRef } from 'react';
import { useField } from 'react-inverted-form';

import {
  Flex,
  FormLabel,
  FormControl,
  VisuallyHidden,
  FormLabelProps,
} from '@chakra-ui/react';

import { Text } from '@ui/typography/Text';

import { InputProps } from './Input';
import { ResizableInput } from './ResizableInput';

interface FormInputProps extends InputProps {
  name: string;
  formId: string;
  label?: string;
  error?: string | null;
  isLabelVisible?: boolean;
  labelProps?: FormLabelProps;
  rightElement?: React.ReactNode;
}

export const FormResizableInput = forwardRef<HTMLInputElement, FormInputProps>(
  (
    {
      name,
      formId,
      label,
      isLabelVisible,
      labelProps,
      rightElement,
      ...props
    }: FormInputProps,
    ref,
  ) => {
    const { getInputProps, renderError, state } = useField(name, formId);

    return (
      <FormControl>
        {isLabelVisible ? (
          <FormLabel {...labelProps}>{label}</FormLabel>
        ) : (
          <VisuallyHidden>
            <FormLabel>{label}</FormLabel>
          </VisuallyHidden>
        )}
        <Flex alignItems='center'>
          <ResizableInput
            ref={ref}
            {...getInputProps()}
            {...props}
            isInvalid={state.meta?.meta?.hasError}
            autoComplete='off'
          />
          {rightElement && rightElement}
        </Flex>
        {renderError((error) => (
          <Text fontSize='xs' color='error.500'>
            {error}
          </Text>
        ))}
      </FormControl>
    );
  },
);
