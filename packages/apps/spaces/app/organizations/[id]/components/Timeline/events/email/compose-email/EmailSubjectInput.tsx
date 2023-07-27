'use client';
import React, { FC } from 'react';
import { Text } from '@ui/typography/Text';
import { Flex } from '@ui/layout/Flex';
import { FormInput } from '@ui/form/Input';

interface EmailSubjectInput {
  fieldName: string;
  formId: string;
}

export const EmailSubjectInput: FC<EmailSubjectInput> = ({
  fieldName,
  formId,
}) => {
  return (
    <Flex alignItems='center' flex={1}>
      <Text as={'span'} color='gray.700' fontWeight={600} mr={1}>
        Subject:
      </Text>
      <FormInput name={fieldName} formId={formId} color='gray.500' height={5} />
    </Flex>
  );
};
