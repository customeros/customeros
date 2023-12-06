'use client';

import { useField } from 'react-inverted-form';

import { Flex } from '@ui/layout/Flex';
import { Text } from '@ui/typography/Text';
import { FormLabel, FormControl } from '@ui/form/FormElement';
import {
  NumberInput,
  blockInvalidChar,
  NumberInputField,
} from '@ui/form/NumberInput';

export interface FormPeriodInputProps {
  name: string;
  label?: string;
  formId: string;
  placeholder?: string;
}

export const FormPeriodInput = ({
  name,
  label,
  formId,
  placeholder = '',
}: FormPeriodInputProps) => {
  const { getInputProps } = useField(name, formId);
  const { value, ...inputProps } = getInputProps();

  return (
    <FormControl w='calc(50% - 8px)'>
      <FormLabel fontWeight={600} color='gray.700' fontSize='sm' mb={-1}>
        {label}
      </FormLabel>
      <NumberInput
        {...inputProps}
        w='full'
        min={2}
        max={100}
        value={value}
        onKeyDown={blockInvalidChar}
      >
        <NumberInputField px={0} autoComplete='off' placeholder={placeholder} />
        {!!value?.length && (
          <Flex position='absolute' left={value?.length * 2} ml='1' top='7px'>
            <Text>years</Text>
          </Flex>
        )}
      </NumberInput>
    </FormControl>
  );
};
