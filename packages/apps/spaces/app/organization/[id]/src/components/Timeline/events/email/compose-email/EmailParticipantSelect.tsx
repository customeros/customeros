'use client';
import React, { FC } from 'react';

import { Flex } from '@ui/layout/Flex';
import { Text } from '@ui/typography/Text';
import { emailRegex } from '@organization/src/components/Timeline/events/email/utils';

import { EmailFormMultiCreatableSelect } from './EmailFormMultiCreatableSelect';
interface EmailParticipantSelect {
  formId: string;
  entryType: string;
  fieldName: string;
  autofocus: boolean;
}

export const EmailParticipantSelect: FC<EmailParticipantSelect> = ({
  entryType,
  fieldName,
  formId,
  autofocus = false,
}) => {
  return (
    <Flex
      alignItems='baseline'
      marginBottom={-1}
      marginTop={0}
      flex={1}
      overflow='visible'
    >
      <Text as={'span'} color='gray.700' fontWeight={600} mr={1}>
        {entryType}:
      </Text>
      <EmailFormMultiCreatableSelect
        autoFocus={autofocus}
        name={fieldName}
        formId={formId}
        placeholder='Enter name or email...'
        noOptionsMessage={() => null}
        allowCreateWhileLoading={false}
        formatCreateLabel={(input) => {
          return input;
        }}
        isValidNewOption={(input) => emailRegex.test(input)}
        getOptionLabel={(d) => {
          if (d?.__isNew__) {
            return `${d.label}`;
          }

          return `${d.label} - ${d.value}`;
        }}
      />
    </Flex>
  );
};
