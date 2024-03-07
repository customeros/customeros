'use client';

import { forwardRef } from 'react';
import { useField } from 'react-inverted-form';

import { FormLabel, FormControl, VisuallyHidden } from '@chakra-ui/react';

import { Input } from '../Input';
import {
  InputGroup,
  InputGroupProps,
  InputLeftElement,
  InputRightElement,
} from './InputGroup';

interface FormInputGroupProps extends InputGroupProps {
  name: string;
  formId: string;
  label?: string;
  autoFocus?: boolean;
  isLabelVisible?: boolean;
  leftElement?: React.ReactNode;
  rightElement?: React.ReactNode;
}

export const FormInputGroup = forwardRef((props: FormInputGroupProps, ref) => {
  const {
    name,
    formId,
    label,
    isLabelVisible,
    leftElement,
    rightElement,
    autoFocus,
    ...rest
  } = props;
  const { getInputProps } = useField(name, formId);

  return (
    <FormControl>
      {isLabelVisible ? (
        <FormLabel>{label}</FormLabel>
      ) : (
        <VisuallyHidden>
          <FormLabel>{label}</FormLabel>
        </VisuallyHidden>
      )}
      <InputGroup {...rest}>
        {leftElement && (
          <InputLeftElement w='4'>{leftElement}</InputLeftElement>
        )}
        <Input
          ref={ref}
          {...getInputProps()}
          pl='30px'
          autoComplete='off'
          {...rest}
          autoFocus={autoFocus}
        />
        {rightElement && <InputRightElement>{rightElement}</InputRightElement>}
      </InputGroup>
    </FormControl>
  );
});
