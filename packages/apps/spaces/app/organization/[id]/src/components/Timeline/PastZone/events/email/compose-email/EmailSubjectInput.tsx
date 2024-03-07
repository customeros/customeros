'use client';
import React, { FC } from 'react';

import { Flex } from '@ui/layout/Flex';
import { Text } from '@ui/typography/Text';
import { FormInput } from '@ui/form/Input';

interface EmailSubjectInput {
  mt?: number;
  formId: string;
  fieldName: string;
}

export const EmailSubjectInput: FC<EmailSubjectInput> = ({
  fieldName,
  formId,
  mt = 0,
}) => {
  return (
    <Flex alignItems='center' flex={1} mt={mt}>
      <Text as={'span'} color='gray.700' fontWeight={600} mr={1}>
        Subject:
      </Text>
      <FormInput
        name={fieldName}
        formId={formId}
        color='gray.500'
        height={5}
        fontSize='inherit'
        border='none'
        _hover={{ border: 'none !important' }}
        _active={{ border: 'none !important' }}
        _visited={{ border: 'none !important' }}
        _focusVisible={{ border: 'none !important' }}
        _focus={{ border: 'none !important' }}
      />
    </Flex>
  );
};
