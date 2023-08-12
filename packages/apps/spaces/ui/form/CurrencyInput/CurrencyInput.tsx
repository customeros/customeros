import React from 'react';
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

export interface CurrencyInputProps
  extends Omit<NumberInputProps, 'onChange' | 'value'> {
  value: string;
  onChange?: (value: string) => void;
  leftElement?: React.ReactNode;
  rightElement?: React.ReactNode;
  label?: string;
  isLabelVisible?: boolean;
}

export const CurrencyInput = ({
  isLabelVisible,
  label,
  leftElement,
  rightElement,
  value,
  onChange,
  ...rest
}: CurrencyInputProps) => {
  const format = (val: string) => `$` + val ?? '';
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
          value={value ? format(value) : ''}
          onChange={(valueString) => {
            // handle weird case of blurring the field with an empty value
            if (valueString === '-9007199254740991') {
              onChange?.('');
              return;
            }
            onChange?.(parse(valueString));
          }}
          _placeholder={{color: 'gray.600'}}
        >
          <NumberInputField
            pl={leftElement ? '30px' : '0'}
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
