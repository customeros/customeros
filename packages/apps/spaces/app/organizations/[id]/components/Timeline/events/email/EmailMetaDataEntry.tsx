'use client';
import React, { FC } from 'react';
import { Text } from '@ui/typography/Text';
import { Flex, Tooltip } from '@chakra-ui/react';
import { InteractionEventParticipant } from '@graphql/types';
import { FormInput } from '@ui/form/Input';
import { getEmailParticipantsNameAndEmail } from '@spaces/utils/getParticipantsName';

interface EmailMetaDataEntry {
  entryType: string;
  content: InteractionEventParticipant[] | string;
}

interface EmailMetaData {
  email: string | null;
  label: string;
}

export const EmailMetaDataEntry: FC<EmailMetaDataEntry> = ({
  entryType,
  content,
}) => {
  const data: boolean | Array<EmailMetaData> =
    typeof content !== 'string' && getEmailParticipantsNameAndEmail(content);

  return (
    <Flex>
      <Text as={'span'} color='#344054' fontWeight={600} mr={1}>
        {entryType}:
      </Text>

      {typeof content === 'string' && (
        <Text as={'span'} color='#667085'>
          {content}
        </Text>
      )}
      {typeof content !== 'string' &&
        !!data &&
        data.map((e, i) => {
          if (!e.label) {
            return (
              <Text
                mr={1}
                as={'span'}
                color='#667085'
                key={`email-participant-tag-${i}-${e.email}`}
              >
                {e.email}
                {i !== data.length - 1 && ','}
              </Text>
            );
          }
          return (
            <Tooltip
              key={`email-participant-tag-${e.label}-${e.email}`}
              label={e.email}
              aria-label={`${e.email}`}
              placement='top'
              zIndex={100}
            >
              <Text
                mr={1}
                as={'span'}
                color='#667085'
                key={`email-participant-tag-${i}-${e.email}`}
              >
                {e.label}
                {i !== data.length - 1 && ','}
              </Text>
            </Tooltip>
          );
        })}
    </Flex>
  );
};

interface EmailMetaDataEntryInput {
  entryType: string;
  formName: string;
  fieldName: string;
}

export const EmailMetaDataEntryInput: FC<EmailMetaDataEntryInput> = ({
  entryType,
  formName,
  fieldName,
}) => {
  return (
    <Flex alignItems='center'>
      <Text as={'span'} color='#344054' fontWeight={600} mr={1}>
        {entryType}:
      </Text>
      <FormInput
        name={fieldName}
        color='#667085'
        borderBottom='none'
        fontWeight={400}
        fontSize='md'
        mr={1}
        minWidth='100%'
        padding={0}
        height={5}
        formId={formName}
      />
    </Flex>
  );
};
