'use client';
import { FC } from 'react';
import { useField } from 'react-inverted-form';
import { InputGroup, InputLeftElement, InputRightElement } from './InputGroup';
import {
  NumberInput,
  NumberInputField,
} from '@ui/form/NumberInput/NumberInput';
import { NumberInputProps } from '@chakra-ui/react';

interface FormNumberInputGroupProps extends NumberInputProps {
  name: string;
  formId: string;
  leftElement?: React.ReactNode;
  rightElement?: React.ReactNode;
}

export const FormNumberInputGroup: FC<FormNumberInputGroupProps> = ({
  name,
  formId,
  leftElement,
  rightElement,
  ...rest
}) => {
  const { getInputProps } = useField(name, formId);

  return (
    <InputGroup>
      {leftElement && <InputLeftElement w='4'>{leftElement}</InputLeftElement>}

      <NumberInput {...rest} {...getInputProps()}>
        <NumberInputField
          pl='30px'
          pr={0}
          autoComplete='off'
          placeholder={rest?.placeholder || ''}
        />
      </NumberInput>
      {rightElement && <InputRightElement>{rightElement}</InputRightElement>}
    </InputGroup>
  );
};
