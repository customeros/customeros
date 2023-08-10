import React from 'react';
import { useField } from 'react-inverted-form';
import {
  FormControl,
  FormLabel,
  NumberInputProps,
  VisuallyHidden,
} from '@chakra-ui/react';
import {
  InputGroup,
  InputLeftElement,
  InputRightElement,
} from '@ui/form/InputGroup/InputGroup';
import {
  NumberInput,
  NumberInputField,
} from '@ui/form/NumberInput/NumberInput';

interface CurrencyInputProps extends NumberInputProps {
  name: string;
  formId: string;
  leftElement?: React.ReactNode;
  rightElement?: React.ReactNode;
  label?: string;
  isLabelVisible?: boolean;
}

export const CurrencyInput: React.FC<CurrencyInputProps> = ({
  formId,
  name,
  isLabelVisible,
  label,
  leftElement,
  rightElement,
  ...rest
}) => {
  const { getInputProps } = useField(name, formId);
  const { value, onChange } = getInputProps();
  const format = (val: number) => `$` + val;
  const parse = (val: string) => val.replace(/^\$/, '');

  return (
    <FormControl>
      {isLabelVisible ? (
        <FormLabel fontWeight={600} color={rest?.color} fontSize='sm' mb={-1}>
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

        <NumberInput
          {...rest}
          {...getInputProps()}
          value={format(value)}
          onChange={(valueString) => onChange(parse(valueString))}
        >
          <NumberInputField
            pl='30px'
            pr={0}
            autoComplete='off'
            placeholder={rest?.placeholder || ''}
          />
        </NumberInput>
        {rightElement && <InputRightElement>{rightElement}</InputRightElement>}
      </InputGroup>
    </FormControl>
  );
};
