'use client';

import { forwardRef } from 'react';
import { useField } from 'react-inverted-form';

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
  leftElement?: React.ReactNode;
  rightElement?: React.ReactNode;
}

export const FormInputGroup = forwardRef((props: FormInputGroupProps, ref) => {
  const { name, formId, leftElement, rightElement, ...rest } = props;
  const { getInputProps } = useField(name, formId);

  return (
    <InputGroup ref={ref} {...rest}>
      {leftElement && <InputLeftElement>{leftElement}</InputLeftElement>}
      <Input {...getInputProps()} {...rest} />
      {rightElement && <InputRightElement>{rightElement}</InputRightElement>}
    </InputGroup>
  );
});
