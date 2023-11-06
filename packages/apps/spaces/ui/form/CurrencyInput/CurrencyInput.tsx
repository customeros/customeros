import React from 'react';

import {
  FormLabel,
  FormControl,
  VisuallyHidden,
  NumberInputProps,
} from '@chakra-ui/react';

import {
  NumberInput,
  NumberInputField,
} from '@ui/form/NumberInput/NumberInput';
import {
  InputGroup,
  InputLeftElement,
  InputRightElement,
} from '@ui/form/InputGroup/InputGroup';

export interface CurrencyInputProps
  extends Omit<NumberInputProps, 'onChange' | 'value'> {
  value: string;
  label?: string;
  isLabelVisible?: boolean;
  leftElement?: React.ReactNode;
  rightElement?: React.ReactNode;
  onChange?: (value: string) => void;
  parseValue?: (val: string) => string;
  formatValue?: (val: string) => string;
}

export const CurrencyInput = React.forwardRef<
  HTMLInputElement,
  CurrencyInputProps
>(
  (
    {
      isLabelVisible,
      label,
      leftElement,
      rightElement,
      value,
      onChange,
      formatValue,
      parseValue,
      ...rest
    },
    ref,
  ) => {
    const handleValueChange = (valueString: string) => {
      // handle weird case of blurring the field with an empty value
      if (valueString === '-9007199254740991') {
        onChange?.('');

        return;
      }
      if (parseValue) {
        onChange?.(parseValue(valueString));

        return;
      }
      onChange?.(valueString);
    };

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
            value={formatValue ? formatValue(value) : value}
            onChange={handleValueChange}
            _placeholder={{ color: 'gray.600' }}
          >
            <NumberInputField
              ref={ref}
              pl={leftElement ? '30px' : '0'}
              pr={0}
              autoComplete='off'
              placeholder={rest?.placeholder || ''}
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
