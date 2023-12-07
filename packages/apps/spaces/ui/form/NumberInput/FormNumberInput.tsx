'use client';

import React, { forwardRef } from 'react';
import { useField } from 'react-inverted-form';

import { FormLabel, FormControl, VisuallyHidden } from '@chakra-ui/react';

import { NumberInputField } from '@ui/form/NumberInput/NumberInput';
import {
  InputGroup,
  InputLeftElement,
  InputRightElement,
  InputLeftElementProps,
  InputRightElementProps,
} from '@ui/form/InputGroup/InputGroup';

import { NumberInput, NumberInputProps } from './NumberInput';

export interface FormNumberInputProps extends NumberInputProps {
  name: string;
  formId: string;
  label?: string;
  isLabelVisible?: boolean;
  leftElement?: React.ReactNode;
  rightElement?: React.ReactNode;
  leftElementProps?: InputLeftElementProps;
  rightElementProps?: InputRightElementProps;
}
export const blockInvalidChar = (e: {
  key: string;
  preventDefault: () => void;
}) => ['e', 'E', '+', '-'].includes(e.key) && e.preventDefault();

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
      leftElementProps,
      rightElementProps,
      ...props
    },
    ref,
  ) => {
    const { getInputProps } = useField(name, formId);

    return (
      <FormControl w={props.w} width={props.width}>
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
            <InputLeftElement w='4' {...leftElementProps}>
              {leftElement}
            </InputLeftElement>
          )}

          <NumberInput
            {...getInputProps()}
            {...props}
            w='full'
            autoComplete='off'
            onKeyDown={blockInvalidChar}
          >
            <NumberInputField
              ref={ref}
              pl={leftElement ? '30px' : '0'}
              pr={0}
              autoComplete='off'
              placeholder={props?.placeholder || ''}
            />
          </NumberInput>

          {rightElement && (
            <InputRightElement {...rightElementProps}>
              {rightElement}
            </InputRightElement>
          )}
        </InputGroup>
      </FormControl>
    );
  },
);
