'use client';

import React, { forwardRef } from 'react';
import { useField } from 'react-inverted-form';

import { FormLabel, FormControl, VisuallyHidden } from '@chakra-ui/react';

import { NumberInputField } from '@ui/form/NumberInput/NumberInput';
import {
  InputGroup,
  InputLeftElement,
  InputRightElement,
} from '@ui/form/InputGroup/InputGroup';

import { NumberInput, NumberInputProps } from './NumberInput';

interface FormNumberInputProps extends NumberInputProps {
  name: string;
  formId: string;
  label?: string;
  isLabelVisible?: boolean;
  leftElement?: React.ReactNode;
  rightElement?: React.ReactNode;
}

export const FormNumberInput = forwardRef<
  HTMLInputElement,
  FormNumberInputProps
>(
  (
    {
      name,
      formId,
      isLabelVisible,
      label = '',
      leftElement,
      rightElement,
      ...props
    },
    ref,
  ) => {
    const { getInputProps } = useField(name, formId);

    return (
      <FormControl>
        {isLabelVisible ? (
          <FormLabel
            fontWeight={600}
            color={props?.color}
            fontSize='sm'
            mb={-1}
          >
            {label}
          </FormLabel>
        ) : (
          <VisuallyHidden>
            <FormLabel>{label}</FormLabel>
          </VisuallyHidden>
        )}
        <InputGroup>
          {leftElement && (
            <InputLeftElement w='4'>{leftElement}</InputLeftElement>
          )}

          <NumberInput {...getInputProps()} {...props} autoComplete='off'>
            <NumberInputField
              ref={ref}
              pl={leftElement ? '30px' : '0'}
              pr={0}
              autoComplete='off'
              placeholder={props?.placeholder || ''}
            />
          </NumberInput>

          {rightElement && (
            <InputRightElement>{rightElement}</InputRightElement>
          )}
        </InputGroup>
      </FormControl>
    );
  },
);
